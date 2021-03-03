package fileman

import (
	"io/fs"
	"io/ioutil"
)

//FileModel is a model with the needed data for template processing
type FileModel struct {
	Name  string
	IsDir bool
}

//FindFiles finds every file in the suplemented path and returns as FileModel
func FindFiles(path string) ([]FileModel, error) {
	var fileInfos []fs.FileInfo
	fileInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return []FileModel{}, err
	}
	var files []FileModel
	for _, fileInfo := range fileInfos {
		file := FileModel{Name: fileInfo.Name(), IsDir: fileInfo.IsDir()}
		files = append(files, file)
	}
	return files, nil
}
