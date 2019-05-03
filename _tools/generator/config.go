package main

import (
	"encoding/json"
	"flag"
	"log"
	"path/filepath"
	"strings"
)

func parseFlag() {
	flag.StringVar(&configFile, `c`, ``, `-c conf.yaml`)

	//DBSettings
	flag.StringVar(&cfg.Username, `u`, `root`, `-u user`)
	flag.StringVar(&cfg.Password, `p`, ``, `-p password`)
	flag.StringVar(&cfg.Host, `h`, `localhost`, `-h host`)
	flag.StringVar(&cfg.Engine, `e`, `mysql`, `-e engine`)
	flag.StringVar(&cfg.Database, `d`, `blog`, `-d database`)
	flag.StringVar(&cfg.Charset, `charset`, `utf8`, `-charset utf8mb4`)
	flag.StringVar(&cfg.Prefix, `pre`, ``, `-pre prefix`)
	flag.StringVar(&cfg.Ignore, `ignore`, ``, `-ignore "^private_"`)
	flag.StringVar(&cfg.Match, `match`, ``, `-match "^publish_"`)
	flag.BoolVar(&cfg.NotGenerated, `notGenerated`, false, `-notGenerated false`)
	flag.StringVar(&cfg.Backup, `backup`, ``, `-backup "./install.0.sql"`)

	//DBSchema
	flag.StringVar(&cfg.SchemaConfig.ImportPath, `import`, ``, `-import github.com/webx-top/project/app/dbschema`)
	flag.StringVar(&cfg.SchemaConfig.PackageName, `pkg`, `dbschema`, `-pkg packageName`)
	flag.StringVar(&cfg.SchemaConfig.SaveDir, `o`, `dbschema`, `-o targetDir`)

	//Model
	flag.StringVar(&cfg.ModelConfig.BaseName, `mbase`, `Base`, `-mbase Base`)
	flag.StringVar(&cfg.ModelConfig.ImportPath, `mimport`, ``, `-mimport github.com/webx-top/project/app/model`)
	flag.StringVar(&cfg.ModelConfig.SaveDir, `mo`, ``, `-mo targetDir`)
	flag.StringVar(&cfg.ModelConfig.PackageName, `mpkg`, `model`, `-mpkg packageName`)

	//Postgres schema
	flag.StringVar(&cfg.Schema, `schema`, `public`, `-schema schemaName`)

	//Time
	flag.StringVar(&autoTime, `autoTime`, `{"update":{"*":["updated"]},"insert":{"*":["created"]}}`, `-autoTime <json-data>`)
	flag.Parse()
}

var cfg = &config{
	SchemaConfig: &SchemaConfig{},
	ModelConfig:  &ModelConfig{},
}
var configFile, autoTime string

type SchemaConfig struct {
	SaveDir     string `json:"saveDir"`
	ImportPath  string `json:"importPath"`
	PackageName string `json:"packageName"`
}

type ModelConfig struct {
	SaveDir     string `json:"saveDir"`
	ImportPath  string `json:"importPath"`
	PackageName string `json:"packageName"`
	BaseName    string `json:"baseName"`
}

type config struct {
	Username       string          `json:"username"`
	Password       string          `json:"password"`
	Host           string          `json:"host"`
	Engine         string          `json:"engine"`
	Database       string          `json:"database"`
	Charset        string          `json:"charset"`
	Prefix         string          `json:"prefix"`
	Ignore         string          `json:"ignore"`
	Match          string          `json:"match"`
	SchemaConfig   *SchemaConfig   `json:"schemaConfig"`
	ModelConfig    *ModelConfig    `json:"modelConfig"`
	Schema         string          `json:"schema"`
	NotGenerated   bool            `json:"notGenerated"`
	AutoTimeFields *AutoTimeFields `json:"autoTime"`
	Backup         string          `json:"backup"`
}

