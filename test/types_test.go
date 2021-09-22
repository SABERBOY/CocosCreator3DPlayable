package test

import (
	"fmt"
	"github.com/Jecced/go-tools/src/fileutil"
	"testing"
)

func TestFileTypes(t *testing.T) {
	dir := "/Users/ankang/develop/projects/git/test/cocos_demo/c3d_demo_01/build/web-mobile"
	types := fileutil.FindAllFileTypes(dir)
	fmt.Println(types)
}
