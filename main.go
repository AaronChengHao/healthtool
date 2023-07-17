package main

import (
	"archive/zip"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"healtool/tool"
	"image/color"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var writeLog = tool.WriteLog
var writeLogFail = tool.WriteLogFail

func main() {
	//os.Setenv("FYNE_FONT", "MSYHL.TTC")
	defer func() {
		if err := recover(); err != nil {
			s := err.(string)
			writeLogFail(s)
		}
	}()

	myApp := app.New()

	tool.FyneWindow = myApp.NewWindow("整合平台工具")

	tool.LogWidget = widget.NewLabel("log info")

	btn1 := widget.NewButton("basic install", func() {
		taskStart()
	})

	btn2 := widget.NewButton("full install", func() {
		writeLog("start execute full install task")
	})

	tool.LogContainer = container.NewScroll(tool.LogWidget)
	tool.LogContainer.SetMinSize(fyne.Size{Height: 200})

	btnContainer := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), btn1, layout.NewSpacer(), btn2, layout.NewSpacer())

	ine := canvas.NewLine(color.White)
	ine.StrokeWidth = 5

	// 使用说明
	explanLabel := widget.NewLabel(`
instructions 
1. dsaffdasfds
2. safsdfsdfsafas
3. sfsdafasdfsd`)

	tool.FyneWindow.SetContent(container.NewVBox(tool.LogContainer, ine, btnContainer, explanLabel))

	tool.FyneWindow.Resize(fyne.Size{Width: 500, Height: 400})

	tool.FyneWindow.ShowAndRun()
}

func taskStart() {
	//downloadZip()
	//downloadChromeInstall()
	openChrome()
}

func openChrome() {
	path, _ := launcher.LookPath()
	u := launcher.New().Bin(path).Delete("--headless").MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	browser.MustPage("https://health.cd120.info")
	browser.MustPage("https://health-zhis.cd120.info")
	dialog.NewConfirm("网址收藏", "请收藏网址", func(b bool) {

	}, tool.FyneWindow).Show()
}

func downloadChromeInstall() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("where", "chrome")
	} else {
		cmd = exec.Command("which", "google-chrome")
	}
	err := cmd.Run()
	if err != nil {
		fmt.Println("Chrome浏览器未安装")
	} else {
		fmt.Println("Chrome浏览器已安装")
	}
}

func downloadZip() {
	writeLog("start download install zip file")
	zipName := "install.zip"
	zipUrl := "https://health.cd120.info/health-web/f/f.zip"

	writeLog("开始下载压缩包文件")
	resp, err := http.Get(zipUrl)
	if err != nil {
		writeLog(err.Error())
		return
	}
	writeLog("压缩包文件下载完成")
	// 创建本地zip文件
	zipFile, err := os.Create(zipName)
	if err != nil {
		writeLog(err.Error())
		return
	}
	writeLog("压缩包文件拷贝到实体文件")
	_, err = io.Copy(zipFile, resp.Body)
	if err != nil {
		writeLog(err.Error())
		return
	}
	// 压缩包
	var src = zipName
	// 获取程序运行目录
	workDir, err := os.Getwd()
	if err != nil {
		writeLog(err.Error())
	}
	var dst = filepath.Join(workDir, "output")

	writeLog("当前工作目录：" + dst)

	writeLog("开始解压")

	if err := UnZip(dst, src); err != nil {
		writeLog(err.Error())
		return
	}
	writeLog("解压完成")
}

func UnZip(dst, src string) (err error) {
	// 打开压缩文件，这个 zip 包有个方便的 ReadCloser 类型
	// 这个里面有个方便的 OpenReader 函数，可以比 tar 的时候省去一个打开文件的步骤
	zr, err := zip.OpenReader(src)
	defer zr.Close()
	if err != nil {
		return
	}

	// 如果解压后不是放在当前目录就按照保存目录去创建目录
	if dst != "" {
		if err := os.MkdirAll(dst, 0755); err != nil {
			return err
		}
	}

	// 遍历 zr ，将文件写入到磁盘
	for _, file := range zr.File {
		//continue
		if strings.HasPrefix(filepath.Base(file.Name), ".") {
			//fmt.Println(file.Name + ": 包含")
			continue
		}

		path := filepath.Join(dst, file.Name)

		// 如果是目录，就创建目录
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, file.Mode()); err != nil {
				return err
			}
			// 因为是目录，跳过当前循环，因为后面都是文件的处理
			continue
		}

		// 获取到 Reader
		fr, err := file.Open()
		if err != nil {
			return err
		}

		// 创建要写出的文件对应的 Write
		fw, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		_, err = io.Copy(fw, fr)
		if err != nil {
			return err
		}

		// 将解压的结果输出
		//fmt.Printf("成功解压 %s ，共写入了 %d 个字符的数据\n", path, n)

		// 因为是在循环中，无法使用 defer ，直接放在最后
		// 不过这样也有问题，当出现 err 的时候就不会执行这个了，
		// 可以把它单独放在一个函数中，这里是个实验，就这样了
		//fw.Close()
		fr.Close()
	}
	return nil
}
