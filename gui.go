package main

import (
	"fmt"
	ui "github.com/VladimirMarkelov/clui"
	term "github.com/nsf/termbox-go"
	"time"

	//term "github.com/nsf/termbox-go"
)

func createDirInfo() DirInfo {
	dir := DirInfo{
		leftPath:  "/Users/juan/a",
		rightPath: "/Users/juan/b",
		children:  make([]*EntryInfo, 0),
	}

	for i := 0; i < 50; i++ {
		dir.AppendEntry(fmt.Sprintf("Item %d", i))
	}

	return dir
}

func createPanel(view ui.Control, dirinfo DirInfo) *ui.TableView {
	panel := ui.CreateTableView(view, 25, 12, 1)

	panel.SetShowLines(true)
	panel.SetShowRowNumber(false)
	panel.SetRowCount(dirinfo.EntryCount())

	cols := []ui.Column{
		ui.Column{Title: "Filename", Width: 25, Alignment: ui.AlignLeft},
		ui.Column{Title: "Hash", Width: 12, Alignment: ui.AlignLeft},
		ui.Column{Title: "Size", Width: 10, Alignment: ui.AlignLeft},
		ui.Column{Title: "Modified", Width: 10, Alignment: ui.AlignLeft},
	}
	panel.SetColumns(cols)
	panel.SetFullRowSelect(true)

	panel.OnDrawCell(func(info *ui.ColumnDrawInfo) {
		if info.RowSelected {
			info.Bg = term.ColorLightGray
			info.Fg = term.ColorRed
		}
		switch info.Col {
		case 0:
			info.Text = dirinfo.GetEntry(info.Row).name
			break
		case 1:
			info.Text = "HASH"
			break
		case 2:
			info.Text = "SIZE"
			break
		case 3:
			info.Text = "MODIFIED"
			break
		}
	})

	return panel
}

func ModifyUI(label *ui.Label) {

	ticker := time.NewTicker(1000 * time.Millisecond)
	done := make(chan bool)
	go func() {
		i:=1
		for {
			select {
			case <-done:
				return
			case _ = <-ticker.C:
				label.SetTitle(fmt.Sprintf("%d sec", i))
				i++
				ui.RefreshScreen()
			}
		}
	}()
}

func main() {
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	dirinfo := createDirInfo()

	window := ui.AddWindow(0, 0, 0, 2, "<c:bright blue>Dirdiff 0.1<c:default>")
	window.SetAlign(ui.AlignCenter)
	window.SetPack(ui.Vertical)
	window.SetGaps(0, 0)
	window.SetPaddings(1, 1)
	window.SetMaximized(true)

	frameTop := ui.CreateFrame(window, ui.AutoSize, ui.AutoSize, ui.BorderThin, 1)
	frameTop.SetPack(ui.Horizontal)
	frameBottom := ui.CreateFrame(window, ui.AutoSize, 1, ui.BorderNone, ui.Fixed)
	frameBottom.SetPack(ui.Horizontal)

	left := createPanel(frameTop, dirinfo)
	right := createPanel(frameTop, dirinfo)

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

	left.OnSelectCell(func(x int, y int) {
		right.SetSelectedRow(y)
		label1.SetTitle(fmt.Sprintf("%d - %d", x, y))
	})

	right.OnSelectCell(func(x int, y int) {
		left.SetSelectedRow(y)
		label2.SetTitle(fmt.Sprintf("%d - %d", x, y))
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

	go ModifyUI(label1)

	ui.MainLoop()
}
