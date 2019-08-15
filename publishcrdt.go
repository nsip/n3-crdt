// publishcrdt.go

package crdt

import (
	"context"

	stan "github.com/nats-io/stan.go"
	"github.com/pkg/errors"
)

func publishCRDT(ctx context.Context, sc stan.Conn, in <-chan CRDTData) (
	<-chan CRDTData, // emits CRDTData objects
	<-chan error, // emits errors encountered to the pipeline manager
	error) { // any error encountered when creating this component

	out := make(chan CRDTData)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		for cd := range in {

			var encodeErr error
			cd.EncodedCRDT, encodeErr = EncodeCRDT(cd.CRDT)
			if encodeErr != nil {
				errc <- errors.Wrap(encodeErr, "error encoding crdt publishCRDT()")
				return
			}

			err := sc.Publish(cd.TopicName, cd.EncodedCRDT)
			if err != nil {
				errc <- errors.Wrap(err, "error publishing crdt publishCRDT()")
				return
			}

			select {
			case out <- cd: // pass the data package on to the next stage
			case <-ctx.Done(): // listen for pipeline shutdown
				return
			}

		}
	}()

	return out, errc, nil

}
