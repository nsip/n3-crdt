// crdtdecode.go

package crdt

import (
	"context"

	"github.com/nsip/vvmap"
	"github.com/pkg/errors"
)

//
// takes binary blob received from stream and re-hydrates into
// a crdt object
//
func crdtDecode(ctx context.Context, in <-chan []byte) (
	<-chan *vvmap.Map, // emits crdt objects
	<-chan error, // emits errors encountered to the pipeline manager
	error) { // any error encountered when creating this component

	out := make(chan *vvmap.Map)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		for blob := range in {

			crdt, err := DecodeCRDT(blob)
			if err != nil {
				errc <- errors.Wrap(err, "unable to decode crdt decodeCRDT():")
				return
			}

			select {
			case out <- crdt: // pass the data package on to the next stage
			case <-ctx.Done(): // listen for pipeline shutdown
				return
			}

		}
	}()

	return out, errc, nil
}
