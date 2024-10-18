// export GOPROXY=direct

package main

import (
	. "fmt"
	"os"

	gw "github.com/xxx/gwasrv"
)

func show(message string) {
	Printf("show: %s\n", message)
}

func main() {

	err := gw.Init("config_tabs.yaml")
	if err != nil {
		Println(err)
		os.Exit(0)
	}

	form1 := gw.WsElemNew("id_3")
	form1.AttachWebSocket(show)

	dd1 := gw.WsElemNew("id_1")
	dd1.AttachWebSocket(show)

	lb1 := gw.WsElemNew("id_6")
	lb1.AttachWebSocket(show)

	h21 := gw.WsElemNew("id_7")
	h21.AttachWebSocket(show)

	sl2 := gw.WsElemNew("id_4")
	sl2.AttachWebSocket(func(message string) {
		Printf("recv: %s\n", message)
		lb1.SetInnerText(message)
	})

	ta1 := gw.WsElemNew("id_5")
	ta1.AttachWebSocket(func(message string) {})

	bt1 := gw.WsElemNew("id_2")
	bt1.AttachWebSocket(func(message string) {
		Printf("recv: %s\n", message)
		err := ta1.WriteTextArea("fgfgfgfgfg")
		if err != nil {
			Println("error ", err)
		}

		h21.SetBackgroundColor("red")
	})

	cv1 := gw.WsElemNew("id_44")
	cv1.AttachWebSocket(func(message string) {})

	buffer := make([]byte, 400*400)
	for i := 0; i < 400*400; i++ {
		buffer[i] = byte(i % 255)
	}
	bt2 := gw.WsElemNew("id_41")
	bt2.AttachWebSocket(func(message string) {
		cv1.DrawCanvas(buffer)
	})

	gw.StartServer()

	for {
	}

}
