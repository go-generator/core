package io

import (
	"archive/zip"
	"fmt"
	"github.com/go-generator/core"
	"os"
)

func ZipFiles(zipWriter *zip.Writer, rootDirectory string, files []core.File) error {
	for _, v := range files {
		fullPath := rootDirectory + string(os.PathSeparator) + v.Name
		zipFile, err := zipWriter.Create(fullPath)
		if err != nil {
			fmt.Println(err)
			return err
		}
		_, err = zipFile.Write([]byte(v.Content))
		if err != nil {
			return err
		}
	}
	return nil
}
