package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	flagConfig = flag.String("config", "demo.json", "Config file to use")
	flagAddr   = flag.String("addr", ":8080", "Service address")
)

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	conf, err := readConfig(*flagConfig)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", conf.serveHome)
	for k, p := range conf.Processes {
		http.HandleFunc("/ws/"+k, p.ServeHTTP)
	}

	h := &http.Server{Addr: *flagAddr, Handler: http.DefaultServeMux}

	// Handle graceful shutdown
	/*
		stop := make(chan os.Signal, 2)
		signal.Notify(stop, os.Interrupt, os.Kill)
		go func() {
			select {
			case sig := <-stop:
				log.Print("Received signal ", sig.String())
				err := conf.Wait()
				if err != nil {
					log.Print(err)
				}
				d := time.Second * 5
				if sig == os.Kill {
					d = time.Second * 15
				}
				wait, cancel := context.WithTimeout(context.Background(), d)
				defer cancel()
				err = h.Shutdown(wait)
				if err != nil {
					log.Print(err)
				}
			}
		}()
	*/

	go func() {
		log.Print("Listening for requests on ", *flagAddr)
		err = h.ListenAndServe()
		log.Print(err)
	}()

	fmt.Println("Press return to stop the service.")
	fmt.Scanln()

	log.Print("Stopping server...")
	err = closeProcesses()
	if err != nil {
		log.Print(err)
	}
	log.Print("Processes stopped.")
	d := time.Second * 5
	wait, cancel := context.WithTimeout(context.Background(), d)
	defer cancel()
	err = h.Shutdown(wait)
	if err != nil {
		log.Print(err)
	}
	log.Print("Done.")
}
