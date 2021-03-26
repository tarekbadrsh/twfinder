package frontend

import (
	"math/rand"
	"strconv"
	"time"
	"twfinder/config"
	"twfinder/gui/server"
	"twfinder/logger"
)

// newStrTxtLblPanel : create new TextBox with lable in Horizontal mode
// strInput : text input for TextBox
func newStrTxtLblPanel(lbltxt string, input *string, isPassword bool) server.Panel {
	pan := server.NewHorizontalPanel()
	lbl := server.NewLabel(lbltxt)
	pan.Add(lbl)

	txtbox := server.NewTextBox(*input)
	if isPassword {
		txtbox = server.NewPasswBox(*input)
	}
	txtbox.AddEHandlerFunc(func(e server.Event) {
		*input = txtbox.Text()
	}, server.ETypeChange)
	pan.Add(txtbox)

	return pan
}

// newIntTxtLblPanel : create new TextBox with lable in Horizontal mode
// intInput : int input for TextBox
func newIntTxtLblPanel(lbltxt string, intInput *int64) server.Panel {
	pan := server.NewHorizontalPanel()

	lbl := server.NewLabel(lbltxt)
	pan.Add(lbl)

	inputStr := strconv.FormatInt(*intInput, 10)
	txtbox := server.NewTextBox(inputStr)
	txtbox.AddEHandlerFunc(func(e server.Event) {
		eventText := txtbox.Text()
		i, err := strconv.ParseInt(eventText, 10, 64)
		if err != nil {
			logger.Errorf("error occurred during conver string to number input:%v  error:%v", eventText, err)
			return
		}
		*intInput = i
	}, server.ETypeChange)
	pan.Add(txtbox)

	return pan
}

// newDatepickerLblPanel : create new datepicker with lable in Horizontal mode
// dateInput : date input for datepicker
func newDatepickerLblPanel(lbltxt string, dateInput *time.Time) server.Panel {
	pan := server.NewHorizontalPanel()
	lbl := server.NewLabel(lbltxt)
	pan.Add(lbl)
	mydatepicker := server.NewDatepicker(*dateInput)
	mydatepicker.AddEHandlerFunc(func(e server.Event) {
		*dateInput = mydatepicker.Date()
	}, server.ETypeChange)
	pan.Add(mydatepicker)
	return pan
}

// newIntTextBoxFromTo : return new two integar only text box
// to handle from/to input
func newIntTextBoxFromTo(panalTitle string, from *int64, to *int64) server.Panel {
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

// newDatepickerFromTo : return new two Datepicker
// to handle from/to input
func newDatepickerFromTo(panalTitle string, from *time.Time, to *time.Time) server.Panel {
	// main
	mainPanal := server.NewVerticalPanel()

	// header
	headerPanel := server.NewHorizontalPanel()
	headerPanellbl := server.NewLabel(panalTitle)
	headerPanel.Add(headerPanellbl)

	// body
	bodyPanelfunc := func(from *time.Time, to *time.Time) server.Panel {
		pan := server.NewVerticalPanel()
		fromPan := newDatepickerLblPanel("From", from)
		pan.Add(fromPan)
		toPan := newDatepickerLblPanel("To", to)
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
		*from = time.Time{}
		*to = time.Time{}
		mainPanal.Remove(bodyPanel)
		headerPanel.Insert(headerPanelAddBtn, headerPanel.CompsCount())
		headerPanel.Remove(headerPanelRemoveBtn)
		e.MarkDirty(mainPanal, headerPanel)
	}, server.ETypeClick)

	mainPanal.Add(headerPanel)

	if !from.IsZero() || !to.IsZero() {
		headerPanel.Add(headerPanelRemoveBtn)
		mainPanal.Add(bodyPanel)
	} else {
		headerPanel.Add(headerPanelAddBtn)
	}

	return mainPanal
}

