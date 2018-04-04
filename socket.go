package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 8192

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Time to wait before force close on connection.
	closeGracePeriod = 10 * time.Second
)

var upgrader = websocket.Upgrader{}

func pumpStdin(ws *websocket.Conn, w io.Writer) {
	defer func() {
		log.Print("Closing stdin websocket")
		ws.Close()
	}()
	ws.SetReadLimit(maxMessageSize)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		message = append(message, '\n')
		if _, err := w.Write(message); err != nil {
			break
		}
	}
}

func pumpStdout(ws *websocket.Conn, r io.Reader, done chan struct{}) {
	defer func() {
	}()
	s := bufio.NewScanner(r)
	for s.Scan() {
		ws.SetWriteDeadline(time.Now().Add(writeWait))
		if err := ws.WriteMessage(websocket.TextMessage, s.Bytes()); err != nil {
			ws.Close()
			break
		}
	}
	if s.Err() != nil {
		log.Println("scan:", s.Err())
	}
	close(done)

	log.Print("Closing stdout websocket")
	ws.SetWriteDeadline(time.Now().Add(writeWait))
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	time.Sleep(closeGracePeriod)
	ws.Close()
}

func ping(ws *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := ws.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(writeWait)); err != nil {
				log.Println("ping:", err)
			}
		case <-done:
			return
		}
	}
}

var (
	// Slice of processes to close on exit. In a real high-capacity site, you'd
	// want to manage this list more carefully. But for our purposes, it's not
	// very significant.
	procsToClose []*process

	// This mutex protects the above slice if multiple callers are trying to use it.
	procsToCloseMutex sync.Mutex
)

func closeProcesses() error {
	procsToCloseMutex.Lock()
	defer procsToCloseMutex.Unlock()

	var merr multiError
	for _, p := range procsToClose {
		err := p.Wait()
		if err != nil {
			merr = append(merr, err)
		}
	}
	if len(merr) > 0 {
		return merr
	}
	return nil
}

func (p *process) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade: ", err)
		return
	}

	// Clone and add to list of processes to close later.
	p = p.Clone()
	procsToCloseMutex.Lock()
	procsToClose = append(procsToClose, p)
	procsToCloseMutex.Unlock()

	defer ws.Close()

	rdr, wtr := io.Pipe()
	err = p.Start(wtr, wtr)
	if err != nil {
		log.Print("stdout: ", err)
		ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
		return
	}
	defer func() {
		err := p.Wait()
		if err != nil {
			log.Print("wait: ", err)
		}
	}()

	stdoutDone := make(chan struct{})
	go pumpStdout(ws, rdr, stdoutDone)
	go ping(ws, stdoutDone)

	pumpStdin(ws, p)

	<-stdoutDone
}
