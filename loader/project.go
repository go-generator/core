package loader

import (
	"bytes"
	"encoding/json"
	"github.com/core-go/io"
	"github.com/go-generator/core"
	"io/ioutil"
	"path/filepath"
)

func LoadProject(filename string) (metadata.Project, error) {
	var input metadata.Project
	byteValue, err := ioutil.ReadFile(filename)
	if err != nil {
		return input, err
	}
	err = json.NewDecoder(bytes.NewBuffer(byteValue)).Decode(&input)
	if err != nil {
		return input, err
	}
	return input, err
}

func LoadProjects(directory string) (map[string]metadata.Project, error) { // map[string]metadata.Project ---> "project name" : metadata project
	projects := make(map[string]metadata.Project)
	names, err := io.ListFileNames(directory)
	if err != nil {
		return nil, err
	}
	for _, n := range names {
		proj, err1 := LoadProject(filepath.Join(directory, n))
		if err1 != nil {
			return nil, err1
		}
		projects[n] = proj
	}
	return projects, err
}
