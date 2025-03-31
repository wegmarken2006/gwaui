// GOOS=js GOARCH=wasm go build -o main.wasm .
// tinygo build -o main.wasm .
// tinygo build --no-debug -o main.wasm .
package main

import (
	"fmt"
	"syscall/js"
	"time"

	"gopkg.in/yaml.v3"
)

type Dom struct {
	doc        js.Value
	wind       js.Value
	console    js.Value
	body       js.Value
	idCount    int
	bodyScript string
	headStyle  string
	Debug      bool
	tabs       []Elem
}

type Elem struct {
	jsValue  js.Value
	children int
	child1   js.Value
	dom      *Dom
	id       string
	ws       js.Value
	canvas   CanvasElem
}
type CanvasElem struct {
	ctx    js.Value
	width  int
	height int
}

func main() {
	c := make(chan struct{}, 0)

	GwauiInit("app", true, true)

	<-c
}

func GwauiInit(titleText string, debug bool, waitForWebSocket bool) Dom {
	dom := Dom{
		doc:     js.Global().Get("document"),
		wind:    js.Global().Get("window"),
		console: js.Global().Get("window").Get("console"),
		body:    js.Global().Get("document").Get("body"),
		idCount: 300,
		Debug:   debug,
	}

	//first element
	title := dom.Title("id_0", titleText)
	title.AppendToBody()

	if waitForWebSocket {
		title.AddWebSocket()
		title.WsReadConfiguration()

		//send to ackonwledge websocket ready
		for ind := 0; ind < 50; ind++ {
			if title.IsWsOpen() {
				title.WsWrite("title")
				break
			} else {
				time.Sleep(50 * time.Millisecond)
			}
		}
	}

	return dom
}

func (dom *Dom) addBodyScript(text string) {
	dom.bodyScript = fmt.Sprintf("%s%s", dom.bodyScript, text)
	scr := dom.doc.Call("createElement", "script")
	scr.Set("innerText", text)
	dom.body.Call("appendChild", scr)
}

func (dom *Dom) addHeadStyle(text string) {
	dom.headStyle = fmt.Sprintf("%s%s", dom.headStyle, text)
	style := dom.doc.Call("createElement", "style")
	style.Set("innerText", text)
	dom.doc.Get("head").Call("appendChild", style)
}

func (dom *Dom) ConsoleLog(text string) {
	dom.console.Call("log", text)
}

func (dom *Dom) getElementByID(id string) js.Value {
	return dom.doc.Call("getElementById", id)
}

func (dom *Dom) createElement(id string) js.Value {
	return dom.doc.Call("createElement", id)
}

func (dom *Dom) newElem(id string, elType string) Elem {
	jsVal := dom.createElement(elType)
	var elem Elem
	if len(id) > 0 {
		elem = Elem{dom: dom, jsValue: jsVal, id: id, children: 0}
	} else {
		elem = Elem{dom: dom, jsValue: jsVal, id: dom.getNewId(), children: 0}
	}

	elem.jsValue.Set("id", elem.id)
	elem.showId()
	return elem
}

func (dom *Dom) Canvas(id string, width int, height int) Elem {
	div := dom.newElem(id, "div")
	div.jsValue.Get("style").Call("setProperty", "text-align", "center")

	cv := dom.newElem("", "canvas")
	widthStr := fmt.Sprintf("%d", width)
	heightStr := fmt.Sprintf("%d", height)
	cv.jsValue.Call("setAttribute", "width", widthStr)
	cv.jsValue.Call("setAttribute", "height", heightStr)
	cv.jsValue.Get("style").Call("setProperty", "border", "1px solid #000000")

	div.canvas.ctx = cv.jsValue.Call("getContext", "2d")
	div.canvas.width = width
	div.canvas.height = height

	div.Append(cv)

	return div
}

/*
func (dom *Dom) Button(id string, text string) Elem {
	elem := dom.newElem(id, "button")
	elem.SetInnerText(text)
	elem.jsValue.Call("setAttribute", "type", "button")
	elem.jsValue.Call("setAttribute", "class", "primary")
	return elem
}
*/

func (dom *Dom) Label(id string, text string) Elem {
	elem := dom.newElem(id, "label")
	elem.SetInnerText(text)
	elem.SetElemSize()
	return elem
}

func (dom *Dom) Input(id string, text string) Elem {
	elem := dom.newElem(id, "input")
	elem.jsValue.Call("setAttribute", "type", "text")
	if len(text) > 0 {
		elem.jsValue.Call("setAttribute", "placeholder", text)
	}
	elem.SetElemSize()
	return elem
}

