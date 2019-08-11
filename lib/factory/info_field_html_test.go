package factory

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
)

func TestFieldHTML(t *testing.T) {
	fields := map[string]*FieldInfo{
		"content":     {Name: "content", DataType: "text", Unsigned: false, PrimaryKey: false, AutoIncrement: false, Min: 0, Max: 0, Precision: 0, MaxSize: 2550, Options: []string{}, DefaultValue: "", Comment: "内容", GoType: "string", GoName: "Content"},
		"description": {Name: "description", DataType: "varchar", Unsigned: false, PrimaryKey: false, AutoIncrement: false, Min: 0, Max: 0, Precision: 0, MaxSize: 255, Options: []string{}, DefaultValue: "", Comment: "简介", GoType: "string", GoName: "Description"},
		"disabled":    {Name: "disabled", DataType: "enum", Unsigned: false, PrimaryKey: false, AutoIncrement: false, Min: 0, Max: 0, Precision: 0, MaxSize: 0, Options: []string{"Y", "N"}, DefaultValue: "N", Comment: "是否禁用", GoType: "string", GoName: "Disabled"},
		"sort":        {Name: "sort", DataType: "int", Unsigned: false, PrimaryKey: false, AutoIncrement: false, Min: -9.999999999e+09, Max: 9.999999999e+09, Precision: 0, MaxSize: 10, Options: []string{}, DefaultValue: "0", Comment: "排序", GoType: "int", GoName: "Sort"},
		"type":        {Name: "type", DataType: "enum", Unsigned: false, PrimaryKey: false, AutoIncrement: false, Min: 0, Max: 0, Precision: 0, MaxSize: 0, Options: []string{"newsid", "prodid", "text", "url", "html", "image", "video", "audio", "file", "json", "list"}, DefaultValue: "text", Comment: "值类型(list-以半角逗号分隔的值列表)", GoType: "string", GoName: "Type"},
	}
	options := echo.H{}
	r := fields[`description`].HTML(`value:description`, options)
	assert.Equal(
		t,
		`<input type="text" class="form-control" name="description" value="value:description" maxlength="255" />`,
		string(r),
	)
	r = fields[`content`].HTML(`content:123`, options)
	assert.Equal(
		t,
		`<textarea class="form-control" name="content"  maxlength="2550">content:123</textarea>`,
		string(r),
	)

	r = fields[`disabled`].HTML(`N`, options)
	assert.Equal(
		t,
		`<div class="radio radio-primary radio-inline">
		<input type="radio" value="Y" name="disabled" id="disabled-Y"> <label for="disabled-Y">Y</label>
	</div><div class="radio radio-primary radio-inline">
		<input type="radio" value="N" name="disabled" id="disabled-N" checked="checked"> <label for="disabled-N">N</label>
	</div>`,
		string(r),
	)

	r = fields[`sort`].HTML(`5000`, options)
	assert.Equal(
		t,
		`<input type="number" class="form-control" name="sort" value="5000" maxlength="10" max="9.999999999e&#43;09" min="-9.999999999e&#43;09" />`,
		string(r),
	)

	r = fields[`type`].HTML(`video`, options)
	//panic(`[` + string(r) + `]`)
	assert.Equal(
		t,
		`<div class="radio radio-primary radio-inline">
		<input type="radio" value="newsid" name="type" id="type-newsid"> <label for="type-newsid">newsid</label>
	</div><div class="radio radio-primary radio-inline">
		<input type="radio" value="prodid" name="type" id="type-prodid"> <label for="type-prodid">prodid</label>
	</div><div class="radio radio-primary radio-inline">
		<input type="radio" value="text" name="type" id="type-text"> <label for="type-text">text</label>
	</div><div class="radio radio-primary radio-inline">
		<input type="radio" value="url" name="type" id="type-url"> <label for="type-url">url</label>
	</div><div class="radio radio-primary radio-inline">
		<input type="radio" value="html" name="type" id="type-html"> <label for="type-html">html</label>
	</div><div class="radio radio-primary radio-inline">
		<input type="radio" value="image" name="type" id="type-image"> <label for="type-image">image</label>
	</div><div class="radio radio-primary radio-inline">
		<input type="radio" value="video" name="type" id="type-video" checked="checked"> <label for="type-video">video</label>
	</div><div class="radio radio-primary radio-inline">
		<input type="radio" value="audio" name="type" id="type-audio"> <label for="type-audio">audio</label>
	</div><div class="radio radio-primary radio-inline">
		<input type="radio" value="file" name="type" id="type-file"> <label for="type-file">file</label>
	</div><div class="radio radio-primary radio-inline">
		<input type="radio" value="json" name="type" id="type-json"> <label for="type-json">json</label>
	</div><div class="radio radio-primary radio-inline">
		<input type="radio" value="list" name="type" id="type-list"> <label for="type-list">list</label>
	</div>`,
		string(r),
	)
}
