package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-generator/core"
	"strings"
)

func ToString(files []metadata.File) (string, error) {
	file, err := json.MarshalIndent(files, "", " ")
	if err != nil {
		return "", err
	}
	return string(file), err
}
func ToOutput(path, directory string, input metadata.Project, templateMap map[string]string, buildModel func(m metadata.Model, types map[string]string, env map[string]interface{}) map[string]interface{}) (metadata.Output, error) {
	var output metadata.Output
	output.Path = path
	output.Directory = directory
	output.Path = strings.TrimSuffix(output.Path, "/")
	outputFiles, err := Generate(input, templateMap, buildModel)
	if err != nil {
		return output, fmt.Errorf("error generating data: %w", err)
	}
	output.Files = outputFiles
	return output, err
}
func GenerateFromFile(templateDir, projectName, projectMetadata string, loadProject func(string) (metadata.Project, error), loadTemplates func(string) (map[string]string, error), initEnv func(map[string]string, string) map[string]string, buildModel func(metadata.Model, map[string]string, map[string]interface{}) map[string]interface{}) (metadata.Output, error) {
	var output metadata.Output
	input, err := loadProject(projectMetadata)
	if err != nil {
		return output, err
	}
	input.Env = initEnv(input.Env, projectName)
	templateMap, err := loadTemplates(templateDir)
	if err != nil {
		return output, err
	}
	output, err = ToOutput("", projectName, input, templateMap, buildModel)
	if err != nil {
		return output, err
	}
	output.Directory = projectName
	return output, nil
}
func GenerateFromString(projectName, jsonInput string, initEnv func(map[string]string, string) map[string]string) error {
	input, err := DecodeProject([]byte(jsonInput), projectName, initEnv)
	if err != nil {
		return err
	}
	input.Env = initEnv(input.Env, projectName)
	return nil
}
func DecodeProject(byteValue []byte, projectName string, initEnv func(map[string]string, string) map[string]string, models...[]metadata.Model) (metadata.Project, error) {
	var input metadata.Project
	err := json.NewDecoder(bytes.NewBuffer(byteValue)).Decode(&input)
	if err != nil {
		return input, err
	}
	if initEnv != nil {
		input.Env = initEnv(input.Env, projectName)
	}
	if len(models) > 0 && models[0] != nil {
		input.Models = models[0]
	}
	return input, err
}
func LoadProject(projectTemplateName, projectName string, templates map[string]string, m []metadata.Model, initEnv func(map[string]string, string) map[string]string) (*metadata.Project, error) {
	if data, ok := templates[projectTemplateName]; ok {
		pr, err := DecodeProject([]byte(data), projectName, initEnv, m)
		if err != nil {
			return nil, err
		}
		return &pr, nil
	}
	return nil, nil
}
func ExportProject(container string, load func(string) (map[string]string, error), projectTemplateName, projectName string, m []metadata.Model, initEnv func(map[string]string, string) map[string]string) (*metadata.Project, error) {
	templates, err := load(container)
	if err != nil {
		return nil, err
	}
	return LoadProject(projectTemplateName, projectName, templates, m, initEnv)
}