package htmlx

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"CocosCreator3DPlayable/internal/assets"

	xhtml "golang.org/x/net/html"
)

type Blocks struct {
	ResourceScript string
	ModuleScripts  []string
	ImportMap      string
	LoaderScript   string
	StartScript    string
}

func Build(inputRoot string, original []byte, blocks Blocks) (string, int, error) {
	doc, err := xhtml.Parse(bytes.NewReader(original))
	if err != nil {
		return "", 0, err
	}

	inlined, err := rewriteNode(inputRoot, doc)
	if err != nil {
		return "", 0, err
	}
	body := findElement(doc, "body")
	if body == nil {
		return "", 0, fmt.Errorf("index.html has no body element")
	}
	removeSystemImportScripts(body)
	appendGeneratedBlocks(body, blocks)

	var out bytes.Buffer
	if err := xhtml.Render(&out, doc); err != nil {
		return "", 0, err
	}
	return out.String(), inlined, nil
}

func rewriteNode(inputRoot string, n *xhtml.Node) (int, error) {
	inlined := 0
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		count, err := rewriteNode(inputRoot, child)
		if err != nil {
			return 0, err
		}
		inlined += count
	}
	if n.Type != xhtml.ElementNode {
		return inlined, nil
	}
	switch strings.ToLower(n.Data) {
	case "script":
		src, ok := attr(n, "src")
		if !ok {
			return inlined, nil
		}
		path, err := assets.Resolve(inputRoot, src)
		if err != nil {
			return 0, err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return 0, err
		}
		removeAttr(n, "src")
		removeChildren(n)
		n.AppendChild(&xhtml.Node{Type: xhtml.TextNode, Data: assets.EscapeRawScript(string(data))})
		inlined++
	case "link":
		rel, _ := attr(n, "rel")
		href, ok := attr(n, "href")
		if !ok || !strings.EqualFold(rel, "stylesheet") {
			return inlined, nil
		}
		path, err := assets.Resolve(inputRoot, href)
		if err != nil {
			return 0, err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return 0, err
		}
		n.Data = "style"
		n.Attr = nil
		removeChildren(n)
		n.AppendChild(&xhtml.Node{Type: xhtml.TextNode, Data: assets.EscapeRawStyle(string(data))})
		inlined++
	}
	return inlined, nil
}

func appendGeneratedBlocks(body *xhtml.Node, blocks Blocks) {
	appendScript(body, blocks.ResourceScript)
	for _, script := range blocks.ModuleScripts {
		appendScript(body, script)
	}
	appendScript(body, "var importMapJson = "+assets.EscapeScriptData(blocks.ImportMap)+";")
	appendScript(body, blocks.LoaderScript)
	appendScript(body, blocks.StartScript)
}

func appendScript(parent *xhtml.Node, text string) {
	if strings.TrimSpace(text) == "" {
		return
	}
	script := &xhtml.Node{Type: xhtml.ElementNode, Data: "script"}
	script.AppendChild(&xhtml.Node{Type: xhtml.TextNode, Data: text})
	parent.AppendChild(&xhtml.Node{Type: xhtml.TextNode, Data: "\n"})
	parent.AppendChild(script)
	parent.AppendChild(&xhtml.Node{Type: xhtml.TextNode, Data: "\n"})
}

func removeSystemImportScripts(root *xhtml.Node) {
	for child := root.FirstChild; child != nil; {
		next := child.NextSibling
		if child.Type == xhtml.ElementNode && strings.EqualFold(child.Data, "script") && nodeTextContains(child, "System.import") {
			root.RemoveChild(child)
		} else {
			removeSystemImportScripts(child)
		}
		child = next
	}
}

func nodeTextContains(n *xhtml.Node, needle string) bool {
	var b strings.Builder
	var walk func(*xhtml.Node)
	walk = func(node *xhtml.Node) {
		if node.Type == xhtml.TextNode {
			b.WriteString(node.Data)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			walk(child)
		}
	}
	walk(n)
	return strings.Contains(b.String(), needle)
}

func findElement(n *xhtml.Node, name string) *xhtml.Node {
	if n.Type == xhtml.ElementNode && strings.EqualFold(n.Data, name) {
		return n
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if found := findElement(child, name); found != nil {
			return found
		}
	}
	return nil
}

func attr(n *xhtml.Node, key string) (string, bool) {
	for _, a := range n.Attr {
		if strings.EqualFold(a.Key, key) {
			return a.Val, true
		}
	}
	return "", false
}

func removeAttr(n *xhtml.Node, key string) {
	filtered := n.Attr[:0]
	for _, a := range n.Attr {
		if !strings.EqualFold(a.Key, key) {
			filtered = append(filtered, a)
		}
	}
	n.Attr = filtered
}

func removeChildren(n *xhtml.Node) {
	for child := n.FirstChild; child != nil; {
		next := child.NextSibling
		n.RemoveChild(child)
		child = next
	}
}
