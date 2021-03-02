package frontend

import (
	"fmt"
	"math/rand"
	"twfinder/config"
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

var twitterConfig = config.Config{}

// NewTextBox : return new text box with two way binding
func NewTextBox(input *string) server.TextBox {
	// func NewTextBox(handlerChange func(input string)) server.TextBox {
	txtbox := server.NewTextBox(*input)
	txtbox.AddSyncOnETypes(server.ETypeKeyUp)
	txtbox.AddEHandlerFunc(func(e server.Event) {
		*input = txtbox.Text()
	}, server.ETypeChange, server.ETypeKeyUp)
	return txtbox
}

// NewArrTextBoxPanal : return new text box with two way binding
func NewArrTextBoxPanal(panalTitle string, inputArr []string) (server.Panel, map[int]string) {

	// Initialize main items
	mainMap := make(map[int]string)
	mainPanal := server.NewVerticalPanel()

	// header
	headerPanel := server.NewHorizontalPanel()
	headerPanellbl := server.NewLabel(panalTitle)
	headerPanel.Add(headerPanellbl)
	headerPanelAddbtn := server.NewButton("+")
	headerPanelAddbtn.AddSyncOnETypes(server.ETypeClick)
	headerPanel.Add(headerPanelAddbtn)
	mainPanal.Insert(headerPanel, mainPanal.CompsCount())

	// printBTN := server.NewButton("Print")
	// printBTN.AddEHandlerFunc(func(e server.Event) { fmt.Println(mainMap) }, server.ETypeClick)
	// mainPanal.Insert(printBTN, mainPanal.CompsCount())

	addInternalPanel := func(txt string, e server.Event) server.Panel {
		internalPanel := server.NewHorizontalPanel()
		elementID := rand.Intn(1000000000000)
		mainMap[elementID] = txt
		elementTxtbox := server.NewTextBox(txt)
		elementTxtbox.AddSyncOnETypes(server.ETypeKeyUp)
		elementTxtbox.AddEHandlerFunc(func(e server.Event) {
			mainMap[elementID] = elementTxtbox.Text()
		}, server.ETypeKeyUp)
		rmvBtn := server.NewButton("-")
		rmvBtn.AddSyncOnETypes(server.ETypeClick)
		rmvBtn.AddEHandlerFunc(func(e server.Event) {
			mainPanal.Remove(internalPanel)
			e.MarkDirty(mainPanal)
			delete(mainMap, elementID)
		}, server.ETypeClick)
		internalPanel.Insert(rmvBtn, 0)
		internalPanel.Insert(elementTxtbox, 0)
		if e != nil {
			e.MarkDirty(internalPanel)
		}
		return internalPanel
	}

	// handle existing values
	for _, v := range inputArr {
		internalPanel := addInternalPanel(v, nil)
		mainPanal.Insert(internalPanel, mainPanal.CompsCount())
	}

	// handel add new panel
	headerPanelAddbtn.AddEHandlerFunc(func(e server.Event) {
		internalPanel := addInternalPanel("", e)
		mainPanal.Insert(internalPanel, mainPanal.CompsCount())
		e.MarkDirty(mainPanal)
	}, server.ETypeClick)

	return mainPanal, mainMap
}

// Gui :
func Gui(name string) server.Window {
	twitterConfig = config.Configuration()

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

	//
	// ---
	//

	win.Add(server.NewLabel("Configuration builder"))

	fp := server.NewHorizontalPanel()
	fp.Add(server.NewLabel("Consumer Key"))
	consumerKey := NewTextBox(&twitterConfig.ConsumerKey)
	fp.Add(consumerKey)
	win.Add(fp)

	fp = server.NewHorizontalPanel()
	fp.Add(server.NewLabel("Consumer Secret"))
	consumerSecret := NewTextBox(&twitterConfig.ConsumerSecret)
	fp.Add(consumerSecret)
	win.Add(fp)

	fp = server.NewHorizontalPanel()
	fp.Add(server.NewLabel("Access Token"))
	accessToken := NewTextBox(&twitterConfig.AccessToken)
	fp.Add(accessToken)
	win.Add(fp)

	fp = server.NewHorizontalPanel()
	fp.Add(server.NewLabel("Access Token Secret"))
	accessTokenSecret := NewTextBox(&twitterConfig.AccessTokenSecret)
	fp.Add(accessTokenSecret)
	win.Add(fp)

	fp = server.NewHorizontalPanel()
	fp.Add(server.NewLabel("Search User"))
	searchUser := server.NewTextBox("")
	fp.Add(searchUser)
	win.Add(fp)

	//
	// ---
	//

	win.Add(server.NewLabel("Search Criteria"))

	//
	// ---
	//

	searchHandlePanal, handleMainMap := NewArrTextBoxPanal("Search Handle Context", twitterConfig.SearchCriteria.SearchHandleContext)
	win.Add(searchHandlePanal)

	//
	// ---
	//

	searchNamePanal, nameMainMap := NewArrTextBoxPanal("Search Name Context", twitterConfig.SearchCriteria.SearchNameContext)
	win.Add(searchNamePanal)

	//
	// ---
	//

	searchBioPanal, bioMainMap := NewArrTextBoxPanal("Search Bio Context", twitterConfig.SearchCriteria.SearchBioContext)
	win.Add(searchBioPanal)

	//
	// ---
	//

	searchLocationPanal, locationMainMap := NewArrTextBoxPanal("Search Location Context", twitterConfig.SearchCriteria.SearchLocationContext)
	win.Add(searchLocationPanal)

	//
	// ---
	//

	fullPrintbtn := server.NewButton("Full Print")
	fullPrintbtn.AddEHandlerFunc(func(e server.Event) {
		twitterConfig.SearchCriteria.SearchHandleContext = nil
		for _, v := range handleMainMap {
			twitterConfig.SearchCriteria.SearchHandleContext = append(twitterConfig.SearchCriteria.SearchHandleContext, v)
		}
		twitterConfig.SearchCriteria.SearchNameContext = nil
		for _, v := range nameMainMap {
			twitterConfig.SearchCriteria.SearchNameContext = append(twitterConfig.SearchCriteria.SearchNameContext, v)
		}
		twitterConfig.SearchCriteria.SearchBioContext = nil
		for _, v := range bioMainMap {
			twitterConfig.SearchCriteria.SearchBioContext = append(twitterConfig.SearchCriteria.SearchBioContext, v)
		}
		twitterConfig.SearchCriteria.SearchLocationContext = nil
		for _, v := range locationMainMap {
			twitterConfig.SearchCriteria.SearchLocationContext = append(twitterConfig.SearchCriteria.SearchLocationContext, v)
		}
		fmt.Println(twitterConfig.SearchCriteria)
	}, server.ETypeClick)
	win.Add(fullPrintbtn)

	//
	// ---
	//

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
