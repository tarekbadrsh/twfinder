// Defines the Datepicker component.

package server

import (
	"net/http"
	"time"
	helper "twfinder/helpers"
)

// Datepicker interface defines a component for date input purpose.
//
// Suggested event type to handle actions: ETypeChange
//
// By default the value of the Datepicker is synchronized with the server
// on ETypeChange event which is when the Datepicker loses focus
// or when the ENTER key is pressed.
// If you want a Datepicker to synchronize values during editing
// add the ETypeKeyUp event type
// to the events on which synchronization happens by calling:
// 		AddSyncOnETypes(ETypeKeyUp)
//
// Default style class: "gui-Datepicker"
type Datepicker interface {
	// Datepicker is a component.
	Comp

	// Date returns the datetime.
	Date() time.Time

	// SetDate sets the date.
	SetDate(date time.Time)

	// Datepicker can be enabled/disabled.
	HasEnabled

	// ReadOnly returns if the Datepicker is read-only.
	ReadOnly() bool

	// SetReadOnly sets if the Datepicker is read-only.
	SetReadOnly(readOnly bool)
}

// Datepicker implementation.
type datepickerImpl struct {
	compImpl       // Component implementation
	hasDateImpl    // Has date implementation
	hasEnabledImpl // Has enabled implementation
}

// NewDatepicker creates a new Datepicker.
func NewDatepicker(data time.Time) Datepicker {
	c := newDatepickerImpl(strEncURIThisV, data)
	c.Style().AddClass("gui-Datepicker")
	return &c
}

// newDatepickerImpl creates a new DatepickerImpl.
func newDatepickerImpl(valueProviderJs []byte, data time.Time) datepickerImpl {
	c := datepickerImpl{newCompImpl(valueProviderJs), newHasDateImpl(data), newHasEnabledImpl()}
	c.AddSyncOnETypes(ETypeChange)
	return c
}

func (d *datepickerImpl) Date() time.Time {
	dd := d.date
	return dd
}

func (d *datepickerImpl) SetText(date time.Time) {
	d.date = date
}

func (d *datepickerImpl) ReadOnly() bool {
	ro := d.Attr("readonly")
	return len(ro) > 0
}

func (d *datepickerImpl) SetReadOnly(readOnly bool) {
	readOnlyTxt := ""
	if readOnly {
		readOnlyTxt = "readonly"
	}
	d.SetAttr("readonly", readOnlyTxt)

}

func (d *datepickerImpl) preprocessEvent(event Event, r *http.Request) {
	// Empty string for datepicker is a valid value.
	// So we have to check whether it is supplied, not just whether its len() > 0
	value := r.FormValue(paramCompValue)
	if len(value) > 0 {
		d.date = helper.StringtoDate(value, "2006-01-02")
	} else {
		// Empty string might be a valid value, if the component value param is present:
		values, present := r.Form[paramCompValue] // Form is surely parsed (we called FormValue())
		if present && len(values) > 0 {
			d.date = helper.StringtoDate(values[0], "2006-01-02")
		}
	}
}

func (d *datepickerImpl) Render(w Writer) {
	d.renderInput(w)
}

var (
	strDate = []byte("date") // "date"
)

// renderInput renders the component as an input HTML tag.
func (d *datepickerImpl) renderInput(w Writer) {
	w.Write(strInputOp)
	w.Write(strDate)
	w.Write(strSize)
	w.Write(strQuote)
	d.renderAttrsAndStyle(w)
	d.renderEnabled(w)
	d.renderEHandlers(w)
	w.Write(strValue)
	d.renderText(w)
	w.Write(strInputCl)
}
