package gwasrv

import (
	"bufio"
	"fmt"
	. "fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v3"
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
}

type TxMessage struct {
	Text            string
	Textarea        string
	BackgroundColor string
	Color           string
	ImageName       string
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
		//defer c.Close() where to
		go func() {
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					Println("read:", err)
					break
				}
				fn(string(message))
				//log.Printf("recv: %s", message)
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
	if wse.gs == nil {
		return fmt.Errorf("no websocket")
	}
	err := wse.gs.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}
	return nil
}

func (wse *WsElem) writeMessage(txMsg TxMessage) error {
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
	txMsg := TxMessage{}
	txMsg.ImageName = pngFile
	err := wse.writeMessage(txMsg)
	return err
}

func (wse *WsElem) SetBackgroundColor(color string) error {
	txMsg := TxMessage{}
	txMsg.BackgroundColor = color
	err := wse.writeMessage(txMsg)
	return err
}

func (wse *WsElem) SetColor(color string) error {
	txMsg := TxMessage{}
	txMsg.Color = color
	err := wse.writeMessage(txMsg)
	return err
}

func (wse *WsElem) SetInnerText(text string) error {
	txMsg := TxMessage{}
	txMsg.Text = text
	err := wse.writeMessage(txMsg)
	return err
}

func (wse *WsElem) WriteTextArea(text string) error {
	txMsg := TxMessage{}
	txMsg.Textarea = text
	err := wse.writeMessage(txMsg)
	return err
}

func Init(yamlName string) (func(string) *WsElem, error) {
	// Read and send Yaml configuration to the client

	yamlFile, err := os.Open(yamlName)

	if err != nil {
		return nil, err
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
			}
		}
	}

	StartServer()

	retFun := func(id string) *WsElem {
		return helems[id]
	}
	return retFun, nil
}

func StartServer() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/", http.FileServer(http.Dir("./static")))

	go func() {
		listener, err := net.Listen("tcp", ":0")
		if err != nil {
			panic(err)
		}

		text := Sprintf("Serving on http://127.0.0.1:%d", listener.Addr().(*net.TCPAddr).Port)
		Println(text)

		panic(http.Serve(listener, nil))
	}()
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