func (dom *Dom) Date(id string) Elem {
	elem := dom.newElem(id, "input")
	elem.jsValue.Call("setAttribute", "type", "date")
	elem.SetElemSize()
	return elem
}

func (dom *Dom) Slider(id string, min int, max int, value int) Elem {
	elem := dom.newElem(id, "input")
	elem.jsValue.Call("setAttribute", "type", "range")
	elem.jsValue.Call("setAttribute", "min", fmt.Sprintf("%d", min))
	elem.jsValue.Call("setAttribute", "max", fmt.Sprintf("%d", max))
	elem.jsValue.Call("setAttribute", "value", fmt.Sprintf("%d", value))
	elem.SetElemSize()
	return elem
}

func (dom *Dom) Switchbox(id string, initial bool) Elem {
	elem := dom.newElem(id, "input")
	elem.jsValue.Call("setAttribute", "type", "checkbox")
	elem.jsValue.Call("setAttribute", "role", "switch")
	if initial {
		elem.jsValue.Call("setAttribute", "checked", "true")
	} else {
		elem.jsValue.Call("setAttribute", "checked", "false")
	}
	elem.SetElemSize()
	return elem
}

func (dom *Dom) Form(id string, btText string) Elem {
	elem := dom.newElem(id, "form")
	i1 := dom.newElem("", "input")
	i1.jsValue.Call("setAttribute", "type", "text")
	b1 := dom.Button("", btText)
	b1.jsValue.Call("setAttribute", "type", "submit")

	grid := dom.GridRow([]Elem{i1, b1})
	elem.Append(grid)
	elem.child1 = i1.jsValue
	elem.children = 1
	elem.SetElemSize()
	return elem
}

func (dom *Dom) Title(id string, title string) Elem {
	elem := dom.newElem(id, "title")
	elem.SetInnerText(title)
	return elem
}

func (dom *Dom) Alert(text string) {
	alert := js.Global().Get("alert")
	alert.Invoke(text)
}

func (dom *Dom) Paragraph(id string, text string) Elem {
	elem := dom.newElem(id, "p")
	elem.SetInnerText(text)
	return elem
}

func (dom *Dom) Empty() Elem {
	elem := dom.newElem("", "p")
	elem.SetInnerText("")
	return elem
}

func (dom *Dom) GridRow(children []Elem) Elem {
	elem := dom.newElem("", "div")
	elem.jsValue.Call("setAttribute", "class", "grid")
	for _, child := range children {
		elem.Append(child)
	}

	return elem
}

func (dom *Dom) Header1(id string, text string) Elem {
	elem := dom.newElem(id, "h1")
	elem.SetInnerText(text)
	return elem
}

func (dom *Dom) Header2(id string, text string) Elem {
	elem := dom.newElem(id, "h2")
	elem.SetInnerText(text)
	elem.SetElemSize()
	elem.jsValue.Get("style").Call("setProperty", "text-align", "center")
	return elem
}

/*
func (dom *Dom) Tab() Elem {
	div1 := dom.newElem("", "div")
	div1.jsValue.Call("setAttribute", "class", "tab")

	return div1
}
*/

func (dom *Dom) Image(id string, fileName string) Elem {
	div1 := dom.newElem(id, "img")
	if len(fileName) > 0 {
		div1.jsValue.Call("setAttribute", "src", fileName)
	}

	return div1
}

/*
func (dom *Dom) Tabcontent(tab Elem, id string, title string) Elem {
	bt := dom.newElem("", "button")
	bt.jsValue.Call("setAttribute", "class", "tablinks outline secondary")
	bt.SetInnerText(title)
	div2 := dom.newElem(id, "div")
	div2.jsValue.Call("setAttribute", "class", "tabcontent")
	div2.child1 = bt.jsValue
	bt.OnClick(func() {
		div2.enableThisTab()
	})
	div2.enableThisTabIfFirst()
	//save Id
	dom.tabs = append(dom.tabs, div2)

	tab.Append(bt)

	return div2
}
*/

func (dom *Dom) Plot(id string) Elem {
	elem := dom.newElem(id, "div")
	elem.jsValue.Get("style").Call("setProperty", "text-align", "center")

	return elem
}

