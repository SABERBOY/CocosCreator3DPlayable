# CocosCreator3D 试玩广告
### [基于该项目修改](https://github.com/Jecced/c3d2one)
### 感谢该大佬的项目


## CocosCreator3D(v3.2.0) 试玩广告制作流程

### 模块裁剪

默认勾选所有的选项,3D物理系统勾选bullet物理,2d勾选Box2D(去除其他则需要修改部分代码)
![](pic/Snipaste_2021-09-22_18-31-00.jpg)

### 使用方法

#### 1.直接运行项目(main.go文件)

```go
var (
mobileDir = "E:\\Projects\\CocosProjects\\WebGameDemo"
outDir = ""
htmlFile = ""
tmpOutFile = ""
)
```

修改路径为当前需要导出试玩广告的项目路径(前提是需要在Cocos导出web-mobile项目)

#### 2.直接编译项目

编译项目为可执行文件(当前命令仅限Windows)

```shell
make build
```

导出为可执行文件,然后命令行执行对应命令(前提是需要在Cocos导出web-mobile项目)

```shell
.\Cocos3DPlayableAdPlugin.exe -path "Cocos项目对应路径"
```

### 导出的文件路径

"Cocos项目对应路径"/build/PS/index.html

### 试玩

查看(当前导出的项目为Cocos官方的示例项目)

```shell
tmpHtml/index.html
```
