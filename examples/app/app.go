// export GOPROXY=direct

package main

import (
	. "fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	webview "github.com/webview/webview_go"
	gw "github.com/wegmarken2006/gwaui/gwasrv"
)

func show(message string) {
	Printf("show: %s\n", message)
}

func testImage(filePath string, width int, height int) {

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Colors are defined by Red, Green, Blue, Alpha uint8 values.
	cyan := color.RGBA{100, 200, 200, 0xff}

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case x < width/2 && y < height/2: // upper left quadrant
				img.Set(x, y, cyan)
			case x >= width/2 && y >= height/2: // lower right quadrant
				img.Set(x, y, color.Black)
			default:
				// Use zero value.
			}
		}
	}

	// Encode as PNG.
	f, _ := os.Create(filePath)
	png.Encode(f, img)
}

func main() {
	args := os.Args

	wv := false
	if len(args) > 1 {
		if args[1] == "wv" {
			wv = true
		}
	}

	getElem, addr, err := gw.Init("config_tabs.yaml")
	if err != nil {
		Println(err)
		os.Exit(0)
	}

	form1 := getElem("id_3")
	dd1 := getElem("id_1")
	lb1 := getElem("id_6")
	h21 := getElem("id_7")
	sl2 := getElem("id_4")
	ta1 := getElem("id_5")
	bt1 := getElem("id_2")
	dt1 := getElem("id_30")
	//cv1 := getElem("id_44")
	img1 := getElem("id_46")
	bt2 := getElem("id_41")
	plt1 := getElem("id_47")
	plt2 := getElem("id_48")
	plt3 := getElem("id_49")
	plt4 := getElem("id_50")

	dt1.AttachWebSocket(show)

	form1.AttachWebSocket(show)

	dd1.SetItemsList([]string{"one", "two", "three"})
	dd1.AttachWebSocket(show)

	sl2.AttachWebSocket(func(message string) {
		Printf("recv: %s\n", message)
		lb1.SetInnerText(message)
	})

	bt1.AttachWebSocket(func(message string) {
		Printf("recv: %s\n", message)
		err := ta1.WriteTextArea("fgfgfgfgfg")
		if err != nil {
			Println("error ", err)
		}

		h21.SetBackgroundColor("red")
		h21.SetColor("white")
	})

	bt2.AttachWebSocket(func(message string) {
		testImage("static/image.png", 400, 400)
		img1.ShowImage("image.png")
	})

	var y2 [][]float64
	x2 := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
	y2 = append(y2, []float64{1.0, 2.0, 4.0, 8.0, 16.0})
	y2 = append(y2, []float64{2.0, 4.0, 8.0, 16.0, 32.0})
	xs2 := []string{"aa", "bb", "cc", "dd", "ee"}
	ys2 := xs2
	pConf := gw.PlotConf{X: x2, Y: y2}
	pConf.X_cat = xs2
	pConf.Y_cat = ys2
	pConf.Title = "prova"
	pConf.Typ = "scatter"
	pConf.Width = 400
	pConf.Height = 400
	pConf.Names = []string{"one", "two"}
	plt1.DrawPlot(&pConf)
	pConf.Typ = "lines"
	plt2.DrawPlot(&pConf)
	pConf.Typ = "bar"
	plt3.DrawPlot(&pConf)
	pConf.Typ = "box"
	plt4.DrawPlot(&pConf)

	if wv {
		w := webview.New(false)
		defer w.Destroy()
		w.SetTitle("Bind Example")
		w.SetSize(1200, 1000, 0)
		w.Navigate(addr)

		w.Run()

	} else {
		text := Sprintf("Serving on %s", addr)
		Println(text)

		gw.WaitKeyFromCOnsole()
	}
}
