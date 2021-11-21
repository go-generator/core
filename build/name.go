package build

import (
	"strings"

	st "github.com/go-generator/core/strings"
	"github.com/stoewer/go-strcase"
)

func BuildNames(name string, options ...func(string) string) map[string]string {
	var toPlural func(string) string
	if len(options) > 0 {
		toPlural = options[0]
	}
	n := make(map[string]string)
	var raw string
	if !strings.Contains(name, "_") {
		raw = strcase.SnakeCase(name) //st.BuildSnakeName(name)
	} else {
		raw = strings.ToLower(name)
		name = strcase.LowerCamelCase(raw) //st.UnBuildSnakeName(raw)
	}
	path := strcase.KebabCase(raw) //strings.Replace(raw, "_", "-", -1)
	n = map[string]string{
		"raw":      raw,
		"path":     path,
		"name":     st.ToCamelCase(name),
		"Name":     st.ToPascalCase(name),
		"NAME":     strings.ToUpper(name),
		"constant": strings.ToUpper(raw),
		"lower":    strings.ToLower(name),
	}
	if toPlural == nil {
		return n
	}
	raws := toPlural(raw)
	paths := strings.Replace(raws, "_", "-", -1)
	names := strcase.LowerCamelCase(raws)
	n["raws"] = raws
	n["paths"] = paths
	n["names"] = st.ToCamelCase(names)
	n["Names"] = st.ToPascalCase(names)
	n["NAMES"] = strings.ToUpper(names)
	n["constants"] = strings.ToUpper(raws)
	n["lowers"] = strings.ToLower(names)
	return n
}
func InitEnv(env map[string]string, projectName string) map[string]string {
	init, ok := env["init"]
	if ok {
		if init == "true" {
			outMap := make(map[string]string)
			for k, v := range env {
				tmp := buildEnvNames(k, v)
				for k1, v1 := range tmp {
					outMap[k1] = v1
				}
			}
			tmp := buildProjectName(projectName)
			for k1, v1 := range tmp {
				outMap[k1] = v1
			}
			return outMap
		}
	}
	env["project"] = projectName
	return env
}
func buildProjectName(name string) map[string]string {
	var raw string
	if !strings.Contains(name, "_") {
		raw = st.BuildSnakeName(name)
	} else {
		raw = strings.ToLower(name)
		name = st.UnBuildSnakeName(raw)
	}
	return map[string]string{
		"project_raw":      raw,
		"project":          st.ToCamelCase(name),
		"Project":          st.ToPascalCase(name),
		"project_lower":    strings.ToLower(name),
		"project_name":     st.BuildSnakeName(name),
		"project_constant": strings.ToUpper(raw),
	}
}

func buildEnvNames(name, v string) map[string]string {
	names := map[string]string{
		name:            v,
		name + "_name":  st.ToCamelCase(v),
		name + "_Name":  st.ToPascalCase(v),
		name + "_NAME":  strings.ToUpper(v),
		name + "_lower": strings.ToLower(v),
	}
	return names
}
