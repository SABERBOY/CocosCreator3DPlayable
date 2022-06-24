package engine

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Jecced/go-tools/src/ak"
	"github.com/Jecced/go-tools/src/fileutil"
	"github.com/Jecced/go-tools/src/strutil"
)

var (
	out     = ""
	content *string
)

func Mix(dir string, html *string) {
	out = dir
	content = html
	register, noJsList := GetAllSystemRegister(out)
	jsMap := map[string]string{}
	for _, jsPath := range register {
		mainIndexJs := strings.Replace(jsPath, out+ak.PS, "", -1)
		jsPath = filepath.ToSlash(jsPath)
		fileName := path.Base(jsPath)
		value, key := jsMap[fileName]
		if mainIndexJs == "index.js" || fileName != "index.js" {
			if key {
				println("fileName:" + value)
			} else {
				jsMap[fileName] = jsPath
			}
		}
	}
	FixWasm(jsMap)
	InsertRegisterJs(register, content, jsMap, noJsList)
	InsertImportMap(content)
	FindAndMixScript(content)
	FindAndMixLink(html)
}

// GetAllSystemRegister 获取所有和System.register相关的js脚本
func GetAllSystemRegister(dir string) ([]string, []string) {
	jsList, err := fileutil.GetFilesBySuffix(dir, ".js")
	if err != nil {
		return nil, nil
	}
	var out []string
	var outNoSys []string
	for _, js := range jsList {
		text, err := fileutil.ReadText(js)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if strings.HasPrefix(text, "System.register(") {
			out = append(out, js)
		} else {
			//E:\Projects\CocosProjects\WebGameDemo\build\PS\src\assets\script\libs\lz-string.min.js
			if strings.Contains(filepath.ToSlash(js), "src/assets/script/libs/") {
				outNoSys = append(outNoSys, js)
			}
		}
	}
	return out, outNoSys
}

// WASM处理
func FixWasm(jsMap map[string]string) {
	for s := range jsMap {
		if strings.Contains(s, "instantiated") {
			content, _ := fileutil.ReadBytes(jsMap[s])
			pat := `fetch\(e\).+assign\(([A-Z]+),n\),t\(\)\}\),n\)\}\(e\)\}\),n\)\}\),n\)`
			re := regexp.MustCompile(pat)
			data := re.ReplaceAll(content, []byte(`let url="cocos-js/"+(e);let wa=base64DecToArr(window.res[url]);WebAssembly.instantiate(wa,o).then((function(e){var n=e.instance.exports;Object.assign($1,n),t()}),n)`))
			fileutil.WriteData(data, jsMap[s])
		}

		if strings.Contains(s, "bullet.wasm") {
			content, _ := fileutil.ReadBytes(jsMap[s])
			pat := `new URL\((".+"),t\.meta\.url\)\.href`
			re := regexp.MustCompile(pat)
			data := re.ReplaceAll(content, []byte("$1"))
			fmt.Println("bullet.wasm:", string(data))
			fileutil.WriteData(data, jsMap[s])
		}
	}
}

func InsertRegisterJs(list []string, html *string, jsMap map[string]string, jsList []string) {
	for _, i2 := range list {
		i2 = strings.Replace(i2, out+ak.PS, "", -1)
		i2 = strings.Replace(i2, "\\", "/", -1)
		pickUpName := pickUpJSFiles(i2)
		insert(html, i2, pickUpName)
	}
	for _, i2 := range jsList {
		i2 = strings.Replace(i2, out+ak.PS, "", -1)
		i2 = strings.Replace(i2, "\\", "/", -1)
		jsFile := out + ak.PS + i2
		content, _ := fileutil.ReadText(jsFile)
		insert := "\n" + `<script>` + content + `</script>` + "\n"
		strutil.InsertString(html, insert, "<!-- SCRIPT -->")
		fmt.Println("libs file:", jsFile)
	}
}

