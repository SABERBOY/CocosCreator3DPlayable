package cocos

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

type TransformResult struct {
	Path     string
	Alias    string
	Content  string
	Warnings []string
}

var (
	appendScriptPattern = regexp.MustCompile(`(?s)var err;.*?appendChild\(script\)`)
	settingsPattern     = regexp.MustCompile(`(?s)var requestSettings.*?requestSettings\(\);`)
	topImportPattern    = regexp.MustCompile(`(?s)(return topLevelImport.+?cc.+?\n.*?\))`)
	wasmFetchPattern    = regexp.MustCompile(`fetch\(e\).+assign\(([A-Z]+),n\),t\(\)\}\),n\)\}\(e\)\}\),n\)\}\),n\)`)
	bulletURLPattern    = regexp.MustCompile(`new URL\((".+?"),t\.meta\.url\)\.href`)
)

func TransformSystemRegister(relPath string, content []byte) TransformResult {
	relPath = strings.ReplaceAll(relPath, "\\", "/")
	text := string(content)
	var warnings []string

	switch {
	case relPath == "index.js":
		text = replaceWithWarning(text, appendScriptPattern, "resolve()", "index.js script loader replacement", &warnings)
		text = strings.ReplaceAll(text, "./application.js", "application.js")
	case relPath == "application.js":
		text = replaceWithWarning(text, settingsPattern, `if (settings.startsWith("./")) {
            settings = settings.substring(2);
          }
          window._CCSettings = window.res[settings];
          resolve();`, "application.js settings loader replacement", &warnings)
		text = regexp.MustCompile(`wait-for-ammo-instantiation(\.js)*`).ReplaceAllString(text, "wait-for-ammo-instantiation.js")
		text = replaceWithWarning(text, topImportPattern, `$1.then(function (engine) {
				cc = engine;
				return regLoading();
				})`, "application.js regLoading hook", &warnings)
	case strings.Contains(relPath, "instantiated"):
		text = replaceWithWarning(text, wasmFetchPattern, `let url="cocos-js/"+(e);let wa=base64DecToArr(window.res[url]);WebAssembly.instantiate(wa,o).then((function(e){var n=e.instance.exports;Object.assign($1,n),t()}),n)`, "instantiated wasm fetch replacement", &warnings)
	case strings.Contains(relPath, "bullet.wasm"):
		text = replaceWithWarning(text, bulletURLPattern, "$1", "bullet wasm URL replacement", &warnings)
	}

	alias := ModuleAlias(relPath)
	if alias != "" && !hasNamedRegister(text) {
		text = nameRegister(text, alias, &warnings)
	}
	return TransformResult{Path: relPath, Alias: alias, Content: text, Warnings: warnings}
}

func ModuleAlias(relPath string) string {
	relPath = strings.ReplaceAll(relPath, "\\", "/")
	switch {
	case relPath == "application.js":
		return "application.js"
	case relPath == "index.js":
		return "index.js"
	case strings.Contains(relPath, "assets/main/index.js"), strings.Contains(relPath, "assets/resources/index.js"):
		return ""
	case strings.Contains(relPath, "cocos-js/cc.js"):
		return "cc"
	case strings.Contains(relPath, "cocos-js/wait-for-ammo-instantiation.js"):
		return "wait-for-ammo-instantiation.js"
	case strings.Contains(relPath, "cocos-js/"):
		return "./" + path.Base(relPath)
	default:
		return "./" + relPath
	}
}

func replaceWithWarning(text string, re *regexp.Regexp, replacement, label string, warnings *[]string) string {
	if !re.MatchString(text) {
		*warnings = append(*warnings, fmt.Sprintf("%s did not match", label))
		return text
	}
	return re.ReplaceAllString(text, replacement)
}

func hasNamedRegister(text string) bool {
	trimmed := strings.TrimLeft(text, " \t\r\n")
	return strings.HasPrefix(trimmed, "System.register(\"") || strings.HasPrefix(trimmed, "System.register('")
}

func nameRegister(text, alias string, warnings *[]string) string {
	idx := strings.Index(text, "System.register(")
	if idx < 0 {
		*warnings = append(*warnings, fmt.Sprintf("System.register call not found for alias %s", alias))
		return text
	}
	insertAt := idx + len("System.register(")
	return text[:insertAt] + fmt.Sprintf("%q,", alias) + text[insertAt:]
}
