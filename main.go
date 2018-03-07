package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	flagConfig = flag.String("config", "demon.json", "Config file to use")
)

func main() {
	flag.Parse()
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)

	conf, err := readConfig(*flagConfig)
	if err != nil {
		log.Fatal(err)
	}

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
}
