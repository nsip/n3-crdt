// updatecrdt.go

package crdt

import (
	"context"

	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"
)

//
// store the crdt in long-term db
//
func saveCRDT(ctx context.Context, wb *badger.WriteBatch, in <-chan CRDTData) (
	<-chan CRDTData, // emits CRDTData objects with updated crdt if changed
	<-chan error, // emits errors encountered to the pipeline manager
	error) { // any error encountered when creating this component

	out := make(chan CRDTData)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		for cd := range in {

			if cd.Updated {
				encoded, err := EncodeCRDT(cd.CRDT)
				if err != nil {
					errc <- errors.Wrap(err, "cannot encode crdt saveCRDT():")
					return
				}
				err = wb.Set([]byte(cd.N3id), encoded)
				if err != nil {
					errc <- errors.Wrap(err, "cannot store crdt saveCRDT():")
					return
				}
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
