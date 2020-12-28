package main

import (
	"fmt"
	ui "github.com/VladimirMarkelov/clui"
	term "github.com/nsf/termbox-go"
	"strconv"
	"time"

	//term "github.com/nsf/termbox-go"
)

func createDirInfo() DirInfo {
	dir := DirInfo{
		LeftPath:  "/Users/juan/a",
		RightPath: "/Users/juan/b",
		Files:     make([]*DirEntry, 0),
	}

	for i := 0; i < 50; i++ {
		dir.AppendEntry(fmt.Sprintf("Item %d", i))
	}

	return dir
}

func createPanel(view ui.Control, state *PanelState, isLeft bool) *ui.TableView {
	panel := ui.CreateTableView(view, 25, 12, 1)

	panel.SetShowLines(true)
	panel.SetShowRowNumber(false)
	panel.SetRowCount(state.currentDirInfo.EntryCount())

	cols := []ui.Column{
		ui.Column{Title: "Filename", Width: 25, Alignment: ui.AlignLeft},
		ui.Column{Title: "Hash", Width: 12, Alignment: ui.AlignLeft},
		ui.Column{Title: "Size", Width: 10, Alignment: ui.AlignLeft},
		ui.Column{Title: "Modified", Width: 10, Alignment: ui.AlignLeft},
	}
	panel.SetColumns(cols)
	panel.SetFullRowSelect(true)

	panel.OnDrawCell(func(info *ui.ColumnDrawInfo) {
		if isLeft {
			info.Bg = term.ColorDefault
		} else {
			info.Bg = term.ColorDefault
		}

		if info.RowSelected {
			info.Bg = term.ColorLightGray
			info.Fg = term.ColorRed
		}
		entry := state.currentDirInfo.GetEntry(info.Row)
		switch info.Col {
		case 0:
			info.Text = entry.Name
			break
		case 1:
			info.Text = entry.GetInfo(isLeft).Hash
			break
		case 2:
			info.Text = strconv.FormatInt(entry.GetInfo(isLeft).Size, 10)
			break
		case 3:
			info.Text = entry.GetInfo(isLeft).LastModification.String()
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
				/*ev := ui.Event{
					// label.Draw() directamente no va
					// ui.EventChanged parece que no esta implementado
					// ui.EventRedraw se puede mandar a toda la pantalla (sin poner Target)
					// si pones el target parece que hace lo mismo :(
					// ui.RefreshScreen() parece que tambien pone el evento
					Type: ui.EventRedraw,
					Target: label,
				}
				ui.PutEvent(ev)
				*/
			}
		}
	}()
}

type PanelState struct {
	currentDirInfo *DirInfo
}

func main() {
	ui.InitLibrary()
	defer ui.DeinitLibrary()

	state := PanelState{
		//currentDirInfo: createDirInfo(),
		currentDirInfo: Prueba2(),
	}

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

	left := createPanel(frameTop, &state, true)
	right := createPanel(frameTop, &state, false)

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

	j := 0

	right.OnSelectCell(func(x int, y int) {
		left.SetSelectedRow(y)
		label2.SetTitle(fmt.Sprintf("%d - %d => %d", x, y, j))
		j++
	})

	/*
	left.OnKeyPress(func(key term.Key) bool {
		if key == term.KeyEnter {
			label2.SetTitle("ENTER PRESS!")
		}
		return false
	})
	*/

	left.OnAction(func(event ui.TableEvent) {
		if event.Action == ui.TableActionEdit {
			entry := state.currentDirInfo.Files[event.Row]
			label2.SetTitle(fmt.Sprintf("TABLE EVENT: %d C:%d R:%d [%s]", event.Action, event.Col, event.Row, entry.Name))
			if entry.Left.Type == DIRECTORY || entry.Right.Type == DIRECTORY || entry.Left.Type == UPDIR || entry.Right.Type == UPDIR {
				//TODO Move this code to a function
				state.currentDirInfo = entry.Info
				left.SetRowCount(state.currentDirInfo.EntryCount())
				left.SetSelectedCol(0)
				left.SetSelectedRow(0)
				right.SetRowCount(state.currentDirInfo.EntryCount())
				right.SetSelectedCol(0)
				right.SetSelectedRow(0)
				ui.RefreshScreen() //TODO Only refresh left and right panel
			}
		}
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

	//Example of goroutine modifying UI
	//go ModifyUI(label1)

	ui.MainLoop()
}
