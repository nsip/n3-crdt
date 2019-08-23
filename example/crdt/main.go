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

	dataFile := "./sample_data/xapi/xapi.json"
	// dataFile := "./sample_data/sif/sif.json"

	err = crdtm.SendFromFile(dataFile)
	if err != nil {
		log.Fatal("Send() Error: ", err)
	}

	log.Println("starting receiver...")
	iterator, err := crdtm.StartReceiver()
	if err != nil {
		log.Fatal("StartReciver() Error: ", err)
	}

	// use the iterator, just audit to log
	go func() {
		count := 0
		for json := range iterator {
			count++
			fmt.Printf("\njson msg recieved:(%d)\n%s\n", count, json)
		}
	}()

	time.Sleep(time.Second * 5)

	crdtm.Close()

}
