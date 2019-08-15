// crdtmanager.go

package crdt

import (
	"log"
	"os"
	"time"

	"github.com/dgraph-io/badger"
	stan "github.com/nats-io/stan.go"
	"github.com/pkg/errors"
)

type CRDTManager struct {
	//
	// the underlying badger k/v stores used by crdt manager
	//
	sdb *badger.DB // used when sending, client specific
	rdb *badger.DB // used when receiving from any/all clients
	//
	// manages parallel async writing to dbs
	//
	swb *badger.WriteBatch // on send
	rwb *badger.WriteBatch // on receive
	//
	// connection to the streaming server
	//
	sc stan.Conn
	//
	// set level of audit ouput, one of: none, basic, high
	//
	AuditLevel string
	//
	// user id to identify who is making changes
	//
	UserId string
	//
	// topic/context stream name used to exchange data
	//
	TopicName string
	//
	// conext cancelFunc used to close the
	// stream-receiver cleanly
	//
	ReceiverCancelFunc func()
}

//
// Open a crdt manager with supporting datastores
// will use the local path
// ./db/crdt/send & ./db/crdt/recv by default
//
func NewManager(userid string, topic string) (*CRDTManager, error) {

	defer timeTrack(time.Now(), "Open()")

	// if no filename provided will create locally
	crdtm, err := openFromFilePath("./db/crdt")
	if err != nil {
		return nil, errors.Wrap(err, "unable to open datastores:")
	}

	// assign user context
	crdtm.UserId = userid
	crdtm.TopicName = topic

	// create streaming server connection
	conn, err := stan.Connect(
		"test-cluster",
		userid,
		stan.NatsURL("nats://localhost:4222"),
	)
	if err != nil {
		return nil, errors.Wrap(err, "unable to find streaming server:")
	}
	crdtm.sc = conn

	return crdtm, nil

}

//
// safely shut down all databases & connections
//
func (crdtm *CRDTManager) Close() {

	defer timeTrack(time.Now(), "Close()")

	err := crdtm.swb.Flush()
	if err != nil {
		log.Println("error flushing send write-batch: ", err)
	}
	err = crdtm.rwb.Flush()
	if err != nil {
		log.Println("error flushing receive write-batch: ", err)
	}
	err = crdtm.sdb.Close()
	if err != nil {
		log.Println("error closing send datastore: ", err)
	}
	err = crdtm.rdb.Close()
	if err != nil {
		log.Println("error closing receive datastore: ", err)
	}
	// shut down the receiver if running
	if crdtm.ReceiverCancelFunc != nil {
		crdtm.StopReceiver()
	}
	// closure for streaming server connection
	crdtm.sc.Close()

}

//
// Opens the two datastores used on send & receive
// in the folder path specified.
//
// under that folder separate datastores will be created
// in subfodlders /send and recv
//
func openFromFilePath(folderPath string) (*CRDTManager, error) {

	sendFolder, recvFolder := "send", "recv"

	send := folderPath + "/" + sendFolder
	err := os.MkdirAll(send, os.ModePerm)
	if err != nil {
		return nil, err
	}

	recv := folderPath + "/" + recvFolder
	err = os.MkdirAll(recv, os.ModePerm)
	if err != nil {
		return nil, err
	}

	options := badger.DefaultOptions(send)
	// options = options.WithSyncWrites(false) // speed optimisation if required
	sdb, err := badger.Open(options)
	if err != nil {
		return nil, err
	}
	swb := sdb.NewWriteBatch()
	log.Println("...send datastore opened")

	options = badger.DefaultOptions(recv)
	// options = options.WithSyncWrites(false) // speed optimisation if required
	rdb, err := badger.Open(options)
	if err != nil {
		return nil, err
	}
	rwb := rdb.NewWriteBatch()
	log.Println("...receive datastore opened")

	return &CRDTManager{
		sdb:        sdb,
		rdb:        rdb,
		swb:        swb,
		rwb:        rwb,
		AuditLevel: "high",
	}, nil
}
