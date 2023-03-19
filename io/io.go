package io

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-generator/core"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func List(path string) ([]string, error) {
	var names []string
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	folder, err := ioutil.ReadDir(absPath)
	if err != nil {
		return names, err
	}
	for _, tmpl := range folder {
		names = append(names, tmpl.Name())
	}
	return names, nil
}

func Load(directory string) (map[string]string, error) {
	tm := make(map[string]string, 0)
	names, er1 := List(directory)
	if er1 != nil {
		return nil, er1
	}
	for _, name := range names {
		content, er2 := ioutil.ReadFile(directory + string(os.PathSeparator) + name)
		if er2 != nil {
			return nil, er2
		}
		tm[name] = string(content)
	}
	return tm, nil
}

func LoadAll(directory string) (map[string]map[string]string, error) {
	templates := make(map[string]map[string]string)
	folders, err := List(directory)
	if err != nil {
		return nil, err
	}
	for _, folder := range folders {
		names, err := List(filepath.Join(directory, folder))
		if err != nil {
			return nil, err
		}
		tm := make(map[string]string, 0)
		for _, n := range names {
			content, err := ioutil.ReadFile(directory + string(os.PathSeparator) + folder + string(os.PathSeparator) + n)
			if err != nil {
				return nil, err
			}
			tm[n] = string(content)
		}
		templates[folder] = tm
	}
	return templates, err
}

func IsValidPath(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func Exec(program string, arguments []string) ([]byte, error) {
	cmd := exec.Command(program, arguments...)
	return cmd.Output()
}

func MkDir(path string) (err error) {
	err = os.MkdirAll(path, 0644)
	if err != nil && os.IsNotExist(err) {
		return
	}
	return
}

func SaveContent(fullName string, content string) error {
	err := os.MkdirAll(filepath.Dir(fullName), os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fullName, []byte(content), os.ModePerm)
}

func Save(fullName string, data []byte) error {
	err := os.MkdirAll(filepath.Dir(fullName), os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fullName, data, os.ModePerm)
}

func SaveFiles(rootDirectory string, files []core.File) error {
	for _, v := range files {
		fullPath := rootDirectory + string(os.PathSeparator) + v.Name
		err := SaveContent(fullPath, v.Content)
		if err != nil {
			return err
		}
	}
	return nil
}

func SaveModels(filePath string, models []core.Model, notAppendExt ...bool) error {
	data, err := json.MarshalIndent(&models, "", " ")
	if err != nil {
		return err
	}
	if !(len(notAppendExt) > 0 && notAppendExt[0]) {
		if filepath.Ext(filePath) != "json" {
			filePath += ".json"
		}
	}
	err = Save(filePath, data)
	if err != nil {
		return err
	}
	return err
}

func SaveProject(filePath string, projectStruct core.Project, notAppendExt ...bool) error {
	data, err := json.MarshalIndent(&projectStruct, "", " ")
	if err != nil {
		return err
	}
	if !(len(notAppendExt) > 0 && notAppendExt[0]) {
		if filepath.Ext(filePath) != "json" {
			filePath += ".json"
		}
	}
	err = Save(filePath, data)
	if err != nil {
		return err
	}
	return err
}

func SaveOutput(directory string, output core.Output) error {
	err := SaveFiles(directory, output.Files)
	if err != nil {
		return fmt.Errorf("error writing files: %w", err)
	}
	return err
}

func SaveConfig(filename string, config core.Config) error {
	var data bytes.Buffer
	err := yaml.NewEncoder(&data).Encode(&config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data.Bytes(), os.ModePerm)
}
