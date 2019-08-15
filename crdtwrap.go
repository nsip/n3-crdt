// crdtwrap.go

package crdt

import (
	"context"

	"github.com/nsip/vvmap"
)

//
// wraps the raw crdt in the CRDTData package to add relevant meta-data
//
func crdtWrap(ctx context.Context, in <-chan *vvmap.Map) (
	<-chan CRDTData, // emits crdtData objects
	<-chan error, // emits errors encountered to the pipeline manager
	error) { // any error encountered when creating this component

	out := make(chan CRDTData)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		for crdt := range in {

			objectId := crdt.Keys()[0]

			cd := CRDTData{
				N3id:    objectId,
				RawData: crdt.Get(objectId).(map[string]interface{}),
				CRDT:    crdt,
				UserId:  string(crdt.ID()),
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
