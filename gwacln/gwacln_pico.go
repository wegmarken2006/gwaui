//go:build pico
// +build pico

package main

func (dom *Dom) Button(id string, text string) Elem {
	elem := dom.newElem(id, "button")
	elem.SetInnerText(text)
	elem.jsValue.Call("setAttribute", "type", "button")
	elem.jsValue.Call("setAttribute", "class", "primary")
	return elem
}

func (dom *Dom) Tab() Elem {
	div1 := dom.newElem("", "div")
	div1.jsValue.Call("setAttribute", "class", "tab")

	return div1
}

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

func (elem *Elem) SetElemSize() {
	//do nothing
}

func (elem *Elem) SetThemeDark() {
}
