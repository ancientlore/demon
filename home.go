package main

import (
	"html/template"
	"net/http"
)

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<title>{{.Title}}</title>
<script type="text/javascript">
var currentStep = 0;
var steps = [{{range $i, $v := .Steps}}{{if $i}},{{end}}
	{
		title: "{{$v.Title}}",
		desc: "{{$v.Desc}}",
		id: "{{$v.ID}}",
		input: "{{$v.Input}}"
	}{{end}}
];

var setStep = function() {
	document.getElementById("step").innerText = (currentStep+1).toString() + ". " + steps[currentStep].title;
	document.getElementById("desc").innerText = steps[currentStep].desc;
	var inp = steps[currentStep].input;
	if (inp != "") {
		inp = "$ " + inp;
	}
	document.getElementById("input").innerText = inp;
}

var prevStep = function() {
	if (currentStep > 0)
		currentStep--;
	setStep();
}
var nextStep = function() {
	if (currentStep < steps.length - 1)
		currentStep++;
	setStep();
}

var runStep = function() {
	var s = steps[currentStep];
	var f = submitters[s.id];
	f(s.input);
}

var submitters = {};
{{range $k, $p := .Processes}}submitters["{{$k}}"] = null;
{{end}}

var smaller = function(s) {
	var item = document.getElementById(s);
	item.classList.remove("larger");
	item.classList.add("smaller");
};

var normal = function(s) {
	var item = document.getElementById(s);
	item.classList.remove("smaller", "larger");
};

var larger = function(s) {
	var item = document.getElementById(s);
	item.classList.remove("smaller");
	item.classList.add("larger");
};

var smallerFrame = function(s) {
	var item = document.getElementById(s);
	item.classList.remove("larger-frame");
	item.classList.add("smaller-frame");
};

var normalFrame = function(s) {
	var item = document.getElementById(s);
	item.classList.remove("smaller-frame", "larger-frame");
};

var largerFrame = function(s) {
	var item = document.getElementById(s);
	item.classList.remove("smaller-frame");
	item.classList.add("larger-frame");
};

