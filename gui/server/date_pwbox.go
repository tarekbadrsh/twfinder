// Defines the TextBox component.

package server

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

	// TextBox has text.
	HasText

	// TextBox can be enabled/disabled.
	HasEnabled

	// ReadOnly returns if the text box is read-only.
	ReadOnly() bool

	// SetReadOnly sets if the text box is read-only.
	SetReadOnly(readOnly bool)

	// Rows returns the number of displayed rows.
	Rows() int

	// SetRows sets the number of displayed rows.
	// rows=1 will make this a simple, one-line input text box,
	// rows>1 will make this a text area.
	SetRows(rows int)

	// Cols returns the number of displayed columns.
	Cols() int

	// SetCols sets the number of displayed columns.
	SetCols(cols int)

	// MaxLength returns the maximum number of characters
	// allowed in the text box.
	// -1 is returned if there is no maximum length set.
	MaxLength() int

	// SetMaxLength sets the maximum number of characters
	// allowed in the text box.
	// Pass -1 to not limit the maximum length.
	SetMaxLength(maxLength int)
}

/*
// PasswBox interface defines a text box for password input purpose.
//
// Suggested event type to handle actions: ETypeChange
//
// By default the value of the PasswBox is synchronized with the server
// on ETypeChange event which is when the PasswBox loses focus
// or when the ENTER key is pressed.
// If you want a PasswBox to synchronize values during editing
// (while you type in characters), add the ETypeKeyUp event type
// to the events on which synchronization happens by calling:
// 		AddSyncOnETypes(ETypeKeyUp)
//
// Default style class: "gui-PasswBox"
type PasswBox interface {
	// PasswBox is a TextBox.
	TextBox
}

// TextBox implementation.
type textBoxImpl struct {
	compImpl       // Component implementation
	hasTextImpl    // Has text implementation
	hasEnabledImpl // Has enabled implementation

	isPassw    bool // Tells if the text box is a password box
	rows, cols int  // Number of displayed rows and columns.
}

var (
	strEncURIThisV = []byte("encodeURIComponent(this.value)") // "encodeURIComponent(this.value)"
)

// NewTextBox creates a new TextBox.
func NewTextBox(text string) TextBox {
	c := newTextBoxImpl(strEncURIThisV, text, false)
	c.Style().AddClass("gui-Datepicker")
	return &c
}

// NewPasswBox creates a new PasswBox.
func NewPasswBox(text string) TextBox {
	c := newTextBoxImpl(strEncURIThisV, text, true)
	c.Style().AddClass("gui-PasswBox")
	return &c
}

// newTextBoxImpl creates a new textBoxImpl.
func newTextBoxImpl(valueProviderJs []byte, text string, isPassw bool) textBoxImpl {
	c := textBoxImpl{newCompImpl(valueProviderJs), newHasTextImpl(text), newHasEnabledImpl(), isPassw, 1, 20}
	c.AddSyncOnETypes(ETypeChange)
	return c
}

func (c *textBoxImpl) ReadOnly() bool {
	ro := c.Attr("readonly")
	return len(ro) > 0
}

func (c *textBoxImpl) SetReadOnly(readOnly bool) {
	if readOnly {
		c.SetAttr("readonly", "readonly")
	} else {
		c.SetAttr("readonly", "")
	}
}

func (c *textBoxImpl) Rows() int {
	return c.rows
}

func (c *textBoxImpl) SetRows(rows int) {
	c.rows = rows
}

func (c *textBoxImpl) Cols() int {
	return c.cols
}

func (c *textBoxImpl) SetCols(cols int) {
	c.cols = cols
}

func (c *textBoxImpl) MaxLength() int {
	if ml := c.Attr("maxlength"); len(ml) > 0 {
		if i, err := strconv.Atoi(ml); err == nil {
			return i
		}
	}
	return -1
}

func (c *textBoxImpl) SetMaxLength(maxLength int) {
	if maxLength < 0 {
		c.SetAttr("maxlength", "")
	} else {
		c.SetAttr("maxlength", strconv.Itoa(maxLength))
	}
}

func (c *textBoxImpl) preprocessEvent(event Event, r *http.Request) {
	// Empty string for text box is a valid value.
	// So we have to check whether it is supplied, not just whether its len() > 0
	value := r.FormValue(paramCompValue)
	if len(value) > 0 {
		c.text = value
	} else {
		// Empty string might be a valid value, if the component value param is present:
		values, present := r.Form[paramCompValue] // Form is surely parsed (we called FormValue())
		if present && len(values) > 0 {
			c.text = values[0]
		}
	}
}

func (c *textBoxImpl) Render(w Writer) {
	if c.rows <= 1 || c.isPassw {
		c.renderInput(w)
	} else {
		c.renderTextArea(w)
	}
}

var (
	strInputOp  = []byte(`<input type="`) // `<input type="`
	strPassword = []byte("password")      // "password"
	strText     = []byte("text")          // "text"
	strSize     = []byte(`" size="`)      // `" size="`
	strValue    = []byte(` value="`)      // ` value="`
	strInputCl  = []byte(`"/>`)           // `"/>`
)

// renderInput renders the component as an input HTML tag.
func (c *textBoxImpl) renderInput(w Writer) {
	w.Write(strInputOp)
	if c.isPassw {
		w.Write(strPassword)
	} else {
		w.Write(strText)
	}
	w.Write(strSize)
	w.Writev(c.cols)
	w.Write(strQuote)
	c.renderAttrsAndStyle(w)
	c.renderEnabled(w)
	c.renderEHandlers(w)

	w.Write(strValue)
	c.renderText(w)
	w.Write(strInputCl)
}

var (
	strTextareaOp   = []byte("<textarea")   // "<textarea"
	strRows         = []byte(` rows="`)     // ` rows="`
	strCols         = []byte(`" cols="`)    // `" cols="`
	strTextAreaOpCl = []byte("\">\n")       // "\">\n"
	strTextAreaCl   = []byte("</textarea>") // "</textarea>"
)

// renderTextArea renders the component as an textarea HTML tag.
func (c *textBoxImpl) renderTextArea(w Writer) {
	w.Write(strTextareaOp)
	c.renderAttrsAndStyle(w)
	c.renderEnabled(w)
	c.renderEHandlers(w)

	// New line char after the <textarea> tag is ignored.
	// So we must render a newline after textarea, else if text value
	// starts with a new line, it will be omitted!
	w.Write(strRows)
	w.Writev(c.rows)
	w.Write(strCols)
	w.Writev(c.cols)
	w.Write(strTextAreaOpCl)

	c.renderText(w)
	w.Write(strTextAreaCl)
}
*/
