package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// 1.新建应用
	a := app.New()
	iconContent, _ := ioutil.ReadFile("icon.png")
	// 2.新建窗口
	w := a.NewWindow("http server")
	w.Resize(fyne.NewSize(300, 120))                       // 设置大小
	w.CenterOnScreen()                                     // 设置位置为屏幕中心
	w.SetFixedSize(true)                                   // 窗口大小固定
	w.SetIcon(fyne.NewStaticResource("icon", iconContent)) // 设置图标
	// 3. 新建组件
	lbl1 := widget.NewLabel("Place the source file in the current folder")
	lbl2 := widget.NewLabel("Port")
	en := widget.NewEntry()
	en.SetText("8080")
	var bt *widget.Button
	bt = widget.NewButton("Start", func() {
		// 如果监听器正在运行，则停止它
		if isServerRunning {
			stopServer()
			bt.SetText("Start")
		} else {
			// 如果监听器没有运行，则启动它
			startServer(en.Text)
			bt.SetText("Stop")
		}
	})
	// 4.设置布局
	vb1 := container.NewCenter(lbl1)
	vb2 := container.NewCenter(lbl2)
	h2 := container.NewGridWithColumns(2, vb2, en)
	c := container.NewVBox(vb1, h2, bt)
	// 5.添加组件
	w.SetContent(c)
	// 6.运行应用
	w.ShowAndRun()
}

var (
	// 用于切换服务器状态的布尔值
	isServerRunning bool

	// 监听器所在的goroutine
	serverRoutine *sync.WaitGroup

	// 监听服务器对象
	serverInstance *http.Server
)

func startServer(port string) {
	// 创建一个新的goroutine来运行监听器
	serverRoutine = &sync.WaitGroup{}
	serverRoutine.Add(1)

	// 创建一个新的ServeMux来处理请求
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("")))

	serverInstance = &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	go func() {
		if err := serverInstance.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("HTTP server ListenAndServe failed: " + err.Error())
		}
		serverRoutine.Done()
	}()

	// 更新服务器状态标志和等待组
	isServerRunning = true
}

func stopServer() {
	// 创建一个新的context及其撤销函数
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 调用Shutdown方法来停止web监听服务
	if err := serverInstance.Shutdown(ctx); err != nil {
		fmt.Println("HTTP server Shutdown failed: " + err.Error())
	}

	// 等待放置所有连接的goroutine完成
	serverRoutine.Wait()

	// 更新服务器状态标志
	isServerRunning = false
}
