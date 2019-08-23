// receive.go

package crdt

import (
	"context"

	"github.com/dgraph-io/badger"
	stan "github.com/nats-io/stan.go"
	"github.com/pkg/errors"
)

//
// creates a pipeline that listens for crdt objects on the
// nominated stream, merges with local versions and outputs
// an iterator of the resulting json objects
//
func runReciever(ctx context.Context, userid string, topicName string, sc stan.Conn, db *badger.DB, wb *badger.WriteBatch, iterator chan []byte) error { // <-chan []byte, // emits json objects

	// monitor all error channels
	var errcList []<-chan error
	var buildError error

	// connect to stream, receive messages (binary crdts)
	// seaparate stream into json objects
	msgOut, errc, err := streamCRDTSource(ctx, userid, topicName, sc)
	if err != nil {
		buildError = errors.Wrap(err, "Error: cannot create stream-crdt source component: ")
		return buildError
	}
	errcList = append(errcList, errc)

	// unpack binary into crdt object
	crdtOut, errc, err := crdtDecode(ctx, msgOut)
	if err != nil {
		buildError = errors.Wrap(err, "Error: cannot create crdt-decode component: ")
		return buildError
	}
	errcList = append(errcList, errc)

	// wrap in CRDTData with meta-data
	cdDataOut, errc, err := crdtWrap(ctx, crdtOut)
	if err != nil {
		buildError = errors.Wrap(err, "Error: cannot create crdt-wrap component: ")
		return buildError
	}
	errcList = append(errcList, errc)

	// locate local vesion & merge with incoming, cache result
	mergeOut, errc, err := crdtMerge(ctx, db, cdDataOut)
	if err != nil {
		buildError = errors.Wrap(err, "Error: cannot create crdt-merge component: ")
		return buildError
	}
	errcList = append(errcList, errc)

	// save the updated crdt
	saveOut, errc, err := saveCRDT(ctx, wb, mergeOut)
	if err != nil {
		buildError = errors.Wrap(err, "Error: cannot create save-crdt component: ")
		return buildError
	}
	errcList = append(errcList, errc)

	// turn crdt back into pure json object
	// iterator is a channel provided by the caller to
	// receive the json objects on
	errc, err = crdtJSON(ctx, iterator, saveOut)
	if err != nil {
		buildError = errors.Wrap(err, "Error: cannot create crdt-JSON component: ")
		return buildError
	}
	errcList = append(errcList, errc)

	err = WaitForPipeline(errcList...)
	if err != nil {
		return err
	}

	err = wb.Flush()

	return err

}
