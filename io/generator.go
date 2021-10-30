package io

import (
	"fmt"
	iou "github.com/core-go/io"
	"github.com/go-generator/core"
	"os"
	"os/exec"
)

func ConvertFile(f metadata.File) iou.File {
	return iou.File{
		Name:    f.Name,
		Content: f.Content,
	}
}

func ConvertFiles(f []metadata.File) []iou.File {
	var out []iou.File
	for _, v := range f {
		out = append(out, ConvertFile(v))
	}
	return out
}

func SaveOutput(directory string, output metadata.Output) error {
	err := iou.SaveFiles(directory, ConvertFiles(output.Files))
	if err != nil {
		return fmt.Errorf("error writing files: %w", err)
	}
	return err
}

func ShellExecutor(program string, arguments []string) ([]byte, error) {
	cmd := exec.Command(program, arguments...)
	return cmd.Output()
}

func CreateDir(path string) (err error) {
	err = os.MkdirAll(path, 0644)
	if err != nil && os.IsNotExist(err) {
		return
	}
	return
}

//func FileDetailsToOutput(content model.FilesDetails, out *model.Output) {
//	var file metadata.File
//	for _, k := range content.Files {
//		out.OutFile = append(output.OutFile, WriteStruct(&k))
//		file.Name = content.Model + "/" + ToLower(k.Name) + ".go_bk"
//		file.Content = k.WriteFile.String()
//		out.Files = append(out.Files, file)
//	}
//}
