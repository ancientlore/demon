# demon

👿

Or, demo-N. A utility to let you orchestrate steps of a demonstration.

`demon` is a simple web server that presents several windows. Each window is either a _site_ (presented in an iframe) or a _process_ attached via websockets. The server is meant to be run locally - it's not designed to secure access to the attached processes.

By default `demon` will look for a demo config called `demo.json` in the current directory. When the demo is done, simply press return and it will attempt a clean exit.

## Configuration

Each demo is configured with a simple JSON file containing _sites_, _processes_, and _steps_.

### Top-Level settings

You can give the demo a title using `title` and define which box holds the steps by setting `stepsPosition`.

### Sites

Sites are specified using a map of a site name and data about the site. `title` is the title for the box, `url` specifies the URL of the site to load in the `iframe`, and `position` specifies where to place the box.

> Note: `url` may optionally be a single environment variable preceded with a `$`, like `"url": "$MYURL"`.

### Processes

Processes are specified using a map of a process name and data about the process. `title` specifies the title for the box, and `position` specifies where to place the box.

`command` is an array containing a command to launch and its arguments. `stdin`, `stdout`, and `stdout` are piped using websockets.

`dir` specifies what directory to start the process in.

`exitInput` allows you send a "last command" before terminating the process, for instance `exit`.

> Note: Any element of the `command` array may be a single environment variable preceded with a `$`. For example, `"command": [ "kubectl", "exec", "-i", "$POD", "/bin/sh"]`.

### Steps

Steps represent input commands that are sent to processes that were started. Steps are listed in a JSON array.

`title` specifies the overall title to show on steps box, while `desc` lets you add (smaller) text with more details.

`input` is optional and specifies what text to send to the process when the RUN button is pressed for the current step. `id` specifies which process gets the command. This is nice because you demo steps can easily switch between different processes.

### Example

    {
        "title": "httpbin demo",
        "stepsPosition": 1,
        "sites": {
            "webnull": {
                "title": "httpbin",
                "url": "http://httpbin.org/",
                "position": 0
            }
        },
        "processes": {
            "bash": {
                "title": "bash",
                "command": ["bash"],
                "dir": "",
                "position": 2,
                "exitInput": "exit"
            }
        },
        "steps": [
            {
                "title": "Demo of httpbin",
                "desc": "This demo shows some things that you can do with httpbin.",
                "id": "",
                "input": ""
            },
            {
                "title": "Get your IP address",
                "desc": "Returns the origin IP address",
                "id": "bash",
                "input": "curl -s http://httpbin.org/ip"
            },
            {
                "title": "Get a UUID4",
                "desc": "Returns a new UUID.",
                "id": "bash",
                "input": "curl -s http://httpbin.org/uuid"
            },
            {
                "title": "Get some XML",
                "desc": "Returns some XML.",
                "id": "bash",
                "input": "curl -s http://httpbin.org/xml"
            },
            {
                "title": "Summary",
                "desc": "Try some commands of your own, based on the documentation above!",
                "id": "",
                "input": ""
            }
        ]
    }

## Screen Shot

![Screen Shot](media/demo_solved.png)
