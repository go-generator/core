package build

import (
	"github.com/go-generator/core"
	"github.com/go-generator/core/generator"
	"strings"
)

func BuildModel(m metadata.Model, types map[string]string, env map[string]interface{}) map[string]interface{} {
	names := generator.BuildNames(m.Name, ToPlural)
	collection := make(map[string]interface{}, 0)
	MergeMap(collection, names)
	table := m.Table
	src := m.Source
	if len(src) == 0 {
		src = m.Name
	}
	if len(table) == 0 {
		table = src
	}
	collection["table"] = table
	collection["source"] = src
	collection["tsId"] = "string"
	collection["goId"] = "string"
	collection["javaId"] = "String"
	collection["netId"] = "string"
	collection["goBsonId"] = "_id,omitempty"
	collection["goGetId"] = "GetRequiredParam"
	collection["goCheckId"] = "len(id) > 0"
	collection["goIdPrefix"] = ""
	collection["goIdType"] = "string"
	collection["env"] = env
	collection["go_id_url"] = "{id}"
	collection["ts_id_url"] = ":id"
	if len(m.Fields) > 0 {
		ck := 0
		goIds := make([]string, 0)
		tsIds := make([]string, 0)
		fields := make([]map[string]interface{}, 0)
		for _, f := range m.Fields {
			sub := make(map[string]interface{}, 0)
			tmp := generator.BuildNames(f.Name)
			MergeMap(sub, tmp)
			sub["simpleTypes"] = f.Type
			t, ok := types[f.Type]
			if ok && len(t) > 0 {
				sub["type"] = t
			} else {
				sub["type"] = f.Type
			}
			column := f.Column
			source := f.Source
			if len(source) == 0 {
				source = f.Name
			}
			if len(column) == 0 {
				column = source
			}
			sub["key"] = f.Key
			sub["source"] = source
			sub["column"] = column
			sub["length"] = f.Length
			sub["env"] = env
			if f.Key {
				ck++
				g, c, p := buildGoGetId(f.Type)
				collection["tsId"] = sub["type"]
				collection["javaId"] = sub["type"]
				collection["netId"] = sub["type"]
				collection["goId"] = sub["type"]
				collection["goGetId"] = g
				collection["goCheckId"] = c
				collection["goIdPrefix"] = p
				collection["goIdType"] = f.Type
				goIds = append(goIds, "{"+tmp["name"]+"}")
				tsIds = append(tsIds, ":"+tmp["name"])
			}
			fields = append(fields, sub)
		}
		if ck > 1 {
			collection["tsId"] = "any"
			collection["javaId"] = names["Table"] + "Id"
			collection["netId"] = names["Table"] + "Id"
			collection["goId"] = "interface{}"
			collection["goBsonId"] = "-"
			collection["goGetId"] = "GetId"
			collection["goCheckId"] = "id != nil"
			collection["goIdPrefix"] = ""
			collection["goIdType"] = "map[string]interface{}"
			collection["go_id_url"] = strings.Join(goIds, "/")
			collection["ts_id_url"] = strings.Join(tsIds, "/")
		}
		collection["fields"] = fields
	}
	return collection
}
func ToPlural(s string) string {
	return s + "s"
}
func MergeMap(m map[string]interface{}, sub map[string]string) {
	for k, v := range sub {
		m[k] = v
	}
}
func buildGoGetId(s string) (string, string, string) {
	if s == "int" {
		return "GetRequiredInt", "id != nil", "*"
	} else if s == "int32" {
		return "GetRequiredInt32", "id != nil", "*"
	} else if s == "int64" {
		return "GetRequiredInt64", "id != nil", "*"
	} else {
		return "GetRequiredParam", "len(id) > 0", ""
	}
}