func (elem *Elem) DrawPlot(pConf *plotConf) {

	typ := pConf.Typ

	// Create the traces
	var xs []interface{}
	var ys [][]interface{}
	var mode string
	var typeS string

	if typ == "lines" || typ == "scatter" {
		for _, item := range pConf.X {
			xs = append(xs, item)
		}
		for _, item := range pConf.Y {
			var ys2 []interface{}
			for _, yelem := range item {
				ys2 = append(ys2, yelem)
			}
			ys = append(ys, ys2)
		}

		if typ == "scatter" {
			mode = "markers"
		} else {
			mode = "lines"
		}
		typeS = "scatter"
	} else if typ == "bar" {
		for _, item := range pConf.X_cat {
			xs = append(xs, item)
		}
		for _, item := range pConf.Y {
			var ys2 []interface{}
			for _, yelem := range item {
				ys2 = append(ys2, yelem)
			}
			ys = append(ys, ys2)
		}
		typeS = typ
		mode = "s"
	} else if typ == "box" {
		for _, item := range pConf.Y {
			var ys2 []interface{}
			for _, yelem := range item {
				ys2 = append(ys2, yelem)
			}
			ys = append(ys, ys2)
		}
		typeS = typ
		mode = ""
	}

	//if no names are passed
	if len(pConf.Names) != len(ys) {
		pConf.Names = nil
		for ind, _ := range ys {
			pConf.Names = append(pConf.Names, fmt.Sprintf("t%d", ind))
		}
	}

	var dataV []interface{}
	if typ == "box" {
		for ind, yelem := range ys {
			data1 := js.ValueOf(map[string]interface{}{
				"y":    yelem,
				"name": pConf.Names[ind],
				"mode": mode,
				"type": typeS,
			})
			dataV = append(dataV, data1)
		}

	} else {
		for ind, yelem := range ys {
			data1 := js.ValueOf(map[string]interface{}{
				"x": xs,
				"y": yelem,
				//"x":    []interface{}{1., 2., 3., 4.},
				//"y":    []interface{}{12., 9., 15., 12.},
				"name": pConf.Names[ind],
				"mode": mode,
				"type": typeS,
			})
			dataV = append(dataV, data1)
		}

	}

	layout1 := js.ValueOf(map[string]interface{}{
		"title":  pConf.Title, //"title",
		"width":  pConf.Width,
		"height": pConf.Height,
	})

	var dataI []interface{}

	for _, item := range dataV {
		dataI = append(dataI, item)
	}

	data := js.ValueOf(dataI)
	//layout := js.ValueOf([]interface{}{layout1})

	// Call the Plotly.newPlot function

	js.Global().Get("Plotly").Call(
		"newPlot",
		elem.jsValue,
		data,
		layout1,
	)
}

/*
func (elem *Elem) enableThisTab() {
	for _, tab := range elem.dom.tabs {
		if tab.id == elem.id {
			tab.jsValue.Get("style").Call("setProperty", "display", "block")
			tab.child1.Call("setAttribute", "class", "tablinks outline active")
		} else {
			tab.jsValue.Get("style").Call("setProperty", "display", "none")
			tab.child1.Call("setAttribute", "class", "tablinks outline secondary")
		}
	}
}
func (elem *Elem) enableThisTabIfFirst() {
	if len(elem.dom.tabs) == 0 {
		elem.jsValue.Get("style").Call("setProperty", "display", "block")
		elem.child1.Call("setAttribute", "class", "tablinks outline active")
	} else {
		elem.jsValue.Get("style").Call("setProperty", "display", "none")
		elem.child1.Call("setAttribute", "class", "tablinks outline secondary")
	}
}
*/

func (dom *Dom) TextArea(id string, lines int, text string) Elem {
	elem := dom.newElem(id, "textarea")
	elem.SetInnerText(text)
	sLines := fmt.Sprintf("%d", lines)
	elem.jsValue.Call("setAttribute", "rows", sLines)
	elem.SetBackgroundColor("black")
	elem.SetColor("white")
	elem.SetElemSize()

	return elem
}

func (dom *Dom) Dropdown(id string, choices []string, defaultInd int) Elem {
	sel := dom.newElem(id, "select")
	sel.SetElemSize()

	op1 := dom.newElem("", "option")
	op1.jsValue.Call("setAttribute", "value", "")

	if len(choices) > 0 && defaultInd < len(choices) {
		sel.jsValue.Call("setAttribute", "aria-label", choices[defaultInd])
		op1.SetInnerText(choices[defaultInd])
	}

	sel.Append(op1)

	for ind, choice := range choices {
		if ind == defaultInd {
			continue
		} else {
			op := dom.newElem("", "option")
			op.SetInnerText(choice)
			sel.Append(op)
		}
	}

	return sel
}

