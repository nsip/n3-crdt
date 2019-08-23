// crdtsource.go

package crdt

import (
	"context"
	"log"

	stan "github.com/nats-io/stan.go"
	"github.com/pkg/errors"
)

//
// crdtsource connects to the streaming server for the given user
// and topic, and retrieves messages containing crdts of
// data updated across the network
//
func streamCRDTSource(ctx context.Context, userid string, topicName string, sc stan.Conn) (
	<-chan []byte, // emits read messages
	<-chan error, // emits errors encountered during receive
	error) { // any error encountered when creating this component

	out := make(chan []byte)
	errc := make(chan error, 1)
	msgchan := make(chan []byte)

	// var sub stan.Subscription

	go func() {
		defer close(out)
		defer close(errc)
		// establish the stream subscription
		sub, err := sc.Subscribe(topicName, func(msg *stan.Msg) {
			msgchan <- msg.Data
		}, stan.DurableName(topicName+userid), stan.DeliverAllAvailable())
		// note order of defers important here,
		// sub will panic if msgchan closed first.
		defer sub.Close()
		defer close(msgchan)
		if err != nil {
			errc <- errors.Wrap(err, "unable to connect to streaming server streamCRDTSource():")
			return
		}

		for msgdata := range msgchan {
			select {
			case out <- msgdata: // pass the data package on to the next stage
			case <-ctx.Done(): // listen for pipeline shutdown
				log.Println("...streamsource got cancel ctx.")
				sub.Close()
				return
			}

		}

	}()

	return out, errc, nil

}
