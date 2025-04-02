//go:build bulma
// +build bulma

package main

import "fmt"

const SIZE = "is-size-4" //"is-medium"

func (dom *Dom) Button(id string, text string) Elem {
	elem := dom.newElem(id, "button")
	elem.SetInnerText(text)
	elem.jsValue.Call("setAttribute", "type", "button")
	class := fmt.Sprintf("button %s %s", "is-primary", SIZE)
	elem.jsValue.Call("setAttribute", "class", class)
	return elem
}

func (dom *Dom) Tab() Elem {
	div1 := dom.newElem("", "div")
	class := fmt.Sprintf("tabs %s", SIZE)
	div1.jsValue.Call("setAttribute", "class", class)
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

func (elem *Elem) SetElemSize() {
	//read current class
	cCurJs := elem.jsValue.Call("getAttribute", "class")
	cCur := cCurJs.String()
	cNew := SIZE
	if len(cCur) > 0 {
		cNew = fmt.Sprintf("%s %s", cCur, SIZE)
	}

	elem.jsValue.Call("setAttribute", "class", cNew)
}

func (elem *Elem) SetThemeDark() {
	htmlJs := elem.dom.doc.Call("getElementById", "id_1000")
	htmlJs.Call("setAttribute", "class", "theme-dark")
}