// newArrTextBoxPanal : return new text box with two way binding
func newArrTextBoxPanal(panalTitle string, inputArr []string) (server.Panel, map[int]string) {

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
	twitterConfig := config.Configuration()
	// Create and build the configuration window
	win := server.NewWindow("configuration", "Configuration - Twitter Finder App")
	win.Style().SetFullWidth()
	win.SetHAlign(server.HACenter)
	win.SetCellPadding(2)

	win.Add(server.NewLabel("Configuration builder"))

	consumerKeyPan := newStrTxtLblPanel("Consumer Key", &twitterConfig.ConsumerKey, true)
	win.Add(consumerKeyPan)

	consumerSecretPan := newStrTxtLblPanel("Consumer Secret", &twitterConfig.ConsumerSecret, true)
	win.Add(consumerSecretPan)

	accessTokenPan := newStrTxtLblPanel("Access Token", &twitterConfig.AccessToken, true)
	win.Add(accessTokenPan)

	accessTokenSecretPan := newStrTxtLblPanel("Access Token Secret", &twitterConfig.AccessTokenSecret, true)
	win.Add(accessTokenSecretPan)

	searchUserPan := newStrTxtLblPanel("Search User", &twitterConfig.SearchUser, false)
	win.Add(searchUserPan)
	//
	// ---
	//
	win.Add(server.NewLabel("Twitter list"))
	twSaveListCb := newCheckPanel("Save the result to list", &twitterConfig.TwitterList.SaveList)
	win.Add(twSaveListCb)
	twitterListNamePan := newStrTxtLblPanel("Name", &twitterConfig.TwitterList.Name, false)
	win.Add(twitterListNamePan)
	twitterListDescPan := newStrTxtLblPanel("Description", &twitterConfig.TwitterList.Description, false)
	win.Add(twitterListDescPan)
	twPublicListCb := newCheckPanel("public list", &twitterConfig.TwitterList.IsPublic)
	win.Add(twPublicListCb)
	//
	// ---
	//
	win.Add(server.NewLabel("Search Criteria"))
	// searchHandlePanal
	searchHandlePanal, handleMainMap := newArrTextBoxPanal("Search Handle Context", twitterConfig.SearchCriteria.SearchHandleContext)
	win.Add(searchHandlePanal)
	// searchNamePanal
	searchNamePanal, nameMainMap := newArrTextBoxPanal("Search Name Context", twitterConfig.SearchCriteria.SearchNameContext)
	win.Add(searchNamePanal)
	// searchBioPanal
	searchBioPanal, bioMainMap := newArrTextBoxPanal("Search Bio Context", twitterConfig.SearchCriteria.SearchBioContext)
	win.Add(searchBioPanal)
	// searchLocationPanal
	searchLocationPanal, locationMainMap := newArrTextBoxPanal("Search Location Context", twitterConfig.SearchCriteria.SearchLocationContext)
	win.Add(searchLocationPanal)
	// followersPanal
	followersPanal := newIntTextBoxFromTo("Followers Count Between", &twitterConfig.SearchCriteria.FollowersCountBetween.From, &twitterConfig.SearchCriteria.FollowersCountBetween.To)
	win.Add(followersPanal)
	// followingPanal
	followingPanal := newIntTextBoxFromTo("Following Count Between", &twitterConfig.SearchCriteria.FollowingCountBetween.From, &twitterConfig.SearchCriteria.FollowingCountBetween.To)
	win.Add(followingPanal)
	// likesPanal
	likesPanal := newIntTextBoxFromTo("Likes Count Between", &twitterConfig.SearchCriteria.LikesCountBetween.From, &twitterConfig.SearchCriteria.LikesCountBetween.To)
	win.Add(likesPanal)
	// tweetsPanal
	tweetsPanal := newIntTextBoxFromTo("Tweets Count Between", &twitterConfig.SearchCriteria.TweetsCountBetween.From, &twitterConfig.SearchCriteria.TweetsCountBetween.To)
	win.Add(tweetsPanal)
	// listsPanal
	listsPanal := newIntTextBoxFromTo("Lists Count Between", &twitterConfig.SearchCriteria.ListsCountBetween.From, &twitterConfig.SearchCriteria.ListsCountBetween.To)
	win.Add(listsPanal)

	// JoinDatePanal
	JoinDatePanal := newDatepickerFromTo("Joined Date Between", &twitterConfig.SearchCriteria.JoinedBetween.From, &twitterConfig.SearchCriteria.JoinedBetween.To)
	win.Add(JoinDatePanal)

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
	saveConfigBtn := server.NewButton("Save & Exit")
	saveConfigBtn.AddEHandlerFunc(func(e server.Event) {
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

		config.SetConfiguration(twitterConfig)
		err := config.SaveConfiguration("")
		if err != nil {
			logger.Error(err)
			return
		}
		e.ReloadWin("home")
		logger.Info("Configuration has been successfully updated")
	}, server.ETypeClick)

	win.Add(saveConfigBtn)
	return win
}
