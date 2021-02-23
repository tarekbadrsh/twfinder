// Built-in static JavaScript codes of Server.

package server

import (
	"strconv"
)

// Static JavaScript resource name
const resNameStaticJs = "server-" + ServerVersion + ".js"

// Static javascript code
var staticJs []byte

func init() {
	// Init staticJs
	staticJs = []byte("" +
		// Param consts
		"var _pEventType='" + paramEventType +
		"',_pCompId='" + paramCompID +
		"',_pCompValue='" + paramCompValue +
		"',_pFocCompId='" + paramFocusedCompID +
		"',_pMouseWX='" + paramMouseWX +
		"',_pMouseWY='" + paramMouseWY +
		"',_pMouseX='" + paramMouseX +
		"',_pMouseY='" + paramMouseY +
		"',_pMouseBtn='" + paramMouseBtn +
		"',_pModKeys='" + paramModKeys +
		"',_pKeyCode='" + paramKeyCode +
		"';\n" +
		// Modifier key masks
		"var _modKeyAlt=" + strconv.Itoa(int(ModKeyAlt)) +
		",_modKeyCtlr=" + strconv.Itoa(int(ModKeyCtrl)) +
		",_modKeyMeta=" + strconv.Itoa(int(ModKeyMeta)) +
		",_modKeyShift=" + strconv.Itoa(int(ModKeyShift)) +
		";\n" +
		// Event response action consts
		"var _eraNoAction=" + strconv.Itoa(eraNoAction) +
		",_eraReloadWin=" + strconv.Itoa(eraReloadWin) +
		",_eraDirtyComps=" + strconv.Itoa(eraDirtyComps) +
		",_eraFocusComp=" + strconv.Itoa(eraFocusComp) +
		";" +
		`

function createXmlHttp() {
	if (window.XMLHttpRequest) // IE7+, Firefox, Chrome, Opera, Safari
		return new XMLHttpRequest();
	else // IE6, IE5
		return new ActiveXObject("Microsoft.XMLHTTP");
}

// Send event
function se(event, etype, compId, compValue) {
	var xhr = createXmlHttp();

	xhr.onreadystatechange = function() {
		if (xhr.readyState == 4 && xhr.status == 200)
			procEresp(xhr);
	}

	xhr.open("POST", _pathEvent, true); // asynch call
	xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");

	var data="";

	if (etype != null)
		data += "&" + _pEventType + "=" + etype;
	if (compId != null)
		data += "&" + _pCompId + "=" + compId;
	if (compValue != null)
		data += "&" + _pCompValue + "=" + compValue;
	if (document.activeElement.id != null && document.activeElement.id !== "")
		data += "&" + _pFocCompId + "=" + document.activeElement.id;

	if (event != null) {
		if (event.clientX != null) {
			// Mouse data
			var x = event.clientX, y = event.clientY;
			// Account for the amount body is scrolled:
			eventDoc = (event.target && event.target.ownerDocument) || document;
			doc = eventDoc.documentElement;
			body = eventDoc.body;
			x += (doc && doc.scrollLeft || body && body.scrollLeft || 0) -
				 (doc && doc.clientLeft || body && body.clientLeft || 0);
			y += (doc && doc.scrollTop  || body && body.scrollTop  || 0) -
				 (doc && doc.clientTop  || body && body.clientTop  || 0 );
			data += "&" + _pMouseWX + "=" + x;
			data += "&" + _pMouseWY + "=" + y;
			var parent = document.getElementById(compId);
			do {
				x -= parent.offsetLeft;
				y -= parent.offsetTop;
			} while (parent = parent.offsetParent);
			data += "&" + _pMouseX + "=" + x;
			data += "&" + _pMouseY + "=" + y;
			data += "&" + _pMouseBtn + "=" + (event.button < 4 ? event.button : 1); // IE8 and below uses 4 for middle btn
		}

		var modKeys;
		modKeys += event.altKey ? _modKeyAlt : 0;
		modKeys += event.ctlrKey ? _modKeyCtlr : 0;
		modKeys += event.metaKey ? _modKeyMeta : 0;
		modKeys += event.shiftKey ? _modKeyShift : 0;
		data += "&" + _pModKeys + "=" + modKeys;
		data += "&" + _pKeyCode + "=" + (event.which ? event.which : event.keyCode);
	}

	xhr.send(data);
}

function procEresp(xhr) {
	var actions = xhr.responseText.split(";");

	if (actions.length == 0) {
		window.alert("No response received!");
		return;
	}
	for (var i = 0; i < actions.length; i++) {
		var n = actions[i].split(",");

		switch (parseInt(n[0])) {
		case _eraDirtyComps:
			for (var j = 1; j < n.length; j++)
				rerenderComp(n[j]);
			break;
		case _eraFocusComp:
			if (n.length > 1)
				focusComp(parseInt(n[1]));
			break;
		case _eraNoAction:
			break;
		case _eraReloadWin:
			if (n.length > 1 && n[1].length > 0)
				window.location.href = _pathApp + n[1];
			else
				window.location.reload(true); // force reload
			break;
		default:
			window.alert("Unknown response code:" + n[0]);
			break;
		}
	}
}

function rerenderComp(compId) {
	var e = document.getElementById(compId);
	if (!e) // Component removed or not visible (e.g. on inactive tab of TabPanel)
		return;

	var xhr = createXmlHttp();

	xhr.onreadystatechange = function() {
		if (xhr.readyState == 4 && xhr.status == 200) {
			// Remember focused comp which might be replaced here:
			var focusedCompId = document.activeElement.id;
			e.outerHTML = xhr.responseText;
			focusComp(focusedCompId);

			// Inserted JS code is not executed automatically, do it manually:
			// Have to "re-get" element by compId!
			var scripts = document.getElementById(compId).getElementsByTagName("script");
			for (var i = 0; i < scripts.length; i++) {
				eval(scripts[i].innerText);
			}
		}
	}

	xhr.open("POST", _pathRenderComp, false); // synch call (if async, browser specific DOM rendering errors may arise)
	xhr.setRequestHeader("Content-type", "application/x-www-form-urlencoded");

	xhr.send(_pCompId + "=" + compId);
}

// Get selected indices (of an HTML select)
function selIdxs(select) {
	var selected = "";

	for (var i = 0; i < select.options.length; i++)
		if(select.options[i].selected)
			selected += i + ",";

	return selected;
}

// Get and update switch button value
function sbtnVal(event, onBtnId, offBtnId) {
	var onBtn = document.getElementById(onBtnId);
	var offBtn = document.getElementById(offBtnId);

	if (onBtn == null)
		return false;

	var value = onBtn == document.elementFromPoint(event.clientX, event.clientY);
	if (value) {
		onBtn.className = "gui-SwitchButton-On-Active";
		offBtn.className = "gui-SwitchButton-Off-Inactive";
	} else {
		onBtn.className = "gui-SwitchButton-On-Inactive";
		offBtn.className = "gui-SwitchButton-Off-Active";
	}

	return value;
}

function focusComp(compId) {
	if (compId != null && compId !== "") {
		var e = document.getElementById(compId);
		if (e) // Else component removed or not visible (e.g. on inactive tab of TabPanel)
			e.focus();
	}
}

function addonload(func) {
	var oldonload = window.onload;
	if (typeof window.onload != 'function') {
		window.onload = func;
	} else {
		window.onload = function() {
			if (oldonload)
				oldonload();
			func();
		}
	}
}

function addonbeforeunload(func) {
	var oldonbeforeunload = window.onbeforeunload;
	if (typeof window.onbeforeunload != 'function') {
		window.onbeforeunload = func;
	} else {
		window.onbeforeunload = function() {
			if (oldonbeforeunload)
				oldonbeforeunload();
			func();
		}
	}
}

var timers = new Object();

function setupTimer(compId, js, timeout, repeat, active, reset) {
	var timer = timers[compId];

	if (timer != null) {
		var changed = timer.js != js || timer.timeout != timeout || timer.repeat != repeat || timer.reset != reset;
		if (!active || changed) {
			if (timer.repeat)
				clearInterval(timer.id);
			else
				clearTimeout(timer.id);
			timers[compId] = null;
		}
		if (!changed)
			return;
	}
	if (!active)
		return;

	// Create new timer
	timers[compId] = timer = new Object();
	timer.js = js;
	timer.timeout = timeout;
	timer.repeat = repeat;
	timer.reset = reset;

	// Start the timer
	if (timer.repeat)
		timer.id = setInterval(js, timeout);
	else
		timer.id = setTimeout(js, timeout);
}

function checkSession(compId) {
	var e = document.getElementById(compId);
	if (!e) // Component removed or not visible (e.g. on inactive tab of TabPanel)
		return;

	var xhr = createXmlHttp();

	xhr.onreadystatechange = function() {
		if (xhr.readyState == 4 && xhr.status == 200) {
			var timeoutSec = parseFloat(xhr.responseText);
			if (timeoutSec < 60)
				e.classList.add("gui-SessMonitor-Expired");
			else
				e.classList.remove("gui-SessMonitor-Expired");
			var cnvtr = window[e.getAttribute("guiJsFuncName")];
			e.children[0].innerText = typeof cnvtr === 'function' ? cnvtr(timeoutSec) : convertSessTimeout(timeoutSec);
		}
	}

	xhr.open("GET", _pathSessCheck, false); // synch call (else we can't catch connection error)
	try {
		xhr.send();
		e.classList.remove("gui-SessMonitor-Error");
	} catch (err) {
		e.classList.add("gui-SessMonitor-Error");
		e.children[0].innerText = "CONN ERR";
	}
}

function convertSessTimeout(sec) {
	if (sec <= 0)
		return "Expired!";
	else if (sec < 60)
			return "<1 min";
	else
		return "~" + Math.round(sec / 60) + " min";
}

// INITIALIZATION

addonload(function() {
	focusComp(_focCompId);
});
`)
}
