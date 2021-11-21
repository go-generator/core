package build

import (
	"github.com/go-generator/core"
	"github.com/go-generator/core/jstypes"
	st "github.com/go-generator/core/strings"
	"regexp"
	"strings"
)

func MergeMap(m map[string]interface{}, sub map[string]string) {
	for k, v := range sub {
		m[k] = v
	}
}
func BuildModel(m metadata.Model, types map[string]string, env map[string]interface{}) map[string]interface{} {
	var re = regexp.MustCompile(`date|datetime|time`)
	names := BuildNames(m.Name, st.ToPlural)
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
	collection["ts_date"] = ""
	collection["ts_number"] = ""
	collection["goId"] = "string"
	collection["javaId"] = "String"
	collection["csId"] = "string"
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
		hasTime := false
		hasNumber := false
		for _, f := range m.Fields {
			sub := make(map[string]interface{}, 0)
			tmp := BuildNames(f.Name)
			MergeMap(sub, tmp)
			x := f.Type
			if re.MatchString(x) {
				collection["ts_date"] = "DateRange, "
				hasTime = true
			}
			sub["simpleTypes"] = f.Type
			t, ok := types[f.Type]
			ut := f.Type
			if ok && len(t) > 0 {
				sub["type"] = t
				ut = t
			} else {
				sub["type"] = f.Type
			}
			jt, ok2 := jstypes.JSTypes[f.Type]
			if ok2 {
				if jt == "number" || jt == "integer" {
					hasNumber = true
					collection["ts_number"] = "NumberRange, "
				}
				sub["jstype"] = jt
			} else {
				sub["jstype"] = f.Type
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
			sub["maxlength"] = f.Length
			sub["env"] = env
			if re.MatchString(x) {
				sub["goFilterType"] = "*TimeRange"
				sub["tsFilterType"] = "Date|DateRange"
				sub["javaFilterType"] = "DateRange"
				sub["csFilterType"] = "DateTimeRange"
			} else if x == "float64" || x == "decimal" || x == "float64[]" || x == "decimal[]" {
				sub["goFilterType"] = "*NumberRange"
				sub["tsFilterType"] = "number|NumberRange"
				sub["javaFilterType"] = "NumberRange"
				sub["csFilterType"] = "NumberRange"
			} else if x == "int64" || x == "int64[]" {
				sub["goFilterType"] = "*Int64Range"
				sub["tsFilterType"] = "number|NumberRange"
				sub["javaFilterType"] = "Int64Range"
				sub["csFilterType"] = "Int64Range"
			} else if x == "float32" || x == "float32[]" {
				sub["goFilterType"] = "*NumberRange"
				sub["tsFilterType"] = "number|NumberRange"
				sub["javaFilterType"] = "FloatRange"
				sub["csFilterType"] = "FloatRange"
				sub["javaFilterType"] = "Int32Range"
				sub["csFilterType"] = "Int32Range"
			} else if x == "int32" || x == "int32[]" {
				sub["goFilterType"] = "*Int32Range"
				sub["tsFilterType"] = "number|NumberRange"
			} else {
				stp := sub["type"]
				sub["goFilterType"] = stp
				sub["tsFilterType"] = stp
				sub["javaFilterType"] = stp
				sub["csFilterType"] = stp
			}
			if f.Key {
				ck++
				g, c, p := buildGoGetId(f.Type)
				sub["bson"] = "-"
				collection["tsId"] = sub["type"]
				collection["javaId"] = sub["type"]
				collection["csId"] = sub["type"]
				collection["goId"] = sub["type"]
				collection["goGetId"] = g
				collection["goCheckId"] = c
				collection["goIdPrefix"] = p
				collection["goIdType"] = f.Type
				goIds = append(goIds, "{"+tmp["name"]+"}")
				tsIds = append(tsIds, ":"+tmp["name"])
			} else {
				sub["bson"] = f.Name + ",omitempty"
			}
			if ut == "float64" || ut == "decimal" || ut == "float32" {
				if f.Scale != nil && *f.Scale > 0 {
					sub["scale"] = *f.Scale
				}
				if f.Precision != nil && *f.Precision > 0 {
					sub["precision"] = *f.Precision
					if f.Scale != nil && *f.Scale > 0 {
						sub["maxlength"] = *f.Precision - *f.Scale + 1
					} else {
						sub["maxlength"] = *f.Precision - *f.Scale
					}
				}
			} else if ut == "int64" {
				sub["scale"] = 0
				sub["maxlength"] = 20
			} else if ut == "int32" {
				sub["scale"] = 0
				sub["maxlength"] = 9
			} else if ut == "int16" {
				sub["scale"] = 0
				sub["maxlength"] = 4
			} else if ut == "int8" {
				sub["scale"] = 0
				sub["maxlength"] = 2
			}
			fields = append(fields, sub)
		}
		collection["time"] = hasTime
		collection["number"] = hasNumber
		if hasTime && hasNumber {
			collection["tsFilterImport"] = "import { DateRange, Filter, NumberRange } from 'onecore';"
		} else if hasTime {
			collection["tsFilterImport"] = "import { DateRange, Filter } from 'onecore';"
		} else if hasNumber {
			collection["tsFilterImport"] = "import { Filter, NumberRange } from 'onecore';"
		} else {
			collection["tsFilterImport"] = "import { Filter } from 'onecore';"
		}
		if ck > 1 {
			collection["tsId"] = "" + names["Name"] + "Id"
			collection["javaId"] = "" + names["Name"] + "Id"
			collection["csId"] = "" + names["Name"] + "Id"
			collection["goId"] = "interface{}"
			collection["goBsonId"] = "-"
			collection["goGetId"] = "GetId"
			collection["goCheckId"] = "id != nil"
			collection["goIdPrefix"] = ""
			collection["goIdType"] = "map[string]interface{}"
			collection["go_id_url"] = strings.Join(goIds, "/")
			collection["ts_id_url"] = strings.Join(tsIds, "/")
		} else if ck == 1 {
			for _, sub := range fields {
				k, ok := sub["key"]
				if ok {
					k0, ok2 := k.(bool)
					if ok2 {
						if k0 {
							sub["bson"] = "_id,omitempty"
							break
						}
					}
				}
			}
		}
		collection["fields"] = fields
	}
	return collection
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
