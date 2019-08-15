// main.go

package main

import (
	"fmt"
	"log"

	crdt "github.com/nsip/n3-crdt"
)

func main() {

	crdtm, err := crdt.NewManager("mattf4", "context1")
	if err != nil {
		log.Fatal(err)
	}
	defer crdtm.Close()

	// dataFile := "./sample_data/xapi/single_xapi.json"
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

	count := 0
	for json := range iterator {
		count++
		fmt.Printf("\njson msg recieved:(%d)\n%s\n", count, json)
		_ = json
		if count == 3200 { // for simple xapi input
			break
		}
	}
	crdtm.StopReceiver()
	log.Println("...reciver stopped.")

}