/*
func (dom *Dom) col() Elem {
	elem := dom.newElem("div")
	elem.jsValue.Call("setAttribute", "class", "column")
	return elem
}

func (dom *Dom) row() Elem {
	elem := dom.newElem("div")
	elem.jsValue.Call("setAttribute", "class", "row")
	return elem
}
*/

func (dom *Dom) getNewId() string {
	id := dom.idCount
	dom.idCount = id + 1
	return fmt.Sprintf("id_%d", id)
}

func (elem *Elem) showId() {
	if elem.dom.Debug {
		elem.jsValue.Call("setAttribute", "title", elem.id)
	}
}

func (elem *Elem) SetInnerText(text string) {
	elem.jsValue.Set("innerText", text)
}

func (elem *Elem) SetBackgroundColor(color string) {
	elem.jsValue.Get("style").Call("setProperty", "background-color", color)
}

func (elem *Elem) SetColor(color string) {
	elem.jsValue.Get("style").Call("setProperty", "color", color)
}

func (elem *Elem) SetItemsList(lst []string) {
	dom := elem.dom
	defaultInd := 0

	for ind, item := range lst {
		op := dom.newElem("", "option")
		op.SetInnerText(item)
		if ind == defaultInd {
			op.jsValue.Call("setAttribute", "selected", "")
		}
		elem.Append(op)
	}

}

func (elem *Elem) setMultipleCols(colsNr int) {
	elem.jsValue.Get("style").Call("setProperty", "float", "left")
	perc := fmt.Sprintf("%d%c", int(100/colsNr), '%')
	elem.jsValue.Get("style").Call("setProperty", "width", perc)
}

func (elem *Elem) setEventListener(fun func()) string {
	name := fmt.Sprintf("fun_%s", elem.id)
	onEventName := fmt.Sprintf("fun_%s()", elem.id)

	jsFun := func(this js.Value, inputs []js.Value) interface{} {
		//elem.dom.consoleLog("Inside func")
		fun()
		return nil
	}

	js.Global().Set(name, js.FuncOf(jsFun))

	return onEventName
}

func (elem *Elem) OnClick(fun func()) {
	elem.jsValue.Call("setAttribute", "onclick", elem.setEventListener(fun))
}

func (elem *Elem) OnInput(fun func()) {
	elem.jsValue.Call("setAttribute", "oninput", elem.setEventListener(fun))
}

func (elem *Elem) OnKeyPress(fun func()) {
	elem.jsValue.Call("setAttribute", "onkeypress", elem.setEventListener(fun))
}

func (elem *Elem) OnChange(fun func()) {
	elem.jsValue.Call("setAttribute", "onchange", elem.setEventListener(fun))
}

func (elem *Elem) OnSubmit(fun func()) {
	name := fmt.Sprintf("fun_%s", elem.id)
	onClickName := fmt.Sprintf("fun_%s(); return false;", elem.id)

	jsFun := func(this js.Value, inputs []js.Value) interface{} {
		//elem.dom.consoleLog("Inside func")
		fun()
		return nil
	}

	js.Global().Set(name, js.FuncOf(jsFun))
	elem.jsValue.Call("setAttribute", "onsubmit", onClickName)
}

func (elem *Elem) AppendToBody() {
	elem.dom.body.Call("appendChild", elem.jsValue)
}

func (elem *Elem) Append(child Elem) {
	elem.jsValue.Call("appendChild", child.jsValue)
}

func (elem *Elem) Value() string {
	if elem.children == 0 {
		return elem.jsValue.Get("value").String()
	} else {
		return elem.child1.Get("value").String()
	}
}

func (elem *Elem) IsChecked() bool {
	return elem.jsValue.Get("checked").Bool()
}

func (elem *Elem) Tooltip(text string) {
	elem.jsValue.Call("setAttribute", "data-tooltip", text)
	elem.jsValue.Call("setAttribute", "data-placement", "right")
}

func (elem *Elem) ShowImage(imageName string) {
	elem.jsValue.Call("setAttribute", "src", imageName)
}

