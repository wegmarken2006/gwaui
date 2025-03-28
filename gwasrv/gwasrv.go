package gwasrv

import (
	"bufio"
	. "fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v3"

	webview "github.com/webview/webview_go"
)

const WEBSOCKET_BUFFER_SIZE = 4096
const BUFFER_SIZE = 4096

var upgrader = websocket.Upgrader{
	ReadBufferSize:  WEBSOCKET_BUFFER_SIZE,
	WriteBufferSize: WEBSOCKET_BUFFER_SIZE,
}

type WsElem struct {
	gs   *websocket.Conn
	addr string
	m    *sync.Mutex
}

func WsElemNew(id string) WsElem {

	addr := Sprintf("/%s", id)
	wsElem := WsElem{addr: addr}

	return wsElem
}

func (wse *WsElem) AttachWebSocket(fn func(message string)) error {

	var conn *websocket.Conn
	var err error
	http.HandleFunc(wse.addr, func(w http.ResponseWriter, r *http.Request) {
		conn, err = upgrader.Upgrade(w, r, nil)
		if err != nil {
			Println("attach error ", err)
			return
		}
		wse.gs = conn

		go func() {
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					//Println("read:", err)
					break
				}
				fn(string(message))
			}
		}()
	})
	return err
}

/*
func (wse *WsElem) InitWriteWs() {
	http.HandleFunc(wse.addr, func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		wse.gs = c
		//defer c.Close() where to?
	})
}
*/

func (wse *WsElem) wsWrite(message []byte) error {
	var err error = nil
	if wse.m == nil {
		var m sync.Mutex
		wse.m = &m
	}

	go func() {
		//wait for available websocket
		for {
			for ind := 0; ind < 250; ind++ {
				if wse.gs != nil {
					break
				}
				time.Sleep(200 * time.Millisecond)
			}
			if wse.gs == nil {
				Println(wse.addr, "still no websocket")
			} else {
				wse.m.Lock()
				err = wse.gs.WriteMessage(websocket.TextMessage, message)
				wse.m.Unlock()
				if err != nil {
					Println(err)
				}
				//Println(wse.addr, "ok websocket")
				return
				//Println(wse.addr, "ok websocket")
			}
		}
	}()
	/*
		if wse.gs == nil {
			return fmt.Errorf("no websocket")
		}
		err := wse.gs.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return err
		}
	*/
	return err
}

func (wse *WsElem) writeMessage(txMsg rxTxMessage) error {
	byteTx, err := yaml.Marshal(txMsg)
	if err != nil {
		return err
	}
	err = wse.wsWrite(byteTx)
	if err != nil {
		return err
	}

	return nil
}

func (wse *WsElem) ShowImage(pngFile string) error {
	txMsg := rxTxMessage{}
	txMsg.ImageName = pngFile
	err := wse.writeMessage(txMsg)
	return err
}

func (wse *WsElem) SetBackgroundColor(color string) error {
	txMsg := rxTxMessage{}
	txMsg.BackgroundColor = color
	err := wse.writeMessage(txMsg)
	return err
}

func (wse *WsElem) SetColor(color string) error {
	txMsg := rxTxMessage{}
	txMsg.Color = color
	err := wse.writeMessage(txMsg)
	return err
}

func (wse *WsElem) SetInnerText(text string) error {
	txMsg := rxTxMessage{}
	txMsg.Text = text
	err := wse.writeMessage(txMsg)
	return err
}

func (wse *WsElem) SetItemsList(lst []string) error {
	txMsg := rxTxMessage{}
	txMsg.ItemList = lst
	err := wse.writeMessage(txMsg)
	return err
}

func (wse *WsElem) WriteTextArea(text string) error {
	txMsg := rxTxMessage{}
	txMsg.Textarea = text
	err := wse.writeMessage(txMsg)
	return err
}

/*
func (wse *WsElem) DrawPlot(pConf *plotConf) error {
	txMsg := rxTxMessage{}
	txMsg.PlotConf = pConf
	err := wse.writeMessage(txMsg)
	return err
}
*/

func (wse *WsElem) DrawPlotLines(x []float64, ys [][]float64, names []string, layout *PlotLayout) error {
	pConf := plotConf{}
	pConf.Typ = "lines"
	pConf.X = x
	pConf.Y = ys
	pConf.Names = names
	pConf.Title = layout.Title
	pConf.Width = layout.Width
	pConf.Height = layout.Height
	txMsg := rxTxMessage{}
	txMsg.PlotConf = &pConf
	err := wse.writeMessage(txMsg)
	return err
}

func (wse *WsElem) DrawPlotScatter(x []float64, ys [][]float64, names []string, layout *PlotLayout) error {
	pConf := plotConf{}
	pConf.Typ = "scatter"
	pConf.X = x
	pConf.Y = ys
	pConf.Names = names
	pConf.Title = layout.Title
	pConf.Width = layout.Width
	pConf.Height = layout.Height
	txMsg := rxTxMessage{}
	txMsg.PlotConf = &pConf
	err := wse.writeMessage(txMsg)
	return err
}

func (wse *WsElem) DrawPlotBars(x []string, ys [][]float64, names []string, layout *PlotLayout) error {
	pConf := plotConf{}
	pConf.Typ = "bar"
	pConf.X_cat = x
	pConf.Y = ys
	pConf.Names = names
	pConf.Title = layout.Title
	pConf.Width = layout.Width
	pConf.Height = layout.Height
	txMsg := rxTxMessage{}
	txMsg.PlotConf = &pConf
	err := wse.writeMessage(txMsg)
	return err
}

