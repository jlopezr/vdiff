package main

import (
	"fmt"
	ui "github.com/VladimirMarkelov/clui"
	term "github.com/nsf/termbox-go"
)

type Column struct {
	position int
	text     string
}

func FormatItem(maxSize int, firstString string, columns ...Column) string {
	return firstString[0:maxSize]
}

func main69() {
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	view := ui.AddWindow(0, 0, 0, 2, "<c:bright blue>Dirdiff 0.1<c:default>")
	view.SetAlign(ui.AlignCenter)
	view.SetPack(ui.Vertical)
	view.SetGaps(0, 0)
	view.SetPaddings(1, 1)
	view.SetMaximized(true)

	frameTop := ui.CreateFrame(view, ui.AutoSize, ui.AutoSize, ui.BorderThin, 1)
	frameTop.SetPack(ui.Horizontal)
	frameBottom := ui.CreateFrame(view, ui.AutoSize, 1, ui.BorderNone, ui.Fixed)
	frameBottom.SetPack(ui.Horizontal)

	left := ui.CreateListBox(frameTop, 1, ui.AutoSize, 1)
	right := ui.CreateListBox(frameTop, 1, ui.AutoSize, 1)

	t := fmt.Sprintf("<%-4.4s>", "PELO")
	left.AddItem(t)
	left.AddItem(t)
	t = fmt.Sprintf("<%-4.4s>", "PE")
	left.AddItem(t)
	left.AddItem(t)

	for i := 0; i <= 50; i++ {
		//txt := FormatItem(6, fmt.Sprintf("ITEM %d  | <c:red>1223239<c:default> | rwx-rwx-rwx ", i))
		//txt := FormatItem(10, fmt.Sprintf("ITEM %d", i), Column{10, "HOLA"}, Column{20, "ADIOS"})
		txt := fmt.Sprintf("%s %d", "ITEM", i)
		left.AddItem(txt)
		right.AddItem(txt)
	}

	label1 := ui.CreateLabel(frameBottom, ui.AutoSize, ui.AutoSize, "HELLO", 1)
	label1.SetAlign(ui.AlignCenter)
	label2 := ui.CreateLabel(frameBottom, ui.AutoSize, ui.AutoSize, "WORLD", 1)
	label2.SetAlign(ui.AlignCenter)

	i := 0
	label1.OnActive(func(active bool) {
		label2.SetTitle(fmt.Sprintf("ACTION %d!", i))
		i++
	})

	label2.OnActive(func(active bool) {
		label1.SetTitle("ACTION 2!")
	})

	left.OnSelectItem(func(ev ui.Event) {
		item := left.SelectedItem()
		txt := fmt.Sprintf("SELECTED <c:red>%d<c:default> ITEM", item)
		label1.SetTitle(txt)
		right.SelectItem(item)
		right.EnsureVisible()
	})

	right.OnSelectItem(func(ev ui.Event) {
		item := right.SelectedItem()
		txt := fmt.Sprintf("SELECTED %d", item)
		label1.SetTitle(txt)
		left.SelectItem(item)
		left.EnsureVisible()
	})

	left.OnKeyPress(func(key term.Key) bool {
		label1.SetTitle(fmt.Sprintf("KEY"))
		return false
	})

	/*
		btnOk := ui.CreateButton(frameBottom, 0, 1, "OK", 1)
		btnOk.SetShadowType(ui.ShadowNone)
		btnOk.OnClick(func(ev ui.Event) {
			label1.SetTitle("OK!!")
		})

		btnQuit := ui.CreateButton(frameBottom, 0, 1, "SALIR", 1)
		btnQuit.OnClick(func(ev ui.Event) {
			go ui.Stop()
		})
	*/

	ui.MainLoop()
}
