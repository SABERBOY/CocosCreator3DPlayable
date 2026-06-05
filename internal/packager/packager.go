package packager

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"CocosCreator3DPlayable/internal/assets"
	"CocosCreator3DPlayable/internal/cocos"
	htmlx "CocosCreator3DPlayable/internal/html"

	"github.com/tdewolff/minify/v2"
	minhtml "github.com/tdewolff/minify/v2/html"
)

type Options struct {
	InputDir   string
	OutputFile string
	Minify     bool
	Force      bool
}

type Report struct {
	OutputFile      string
	ResourceCount   int
	ScriptCount     int
	InlinedHTMLRefs int
	SkippedSymlinks int
	Warnings        []string
}

func Run(ctx context.Context, opts Options) (Report, error) {
	if err := ctx.Err(); err != nil {
		return Report{}, err
	}
	inputDir, outputFile, err := normalizeOptions(opts)
	if err != nil {
		return Report{}, err
	}

	indexPath := filepath.Join(inputDir, "index.html")
	indexHTML, err := os.ReadFile(indexPath)
	if err != nil {
		return Report{}, fmt.Errorf("read index.html: %w", err)
	}

	inventory, err := assets.Scan(inputDir)
	if err != nil {
		return Report{}, err
	}
	resourceScript, err := assets.ResourceScript(inventory.Assets)
	if err != nil {
		return Report{}, err
	}
	moduleScripts, warnings, err := buildModuleScripts(inputDir, inventory)
	if err != nil {
		return Report{}, err
	}
	importMap, err := readImportMap(inputDir)
	if err != nil {
		return Report{}, err
	}

	html, inlinedRefs, err := htmlx.Build(inputDir, indexHTML, htmlx.Blocks{
		ResourceScript: resourceScript,
		ModuleScripts:  moduleScripts,
		ImportMap:      importMap,
		LoaderScript:   assets.EscapeRawScript(cocos.LoaderScript),
		StartScript:    "System.import('index.js').catch(function (err) { console.error(err); });",
	})
	if err != nil {
		return Report{}, err
	}
	if opts.Minify {
		html, err = minifyHTML(html)
		if err != nil {
			return Report{}, err
		}
	}
	if err := writeOutput(outputFile, []byte(html), opts.Force); err != nil {
		return Report{}, err
	}

	return Report{
		OutputFile:      outputFile,
		ResourceCount:   len(inventory.Assets),
		ScriptCount:     len(moduleScripts),
		InlinedHTMLRefs: inlinedRefs,
		SkippedSymlinks: inventory.SkippedSymlinks,
		Warnings:        warnings,
	}, nil
}

func normalizeOptions(opts Options) (string, string, error) {
	if strings.TrimSpace(opts.InputDir) == "" {
		return "", "", errors.New("--input is required")
	}
	if strings.TrimSpace(opts.OutputFile) == "" {
		return "", "", errors.New("--output is required")
	}
	inputDir, err := filepath.Abs(opts.InputDir)
	if err != nil {
		return "", "", err
	}
	info, err := os.Stat(inputDir)
	if err != nil {
		return "", "", fmt.Errorf("input directory: %w", err)
	}
	if !info.IsDir() {
		return "", "", fmt.Errorf("input is not a directory: %s", inputDir)
	}
	outputFile, err := filepath.Abs(opts.OutputFile)
	if err != nil {
		return "", "", err
	}
	if inside(inputDir, outputFile) {
		return "", "", fmt.Errorf("output file must not be inside input directory: %s", outputFile)
	}
	return inputDir, outputFile, nil
}

func buildModuleScripts(inputDir string, inventory assets.Inventory) ([]string, []string, error) {
	var scripts []string
	var warnings []string
	for _, rel := range append(append([]string{}, inventory.SystemScripts...), inventory.LibraryScripts...) {
		path, err := assets.Resolve(inputDir, rel)
		if err != nil {
			return nil, nil, err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, nil, err
		}
		if strings.HasSuffix(strings.ToLower(rel), ".js") && contains(inventory.SystemScripts, rel) {
			result := cocos.TransformSystemRegister(rel, data)
			warnings = append(warnings, result.Warnings...)
			scripts = append(scripts, assets.EscapeRawScript(result.Content))
		} else {
			scripts = append(scripts, assets.EscapeRawScript(string(data)))
		}
	}
	return scripts, warnings, nil
}

func readImportMap(inputDir string) (string, error) {
	path, err := assets.Resolve(inputDir, "src/import-map.json")
	if err != nil {
		if os.IsNotExist(err) {
			return "{}", nil
		}
		return "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(strings.TrimSpace(string(data)), "{") {
		return "", fmt.Errorf("src/import-map.json is not a JSON object")
	}
	return string(data), nil
}

func writeOutput(outputFile string, data []byte, force bool) error {
	if _, err := os.Stat(outputFile); err == nil && !force {
		return fmt.Errorf("output file exists; pass --force to overwrite: %s", outputFile)
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(outputFile), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(outputFile), "."+filepath.Base(outputFile)+".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if force {
		if err := os.Remove(outputFile); err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	return os.Rename(tmpName, outputFile)
}

func minifyHTML(input string) (string, error) {
	m := minify.New()
	m.AddFunc("text/html", minhtml.Minify)
	return m.String("text/html", input)
}

func inside(root, target string) bool {
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return false
	}
	return rel == "." || (!strings.HasPrefix(rel, ".."+string(os.PathSeparator)) && rel != ".." && !filepath.IsAbs(rel))
}

func contains(list []string, value string) bool {
	for _, item := range list {
		if item == value {
			return true
		}
	}
	return false
}