func pickUpJSFiles(file string) string {
	if file == "application.js" {
		return "application.js"
	}
	if file == "index.js" {
		return "index.js"
	}
	if strings.Contains(file, "assets/main/index.js") || strings.Contains(file, "assets/resources/index.js") {
		return file
	}
	if strings.Contains(file, "cocos-js") {
		if strings.Contains(file, "cc.js") {
			return "cc"
		}
		if strings.Contains(file, "wait-for-ammo-instantiation.js") {
			return "wait-for-ammo-instantiation.js"
		}
		baseFile := path.Base(file)
		return "./" + baseFile
	}

	return "./" + file
}

func insert(html *string, file string, pickUpName string) {
	fmt.Printf("insert js file:%s\n", file)
	jsFile := out + ak.PS + file
	content, _ := fileutil.ReadText(jsFile)
	st := strings.Index(content, "(")
	if strings.Contains(file, "assets/main/index.js") || strings.Contains(file, "assets/resources/index.js") {
		// content = content[:st+1] /*+ "\"" + file + "\","*/ + content[st+1:]
	} else {
		/*slashJsPath := filepath.ToSlash(file)
		//fileName := path.Base(slashJsPath)
		filenameall := path.Base(slashJsPath)
		filesuffix := path.Ext(slashJsPath)
		fileName := filenameall[0 : len(filenameall)-len(filesuffix)]
		content = content[:st+1] + "\"" + fileName + "\"," + content[st+1:]*/
		content = content[:st+1] + "\"" + pickUpName + "\"," + content[st+1:]
	}
	/*	for s := range jsMap {
		var match = "./" + s
		if strings.Contains(content, match) {
			filenameall := path.Base(s)
			filesuffix := path.Ext(s)
			pick := filenameall[0 : len(filenameall)-len(filesuffix)]
			content = strings.Replace(content, match, pick, -1)
		}
	}*/
	insert := "\n" + `<script>` + content + `</script>` + "\n"
	strutil.InsertString(html, insert, "<!-- SCRIPT -->")
	fmt.Println("RegisterJs file:", jsFile)
	//_ = os.Remove(jsFile)
}

// InsertImportMap 插入import-map.json到html内
func InsertImportMap(html *string) {
	file := out + ak.PS + "src/import-map.json"
	text, _ := fileutil.ReadText(file)
	insert := "<script> var importMapJson=" + text + "</script>\n"
	strutil.InsertString(html, insert, "<!-- JSON IMPORT-MAP -->")
	//_ = os.Remove(file)
}

// FindAndMixScript 找到所有的script标签, 用src中的内容替换方块
func FindAndMixScript(html *string) {
	scripts := strutil.MatchString(*html, "<script src=", "</script>", true)
	for _, script := range scripts {
		ReplaceScriptBlock(html, script)
	}
}

func ReplaceScriptBlock(html *string, script string) {
	src := strutil.MatchStringFirst(script, "src=\"", "\"", false)
	path := out + ak.PS + src
	content, _ := fileutil.ReadText(path)
	insert := fmt.Sprintf(
		"\n<script>\n%s\n</script>\n",
		content,
	)
	*html = strings.Replace(*html, script, insert, -1)
	fmt.Println("file:", path)
	//os.Remove(path)
}

// FindAndMixLink 找到所有link标签, 替换成stype标签
func FindAndMixLink(html *string) {
	links := strutil.MatchString(*html, "  <link rel=\"", ".css\"/>", true)
	for _, link := range links {
		ReplaceStyleBlock(html, link)
	}
}

func ReplaceStyleBlock(html *string, script string) {
	src := strutil.MatchStringFirst(script, "href=\"", "\"", false)
	path := out + ak.PS + src
	content, _ := fileutil.ReadText(path)
	insert := fmt.Sprintf(
		"\n<style>\n%s\n</style>\n",
		content,
	)
	*html = strings.Replace(*html, script, insert, -1)
	fmt.Println("style file:", path)
	//os.Remove(path)
}
