package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"proxy-node2more/config"
	"proxy-node2more/utils"
	"regexp"
	"strconv"
	"strings"
)

func main() {

	app := app.New()
	theme := &ChineseWordTheme{}
	theme.SetFonts("./assets/font/msyh.ttf", "")
	app.Settings().SetTheme(theme)

	window := app.NewWindow("CDN节点快速提取替换工具")
	window.Resize(fyne.NewSize(875, 590))
	//全局配置
	var globalConfig = config.AllConfig{
		InputNodeStr:   nil,
		CDNName:        0,
		CustomCDNIp:    nil,
		GetMethodName:  0,
		WantedNodeNum:  100,
		OutPutNodeList: nil,
	}

	//输入节点
	nodeinput := widget.NewMultiLineEntry()
	nodeinput.SetPlaceHolder("请输入Vmess/Trojan/Vless节点分享链接...\n请确保输入的节点已经套用了CDN,并且能正常使用.")
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
		fmt.Println("CDN Provider你选择了: ", value)
		var tempValue config.CdnProvider = 0
		if value == "Cloudflare" {
			tempValue = config.CDNCloudflare
			if !customCdnIpInput.Disabled() {
				customCdnIpInput.Disable()
				customIpContainer.Refresh()
			}
		}
		if value == "CloudFront" {
			tempValue = config.CDNCloudFront
			if !customCdnIpInput.Disabled() {
				customCdnIpInput.Disable()
				customIpContainer.Refresh()
			}
		}
		if value == "Gcore" {
			tempValue = config.CDNGcore
			if !customCdnIpInput.Disabled() {
				customCdnIpInput.Disable()
				customIpContainer.Refresh()
			}
		}
		if value == "自定义" {
			tempValue = config.CDNOther
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
		fmt.Println("你选择了: ", value)
		var tempValue config.GetMethod = 0
		if value == "顺序" {
			tempValue = config.GetMethodSequance
		}
		if value == "随机" {
			tempValue = config.GetMethodRandom
		}
		globalConfig.GetMethodName = tempValue
	})
	getMethodList.Selected = "顺序"
	getMethodList.Alignment = fyne.TextAlignCenter
	//获取节点数
	getNodeNumLabel := widget.NewLabel("获取节点数(一个原始节点): ")
	inputNodeNum := widget.NewEntry()
	inputNodeNum.SetText("100")
	inputNodeNum.SetMinRowsVisible(10)
	inputNodeNum.OnChanged = func(value string) {
		fmt.Println("GetNodeNum 你输入了: ", value)
		number, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("格式输入错误: ", err.Error())
			return
		}
		globalConfig.WantedNodeNum = number
	}
	//提交按钮
	subButton := widget.NewButton("↓点击提取节点↓", func() {

	})

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

	subButton.OnTapped = func() {
		text := textarea.Text
		//清空
		if text != "" {
			textarea.SetText("")
		}
		fmt.Println("提交按钮已点击")
		//获取输入框内容
		var nodeInpputStr = nodeinput.Text
		//nodeInpputStr = strings.ReplaceAll(nodeInpputStr, " ", "")
		//nodeInpputStr = strings.ReplaceAll(nodeInpputStr, "\t", "")
		//nodeInpputStr = strings.ReplaceAll(nodeInpputStr, "\n", "")
		//process nodeinputStr
		var nodes []string

		if strings.Contains(nodeInpputStr, "vmess://") {
			pattern := "vmess://[\\w\\d\\+\\/=]+"
			// Compile the regular expression pattern
			regex := regexp.MustCompile(pattern)
			// Find all the matches in the input string
			matches := regex.FindAllString(nodeInpputStr, -1)
			// Print the matches
			for _, match := range matches {
				//fmt.Printf("Match %d: %s\n", i+1, match)
				nodes = append(nodes, match)
				//移除
				nodeInpputStr = strings.ReplaceAll(nodeInpputStr, match, "")
			}
		}
		if strings.Contains(nodeInpputStr, "vless://") {
			re := regexp.MustCompile(`vless://[^\s]+`)

			matches := re.FindAllString(nodeInpputStr, -1)
			// Print the matches
			for _, match := range matches {
				//fmt.Printf("Match %d: %s\n", i+1, match)
				nodes = append(nodes, match)
				//移除
				nodeInpputStr = strings.ReplaceAll(nodeInpputStr, match, "")
			}
		}
		if strings.Contains(nodeInpputStr, "trojan") {
			//extract trojan nodes
			re := regexp.MustCompile(`trojan://[^\s]+`)

			matches := re.FindAllString(nodeInpputStr, -1)
			for _, match := range matches {
				//fmt.Printf("Match %d: %s\n", i+1, match)
				nodes = append(nodes, match)
				//移除
				nodeInpputStr = strings.ReplaceAll(nodeInpputStr, match, "")
			}
		}

		globalConfig.InputNodeStr = nodes
		fmt.Println(nodeInpputStr)
		//获取cdn类型
		nodesResult, err := utils.CaculateNodesResult(&globalConfig)
		if err != nil {
			fmt.Println("程序运行出错: ", err.Error())
		}

		nodeList := nodesResult.OutPutNodeList
		resultStr := strings.Join(nodeList, "\n")
		//更新界面
		textarea.SetText(resultStr)
	}

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
