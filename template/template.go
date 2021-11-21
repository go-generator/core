package template

import (
	"strings"
	"text/template"

	st "github.com/go-generator/core/strings"
)

func MakeFuncMap() template.FuncMap {
	funcMap := make(template.FuncMap, 0)
	funcMap["lower"] = strings.ToLower
	funcMap["upper"] = strings.ToUpper
	funcMap["snake"] = st.BuildSnakeName
	funcMap["unsnake"] = st.UnBuildSnakeName
	funcMap["plural"] = st.ToPlural
	funcMap["singular"] = st.ToSingular
	funcMap["camel"] = st.ToCamelCase
	funcMap["pascal"] = st.ToPascalCase
	funcMap["go_driver"] = st.ImportDriver
	return funcMap
}
