// GOOS=js GOARCH=wasm go build -o main.wasm test1.go
// tinygo build -o main.wasm test1.go
package main

import (
	gw "github.com/xxx/gwacln"
)

func main() {

	c := make(chan struct{}, 0)

	// read configuration sent by the server or define the GUI here
	gw.GwauiInit("test", true, true)

	//auto contained version, requires a server
	/*
		dom := gw.GwauiInit("test", true, false)

		dom.ConsoleLog("GUI started")

		tb1 := dom.Tab()
		tb1.AppendToBody()

		tbc1 := dom.Tabcontent(tb1, "", "tab_1")
		tbc2 := dom.Tabcontent(tb1, "", "tab_2")

		p1 := dom.Paragraph("", "")
		p2 := dom.Paragraph("", "")
		h21 := dom.Header2("", "header2")
		gr1 := dom.GridRow([]gw.Elem{p1, h21, p2})

		dd1 := dom.Dropdown("", []string{"one", "two", "three", "four"}, 1)

		//dd1.WsReadForDropDown()
		dd1.OnChange(func() {
			value := dd1.Value()
			dom.ConsoleLog(value)
		})

		b1 := dom.Button("", "Press")
		b1.OnClick(func() {
			h21.SetBackgroundColor("green")
		})
		gr2 := dom.GridRow([]gw.Elem{dd1, b1})

		ck1 := dom.Switchbox("", true)
		ck1.OnClick(func() {
			value := ck1.IsChecked()
			if value {
				dom.ConsoleLog("true")
			} else {
				dom.ConsoleLog("false")
			}
		})

		l1 := dom.Label("", "empty")

		i1 := dom.Input("")
		i1.OnInput(func() {
			value := i1.Value()
			dom.ConsoleLog(value)
			l1.SetInnerText(value)
		})
		gr3 := dom.GridRow([]gw.Elem{ck1, l1, i1})

		f1 := dom.Form("", "Submit")
		f1.OnSubmit(func() {
			value := f1.Value()
			dom.ConsoleLog(value)
		})

		sl1 := dom.Slider("", 10, 100, 30)
		l2 := dom.Label("", "30")

		sl1.OnChange(func() {
			value := sl1.Value()
			l2.SetInnerText(value)
		})
		gr4 := dom.GridRow([]gw.Elem{f1, sl1, l2})

		tbc1.Append(gr1)
		tbc1.Append(gr2)
		tbc2.Append(gr3)
		tbc2.Append(gr4)
		tbc1.AppendToBody()
		tbc2.AppendToBody()
	*/

	<-c
}
