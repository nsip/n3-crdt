// findcrdt.go

package crdt

import (
	"context"

	"github.com/dgraph-io/badger"
	"github.com/nsip/vvmap"
	"github.com/tidwall/buntdb"
)

//
// sets the new version of content in the crdt, crdt will be updated if
// found in cache or long-term storage, if not found new crdt is created.
//
func updateCRDT(ctx context.Context, db *badger.DB, in <-chan CRDTData) (
	<-chan CRDTData, // emits CRDTData objects with located crdt if exists
	<-chan error, // emits errors encountered to the pipeline manager
	error) { // any error encountered when creating this component

	out := make(chan CRDTData)
	errc := make(chan error, 1)

	// keep a temporary cache of crdts in case multiple revisions
	// to the same object are sent through
	cache, err := buntdb.Open(":memory:")
	if err != nil {
		return nil, nil, err
	}

	go func() {
		defer close(out)
		defer close(errc)
		defer cache.Close()

		for cd := range in {

			var crdt *vvmap.Map
			var found bool
			// find cached or stored crdt for this object
			if crdt, found = findCached(cd, cache); !found {
				crdt, _ = findStored(cd, db)
			}

			// update the data
			crdt.Set(cd.N3id, Flatten(cd.RawData))

			err := updateCache(cd, crdt, cache)
			if err != nil {
				errc <- err
				return
			}

			cd.CRDT = crdt
			cd.Updated = true // for now always, in future may check if data has been sent before

			select {
			case out <- cd: // pass the data package on to the next stage
			case <-ctx.Done(): // listen for pipeline shutdown
				return
			}

		}
	}()

	return out, errc, nil

}
