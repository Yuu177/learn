[TOC]

# PlantUML

> UML-Unified Modeling Language [统一建模语言](https://baike.baidu.com/item/统一建模语言/3160571?fromModule=lemma_inlink)。这个文档主要记录画图的一些坑。

**PlantUML** 是一个可以让你快速编写 UML 图的组件。官方文档地址：https://plantuml.com/zh/

也可以在 vscode 中安装 plantuml 插件。

## 内容格式化

快捷键 `Shift + Alt + F`  或者右键选择格式化文档。

## 画 UML 图

### 泳道图

https://plantuml.com/zh/activity-diagram-beta

#### 组合

> 分组, 分区, 包, 矩形或卡片式

你可以通过定义分组活动:

- group
- partition
- package
- rectangle
- card

```
@startuml
start
group 分组
  :Activity;
end group
floating note: 分组备注

partition 分区 {
  :Activity;
}
floating note: 分区备注

package 包 {
  :Activity;
}
floating note: 包备注

rectangle 矩形 {
  :Activity;
}
floating note: 矩形备注

card 卡片式 {
  :Activity;
}
floating note: 卡片式备注
end
@enduml

```



![](//www.plantuml.com/plantuml/png/SoWkIImgAStDuG8pkDAByaiB59vsj3tVtSAbe63bc5oIMPPPKcdDbPgNeW2MvKhBoKyioSnBLyZBBqcrWYf-kgJzsUOLN5m5G5CoIpBpyq3YJtjsALIZ6bEBfXsg3A4zEJinFLNXQKyh8PqWDJ1jHQd99ObvwJcf2i_dhtowTn4XlL1bCEt9YKKf2azx5pxlR7-wfv-GPeHAY7vGq70v00bWC080)

这里都可以分组，但是不建议用 `group` 来分组，因为这里有个 bug，格式化文档的时候会排版会有问题。

## 命令行

> 使用命令行运行 PlantUML，有可能写 Makefile 生成图片的时候会使用到。

https://plantuml.com/zh/command-line

- Makefile

```makefile
# 生成图片并存储在 Makefile 的同一目录下

makeFilePath:=$(shell pwd)/$(lastword $(MAKEFILE_LIST))
rootPath=$(shell dirname $(makeFilePath)) # 获取绝对路径

all:
	@echo $(rootPath)
	@java -jar plantuml.jar -o "$(rootPath)/images" "./test/**.puml"
```

[下载 plantuml.jar](https://plantuml.com/zh/download)