package db

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Load(fileName string, config *Database) error {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	err2 := yaml.Unmarshal(file, config)
	if err2 == nil {
		l0 := len(config.DatabaseChangeLog)
		for i := 0; i < l0; i++ {
			l1 := len(config.DatabaseChangeLog[i].Changes)
			for j := 0; j < l1; j++ {
				l2 := len(config.DatabaseChangeLog[i].Changes[j].CreateTable.Columns)
				for k := 0; k < l2; k++ {
					config.DatabaseChangeLog[i].Changes[j].CreateTable.Columns[k].No = k
				}
			}
		}
	}
	return err2
}

func Save(fullName string, content string) error {
	err := os.MkdirAll(filepath.Dir(fullName), os.ModePerm)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fullName, []byte(content), os.ModePerm)
}
