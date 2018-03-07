package main

import (
	"html/template"
	"net/http"
)

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<title>demon</title>
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
    overflow: hidden;
    padding: 0;
    margin: 0;
    width: 100%;
    height: 100%;
    background: gray;
}
.log {
    background: white;
    margin: 0;
    padding: 0.5em 0.5em 0.5em 0.5em;
    top: 0.5em;
    left: 0.5em;
    right: 0.5em;
    bottom: 3em;
	overflow: auto;
	height: 200px;
}
.log pre {
  margin: 0;
  font-family: "Go Mono", "Consolas", monospace;
  font-size: 10pt;
}
.form {
    padding: 0 0.5em 0 0.5em;
    margin: 0;
    bottom: 1em;
    left: 0px;
    width: 100%;
    overflow: hidden;
}
</style>
</head>
<body>
{{range $k, $p := .Processes}}
<h1>{{$p.Title}}</h1>
<div id="{{$k}}_log" class="log"></div>
<form id="{{$k}}_form" class="form">
    <input type="submit" value="Send" />
    <input type="text" id="{{$k}}_msg" size="64"/>
</form>
{{end}}
</body>
</html>`

var tpl = template.Must(template.New("index").Parse(indexHTML))

func (c *config) serveHome(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, c)
}
