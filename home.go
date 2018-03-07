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
window.onload = function () {
	var f = function(k) {
		var conn;
		var msg = document.getElementById(k+"_msg");
		var log = document.getElementById(k+"_log");
		function appendLog(item) {
			var doScroll = log.scrollTop > log.scrollHeight - log.clientHeight - 1;
			log.appendChild(item);
			if (doScroll) {
				log.scrollTop = log.scrollHeight - log.clientHeight;
			}
		}
		document.getElementById(k+"_form").onsubmit = function () {
			if (!conn) {
				return false;
			}
			if (!msg.value) {
				return false;
			}
			var item = document.createElement("div");
			item.innerHTML = "<pre>$ <b>"+msg.value+"</b></pre>";
			appendLog(item);
			conn.send(msg.value);
			msg.value = "";
			return false;
		};
		if (window["WebSocket"]) {
			conn = new WebSocket("ws://" + document.location.host + "/ws/"+k);
			conn.onclose = function (evt) {
				var item = document.createElement("div");
				item.innerHTML = "<b>Connection closed.</b>";
				appendLog(item);
			};
			conn.onmessage = function (evt) {
				var messages = evt.data.split('\n');
				for (var i = 0; i < messages.length; i++) {
					var item = document.createElement("div");
					var s = messages[i];
					if (s.length == 0) {
						s = " ";
					}
					item.innerHTML = "<pre>"+s+"</pre>";
					appendLog(item);
				}
			};
		} else {
			var item = document.createElement("div");
			item.innerHTML = "<b>Your browser does not support WebSockets.</b>";
			appendLog(item);
		}
	};

	{{ range $k, $p := .Processes}}
	f("{{$k}}");
	{{end}}
};
</script>
<style type="text/css">
html {
    overflow: hidden;
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
	overflow: auto;
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
.log iframe {
	height: 100%;
	width: 100%;
	margin: 0;
	padding: 0;
	border: none;
}
.log {
	text-align: left;
    background: white;
    margin: 0;
	padding: 0.5em 0.5em 0.5em 0.5em;
	top: 3em;
	left: 0.5em;
	right: 0.5em;
	bottom: 3em;
	overflow: auto;
	position: absolute;  
}
.log pre {
  	margin: 0;
  	font-family: "Go Mono", "Consolas", monospace;
  	font-size: 10pt;
}
.form {
    padding: 0.5em 0.5em 0.5em 0.5em;
    margin: 0;
    overflow: hidden;
    bottom: 1em;
	left: 0px;
	position: absolute;
}
</style>
</head>
<body>
<h1>{{.Title}}</h1>
{{ range $k, $s := .Sites}}
<div class="group">
	<h1>{{$s.Title}}</h1>
	<div class="log">
		<iframe src="{{$s.URL}}"></iframe>
	</div>
</div>
{{end}}
{{range $k, $p := .Processes}}
<div class="group">
	<h1>{{$p.Title}}</h1>
	<div id="{{$k}}_log" class="log"></div>
	<form id="{{$k}}_form" class="form">
    	<input type="submit" value="Send" />
    	<input type="text" id="{{$k}}_msg" size="64" />
	</form>
</div>
{{end}}
</body>
</html>`

var tpl = template.Must(template.New("index").Parse(indexHTML))

func (c *config) serveHome(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, c)
}
