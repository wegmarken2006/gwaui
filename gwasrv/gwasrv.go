package gwasrv

import (
	"fmt"
	. "fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v3"
)

const WEBSOCKET_BUFFER_SIZE = 4096
const BUFFER_SIZE = 4096

var upgrader = websocket.Upgrader{
	ReadBufferSize:  WEBSOCKET_BUFFER_SIZE,
	WriteBufferSize: WEBSOCKET_BUFFER_SIZE,
}

const STARTING_PORT = 9001

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

func Init(yamlName string) error {
	// Read and send Yaml configuration to the client

	yamlFile, err := os.Open(yamlName)

	if err != nil {
		return err
	}

	defer yamlFile.Close()

	byteValue, _ := io.ReadAll(yamlFile)

	//client reading configuration attached to the "title" element defined in static/index.html
	title := WsElem{addr: "/id_0"}
	title.AttachWebSocket(func(message string) {
		//Println(message)
		title.wsWrite(byteValue)
	})

	return nil
}

func StartServer() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.Handle("/", http.FileServer(http.Dir("./static")))

	port := STARTING_PORT
	go func() {
		for {
			portStr := Sprintf(":%d", port)
			Println(portStr)
			text := Sprintf("Serving on http://localhost%s", portStr)
			Println(text)
			err := http.ListenAndServe(portStr, nil)
			if err != nil {
				Println(err)
				port += 1
			} else {
				break
			}
		}
	}()
}
