# CocosCreator3DPlayable

把 Cocos Creator 3.x 导出的 `web-mobile` 目录打包成一个可投放的 playable-ad HTML 文件。

## 安装和运行

```powershell
go run ./cmd/cocos-playable pack `
  --input "C:\path\to\project\build\web-mobile" `
  --output "C:\path\to\playable.html"
```

覆盖已有文件时显式加 `--force`：

```powershell
go run ./cmd/cocos-playable pack --input build/web-mobile --output dist/playable.html --force
```

压缩输出：

```powershell
go run ./cmd/cocos-playable pack --input build/web-mobile --output dist/playable.html --minify
```

构建可执行文件：

```powershell
go build -o out/cocos-playable.exe ./cmd/cocos-playable
```

## 当前行为

- 输入必须直接指向 Cocos Creator 导出的 `web-mobile` 目录。
- 输出是一个单独的 HTML 文件，不再使用旧的 `-path` 和 `build/PS/index.html` 契约。
- 不复制项目目录，不递归清空输出目录，只写用户指定的输出文件。
- 默认拒绝 `../`、绝对路径、外部 URL 和 symlink 资源引用，防止打包时读取输入目录外的文件。
- JSON 和 CCON 资源会校验为合法 JSON；文本和二进制资源会安全序列化进 `window.res`。
- 生成脚本会转义 `<`、`>`、`&`、U+2028、U+2029、`</script>` 和 `</style>` 相关上下文，避免资源内容破坏 HTML 结构。

## Cocos 3.x 适配策略

打包器按规则处理 Cocos 3.x 常见导出结构：

- 给匿名 `System.register(...)` 模块补充稳定模块名。
- 内联 `script src` 和 stylesheet 链接。
- 将 `src/import-map.json`、文本资源和二进制资源放入 `window.res`。
- 注入 playable 资源加载器，注册图片、JSON、文本、音频、字体、WASM、CCON/CCONB 等 downloader。
- 对 `index.js`、`application.js`、WASM 初始化和 bullet WASM 路径做兼容转换。

如果新的 Cocos 3.x 导出结构导致转换规则没有命中，命令会输出 warning。把真实导出裁剪成 fixture 后放进 `testdata/cocos3x/`，再增加回归测试。

## 验证

```powershell
go test ./...
go test -race ./...
go vet ./...
go build ./...
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

CI 使用 Go 1.25.11 运行 build、test、race、vet 和 govulncheck。

## 安全说明

这个工具处理本地文件系统和会发布到广告环境的 HTML，因此输入目录应来自可信的 Cocos 构建链。当前实现仍按安全默认值处理半可信输入：

- 所有 HTML 引用资源都必须解析在输入目录内。
- symlink 默认跳过或拒绝。
- 输出文件使用同目录临时文件写入，再原子 rename。
- 未加 `--force` 时不会覆盖已有输出。

历史版本说明保留在 `version/` 下，仅用于理解旧方案和 Cocos 3.2/3.4 的背景；旧文档里的 `-path` 用法已经废弃。
