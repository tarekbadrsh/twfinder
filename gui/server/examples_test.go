package server_test

import (
	"twfinder/gui/server"
)

// Example code determining which button was clicked.
func ExampleButton() {
	b := server.NewButton("Click me")
	b.AddEHandlerFunc(func(e server.Event) {
		if e.MouseBtn() == server.MouseBtnMiddle {
			// Middle click
		}
	}, server.ETypeClick)
}

// Example code determining what kind of key is involved.
func ExampleTextBox() {
	b := server.NewTextBox("")
	b.AddSyncOnETypes(server.ETypeKeyUp) // This is here so we will see up-to-date value in the event handler
	b.AddEHandlerFunc(func(e server.Event) {
		if e.ModKey(server.ModKeyShift) {
			// SHIFT is pressed
		}

		c := e.KeyCode()
		switch {
		case c == server.KeyEnter: // Enter
		case c >= server.Key0 && c <= server.Key9:
			fallthrough
		case c >= server.KeyNumpad0 && c <= server.KeyNumpad9: // Number
		case c >= server.KeyA && c <= server.KeyZ: // Letter
		case c >= server.KeyF1 && c <= server.KeyF12: // Function key
		}
	}, server.ETypeKeyUp)
}
