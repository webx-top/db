package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Replace struct {
	Old      string
	New      string
	Regexp   *regexp.Regexp
	FileRule *regexp.Regexp
}

var replaces = []*Replace{
	&Replace{` // import "upper.io/db.v2"`, ``, nil, nil},
	&Replace{`"upper.io/db.v2"`, `"github.com/webx-top/db"`, nil, nil},
	&Replace{`"upper.io/db.v2/`, `"github.com/webx-top/db/`, nil, nil},
	&Replace{"",
		"${1}case `!=`, `<>`:\n\t\t\t\top = `$$ne`\n\t\t\t",
		regexp.MustCompile("([\\s]+op \\= \\`\\$gte\\`[\\s]+)"),
		regexp.MustCompile(`mongo[/\\]collection\.go$`),
	},
	&Replace{"conds[chunks[0]] = bson.M{op: value}",
		`
			if v, y := conds[chunks[0]]; y {
				if bsonM, ok := v.(bson.M); ok {
					bsonM[op] = value
					continue
				}
			}
			conds[chunks[0]] = bson.M{op: value}
`,
		nil,
		regexp.MustCompile(`mongo[/\\]collection\.go$`),
	},
	&Replace{``, ``,
		regexp.MustCompile(regexp.QuoteMeta(`if c.Database == "" {`) + "[\\s]+" + regexp.QuoteMeta(`return ""`) + "[\\s]+\\}"),
		regexp.MustCompile(`mysql[/\\]collection\.go$`),
	},
	&Replace{`if c.Database == "" {`, `if false {`,
		nil,
		regexp.MustCompile(`mysql[/\\]collection\.go$`),
	},
	&Replace{``,
		`
	if iter.Next() {
		var name sql.NullString
		err := iter.Scan(&name)
		return name.String, err
	}
`,
		regexp.MustCompile("[\\s]+if iter\\.Next\\(\\) \\{[\\s]+var name string[\\s]+err := iter\\.Scan\\(&name\\)[\\s]+return name, err[\\s]+\\}"),
		regexp.MustCompile(`mysql[/\\]database\.go$`),
	},

	&Replace{"connTimeout",
		"ConnTimeout",
		nil,
		regexp.MustCompile(`mongo[/\\]database\.go$`),
	},
}

func main() {
	root := filepath.Join(os.Getenv(`GOPATH`), `src`, `github.com/admpub/db`)
	save := filepath.Join(os.Getenv(`GOPATH`), `src`, `github.com/webx-top/db`)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if info.Name() == `_tools` || strings.HasPrefix(info.Name(), `.`) {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasPrefix(info.Name(), `.`) {
			return nil
		}
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		content := string(b)
		for _, re := range replaces {
			if re.FileRule != nil {
				if re.FileRule.MatchString(path) == false {
					continue
				}
				//panic(path)
			}
			if re.Regexp == nil {
				content = strings.Replace(content, re.Old, re.New, -1)
			} else {
				fmt.Printf("%#v\n", re.Regexp.FindAllString(content, -1))
				content = re.Regexp.ReplaceAllString(content, re.New)
			}
		}
		saveAs := strings.TrimPrefix(path, root)
		saveAs = filepath.Join(save, saveAs)
		err = os.MkdirAll(filepath.Dir(saveAs), os.ModePerm)
		if err == nil {
			file, err := os.Create(saveAs)
			if err == nil {
				_, err = file.WriteString(content)
			}
		}
		if err != nil {
			return err
		}
		fmt.Println(`Autofix ` + path + `.`)
		return nil
	})
	defer time.Sleep(5 * time.Minute)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(`Autofix complete.`)
}
