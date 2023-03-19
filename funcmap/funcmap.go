package funcmap

import (
	"strings"
	"text/template"

	"github.com/gertd/go-pluralize"
	st "github.com/go-generator/core/strings"
)

func MakeFuncMap() template.FuncMap {
	pluralize := pluralize.NewClient()
	funcMap := make(template.FuncMap, 0)
	funcMap["lower"] = strings.ToLower
	funcMap["upper"] = strings.ToUpper
	funcMap["snake"] = st.BuildSnakeName
	funcMap["unsnake"] = st.UnBuildSnakeName
	funcMap["plural"] = pluralize.Plural     //st.ToPlural
	funcMap["singular"] = pluralize.Singular //st.ToSingular
	funcMap["camel"] = st.ToCamelCase
	funcMap["pascal"] = st.ToPascalCase
	funcMap["go_driver"] = st.ImportDriver
	funcMap["go_mod_import"] = st.ImportGoMod
	return funcMap
}
