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
	entities, collections := InitProject(project, buildModel, env)
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
		var cols []map[string]interface{}
		for i := range collections {
			if InCollection(project.Collection, collections[i]["Name"].(string)) {
				cols = append(cols, collections[i])
			}
		}
		m["collections"] = cols
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
		for _, v := range entities {
			count := 0
			for _, k := range v.Model.Fields {
				if k.Key {
					count++
				}
			}
			if count > 1 {
				if str, ok := templates[e.Name]; ok {
					text, err2 := parsing(str, v.Params, "entity_"+e.Name, funcMap)
					if err2 != nil {
						return nil, fmt.Errorf("generating model file error: %w", err2)
					}
					entityPath, err3 := generateFilePath(e.File, v.Params, funcMap)
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
			} else {
				if e.Model || InCollection(project.Collection, v.Model.Name) && e.Name != "id" {
					if str, ok := templates[e.Name]; ok {
						text, err2 := parsing(str, v.Params, "entity_"+e.Name, funcMap)
						if err2 != nil {
							return nil, fmt.Errorf("generating model file error: %w", err2)
						}
						entityPath, err3 := generateFilePath(e.File, v.Params, funcMap)
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
		}
	}
	return outputFile, err
}

func InCollection(collection []string, name string) bool {
	for _, c := range collection {
		if strings.Compare(strings.ToLower(name), strings.ToLower(c)) == 0 {
			return true
		}
	}
	return false
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

type Entity struct {
	Model  metadata.Model         `mapstructure:"model" json:"model,omitempty" gorm:"column:model" bson:"model,omitempty" dynamodbav:"model,omitempty" firestore:"model,omitempty"`
	Params map[string]interface{} `mapstructure:"params" json:"params,omitempty" gorm:"column:params" bson:"params,omitempty" dynamodbav:"params,omitempty" firestore:"params,omitempty"`
}

func InitProject(project metadata.Project, buildModel func(m metadata.Model, types map[string]string, env map[string]interface{}) map[string]interface{}, options ...map[string]interface{}) ([]Entity, []map[string]interface{}) {
	var entities []Entity
	var collections []map[string]interface{}
	var env map[string]interface{}
	if len(options) > 0 && options[0] != nil {
		env = options[0]
	} else {
		env = ParseEnv(project.Env)
	}
	if project.Models == nil || len(project.Models) == 0 {
		if project.Collection != nil && len(project.Collection) > 0 {
			var models []metadata.Model
			for _, c := range project.Collection {
				m := metadata.Model{Name: c}
				models = append(models, m)
			}
			project.Models = models
		}
	}
	for _, m := range project.Models {
		count := 0
		for _, k := range m.Fields {
			if k.Key {
				count++
			}
		} //check composite key
		model := buildModel(m, project.Types, env)
		entity := Entity{Model: m, Params: model}
		hasTime := HasTime(project.Models, m.Name)
		if _, ok := project.Env["layer"]; !ok {
			for _, f := range m.Arrays {
				if HasTime(project.Models, f.Ref) {
					hasTime = true
					break
				}
			}
		}
		if m.Ones != nil && len(m.Ones) > 0 {
			var ones []map[string]interface{}
			for _, c := range m.Ones {
				s := GetModel(project.Models, c.Model)
				if s != nil {
					sub := buildModel(*s, project.Types, env)
					ones = append(ones, sub)
				}
			}
			model["ones"] = ones
		}
		if m.Models != nil && len(m.Models) > 0 {
			var models []map[string]interface{}
			for _, c := range m.Models {
				s := GetModel(project.Models, c.Model)
				if s != nil {
					sub := buildModel(*s, project.Types, env)
					models = append(models, sub)
				}
			}
			model["models"] = models
			if m.Arrays == nil || len(m.Arrays) == 0 {
				model["leaf"] = true
			}
		}
		if m.Arrays != nil && len(m.Arrays) > 0 {
			var arrays []map[string]interface{}
			for _, c := range m.Arrays {
				s := GetModel(project.Models, c.Model)
				if s != nil {
					sub := buildModel(*s, project.Types, env)
					if s.Arrays == nil || len(s.Arrays) == 0 {
						sub["leaf"] = true
					}
					sub["parent"] = m.Name
					sub["link"] = c.Fields
					arrays = append(arrays, sub)
				}
			}
			model["arrays"] = arrays
		}
		model["time"] = hasTime
		if count > 1 {
			model["composite"] = true
		} else {
			model["composite"] = false
		}
		entities = append(entities, entity)
		collections = append(collections, model)
	}
	return entities, collections
}

func HasTime(models []metadata.Model, name string) bool {
	for i := range models {
		if models[i].Table == name || models[i].Name == name {
			for _, f := range models[i].Fields {
				if strings.Contains(f.Type, "date") ||
					strings.Contains(f.Type, "time") {
					return true
				}
			}
		}
	}
	return false
}

func GetModel(models []metadata.Model, name string) *metadata.Model {
	for _, m := range models {
		if m.Name == name {
			return &m
		}
	}
	return nil
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
