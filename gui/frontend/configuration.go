package frontend

import (
	"fmt"
	"math/rand"
	"strconv"
	"twfinder/config"
	"twfinder/gui/server"
)

var twitterConfig = config.Config{}

// newStrTxtLblPanel : create new TextBox with lable in Horizontal mode
// strInput : text input for TextBox
func newStrTxtLblPanel(lbltxt string, input *string) server.Panel {
	pan := server.NewHorizontalPanel()

	lbl := server.NewLabel(lbltxt)
	pan.Add(lbl)

	txtbox := server.NewTextBox(*input)
	txtbox.AddEHandlerFunc(func(e server.Event) {
		*input = txtbox.Text()
	}, server.ETypeChange)
	pan.Add(txtbox)

	return pan
}

// newIntTxtLblPanel : create new TextBox with lable in Horizontal mode
// intInput : int input for TextBox
func newIntTxtLblPanel(lbltxt string, input *int64) server.Panel {
	pan := server.NewHorizontalPanel()

	lbl := server.NewLabel(lbltxt)
	pan.Add(lbl)

	inputStr := strconv.FormatInt(*input, 10)
	txtbox := server.NewTextBox(inputStr)
	txtbox.AddEHandlerFunc(func(e server.Event) {
		i, _ := strconv.ParseInt(txtbox.Text(), 10, 64)
		*input = i
	}, server.ETypeChange)
	pan.Add(txtbox)

	return pan
}

// NewTextBoxFromTo : return new two integar only text box
// to handle from/to input
func NewTextBoxFromTo(panalTitle string, from *int64, to *int64) server.Panel {
	// main
	mainPanal := server.NewVerticalPanel()

	// header
	headerPanel := server.NewHorizontalPanel()
	headerPanellbl := server.NewLabel(panalTitle)
	headerPanel.Add(headerPanellbl)

	// body
	bodyPanelfunc := func(from *int64, to *int64) server.Panel {
		pan := server.NewVerticalPanel()
		fromPan := newIntTxtLblPanel("From", from)
		pan.Add(fromPan)
		toPan := newIntTxtLblPanel("To", to)
		pan.Add(toPan)
		return pan
	}
	bodyPanel := bodyPanelfunc(from, to)

	// btns
	headerPanelAddBtn := server.NewButton("+")
	headerPanelRemoveBtn := server.NewButton("-")

	// addbtn
	headerPanelAddBtn.AddSyncOnETypes(server.ETypeClick)
	headerPanelAddBtn.AddEHandlerFunc(func(e server.Event) {
		bodyPanel = bodyPanelfunc(from, to)
		mainPanal.Add(bodyPanel)
		headerPanel.Insert(headerPanelRemoveBtn, headerPanel.CompsCount())
		headerPanel.Remove(headerPanelAddBtn)
		e.MarkDirty(mainPanal, headerPanel)
	}, server.ETypeClick)

	// removebtn
	headerPanelRemoveBtn.AddSyncOnETypes(server.ETypeClick)
	headerPanelRemoveBtn.AddEHandlerFunc(func(e server.Event) {
		// update values
		*from = 0
		*to = 0
		mainPanal.Remove(bodyPanel)
		headerPanel.Insert(headerPanelAddBtn, headerPanel.CompsCount())
		headerPanel.Remove(headerPanelRemoveBtn)
		e.MarkDirty(mainPanal, headerPanel)
	}, server.ETypeClick)

	mainPanal.Add(headerPanel)

	if *from > 0 || *to > 0 {
		headerPanel.Add(headerPanelRemoveBtn)
		mainPanal.Add(bodyPanel)
	} else {
		headerPanel.Add(headerPanelAddBtn)
	}

	return mainPanal
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

	bodyPanel := func(txt string, e server.Event) server.Panel {
		bodyPanel := server.NewHorizontalPanel()
		elementID := rand.Intn(1000000000000)
		mainMap[elementID] = txt
		elementTxtbox := server.NewTextBox(txt)
		elementTxtbox.AddEHandlerFunc(func(e server.Event) {
			mainMap[elementID] = elementTxtbox.Text()
		}, server.ETypeChange)
		rmvBtn := server.NewButton("-")
		rmvBtn.AddSyncOnETypes(server.ETypeClick)
		rmvBtn.AddEHandlerFunc(func(e server.Event) {
			mainPanal.Remove(bodyPanel)
			e.MarkDirty(mainPanal)
			delete(mainMap, elementID)
		}, server.ETypeClick)
		bodyPanel.Insert(rmvBtn, 0)
		bodyPanel.Insert(elementTxtbox, 0)
		if e != nil {
			e.MarkDirty(bodyPanel)
		}
		return bodyPanel
	}

	// handle existing values
	for _, v := range inputArr {
		bodyPanel := bodyPanel(v, nil)
		mainPanal.Insert(bodyPanel, mainPanal.CompsCount())
	}

	// handel add new panel
	headerPanelAddbtn.AddEHandlerFunc(func(e server.Event) {
		bodyPanel := bodyPanel("", e)
		mainPanal.Insert(bodyPanel, mainPanal.CompsCount())
		e.MarkDirty(mainPanal)
	}, server.ETypeClick)

	return mainPanal, mainMap
}

