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
func Gui() server.Window {
	// Create and build a window
	win := server.NewWindow("main", "Test GUI Window")
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

	//
	// ---
	//
	win.Add(server.NewLabel("Configuration builder"))

	fp := server.NewHorizontalPanel()
	fp.Add(server.NewLabel("Consumer Key"))
	consumerKey := server.NewTextBox("")
	fp.Add(consumerKey)
	win.Add(fp)

	fp = server.NewHorizontalPanel()
	fp.Add(server.NewLabel("Consumer Secret"))
	consumerSecret := server.NewTextBox("")
	fp.Add(consumerSecret)
	win.Add(fp)

	fp = server.NewHorizontalPanel()
	fp.Add(server.NewLabel("Access Token"))
	accessToken := server.NewTextBox("")
	fp.Add(accessToken)
	win.Add(fp)

	fp = server.NewHorizontalPanel()
	fp.Add(server.NewLabel("Access Token Secret"))
	accessTokenSecret := server.NewTextBox("")
	fp.Add(accessTokenSecret)
	win.Add(fp)

	fp = server.NewHorizontalPanel()
	fp.Add(server.NewLabel("Search User"))
	searchUser := server.NewTextBox("")
	fp.Add(searchUser)
	win.Add(fp)

	win.Add(server.NewLabel("Search Criteria"))
	win.Add(server.NewLabel("Search Handle Context"))
	btnSearchAdd := server.NewButton("+")
	win.Add(btnSearchAdd)
	fp = server.NewHorizontalPanel()
	btnSearchAdd.AddEHandlerFunc(func(e server.Event) {
		// Create and add a new button...
		btnSearchRemove := server.NewButton("-")
		fp.Insert(btnSearchRemove, 0)
		searchHandleContext := server.NewTextBox("")
		fp.Insert(searchHandleContext, 0)
		e.MarkDirty(fp)

	}, server.ETypeClick)

	searchHandleContext := server.NewTextBox("")
	fp.Add(searchHandleContext)
	btnSearchRemove := server.NewButton("-")
	fp.Add(btnSearchRemove)

	win.Add(server.NewLabel("Search Criteria"))
	win.Add(server.NewLabel("Search Handle Context"))
	mybtn := server.NewButton("+")
	win.Add(mybtn)
	mybtnsPanel := server.NewNaturalPanel()
	mybtn.AddEHandlerFunc(func(e server.Event) {
		// Create and add a new button...
		searchHandleContext := server.NewTextBox("")
		mybtnsPanel.Insert(searchHandleContext, 0)
		newbtn := server.NewButton("-")
		newbtn.AddEHandlerFunc(func(e server.Event) {
			searchHandleContext.Parent().Remove(searchHandleContext)
			e.MarkDirty(searchHandleContext)
			mybtnsPanel.Remove(newbtn)
			e.MarkDirty(mybtnsPanel)
		}, server.ETypeClick)
		mybtnsPanel.Insert(newbtn, 0)
		e.MarkDirty(mybtnsPanel)
	}, server.ETypeClick)
	win.Add(mybtnsPanel)

	fp = server.NewHorizontalPanel()
	fp.Add(server.NewLabel("Search Name Context"))
	searchNameContext := server.NewTextBox("")
	fp.Add(searchNameContext)
	win.Add(fp)

	fp = server.NewHorizontalPanel()
	fp.Add(server.NewLabel("Search Bio Context"))
	searchBioContext := server.NewTextBox("")
	fp.Add(searchBioContext)
	win.Add(fp)

	fp = server.NewHorizontalPanel()
	fp.Add(server.NewLabel("Search Location Context"))
	searchLocationContext := server.NewTextBox("")
	fp.Add(searchLocationContext)
	win.Add(fp)

	fp = server.NewHorizontalPanel()
	fp.Add(server.NewLabel("Followers Count Between"))
	followersCountFrom := server.NewTextBox("")
	fp.Add(followersCountFrom)
	followersCountTo := server.NewTextBox("")
	fp.Add(followersCountTo)
	win.Add(fp)

	fp = server.NewHorizontalPanel()
	fp.Add(server.NewLabel("Following Count Between"))
	followingCountFrom := server.NewTextBox("")
	fp.Add(followingCountFrom)
	followingCountTo := server.NewTextBox("")
	fp.Add(followingCountTo)
	win.Add(fp)

	fp = server.NewHorizontalPanel()
	fp.Add(server.NewLabel("Likes Count Between"))
	likesCountFrom := server.NewTextBox("")
	fp.Add(likesCountFrom)
	likesCountTo := server.NewTextBox("")
	fp.Add(likesCountTo)
	win.Add(fp)

	fp = server.NewHorizontalPanel()
	fp.Add(server.NewLabel("Tweets Count Between"))
	tweetsCountFrom := server.NewTextBox("")
	fp.Add(tweetsCountFrom)
	tweetsCountTo := server.NewTextBox("")
	fp.Add(tweetsCountTo)
	win.Add(fp)

	fp = server.NewHorizontalPanel()
	fp.Add(server.NewLabel("Lists Count Between"))
	listsCountFrom := server.NewTextBox("")
	fp.Add(listsCountFrom)
	listsCountTo := server.NewTextBox("")
	fp.Add(listsCountTo)
	win.Add(fp)

	fp = server.NewHorizontalPanel()
	fp.Add(server.NewLabel("Joined Count Between"))
	JoinedCountFrom := server.NewTextBox("2000-09-22T12:42:31Z")
	fp.Add(JoinedCountFrom)
	JoinedCountTo := server.NewTextBox("2018-09-22T12:42:31Z")
	fp.Add(JoinedCountTo)
	win.Add(fp)

	fp = server.NewHorizontalPanel()
	verifiedCb := server.NewCheckBox("Verified")
	fp.Add(verifiedCb)
	win.Add(fp)

	fp = server.NewHorizontalPanel()
	followingCb := server.NewCheckBox("Following")
	fp.Add(followingCb)
	win.Add(fp)

	fp = server.NewHorizontalPanel()
	followersCb := server.NewCheckBox("Followers")
	fp.Add(followersCb)
	win.Add(fp)

	fp = server.NewHorizontalPanel()
	recursiveCb := server.NewCheckBox("Recursive")
	fp.Add(recursiveCb)
	win.Add(fp)

	fp = server.NewHorizontalPanel()
	recursiveSuccessUsersOnlyCb := server.NewCheckBox("Recursive Success Users Only")
	fp.Add(recursiveSuccessUsersOnlyCb)
	win.Add(fp)

	return win
}
