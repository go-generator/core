package io

import (
	"fmt"
	"github.com/go-generator/core"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

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
		err := SaveFile(fullPath, v.Content)
		if err != nil {
			return err
		}
	}
	return nil
}
func SaveFile(fullName string, content string) error {
	err := os.MkdirAll(filepath.Dir(fullName), os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fullName, []byte(content), os.ModePerm)
}