// newCheckPanel : create new checkbos with lable in Horizontal mode
// state : State of the checkbox
func newCheckPanel(lbltxt string, state *bool) server.Panel {
	pan := server.NewHorizontalPanel()
	lbl := server.NewLabel(lbltxt)
	pan.Add(lbl)
	checkBox := server.NewCheckBox("")
	checkBox.SetState(*state)
	checkBox.AddEHandlerFunc(func(e server.Event) {
		*state = checkBox.State()
	}, server.ETypeChange, server.ETypeKeyUp)
	pan.Add(checkBox)
	return pan
}

// ConfigWin : build configuration window with all required elements
func ConfigWin() server.Window {
	twitterConfig = config.Configuration()

	// Create and build a window
	win := server.NewWindow("Configuration", "Configuration")
	win.Style().SetFullWidth()
	win.SetHAlign(server.HACenter)
	win.SetCellPadding(2)

	win.Add(server.NewLabel("Configuration builder"))

	consumerKeyPan := newStrTxtLblPanel("Consumer Key", &twitterConfig.ConsumerKey)
	win.Add(consumerKeyPan)

	consumerSecretPan := newStrTxtLblPanel("Consumer Secret", &twitterConfig.ConsumerSecret)
	win.Add(consumerSecretPan)

	accessTokenPan := newStrTxtLblPanel("Access Token", &twitterConfig.AccessToken)
	win.Add(accessTokenPan)

	accessTokenSecretPan := newStrTxtLblPanel("Access Token Secret", &twitterConfig.AccessTokenSecret)
	win.Add(accessTokenSecretPan)

	searchUserPan := newStrTxtLblPanel("Search User", &twitterConfig.SearchUser)
	win.Add(searchUserPan)
	//
	// ---
	//
	win.Add(server.NewLabel("Search Criteria"))
	// searchHandlePanal
	searchHandlePanal, handleMainMap := NewArrTextBoxPanal("Search Handle Context", twitterConfig.SearchCriteria.SearchHandleContext)
	win.Add(searchHandlePanal)
	// searchNamePanal
	searchNamePanal, nameMainMap := NewArrTextBoxPanal("Search Name Context", twitterConfig.SearchCriteria.SearchNameContext)
	win.Add(searchNamePanal)
	// searchBioPanal
	searchBioPanal, bioMainMap := NewArrTextBoxPanal("Search Bio Context", twitterConfig.SearchCriteria.SearchBioContext)
	win.Add(searchBioPanal)
	// searchLocationPanal
	searchLocationPanal, locationMainMap := NewArrTextBoxPanal("Search Location Context", twitterConfig.SearchCriteria.SearchLocationContext)
	win.Add(searchLocationPanal)
	// followersPanal
	followersPanal := NewTextBoxFromTo("Followers Count Between", &twitterConfig.SearchCriteria.FollowersCountBetween.From, &twitterConfig.SearchCriteria.FollowersCountBetween.To)
	win.Add(followersPanal)
	// followingPanal
	followingPanal := NewTextBoxFromTo("Following Count Between", &twitterConfig.SearchCriteria.FollowingCountBetween.From, &twitterConfig.SearchCriteria.FollowingCountBetween.To)
	win.Add(followingPanal)
	// likesPanal
	likesPanal := NewTextBoxFromTo("Likes Count Between", &twitterConfig.SearchCriteria.LikesCountBetween.From, &twitterConfig.SearchCriteria.LikesCountBetween.To)
	win.Add(likesPanal)
	// tweetsPanal
	tweetsPanal := NewTextBoxFromTo("Tweets Count Between", &twitterConfig.SearchCriteria.TweetsCountBetween.From, &twitterConfig.SearchCriteria.TweetsCountBetween.To)
	win.Add(tweetsPanal)
	// listsPanal
	listsPanal := NewTextBoxFromTo("Lists Count Between", &twitterConfig.SearchCriteria.ListsCountBetween.From, &twitterConfig.SearchCriteria.ListsCountBetween.To)
	win.Add(listsPanal)

	// fp = server.NewHorizontalPanel()
	// fp.Add(server.NewLabel("Joined Count Between"))
	// JoinedCountFrom := server.NewTextBox("2000-09-22T12:42:31Z")
	// fp.Add(JoinedCountFrom)
	// JoinedCountTo := server.NewTextBox("2018-09-22T12:42:31Z")
	// fp.Add(JoinedCountTo)
	// win.Add(fp)

	// verifiedCb
	verifiedCb := newCheckPanel("Verified", &twitterConfig.SearchCriteria.Verified)
	win.Add(verifiedCb)
	//
	// ---
	//
	// followingCb
	followingCb := newCheckPanel("Following", &twitterConfig.Following)
	win.Add(followingCb)
	// followersCb
	followersCb := newCheckPanel("Followers", &twitterConfig.Followers)
	win.Add(followersCb)
	// recursiveCb
	recursiveCb := newCheckPanel("Recursive", &twitterConfig.Recursive)
	win.Add(recursiveCb)
	// recursiveSuccessUsersOnlyCb
	recursiveSuccessUsersOnlyCb := newCheckPanel("Recursive Success Users Only", &twitterConfig.RecursiveSuccessUsersOnly)
	win.Add(recursiveSuccessUsersOnlyCb)
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

		fmt.Println(twitterConfig)
	}, server.ETypeClick)
	win.Add(fullPrintbtn)

	return win
}
