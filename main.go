package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
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
	stop := make(chan os.Signal, 2)
	signal.Notify(stop, os.Interrupt, os.Kill)
	go func() {
		select {
		case sig := <-stop:
			log.Print("Received signal ", sig.String())
			d := time.Second * 5
			if sig == os.Kill {
				d = time.Second * 15
			}
			conf.Wait()
			wait, cancel := context.WithTimeout(context.Background(), d)
			defer cancel()
			err := h.Shutdown(wait)
			if err != nil {
				log.Print(err)
			}
		}
	}()

	log.Fatal(h.ListenAndServe())

	/*
		err = conf.Start()
		if err != nil {
			log.Fatal(err)
		}
		defer conf.Wait()

		for i := range conf.Steps {
			fmt.Println()
			fmt.Println(conf.Steps[i].Title)
			if conf.Steps[i].Desc != "" {
				fmt.Println(conf.Steps[i].Desc)
			}
			if conf.Steps[i].Interrupt {
				fmt.Println("INTERRUPT")
				err = conf.Processes[conf.Steps[i].ID].Interrupt()
				if err != nil {
					log.Println(err)
				}
			}
			if conf.Steps[i].Input != "" {
				fmt.Println(conf.Steps[i].Input)
				_, err = conf.Processes[conf.Steps[i].ID].Write([]byte(conf.Steps[i].Input + "\n"))
				if err != nil {
					log.Println(err)
				}
			}
		}
	*/
}
