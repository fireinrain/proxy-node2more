package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"strings"
)

// AllConfig 配置枚举
type AllConfig struct {
	//输入的节点切片
	InputNodeList []string
	//cdn提供商
	CDNName CdnProvider
	//获取方式
	GetMethodName GetMethod
	//获取的节点数
	WantedNodeNum int

	//输出的节点切片
	OutPutNodeList []string
}

// CdnProvider cdn提供商枚举
type CdnProvider int

const (
	CDNCloudflare CdnProvider = iota
	CDNCloudFront
	CDNGcore
	CDNOther
)

// GetMethod 获取方式
type GetMethod int

const (
	GetMethodSequance GetMethod = iota
	GetMethodRandom
)

func main() {

	app := app.New()
	theme := &ChineseWordTheme{}
	theme.SetFonts("./assets/font/msyh.ttf", "")
	app.Settings().SetTheme(theme)

	window := app.NewWindow("CDN节点快速提取替换工具")
	window.Resize(fyne.NewSize(875, 590))
	//全局配置
	var globalConfig = AllConfig{
		InputNodeList:  nil,
		CDNName:        0,
		GetMethodName:  0,
		WantedNodeNum:  100,
		OutPutNodeList: nil,
	}

	//输入节点
	nodeinput := widget.NewMultiLineEntry()
	nodeinput.SetPlaceHolder("请输入Vmess/Trojan节点(可多行)...")
	nodeinput.SetMinRowsVisible(8)
	inputNodeLabel := widget.NewLabel("原始节点: ")
	inputNodeLabel.Resize(fyne.NewSize(400, 400))
	inputContainer := container.New(layout.NewFormLayout(), inputNodeLabel, nodeinput)

	//自定义CDN ip列表输入框
	customCdnIpLabel := widget.NewLabel("CDN IP列表: ")
	customCdnIpInput := widget.NewMultiLineEntry()
	customCdnIpInput.SetPlaceHolder("192.168.1.1\n172.3.5.2")
	customCdnIpInput.SetMinRowsVisible(8)
	customIpContainer := container.New(layout.NewFormLayout(), customCdnIpLabel, customCdnIpInput)
	//默认设置为不可输入
	customCdnIpInput.Disable()

	//选择栏
	cdnProviderLabel := widget.NewLabel("CDN提供商: ")
	var cdnProvider = []string{"Cloudflare", "CloudFront", "Gcore", "自定义"}
	cdnSelectList := widget.NewSelect(cdnProvider, func(value string) {
		println("CDN Provider你选择了: ", value)
		var tempValue CdnProvider = 0
		if value == "Cloudflare" {
			tempValue = CDNCloudflare
			if !customCdnIpInput.Disabled() {
				customCdnIpInput.Disable()
				customIpContainer.Refresh()
			}
		}
		if value == "CloudFront" {
			tempValue = CDNCloudFront
			if !customCdnIpInput.Disabled() {
				customCdnIpInput.Disable()
				customIpContainer.Refresh()
			}
		}
		if value == "Gcore" {
			tempValue = CDNGcore
			if !customCdnIpInput.Disabled() {
				customCdnIpInput.Disable()
				customIpContainer.Refresh()
			}
		}
		if value == "自定义" {
			tempValue = CDNOther
			if customCdnIpInput.Disabled() {
				customCdnIpInput.Enable()
				customIpContainer.Refresh()

			}
		}
		globalConfig.CDNName = tempValue
	})
	cdnSelectList.Selected = "Cloudflare"
	cdnSelectList.Alignment = fyne.TextAlignCenter
	//获取方式
	getMethodLabel := widget.NewLabel("获取方式: ")
	var getMethodProvider = []string{"顺序", "随机"}
	getMethodList := widget.NewSelect(getMethodProvider, func(value string) {
		println("你选择了: ", value)
		var tempValue GetMethod = 0
		if value == "顺序" {
			tempValue = GetMethodSequance
		}
		if value == "随机" {
			tempValue = GetMethodRandom
		}
		globalConfig.GetMethodName = tempValue
	})
	getMethodList.Selected = "顺序"
	getMethodList.Alignment = fyne.TextAlignCenter
	//获取节点数
	getNodeNumLabel := widget.NewLabel("获取节点数: ")
	inputNodeNum := widget.NewEntry()
	inputNodeNum.SetText("100")
	inputNodeNum.SetMinRowsVisible(10)
	inputNodeNum.OnChanged = func(value string) {
		println("GetNodeNum 你输入了: ", value)

	}
	//提交按钮
	subButton := widget.NewButton("↓点击提取节点↓", func() {

	})
	subButton.OnTapped = func() {
		println("提交按钮已点击")
	}

	selectContainer := container.New(layout.NewHBoxLayout(),
		cdnProviderLabel, cdnSelectList,
		getMethodLabel, getMethodList,
		getNodeNumLabel, inputNodeNum, subButton)

	//输出
	textarea := widget.NewMultiLineEntry()
	textarea.SetMinRowsVisible(8)
	outPutNodes := widget.NewLabel("节点列表: ")
	outPutContainer := container.New(layout.NewFormLayout(), outPutNodes, textarea)

	window.SetContent(container.New(layout.NewVBoxLayout(),
		inputContainer,
		customIpContainer,
		selectContainer,
		outPutContainer,
	))

	window.ShowAndRun()
}

// ChineseWordTheme 自定义主题支持中文显示
type ChineseWordTheme struct {
	regular, bold, italic, boldItalic, monospace fyne.Resource
}

func (cwt *ChineseWordTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (cwt *ChineseWordTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (cwt *ChineseWordTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace {
		return cwt.monospace
	}
	if style.Bold {
		if style.Italic {
			return cwt.boldItalic
		}
		return cwt.bold
	}
	if style.Italic {
		return cwt.italic
	}
	return cwt.regular
}

func (cwt *ChineseWordTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}

func (cwt *ChineseWordTheme) SetFonts(regularFontPath string, monoFontPath string) {
	cwt.regular = theme.TextFont()
	cwt.bold = theme.TextBoldFont()
	cwt.italic = theme.TextItalicFont()
	cwt.boldItalic = theme.TextBoldItalicFont()
	cwt.monospace = theme.TextMonospaceFont()

	if regularFontPath != "" {
		cwt.regular = loadCustomFont(regularFontPath, "Regular", cwt.regular)
		cwt.bold = loadCustomFont(regularFontPath, "Bold", cwt.bold)
		cwt.italic = loadCustomFont(regularFontPath, "Italic", cwt.italic)
		cwt.boldItalic = loadCustomFont(regularFontPath, "BoldItalic", cwt.boldItalic)
	}
	if monoFontPath != "" {
		cwt.monospace = loadCustomFont(monoFontPath, "Regular", cwt.monospace)
	} else {
		cwt.monospace = cwt.regular
	}
}

func loadCustomFont(env, variant string, fallback fyne.Resource) fyne.Resource {
	variantPath := strings.Replace(env, "Regular", variant, -1)

	res, err := fyne.LoadResourceFromPath(variantPath)
	if err != nil {
		fyne.LogError("Error loading specified font", err)
		return fallback
	}

	return res
}
