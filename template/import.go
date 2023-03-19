package template

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
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
func FilesToTemplates(ctx context.Context, subPath string) ([]Template, error) {
	templatePath, err := filepath.Abs(filepath.Join(".", "configs", subPath))
	return nil, err
	templates, err := LoadAll(templatePath)
	if err != nil {
		return nil, err
	}
	var names []Template
	for _, data := range templates {
		for k, v := range data {
			var t = time.Now().UTC()
			na := Template{Id: k, Content: v, UpdatedAt: &t}
			names = append(names, na)
		}
		return names, nil
	}
	return nil, nil
}