func (cfg *config) Check() {
	if len(cfg.ModelConfig.SaveDir) > 0 {
		cfg.ModelConfig.SaveDir = strings.TrimSuffix(cfg.ModelConfig.SaveDir, `/`)
		cfg.ModelConfig.SaveDir = strings.TrimSuffix(cfg.ModelConfig.SaveDir, `\`)
		if len(cfg.ModelConfig.PackageName) == 0 {
			cfg.ModelConfig.PackageName = filepath.Base(cfg.ModelConfig.SaveDir)
			switch cfg.ModelConfig.PackageName {
			case `.`, `/`, `\`:
				cfg.ModelConfig.PackageName = `model`
			}
		}
	}

	if cfg.AutoTimeFields == nil && len(autoTime) > 0 {
		cfg.AutoTimeFields = &AutoTimeFields{}
		cfg.AutoTimeFields.Parse(autoTime)
	}
}

type AutoTimeFields struct {
	//Update update操作时，某个字段自动设置为当前时间（map的键和值分别为表名称和字段名称。当表名称设置为“*”时，代表所有表中的这个字段）
	Update map[string][]string `json:"update"`

	//Insert insert操作时，某个字段自动设置为当前时间（map的键和值分别为表名称和字段名称。当表名称设置为“*”时，代表所有表中的这个字段）
	Insert map[string][]string `json:"insert"`
}

func (a *AutoTimeFields) Parse(autoTime string) {
	// JSON
	if (autoTime)[0] == '{' {
		err := json.Unmarshal([]byte(autoTime), a)
		if err != nil {
			log.Fatal(err)
		}
		return
	}
	// update(*:updated)/insert(*:created) 括号内的格式：<表1>:<字段1>,<字段2>,<...字段N>;<表2>:<字段1>,<字段2>,<...字段N>
	a.Update = make(map[string][]string)
	a.Insert = make(map[string][]string)
	for _, par := range strings.Split(autoTime, `/`) {
		par = strings.TrimSpace(par)
		switch {
		case strings.HasPrefix(par, `update(`):
			a.parseUpdateTime(par)

		case strings.HasPrefix(par, `insert(`):
			a.parseInsertTime(par)
		}
	}
}

func (a *AutoTimeFields) parseUpdateTime(par string) {
	par = strings.TrimPrefix(par, `update(`)
	par = strings.TrimSuffix(par, `)`)
	for _, item := range strings.Split(par, `;`) {
		t := strings.SplitN(item, `:`, 2)
		if len(t) < 2 {
			continue
		}
		t[0] = strings.TrimSpace(t[0])
		t[1] = strings.TrimSpace(t[1])
		if len(t[0]) == 0 || len(t[1]) == 0 {
			continue
		}
		if _, ok := a.Update[t[0]]; !ok {
			a.Update[t[0]] = []string{}
		}
		for _, field := range strings.Split(t[1], `,`) {
			field = strings.TrimSpace(field)
			if len(field) == 0 {
				continue
			}
			a.Update[t[0]] = append(a.Update[t[0]], field)
		}
	}
}

func (a *AutoTimeFields) parseInsertTime(par string) {
	par = strings.TrimPrefix(par, `insert(`)
	par = strings.TrimSuffix(par, `)`)
	for _, item := range strings.Split(par, `;`) {
		t := strings.SplitN(item, `:`, 2)
		if len(t) < 2 {
			continue
		}
		t[0] = strings.TrimSpace(t[0])
		t[1] = strings.TrimSpace(t[1])
		if len(t[0]) == 0 || len(t[1]) == 0 {
			continue
		}
		if _, ok := a.Insert[t[0]]; !ok {
			a.Insert[t[0]] = []string{}
		}
		for _, field := range strings.Split(t[1], `,`) {
			field = strings.TrimSpace(field)
			if len(field) == 0 {
				continue
			}
			a.Insert[t[0]] = append(a.Insert[t[0]], field)
		}
	}
}
