// receivefromstream.go

package crdt

import (
	"context"
	"log"
)

//
// starts a stream listener/processor for the topic associated
// with this manager
//
func (crdtm *CRDTManager) StartReceiver() (<-chan []byte, error) {

	ctx, cancelFunc := context.WithCancel(context.Background())
	crdtm.ReceiverCancelFunc = cancelFunc
	iterator := make(chan []byte, 0)

	go func() {
		defer close(iterator)
		err := runReciever(ctx, crdtm.UserId, crdtm.TopicName, crdtm.sc, crdtm.rdb, crdtm.rwb, iterator)
		if err != nil {
			log.Println("ReceiverError: ", err)
			return
		}
	}()

	return iterator, nil

}

//
// shuts down the receiver gracefully
//
func (crdtm *CRDTManager) StopReceiver() {
	// flush write buffers
	crdtm.rwb.Flush()
}
