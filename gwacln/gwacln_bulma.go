//go:build bulma
// +build bulma

package main

func (dom *Dom) Button(id string, text string) Elem {
	elem := dom.newElem(id, "button")
	elem.SetInnerText(text)
	elem.jsValue.Call("setAttribute", "type", "button")
	elem.jsValue.Call("setAttribute", "class", "button is-primary")
	return elem
}

func (dom *Dom) Tab() Elem {
	div1 := dom.newElem("", "div")
	div1.jsValue.Call("setAttribute", "class", "tabs is-medium")
	ul1 := dom.newElem("", "ul")
	div1.Append(ul1)
	div1.child1 = ul1.jsValue

	return div1
}

func (dom *Dom) Tabcontent(tab Elem, id string, title string) Elem {
	li := dom.newElem("", "li")
	a := dom.newElem("", "a")
	//bt.jsValue.Call("setAttribute", "class", "tablinks outline secondary")
	li.Append(a)
	a.SetInnerText(title)
	div2 := dom.newElem(id, "div")
	div2.jsValue.Call("setAttribute", "class", "tabcontent")
	div2.child1 = li.jsValue
	li.OnClick(func() {
		div2.enableThisTab()
	})
	div2.enableThisTabIfFirst()
	//save Id
	dom.tabs = append(dom.tabs, div2)
	tab.child1.Call("appendChild", li.jsValue)

	return div2
}

func (elem *Elem) enableThisTab() {
	for _, tab := range elem.dom.tabs {
		if tab.id == elem.id {
			tab.jsValue.Get("style").Call("setProperty", "display", "block")
			tab.child1.Call("setAttribute", "class", "is-active")
		} else {
			tab.jsValue.Get("style").Call("setProperty", "display", "none")
			tab.child1.Call("setAttribute", "class", "")
		}
	}
}
func (elem *Elem) enableThisTabIfFirst() {
	if len(elem.dom.tabs) == 0 {
		elem.jsValue.Get("style").Call("setProperty", "display", "block")
		elem.child1.Call("setAttribute", "class", "is-active")
	} else {
		elem.jsValue.Get("style").Call("setProperty", "display", "none")
		elem.child1.Call("setAttribute", "class", "")
	}
}
