package io

import (
	"encoding/json"
	"fmt"
	"github.com/go-generator/core"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func SaveModels(models []metadata.Model, filePath string, notAppendExt...bool) error {
	data, err := json.MarshalIndent(&models, "", " ")
	if err != nil {
		return err
	}
	if !(len(notAppendExt) > 0 && notAppendExt[0]) {
		if filepath.Ext(filePath) != "json" {
			filePath += ".json"
		}
	}
	err = SaveFile(filePath, data)
	if err != nil {
		return err
	}
	return err
}
func SaveProject(projectStruct metadata.Project, filePath string, notAppendExt...bool) error {
	data, err := json.MarshalIndent(&projectStruct, "", " ")
	if err != nil {
		return err
	}
	if !(len(notAppendExt) > 0 && notAppendExt[0]) {
		if filepath.Ext(filePath) != "json" {
			filePath += ".json"
		}
	}
	err = SaveFile(filePath, data)
	if err != nil {
		return err
	}
	return err
}
func Save(directory string, output metadata.Output) error {
	err := SaveFiles(directory, output.Files)
	if err != nil {
		return fmt.Errorf("error writing files: %w", err)
	}
	return err
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
func SaveFiles(rootDirectory string, files []metadata.File) error {
	for _, v := range files {
		fullPath := rootDirectory + string(os.PathSeparator) + v.Name
		err := SaveContent(fullPath, v.Content)
		if err != nil {
			return err
		}
	}
	return nil
}
func SaveFile(fullName string, data []byte) error {
	err := os.MkdirAll(filepath.Dir(fullName), os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fullName, data, os.ModePerm)
}
func SaveContent(fullName string, content string) error {
	err := os.MkdirAll(filepath.Dir(fullName), os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fullName, []byte(content), os.ModePerm)
}