window.onload = function () {
	var f = function(k) {
		var conn;
		var msg = document.getElementById(k+"_msg");
		var log = document.getElementById(k+"_log");
		function appendLog(item) {
			var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
			log.appendChild(item);
			if (log.childElementCount > 120) {
				log.removeChild(log.childNodes[0]);
			}
			if (doScroll) {
				log.scrollTop = log.scrollHeight - log.clientHeight;
			}
		}
		submitters[k] = function (val) {
			if (!conn) {
				return false;
			}
			if (!val) {
				return false;
			}
			var item = document.createElement("div");
			item.innerHTML = "<pre>$ <b>"+val+"</b></pre>";
			appendLog(item);
			conn.send(val); // HL
			return false;
		};
		document.getElementById(k+"_form").onsubmit = function () {
			if (!msg.value) {
				return false;
			}
			var r = submitters[k](msg.value);
			msg.value = "";
			return r;
		};

		if (window["WebSocket"]) {
			conn = new WebSocket("ws://" + document.location.host + "/ws/"+k); // HL
			conn.onclose = function (evt) {
				var item = document.createElement("div");
				item.innerHTML = "<b>Connection closed.</b>";
				appendLog(item);
			};
			conn.onmessage = function (evt) { // HL
				var messages = evt.data.split('\n');
				for (var i = 0; i < messages.length; i++) {
					var item = document.createElement("div");
					var s = messages[i];
					if (s.length == 0) {
						s = " ";
					}
					var pre = document.createElement("pre");
					pre.innerText = s;
					item.appendChild(pre);
					appendLog(item);
				}
			};
		} else {
			var item = document.createElement("div");
			item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
			appendLog(item);
		}
	};

	{{ range $k, $p := .Processes}}	f("{{$k}}");
{{end}}

	prevStep();
};
</script>
<style type="text/css">
html {
}
body {
	background: #e2e1e0;
	text-align: center;
	font-family: "Go", "Arial", sans-serif;
}
h1 {
	color: #0074D9;
	font-size: 16pt;
}
.group {
	background: #fff;
	border-radius: 2px;
	display: inline-block;
	height: 300px;
	margin: 1rem;
	position: relative;
	width: 400px;
	resize: both;
	overflow: hidden;
	box-shadow: 0 1px 3px rgba(0,0,0,0.12), 0 1px 2px rgba(0,0,0,0.24);
	transition: all 0.3s cubic-bezier(.25,.8,.25,1);
	vertical-align: top;
}
.group:hover {
	box-shadow: 0 14px 28px rgba(0,0,0,0.25), 0 10px 10px rgba(0,0,0,0.22);
}
.group h1 {
	font-size: 12pt;
}
.group h2 {
	font-size: 10pt;
}
.group p {
	font-size: 8pt;
}
.fr iframe {
	height: 100%;
	width: 100%;
	margin: 0;
	padding: 0;
	border: none;
}
.fr {
	text-align: left;
    background: white;
    margin: 0;
	padding: 0.5em 0.5em 0.5em 0.5em;
	top: 2em;
	left: 0.5em;
	right: 0.5em;
	bottom: 0.5em;
	overflow: hidden;
	position: absolute;  
}
.frwrapper {
	padding: 0;
	margin: 0;
	height: 100%;
	width: 100%;
}
.log {
	text-align: left;
    background: white;
    margin: 0;
	padding: 0.5em 0.5em 0.5em 0.5em;
	top: 2em;
	left: 0.5em;
	right: 0.5em;
	bottom: 2em;
	overflow: auto;
	position: absolute;
}
.log pre {
  	margin: 0;
  	font-family: "Go Mono", "Consolas", monospace;
	font-size: 10pt;
}
.group pre {
	font-family: "Go Mono", "Consolas", monospace;
	font-size: 10pt;
}
.larger pre {
	font-size: 12pt;
}
.smaller pre {
	font-size: 8pt;
}
#input {
	white-space: pre-wrap;       /* css-3 */
	white-space: -moz-pre-wrap;  /* Mozilla, since 1999 */
	white-space: -pre-wrap;      /* Opera 4-6 */
	white-space: -o-pre-wrap;    /* Opera 7 */
	word-wrap: break-word;       /* Internet Explorer 5.5+ */
}
.form {
    padding: 0.5em 0.5em 0.5em 0.5em;
    margin: 0;
    overflow: hidden;
    bottom: 0.2em;
	left: 0px;
	position: absolute;
}
.controls {
	text-align: right;
    margin: 0;
	padding: 0.1em 0.1em 0.1em 0.1em;
	right: 0.2em;
	top: 0.2em;
	overflow: hidden;
	position: absolute;  
}
.smaller-frame {
	zoom: 0.75;
	-moz-transform: scale(0.75);
	-moz-transform-origin: 0 0;
	-o-transform: scale(0.75);
	-o-transform-origin: 0 0;
	-webkit-transform: scale(0.75);
	-webkit-transform-origin: 0 0;
	width: 134%;
	height: 134%;
}
.larger-frame {
	zoom: 1.25;
	-moz-transform: scale(1.25);
	-moz-transform-origin: 0 0;
	-o-transform: scale(1.25);
	-o-transform-origin: 0 0;
	-webkit-transform: scale(1.25);
	-webkit-transform-origin: 0 0;
	width: 80%;
	height: 80%;
}
</style>
</head>
<body>
<h1>{{.Title}}</h1>{{with $c := .}}{{range $i := .Loop }}
{{if eq $c.StepsPosition $i}}<div class="group">
	<h1 id="step"></h1>
	<p id="desc"></p>
	<pre id="input"></pre>
	<div>
		<input type="button" value="<" onclick="javascript:prevStep()"/>&nbsp;
		<input type="button" value="RUN" onclick="javascript:runStep()"/>&nbsp;
		<input type="button" value=">" onclick="javascript:nextStep()"/>
	</div>
</div>
{{end}}
{{ range $k, $s := $c.Sites}}{{if eq $i $s.Position}}
<div class="group">
	<h1>{{$s.Title}}</h1>
	<div class="fr">
		<div class="frwrapper" id="{{$k}}_fr">
			<iframe src="{{$s.URL}}"></iframe>
		</div>
	</div>
	<div class="controls">
		<input type="button" value="-" onclick="javascript:smallerFrame('{{$k}}_fr')"/>
		<input type="button" value="0" onclick="javascript:normalFrame('{{$k}}_fr')"/>
		<input type="button" value="+" onclick="javascript:largerFrame('{{$k}}_fr')"/>
	</div>
</div>
{{end}}{{end}}
{{range $k, $p := $c.Processes}}{{if eq $i $p.Position}}
<div class="group">
	<h1>{{$p.Title}}</h1>
	<div id="{{$k}}_log" class="log"></div>
	<form id="{{$k}}_form" class="form">
    	<input type="submit" value="Send" />&nbsp;
    	<input type="text" id="{{$k}}_msg" size="50" />
	</form>
	<div class="controls">
		<input type="button" value="-" onclick="javascript:smaller('{{$k}}_log')"/>
		<input type="button" value="0" onclick="javascript:normal('{{$k}}_log')"/>
		<input type="button" value="+" onclick="javascript:larger('{{$k}}_log')"/>
	</div>
</div>
{{end}}{{end}}
{{end}}{{end}}
</body>
</html>`

var tpl = template.Must(template.New("index").Parse(indexHTML))

func (c *config) serveHome(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, c)
}
