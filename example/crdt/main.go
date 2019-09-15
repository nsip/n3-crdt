// main.go

package main

import (
	"fmt"
	"log"
	"time"

	crdt "github.com/nsip/n3-crdt"
)

func main() {

	crdtm, err := crdt.NewCRDTManager("mattf01", "context1")
	if err != nil {
		log.Fatal(err)
	}

	// dataFile := "./sample_data/xapi/xapi.json"
	// dataFile := "./sample_data/sif/sif.json"

	//
	// load sample data
	//
	sampleDataPaths := []string{
		// "./sample_data/xapi/xapi.json",
		"./sample_data/naplan/sif.json",
		// "./sample_data/subjects/subjects.json",
		// "./sample_data/lessons/lessons.json",
		// "./sample_data/curriculum/overview.json",
		// "./sample_data/curriculum/content.json",
		// "./sample_data/sif/sif.json",
	}

	for _, path := range sampleDataPaths {
		err := crdtm.SendFromFile(path)
		if err != nil {
			log.Fatal("Send() Error: ", err)
		}
	}

	log.Println("starting receiver...")
	iterator, err := crdtm.StartReceiver()
	if err != nil {
		log.Fatal("StartReciver() Error: ", err)
	}

	// use the iterator, just audit to log
	count := 0
	for json := range iterator {
		count++
		fmt.Printf("\njson msg received:(%d)\t%s\n%s\n", count, time.Now().Format(time.RFC850), json)
		// _ = json
	}
	log.Println("...iterator closed (", count, ") messages")

	crdtm.Close()

	// time for graceful shutdown
	// time.Sleep(time.Second * 3)

}
