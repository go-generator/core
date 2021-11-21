package generator

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/go-generator/core"
	"github.com/go-generator/core/build"
	"github.com/go-generator/core/types"
)

func GenerateFiles(projectName, projectJson string, projectTemplate map[string]map[string]string, funcMap template.FuncMap, options ...map[string]map[string]string) ([]metadata.File, error) {
	prj, err := DecodeProject([]byte(projectJson), projectName, build.InitEnv)
	if err != nil {
		return nil, err
	}
	_, ok := prj.Env["go_module"]
	if ok && projectName != "" {
		prj.Env["go_module"] = projectName
	}
	if !(prj.Types != nil && len(prj.Types) > 0) {
		if len(options) > 0 && options[0] != nil {
			prj.Types = options[0][prj.Language]
		} else {
			prj.Types = types.Types[prj.Language]
		}
	}
	return Generate(prj, projectTemplate[prj.Language], funcMap, build.BuildModel)
}
func Generate(
	project metadata.Project,
	templates map[string]string,
	funcMap template.FuncMap,
	buildModel func(m metadata.Model, types map[string]string, env map[string]interface{}) map[string]interface{},
	options ...func(map[string]string) map[string]interface{},
) ([]metadata.File, error) {
	var outputFile []metadata.File
	var err error
	pathSeparator := string(os.PathSeparator)
	var parseEnv func(map[string]string) map[string]interface{}
	if len(options) > 0 && options[0] != nil {
		parseEnv = options[0]
	} else {
		parseEnv = ParseEnv
	}
	env := parseEnv(project.Env)
	collections := InitProject(project, buildModel, env)
	for _, v := range project.Statics {
		m := make(map[string]interface{}, 0)
		m["env"] = env
		v.File, err = parsing(v.File, m, "static_"+v.Name, funcMap)
		if err != nil {
			return nil, fmt.Errorf("generating static file error: %w", err)
		}
		if s, ok := templates[v.Name]; ok {
			text, err1 := parsing(s, m, "static_"+v.Name, funcMap)
			if err1 != nil {
				return nil, fmt.Errorf("generating static file content error: %w", err1)
			}
			if v.Replace {
				if strings.Contains(text, "{|") {
					text = strings.Replace(text, "{|", "{", -1)
				}
				if strings.Contains(text, "|}") {
					text = strings.Replace(text, "|}", "}", -1)
				}
			}
			outputFile = append(outputFile, metadata.File{Name: v.File, Content: text})
		}
	}
	for _, a := range project.Arrays {
		m := make(map[string]interface{}, 0)
		m["env"] = env
		m["collections"] = collections
		if str, ok := templates[a.Name]; ok {
			text, err2 := parsing(str, m, "array_"+a.Name, funcMap)
			if err2 != nil {
				return nil, fmt.Errorf("generating model file error: %w", err2)
			}
			entityPath, err3 := generateFilePath(a.File, m, funcMap)
			if err3 != nil {
				return nil, fmt.Errorf("generating file path error: %w", err3)
			}
			entityPath = strings.ReplaceAll(entityPath, "/", pathSeparator)
			if a.Replace {
				if strings.Contains(text, "{|") {
					text = strings.Replace(text, "{|", "{", -1)
				}
				if strings.Contains(text, "|}") {
					text = strings.Replace(text, "|}", "}", -1)
				}
			}
			outputFile = append(outputFile, metadata.File{
				Name:    entityPath,
				Content: text,
			})
		} else {
			return nil, errors.New("template must be string")
		}
	}
	for _, e := range project.Entities {
		for _, v := range collections {
			if str, ok := templates[e.Name]; ok {
				text, err2 := parsing(str, v, "entity_"+e.Name, funcMap)
				if err2 != nil {
					return nil, fmt.Errorf("generating model file error: %w", err2)
				}
				entityPath, err3 := generateFilePath(e.File, v, funcMap)
				if err3 != nil {
					return nil, fmt.Errorf("generating file path error: %w", err3)
				}
				entityPath = strings.ReplaceAll(entityPath, "/", pathSeparator)
				if e.Replace {
					if strings.Contains(text, "{|") {
						text = strings.Replace(text, "{|", "{", -1)
					}
					if strings.Contains(text, "|}") {
						text = strings.Replace(text, "|}", "}", -1)
					}
				}
				outputFile = append(outputFile, metadata.File{
					Name:    entityPath,
					Content: text,
				})
			} else {
				return nil, errors.New("template must be string")
			}
		}
	}
	return outputFile, err
}
func parsing(t string, m map[string]interface{}, name string, funcMap template.FuncMap) (string, error) {
	strBld := &strings.Builder{}
	tmp, err := template.New(name).Funcs(funcMap).Parse(t)
	if err != nil {
		return "", err
	}
	err = tmp.Execute(strBld, m)
	if err != nil {
		return "", err
	}
	return strBld.String(), err
}
func generateFilePath(path string, m map[string]interface{}, funcMap template.FuncMap) (string, error) {
	strBld := strings.Builder{}
	tmp, err := template.New(path).Funcs(funcMap).Parse(path)
	if err != nil {
		return "", err
	}
	err = tmp.Execute(&strBld, m)
	if err != nil {
		return "", err
	}
	filePath := strBld.String()
	filePath = strings.ReplaceAll(filePath, "/", string(os.PathSeparator))
	return filePath, err
}
func InitProject(project metadata.Project, buildModel func(m metadata.Model, types map[string]string, env map[string]interface{}) map[string]interface{}, options ...map[string]interface{}) []map[string]interface{} {
	var collections []map[string]interface{}
	var env map[string]interface{}
	if len(options) > 0 && options[0] != nil {
		env = options[0]
	} else {
		env = ParseEnv(project.Env)
	}
	for _, m := range project.Models {
		model := buildModel(m, project.Types, env)
		collections = append(collections, model)
	}
	return collections
}
func ParseEnv(env map[string]string) map[string]interface{} {
	res := make(map[string]interface{}, 0)
	res["layer"] = false
	for k, v := range env {
		if k == "layer" && v == "true" {
			res[k] = true
		} else {
			res[k] = v
		}
	}
	return res
}
