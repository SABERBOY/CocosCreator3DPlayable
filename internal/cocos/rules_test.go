package cocos

import (
	"strings"
	"testing"
)

func TestTransformSystemRegisterAddsAlias(t *testing.T) {
	result := TransformSystemRegister("index.js", []byte(`System.register([], function () { return {}; });`))
	if !strings.Contains(result.Content, `System.register("index.js",[]`) {
		t.Fatalf("missing named System.register alias: %s", result.Content)
	}
}

func TestTransformApplicationDoesNotDoubleAppendAmmoExtension(t *testing.T) {
	content := `System.register([], function () {
var waitForAmmo = "wait-for-ammo-instantiation.js";
return topLevelImport('cc')
})`
	result := TransformSystemRegister("application.js", []byte(content))
	if strings.Contains(result.Content, "wait-for-ammo-instantiation.js.js") {
		t.Fatalf("double-appended ammo extension: %s", result.Content)
	}
}
