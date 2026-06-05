package assets

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Kind int

const (
	KindJSON Kind = iota
	KindText
	KindBinary
)

type Asset struct {
	Path string
	Kind Kind
	Data []byte
}

type Inventory struct {
	Assets          []Asset
	SystemScripts   []string
	LibraryScripts  []string
	SkippedSymlinks int
}

var jsonExts = map[string]bool{
	".json": true,
	".ccon": true,
}

var textExts = map[string]bool{
	".txt": true,
}

var binaryExts = map[string]bool{
	".bin":   true,
	".png":   true,
	".jpg":   true,
	".jpeg":  true,
	".mp3":   true,
	".ogg":   true,
	".wav":   true,
	".m4a":   true,
	".wasm":  true,
	".ttf":   true,
	".cconb": true,
}

func Scan(root string) (Inventory, error) {
	var inv Inventory
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			return nil
		}
		if d.Type()&os.ModeSymlink != 0 {
			inv.SkippedSymlinks++
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}

		rel, err := Rel(root, path)
		if err != nil {
			return err
		}
		ext := strings.ToLower(filepath.Ext(rel))
		if ext == ".js" {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			text := strings.TrimSpace(string(data))
			if strings.HasPrefix(text, "System.register(") {
				inv.SystemScripts = append(inv.SystemScripts, rel)
			} else if strings.Contains(rel, "src/assets/script/libs/") {
				inv.LibraryScripts = append(inv.LibraryScripts, rel)
			}
			return nil
		}

		kind, ok := classify(ext)
		if !ok {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if kind == KindJSON && !json.Valid(data) {
			return fmt.Errorf("%s is not valid JSON", rel)
		}
		inv.Assets = append(inv.Assets, Asset{Path: rel, Kind: kind, Data: data})
		return nil
	})
	if err != nil {
		return Inventory{}, err
	}

	sort.Slice(inv.Assets, func(i, j int) bool { return inv.Assets[i].Path < inv.Assets[j].Path })
	sort.Strings(inv.SystemScripts)
	sort.Strings(inv.LibraryScripts)
	return inv, nil
}

func ResourceScript(assets []Asset) (string, error) {
	var b strings.Builder
	b.WriteString("window.res = Object.create(null);\n")
	for _, asset := range assets {
		key, err := jsonString(asset.Path)
		if err != nil {
			return "", err
		}
		value, err := assetLiteral(asset)
		if err != nil {
			return "", err
		}
		b.WriteString("window.res[")
		b.WriteString(key)
		b.WriteString("] = ")
		b.WriteString(value)
		b.WriteString(";\n")
	}
	return EscapeScriptData(b.String()), nil
}

func Resolve(root, ref string) (string, error) {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return "", errors.New("empty asset reference")
	}
	if strings.Contains(ref, "://") || strings.HasPrefix(ref, "//") || strings.HasPrefix(ref, "data:") {
		return "", fmt.Errorf("external asset reference %q is not supported", ref)
	}
	if i := strings.IndexAny(ref, "?#"); i >= 0 {
		ref = ref[:i]
	}
	ref = strings.ReplaceAll(ref, "\\", "/")
	clean := filepath.Clean(filepath.FromSlash(ref))
	if filepath.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("asset reference %q escapes input directory", ref)
	}

	target := filepath.Join(root, clean)
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	targetAbs, err := filepath.Abs(target)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(rootAbs, targetAbs)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) || filepath.IsAbs(rel) {
		return "", fmt.Errorf("asset reference %q escapes input directory", ref)
	}
	info, err := os.Lstat(targetAbs)
	if err != nil {
		return "", err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("asset reference %q points to a symlink", ref)
	}
	return targetAbs, nil
}

func Rel(root, target string) (string, error) {
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return "", err
	}
	if rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) || filepath.IsAbs(rel) {
		return "", fmt.Errorf("%s is outside %s", target, root)
	}
	return filepath.ToSlash(rel), nil
}

func EscapeScriptData(s string) string {
	replacer := strings.NewReplacer(
		"<", "\\u003c",
		">", "\\u003e",
		"&", "\\u0026",
		"\u2028", "\\u2028",
		"\u2029", "\\u2029",
	)
	return replacer.Replace(s)
}

func EscapeRawScript(s string) string {
	return escapeRawTag(s, "script")
}

func EscapeRawStyle(s string) string {
	return escapeRawTag(s, "style")
}

func classify(ext string) (Kind, bool) {
	switch {
	case jsonExts[ext]:
		return KindJSON, true
	case textExts[ext]:
		return KindText, true
	case binaryExts[ext]:
		return KindBinary, true
	default:
		return KindText, false
	}
}

func assetLiteral(asset Asset) (string, error) {
	switch asset.Kind {
	case KindJSON:
		return EscapeScriptData(string(asset.Data)), nil
	case KindText:
		return jsonString(string(asset.Data))
	case KindBinary:
		return jsonString(base64.StdEncoding.EncodeToString(asset.Data))
	default:
		return "", fmt.Errorf("unknown asset kind for %s", asset.Path)
	}
}

func jsonString(s string) (string, error) {
	data, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return EscapeScriptData(string(data)), nil
}

func escapeRawTag(s, tag string) string {
	lower := strings.ToLower(s)
	needle := "</" + tag
	var b strings.Builder
	last := 0
	for {
		i := strings.Index(lower[last:], needle)
		if i < 0 {
			b.WriteString(s[last:])
			break
		}
		i += last
		b.WriteString(s[last:i])
		b.WriteString("<\\/")
		b.WriteString(s[i+2 : i+2+len(tag)])
		last = i + 2 + len(tag)
	}
	return b.String()
}
