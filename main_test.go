package main

import (
	"io/ioutil"
	"os"
	"testing"
)

// 删除所有空文件夹
func TestDelEmptyDir(t *testing.T) {
	dir := "/Library/WebServer/Documents/game/test"

	DelEmptyDir(dir)

}
func DelEmptyDir(dir string) (bool, error) {
	list, err := ioutil.ReadDir(dir)
	if err != nil {
		return false, err
	}

	// 确定文件夹内部文件数量
	l := len(list)
	path := ""

	for _, info := range list {
		path = dir + string(os.PathSeparator) + info.Name()
		if !info.IsDir() && info.Name() == ".DS_Store" {
			err := os.Remove(path)
			if err != nil {
				continue
			}
			l--
			continue
		}
		if !info.IsDir() {
			continue
		}
		b, err := DelEmptyDir(path)
		if err != nil {
			continue
		}
		if b {
			err := os.Remove(path)
			if err != nil {
				continue
			}
			l--
		}
	}

	return l == 0, nil
}
