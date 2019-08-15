// crdtjson.go

package crdt

import (
	"context"

	"github.com/pkg/errors"
	"github.com/tidwall/sjson"
)

//
// crdtJSON takes a feed of previously merged crdts, and turns them
// into the json data of the object, returned as a byte array.
//
func crdtJSON(ctx context.Context, iterator chan []byte, in <-chan CRDTData) (
	<-chan error, // emits errors encountered to the pipeline manager
	error) { // any error encountered when creating this component

	errc := make(chan error, 1)

	go func() {
		defer close(errc)

		for cd := range in {

			var jsonDoc []byte
			for _, id := range cd.CRDT.Keys() {
				tuples, ok := cd.CRDT.Get(id).(map[string]interface{})
				if !ok {
					errc <- errors.New("unexpected content in crdt crdtJSON():")
					return
				}
				for k, v := range tuples {
					var err error
					jsonDoc, err = sjson.SetBytes(jsonDoc, k, v)
					if err != nil {
						errc <- errors.Wrap(err, "unable to assign k/v pair crdtJSON():")
						return
					}
				}
			}

			select {
			case iterator <- jsonDoc: // pass the json data out to the consumer
			case <-ctx.Done(): // listen for pipeline shutdown
				return
			}

		}
	}()

	return errc, nil

}
