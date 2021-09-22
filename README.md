# CocosCreator3D 试玩广告
## CocosCreator3D(v3.2.0) 试玩广告制作流程
### 模块裁剪
默认勾选所有的选项,3D物理系统勾选bullet物理,2d勾选Box2D(去除其他则需要修改部分代码)
![](pic/Snipaste_2021-09-22_18-31-00.jpg)
### 使用方法
#### Example

```go
var (
mobileDir  = "E:\\Projects\\CocosProjects\\WebGameDemo\\build\\web-mobile"
outDir     = "E:\\Projects\\CocosProjects\\WebGameDemo\\build\\PS"
htmlFile   = outDir + ak.PS + "index.html"
tmpOutFile = outDir + ak.PS + "/index.html"
)
```
修改路径为