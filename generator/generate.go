package generator

import (
	"errors"
	"fmt"
	"github.com/go-generator/core"
	"os"
	"strings"
	"text/template"
)

func Generate(project metadata.Project, templates map[string]interface{}, buildModel func(m metadata.Model, types map[string]string, env map[string]interface{}) map[string]interface{}, options...func(map[string]string) map[string]interface{}) ([]metadata.File, error) {
	var (
		outputFile []metadata.File
		err        error
	)
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
		v.File, err = parsing(v.File, m, "static_"+v.Name)
		if err != nil {
			return nil, fmt.Errorf("generating static file error: %w", err)
		}
		if s, ok := templates[v.Name].(string); ok {
			text, err1 := parsing(s, m, "static_"+v.Name)
			if err1 != nil {
				return nil, fmt.Errorf("generating static file content error: %w", err1)
			}
			outputFile = append(outputFile, metadata.File{Name: v.File, Content: text})
		}
	}
	for _, a := range project.Arrays {
		m := make(map[string]interface{}, 0)
		m["env"] = env
		m["collections"] = collections
		if str, ok := templates[a.Name].(string); ok {
			text, err2 := parsing(str, m, "array_"+a.Name)
			if err2 != nil {
				return nil, fmt.Errorf("generating model file error: %w", err2)
			}
			entityPath, err3 := generateFilePath(a.File, m)
			if err3 != nil {
				return nil, fmt.Errorf("generating file path error: %w", err3)
			}
			entityPath = strings.ReplaceAll(entityPath, "/", pathSeparator)
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
			if str, ok := templates[e.Name].(string); ok {
				text, err2 := parsing(str, v, "entity_"+e.Name)
				if err2 != nil {
					return nil, fmt.Errorf("generating model file error: %w", err2)
				}
				entityPath, err3 := generateFilePath(e.File, v)
				if err3 != nil {
					return nil, fmt.Errorf("generating file path error: %w", err3)
				}
				entityPath = strings.ReplaceAll(entityPath, "/", pathSeparator)
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
func parsing(tmpl string, m map[string]interface{}, templateName string) (string, error) {
	strBld := &strings.Builder{}
	tmp, err := template.New(templateName).Parse(tmpl)
	if err != nil {
		return "", err
	}
	err = tmp.Execute(strBld, m)
	if err != nil {
		return "", err
	}
	return strBld.String(), err
}
func generateFilePath(path string, m map[string]interface{}) (string, error) {
	strBld := strings.Builder{}
	tmp, err := template.New(path).Parse(path)
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
func InitProject(project metadata.Project, buildModel func(m metadata.Model, types map[string]string, env map[string]interface{}) map[string]interface{}, options...map[string]interface{}) []map[string]interface{} {
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
