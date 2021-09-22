package res

import (
	"fmt"
	"github.com/Jecced/go-tools/src/ak"
	"github.com/Jecced/go-tools/src/commutil"
	"github.com/Jecced/go-tools/src/fileutil"
	"github.com/Jecced/go-tools/src/strutil"
	"strings"
)

var (
	out     = ""
	content *string
	resp    = "window.res = {}\n"
)

func Mix(dir string, html *string) {
	out = dir
	content = html

	text := []string{
		".json",
		".txt",
	}

	for _, t := range text {
		TextMix(t)
	}

	bin := []string{
		".bin",
		".png",
		".jpg",
		".mp3",
		".wasm",
		".ttf",
		".cconb",
	}

	for _, s := range bin {
		BinaryMix(s)
	}

	insert := fmt.Sprintf(
		"\n<script>%s</script>\n",
		resp,
	)
	strutil.InsertString(content, insert, "<!-- RESOURCE -->")
}

func BinaryMix(suffix string) {
	files, _ := fileutil.GetFilesBySuffix(out, suffix)
	for _, file := range files {
		dealBinary(file)
	}
}
func dealBinary(file string) {
	key := strings.Replace(file, out, "", -1)
	if strings.HasPrefix(key, ak.PS) {
		key = key[1:]
	}
	bytes, _ := fileutil.ReadBytes(file)
	base64 := commutil.EncodeBase64(&bytes)
	z := fmt.Sprintf(
		"window.res[\"%s\"]=\"%s\"",
		strings.Replace(key, "\\", "/", -1),
		base64,
	)
	resp = resp + "\n" + z + ""
	//os.Remove(file)
	fmt.Println("删除", file)
}

func TextMix(suffix string) {
	files, _ := fileutil.GetFilesBySuffix(out, suffix)
	for _, file := range files {
		dealText(file)
	}
}

func dealText(file string) {
	key := strings.Replace(file, out, "", -1)
	if strings.HasPrefix(key, ak.PS) {
		key = key[1:]
	}
	c, _ := fileutil.ReadText(file)
	// res/import/90/90115bf8-2d9c-444e-919c-50b73e73c5c9.bin
	z := fmt.Sprintf(
		"window.res[\"%s\"]=%s",
		strings.Replace(key, "\\", "/", -1),
		c,
	)
	resp = resp + "\n" + z + ""
	//os.Remove(file)
	fmt.Println("删除", file)
}
