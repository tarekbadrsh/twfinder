// Built-in static CSS themes of Server.

package server

// Built-in CSS themes.
const (
	ThemeDefault = "default" // Default CSS theme
	ThemeDebug   = "debug"   // Debug CSS theme, useful for developing/debugging purposes.
)

// resNameStaticCSS returns the CSS resource name
// for the specified CSS theme.
func resNameStaticCSS(theme string) string {
	// E.g. "server-default-0.8.0.css"
	return "server-" + theme + "-" + ServerVersion + ".css"
}

var staticCSS = make(map[string][]byte)

func init() {
	staticCSS[resNameStaticCSS(ThemeDefault)] = []byte("" +
		`
.guiimg-collapsed {background-image:url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAATUlEQVQ4y83RsQkAMAhEURNc+iZw7KQNgnjGRlv5D0SRMQPgADjVbr3AuzCz1QJYKAUyiAYiqAx4aHe/p9XAn6C/IQ1kb9TfMATYcM5cL5cg3qDaS5UAAAAASUVORK5CYII=)}
.guiimg-expanded {background-image:url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQCAYAAAAf8/9hAAAATElEQVQ4y2NgGGjACGNUVlb+J0Vje3s7IwMDAwMT1VxAiitgtlPfBcS4Atl22rgAnyvQbaedC7C5ApvtVHEBXlBZWfmfUKwwMQx5AADNQhjmAryM3wAAAABJRU5ErkJggg==)}

.guiimg-collapsed, .guiimg-expanded {background-position:0px 0px; background-repeat:no-repeat}

body {font-family:Arial}

.gui-Window {}

.gui-Panel {}

.gui-Table {}

.gui-Label {}

.gui-Link {}

.gui-Image {}

.gui-Button {}

.gui-CheckBox {}
.gui-CheckBox-Disabled {color:#888}

.gui-RadioButton {}
.gui-RadioButton-Disabled {color:#888}

.gui-ListBox {}

.gui-TextBox {}

.gui-PasswBox {}

.gui-HTML {}

.gui-SwitchButton {}
.gui-SwitchButton-On-Active {background:#00a000; color:#d0ffd0}
.gui-SwitchButton-Off-Active {background:#d03030; color:#ffd0d0}
.gui-SwitchButton-On-Inactive, .gui-SwitchButton-Off-Inactive {background:#606060; color:#909090}
.gui-SwitchButton-On-Inactive:enabled, .gui-SwitchButton-Off-Inactive:enabled {cursor:pointer}
.gui-SwitchButton-On-Active, .gui-SwitchButton-Off-Active, .gui-SwitchButton-On-Inactive, .gui-SwitchButton-Off-Inactive {margin:0px;border: 0px; width:100%}
.gui-SwitchButton-On-Active:disabled, .gui-SwitchButton-Off-Active:disabled, .gui-SwitchButton-On-Inactive:disabled, .gui-SwitchButton-Off-Inactive:disabled {color:black}

.gui-Expander {}
.gui-Expander-Header, .gui-Expander-Header-Expanded {cursor:pointer}
.gui-Expander-Header, .gui-Expander-Header-Expanded, .gui-Expander-Content {padding-left:19px}

.gui-TabBar {}
.gui-TabBar-Top {padding:0px 5px 0px 5px; border-bottom:5px solid #8080f8}
.gui-TabBar-Bottom {padding:0px 5px 0px 5px; border-top:5px solid #8080f8}
.gui-TabBar-Left {padding:5px 0px 5px 0px; border-right:5px solid #8080f8}
.gui-TabBar-Right {padding:5px 0px 5px 0px; border-left:5px solid #8080f8}
.gui-TabBar-NotSelected {padding-left:5px; padding-right:5px; border:1px solid white  ; background:#c0c0ff; cursor:default}
.gui-TabBar-Selected    {padding-left:5px; padding-right:5px; border:1px solid #8080f8; background:#8080f8; cursor:default}
.gui-TabPanel {}
.gui-TabPanel-Content {border:1px solid #8080f8; width:100%; height:100%}

.gui-SessMonitor {}
.gui-SessMonitor-Expired, .gui-SessMonitor-Error {color:red}
`)

	staticCSS[resNameStaticCSS(ThemeDebug)] = []byte(string(staticCSS[resNameStaticCSS(ThemeDefault)]) +
		`
.gui-Window td, .gui-Table td, .gui-Panel td, .gui-TabPanel td {border:1px solid black}
`)
}
