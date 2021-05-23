package frontend

import (
	"fmt"
	"twfinder/gui/server"
)

type myButtonHandler struct {
	counter int
	text    string
}

func (h *myButtonHandler) HandleEvent(e server.Event) {
	if b, isButton := e.Src().(server.Button); isButton {
		b.SetText(b.Text() + h.text)
		h.counter++
		b.SetToolTip(fmt.Sprintf("You've clicked %d times!", h.counter))
		e.MarkDirty(b)
	}
}

// Gui :
func Gui(name string) server.Window {

	// Create and build a window
	// win := server.NewWindow("main", "Test GUI Window")
	win := server.NewWindow(name, "Test GUI Window")
	win.Style().SetFullWidth()
	win.SetHAlign(server.HACenter)
	win.SetCellPadding(2)

	// Button which changes window content
	// win.Add(server.NewLabel("I'm a label! Try clicking on the button=>"))
	// btn := server.NewButton("Click me")
	// btn.AddEHandler(&myButtonHandler{text: ":-)"}, server.ETypeClick)
	// win.Add(btn)
	// btnsPanel := server.NewNaturalPanel()
	// btn.AddEHandlerFunc(func(e server.Event) {
	// 	// Create and add a new button...
	// 	newbtn := server.NewButton(fmt.Sprintf("Extra #%d", btnsPanel.CompsCount()))
	// 	newbtn.AddEHandlerFunc(func(e server.Event) {
	// 		btnsPanel.Remove(newbtn) // ...which removes itself when clicked
	// 		e.MarkDirty(btnsPanel)
	// 	}, server.ETypeClick)
	// 	btnsPanel.Insert(newbtn, 0)
	// 	e.MarkDirty(btnsPanel)
	// }, server.ETypeClick)
	// win.Add(btnsPanel)

	// ListBox examples
	p := server.NewHorizontalPanel()
	p.Style().SetBorder2(1, server.BrdStyleSolid, server.ClrBlack)
	p.SetCellPadding(2)
	p.Add(server.NewLabel("A drop-down list being"))
	widelb := server.NewListBox([]string{"50", "100", "150", "200", "250"})
	widelb.Style().SetWidth("50")
	widelb.AddEHandlerFunc(func(e server.Event) {
		widelb.Style().SetWidth(widelb.SelectedValue() + "px")
		e.MarkDirty(widelb)
	}, server.ETypeChange)
	p.Add(widelb)
	p.Add(server.NewLabel("pixel wide. And a multi-select list:"))
	listBox := server.NewListBox([]string{"First", "Second", "Third", "Forth", "Fifth", "Sixth"})
	listBox.SetMulti(true)
	listBox.SetRows(4)
	p.Add(listBox)
	countLabel := server.NewLabel("Selected count: 0")
	listBox.AddEHandlerFunc(func(e server.Event) {
		countLabel.SetText(fmt.Sprintf("Selected count: %d", len(listBox.SelectedIndices())))
		e.MarkDirty(countLabel)
	}, server.ETypeChange)
	p.Add(countLabel)
	win.Add(p)

	// Self-color changer check box
	greencb := server.NewCheckBox("I'm a check box. When checked, I'm green!")
	greencb.AddEHandlerFunc(func(e server.Event) {
		if greencb.State() {
			greencb.Style().SetBackground(server.ClrGreen)
		} else {
			greencb.Style().SetBackground("")
		}
		e.MarkDirty(greencb)
	}, server.ETypeClick)
	win.Add(greencb)

	// TextBox with echo
	p = server.NewHorizontalPanel()
	p.Add(server.NewLabel("Enter your name:"))
	tb := server.NewTextBox("")
	tb.AddSyncOnETypes(server.ETypeKeyUp)
	p.Add(tb)
	p.Add(server.NewLabel("You entered:"))
	nameLabel := server.NewLabel("")
	nameLabel.Style().SetColor(server.ClrRed)
	tb.AddEHandlerFunc(func(e server.Event) {
		nameLabel.SetText(tb.Text())
		e.MarkDirty(nameLabel)
	}, server.ETypeChange, server.ETypeKeyUp)
	p.Add(nameLabel)
	win.Add(p)

	return win
}
