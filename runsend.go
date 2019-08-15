// runsend.go

package crdt

import (
	"context"
	"fmt"
	"io"

	"github.com/dgraph-io/badger"
	stan "github.com/nats-io/stan.go"
	"github.com/pkg/errors"
)

//
// processes the stream of json objects from the given reader;
// converts each object into a crdt, compares with any exisitng
// local version of the same object and performs an update and sends
// the resultiing object on.
//
//
// Sends the resulting object to the
// streaming provider, nats in the case of N3.
//
func runSendWithReader(sdb *badger.DB, swb *badger.WriteBatch, userId string, topicName string, sc stan.Conn, r io.Reader, auditLevel string) error {

	// set up a context to manage send pipeline
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	// monitor all error channels
	var errcList []<-chan error

	//
	// build the pipleine by connecting all stages
	//

	// seaparate stream into json objects
	jsonOut, errc, err := jsonReaderSource(ctx, r)
	if err != nil {
		return errors.Wrap(err, "Error: cannot create json-reader source component: ")
	}
	errcList = append(errcList, errc)

	// classify the json data & add meta-data
	classOut, errc, err := objectClassifier(ctx, userId, topicName, jsonOut)
	if err != nil {
		return errors.Wrap(err, "Error: cannot create object-classifier component: ")
	}
	errcList = append(errcList, errc)

	// find existing crdt
	findOut, errc, err := updateCRDT(ctx, sdb, classOut)
	if err != nil {
		return errors.Wrap(err, "Error: cannot create find-crdt component: ")
	}
	errcList = append(errcList, errc)

	// publish crdt
	publishOut, errc, err := publishCRDT(ctx, sc, findOut)
	if err != nil {
		return errors.Wrap(err, "Error: cannot create publish-crdt component: ")
	}
	errcList = append(errcList, errc)

	// save the crdt
	saveOut, errc, err := saveCRDT(ctx, swb, publishOut)
	if err != nil {
		return errors.Wrap(err, "Error: cannot create save-crdt component: ")
	}
	errcList = append(errcList, errc)

	// minimal audit log
	for cd := range saveOut {
		fmt.Printf("sent & saved:\t%s, v: %v id: %s\n", cd.N3id, cd.CRDT.Version(), cd.CRDT.ID())
	}

	return WaitForPipeline(errcList...)

}