/*
func (elem *Elem) PutImage(pngName string) {

	elem.canvas.ctx.Call("clearRect", 0, 0, elem.canvas.width, elem.canvas.height)
	//imageData := elem.canvas.ctx.Call("getImageData", 0, 0, elem.canvas.width, elem.canvas.height)
	imageData := elem.canvas.ctx.Call("createImageData", elem.canvas.width, elem.canvas.height)
	data := imageData.Get("data")
	len := data.Get("length").Int()

	goData := make([]uint8, len)

	for i := 0; i < len; i++ {
		goData[i] = 255 - goData[i]     // red
		goData[i+1] = 255 - goData[i+1] // green
		goData[i+2] = 255 - goData[i+2] // blue

	}
	js.CopyBytesToJS(data, goData)

	elem.canvas.ctx.Call("putImageData", imageData, 0, 0)
}
*/

func (elem *Elem) AddWebSocket() {
	id := elem.id
	host := elem.dom.doc.Get("location").Get("host")
	addr := fmt.Sprintf("ws://%s/%s", host, id)
	ws := js.Global().Get("WebSocket").New(addr)

	elem.dom.ConsoleLog(addr) //TODO
	elem.ws = ws
}

func (elem *Elem) WsRead() {

	jsFun := func(this js.Value, inputs []js.Value) interface{} {
		var rxMsg rxTxMessage

		event := inputs[0]

		edata := event.Get("data").String()
		yaml.Unmarshal([]byte(edata), &rxMsg)

		if len(rxMsg.BackgroundColor) > 0 {
			elem.SetBackgroundColor(rxMsg.BackgroundColor)
		}
		if len(rxMsg.Color) > 0 {
			elem.SetColor(rxMsg.Color)
		}

		if len(rxMsg.Text) > 0 {
			elem.SetInnerText(rxMsg.Text)
		}
		if len(rxMsg.Textarea) > 0 {
			textValue := elem.Value() //current content
			str := fmt.Sprintf("%s%s", textValue, rxMsg.Textarea)

			diff := len(str) - 4096
			if diff > 0 {
				textValue = str[:diff] //str.slice(diff);
			} else {
				textValue = str
			}
			elem.SetInnerText(textValue)
		}
		if len(rxMsg.ItemList) > 0 {
			elem.SetItemsList(rxMsg.ItemList)
		}
		if len(rxMsg.ImageName) > 0 {
			elem.ShowImage(rxMsg.ImageName)
		}
		if rxMsg.PlotConf != nil {
			elem.DrawPlot(rxMsg.PlotConf)
		}
		if len(rxMsg.FuncToCall) > 0 {
			if rxMsg.FuncToCall == "SetThemeDark" {
				elem.SetThemeDark()
			}
			// refletion doesn't work
			//fn := reflect.ValueOf(elem).MethodByName(rxMsg.FuncToCall)
			//fn.Call(nil)
		}

		return nil
	}

	elem.ws.Call("addEventListener", "message", js.FuncOf(jsFun))
}

/*
func (elem *Elem) WsReadForTextArea() {

		jsFun := func(this js.Value, inputs []js.Value) interface{} {
			event := inputs[0]
			edata := event.Get("data").String()

			textValue := elem.Value() //current content
			str := fmt.Sprintf("%s%s", textValue, edata)

			diff := len(str) - 4096
			if diff > 0 {
				textValue = str[:diff] //str.slice(diff);
			} else {
				textValue = str
			}

			elem.SetInnerText(textValue)
			return nil
		}

		elem.ws.Call("addEventListener", "message", js.FuncOf(jsFun))
	}

func (elem *Elem) WsReadForGenericText() {

		jsFun := func(this js.Value, inputs []js.Value) interface{} {
			event := inputs[0]
			edata := event.Get("data").String()

			elem.SetInnerText(edata)
			return nil
		}

		elem.ws.Call("addEventListener", "message", js.FuncOf(jsFun))
	}
*/

