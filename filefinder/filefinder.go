package filefinder

import (
	"io/ioutil"
)

//FindFiles finds every filename in the suplemented path
func FindFiles(path string) ([]string, error) {
	var paths []string
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return paths, err
	}
	for _, file := range files {
		paths = append(paths, file.Name())
	}
	return paths, nil
}