func (wse *WsElem) DrawPlotBox(ys [][]float64, names []string, layout *PlotLayout) error {
	pConf := plotConf{}
	pConf.Typ = "box"
	pConf.Y = ys
	pConf.Names = names
	pConf.Title = layout.Title
	pConf.Width = layout.Width
	pConf.Height = layout.Height
	txMsg := rxTxMessage{}
	txMsg.PlotConf = &pConf
	err := wse.writeMessage(txMsg)
	return err
}

func Init(yamlName string) (func(string) *WsElem, string, error) {
	// Read and send Yaml configuration to the client

	if _, err := os.Stat("./static"); os.IsNotExist(err) {
		Println("Folder ./static missing.")
		os.Exit(0)
	}

	yamlFile, err := os.Open(yamlName)

	if err != nil {
		return nil, "", err
	}

	defer yamlFile.Close()

	byteValue, _ := io.ReadAll(yamlFile)

	//client reading configuration attached to the "title" element defined in static/index.html
	//title := WsElem{addr: "/id_0"}
	title := WsElemNew("id_0")
	title.AttachWebSocket(func(message string) {
		//Println(message)
		title.wsWrite(byteValue)
	})

	var guiDescr []GuiDescr
	yaml.Unmarshal([]byte(byteValue), &guiDescr)

	helems := make(map[string]*WsElem)

	for _, tabs := range guiDescr {
		//for _, grid := range row.GridRow {
		rows := tabs.Tab.Row
		for _, row := range rows {
			for _, grid := range row.GridRow {
				if len(grid.Button.Id) > 0 {
					bt := WsElemNew(grid.Button.Id)
					helems[grid.Button.Id] = &bt
				}
				if len(grid.Dropdown.Id) > 0 {
					dd := WsElemNew(grid.Dropdown.Id)
					helems[grid.Dropdown.Id] = &dd
				}
				if len(grid.Form.Id) > 0 {
					fm := WsElemNew(grid.Form.Id)
					helems[grid.Form.Id] = &fm
				}
				if len(grid.Input.Id) > 0 {
					ip := WsElemNew(grid.Input.Id)
					helems[grid.Input.Id] = &ip
				}
				if len(grid.Date.Id) > 0 {
					dt := WsElemNew(grid.Date.Id)
					helems[grid.Date.Id] = &dt
				}
				if len(grid.Slider.Id) > 0 {
					sl := WsElemNew(grid.Slider.Id)
					helems[grid.Slider.Id] = &sl
				}
				if len(grid.Textarea.Id) > 0 {
					ta := WsElemNew(grid.Textarea.Id)
					ta.AttachWebSocket(func(message string) {})
					helems[grid.Textarea.Id] = &ta
				}
				if len(grid.Label.Id) > 0 {
					lb := WsElemNew(grid.Label.Id)
					if grid.Label.Mutable {
						lb.AttachWebSocket(func(message string) {})
					}
					helems[grid.Label.Id] = &lb
				}
				if len(grid.H2.Id) > 0 {
					h2 := WsElemNew(grid.H2.Id)
					if grid.H2.Mutable {
						h2.AttachWebSocket(func(message string) {})
					}
					helems[grid.H2.Id] = &h2
				}
				if len(grid.Canvas.Id) > 0 {
					cv := WsElemNew(grid.Canvas.Id)
					cv.AttachWebSocket(func(message string) {})
					helems[grid.Canvas.Id] = &cv

				}
				if len(grid.Image.Id) > 0 {
					im := WsElemNew(grid.Image.Id)
					im.AttachWebSocket(func(message string) {})
					helems[grid.Image.Id] = &im
				}
				if len(grid.Plot.Id) > 0 {
					plt := WsElemNew(grid.Plot.Id)
					plt.AttachWebSocket(func(message string) {})
					helems[grid.Plot.Id] = &plt
				}

			}
		}
	}

	addr := StartServer()

	retFun := func(id string) *WsElem {
		elem := helems[id]
		if elem == nil {
			Println(id, "not in yaml file")
			os.Exit(0)
		}
		return elem
	}
	return retFun, addr, nil
}

func Run(addr string, titleStr string, width int, height int, wview bool) {
	if wview {
		w := webview.New(false)
		defer w.Destroy()
		w.SetTitle(titleStr)
		w.SetSize(width, height, 0)
		w.Navigate(addr)

		w.Run()
	} else {
		text := Sprintf("Serving on %s", addr)
		Println(text)

		WaitKeyFromCOnsole()
	}
}

func StartServer() string {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/", http.FileServer(http.Dir("./static")))

	ch := make(chan string)

	go func() {
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			panic(err)
		}

		addrStr := Sprintf("http://127.0.0.1:%d", listener.Addr().(*net.TCPAddr).Port)
		ch <- addrStr
		//text := Sprintf("Serving on %s", addrStr)
		//Println(text)

		panic(http.Serve(listener, nil))
	}()

	addrStr := <-ch

	return addrStr
}

func WaitKeyFromCOnsole() {
	//Wait for a key press
	reader := bufio.NewReader(os.Stdin)
	Println("Press:\n q<Enter> to exit")
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			os.Exit(0)
		}
		if strings.HasPrefix(text, "q") || strings.HasPrefix(text, "Q") {
			os.Exit(0)
		}
	}
}