func (elem *Elem) WsReadConfiguration() {

	jsFun := func(this js.Value, inputs []js.Value) interface{} {
		event := inputs[0]
		edata := event.Get("data").String()
		var guiDescr []GuiDescr
		yaml.Unmarshal([]byte(edata), &guiDescr)
		tabValid := false
		var tb Elem

		for _, tabs := range guiDescr {

			//for _, grid := range row.GridRow {
			rows := tabs.Tab.Row
			var tbc Elem

			if len(tabs.Tab.Id) > 0 && !tabValid {
				tb = elem.dom.Tab()
				tb.AppendToBody()
				tabValid = true
			}
			if len(tabs.Tab.Id) > 0 {
				tbc = elem.dom.Tabcontent(tb, tabs.Tab.Id, tabs.Tab.Text)
			}

			for _, row := range rows {
				elems := []Elem{}
				for _, grid := range row.GridRow {
					if len(grid.Button.Id) > 0 {
						bt := elem.dom.Button(grid.Button.Id, grid.Button.Text)
						bt.AddWebSocket()
						bt.OnClick(func() { bt.WsWrite("pressed") })
						elems = append(elems, bt)
					}
					if len(grid.Dropdown.Id) > 0 {
						dd := elem.dom.Dropdown(grid.Dropdown.Id, grid.Dropdown.Items, grid.Dropdown.DefaultInd)
						dd.AddWebSocket()
						dd.WsRead()
						dd.OnChange(func() {
							value := dd.Value()
							dd.WsWrite(value)
						})
						elems = append(elems, dd)
					}
					if len(grid.Form.Id) > 0 {
						fm := elem.dom.Form(grid.Form.Id, grid.Form.Text)
						fm.AddWebSocket()
						fm.OnSubmit(func() {
							value := fm.Value()
							fm.WsWrite(value)
						})
						elems = append(elems, fm)
					}
					if len(grid.Input.Id) > 0 {
						ip := elem.dom.Input(grid.Input.Id, grid.Input.Text)
						ip.AddWebSocket()
						ip.OnChange(func() {
							value := ip.Value()
							ip.WsWrite(value)
						})
						elems = append(elems, ip)
					}
					if len(grid.Date.Id) > 0 {
						dt := elem.dom.Date(grid.Date.Id)
						dt.AddWebSocket()
						dt.OnChange(func() {
							value := dt.Value()
							dt.WsWrite(value)
						})
						elems = append(elems, dt)
					}
					if len(grid.Slider.Id) > 0 {
						sl := elem.dom.Slider(grid.Slider.Id, grid.Slider.MinMaxIni[0], grid.Slider.MinMaxIni[1], grid.Slider.MinMaxIni[2])
						sl.AddWebSocket()
						sl.OnChange(func() {
							value := sl.Value()
							sl.WsWrite(value)
						})
						elems = append(elems, sl)
					}
					if len(grid.Textarea.Id) > 0 {
						ta := elem.dom.TextArea(grid.Textarea.Id, grid.Textarea.Lines, grid.Textarea.Text)
						ta.AddWebSocket()
						//ta.WsReadForTextArea()
						ta.WsRead()
						elems = append(elems, ta)
					}
					if len(grid.Label.Id) > 0 {
						lb := elem.dom.Label(grid.Label.Id, grid.Label.Text)
						if grid.Label.Mutable {
							lb.AddWebSocket()
							lb.WsRead()
						}
						elems = append(elems, lb)
					}
					if len(grid.H2.Id) > 0 {
						h2 := elem.dom.Header2(grid.H2.Id, grid.H2.Text)
						if grid.H2.Mutable {
							h2.AddWebSocket()
							h2.WsRead()
						}
						elems = append(elems, h2)
					}
					if grid.Paragraph != nil {
						id := ""
						if len(grid.Paragraph.Id) > 0 {
							id = grid.Paragraph.Id
						}
						par := elem.dom.Paragraph(id, grid.Paragraph.Text)
						elems = append(elems, par)
					}
					if len(grid.Canvas.Id) > 0 {
						cv := elem.dom.Canvas(grid.Canvas.Id, grid.Canvas.Width, grid.Canvas.Height)
						cv.AddWebSocket()
						cv.WsRead()
						elems = append(elems, cv)
					}
					if len(grid.Image.Id) > 0 {
						img := elem.dom.Image(grid.Image.Id, "")
						img.AddWebSocket()
						img.WsRead()
						elems = append(elems, img)
					}
					if len(grid.Plot.Id) > 0 {
						plt := elem.dom.Plot(grid.Plot.Id)
						plt.AddWebSocket()
						plt.WsRead()
						elems = append(elems, plt)
					}
				}
				gd := elem.dom.GridRow(elems)
				if tabValid {
					tbc.Append(gd)
					tbc.AppendToBody()
				} else {
					gd.AppendToBody()
				}
			}
		}

		return nil
	}

	elem.ws.Call("addEventListener", "message", js.FuncOf(jsFun))
}

func (elem *Elem) IsWsOpen() bool {
	status := elem.ws.Get("readyState").Int()
	if status != 1 {
		return false
	} else {
		// WebSocket.OPEN
		return true
	}
}

func (elem *Elem) WsWrite(toWrite string) {
	elem.ws.Call("send", toWrite)
}
