// crdtmerge.go

package crdt

import (
	"context"

	"github.com/dgraph-io/badger"
	"github.com/nsip/vvmap"
	"github.com/tidwall/buntdb"
)

//
// takes the crdt from the stream and merges with local version found
// either in cache or long-term storage.
//
func crdtMerge(ctx context.Context, db *badger.DB, in <-chan CRDTData) (
	<-chan CRDTData, // emits CRDTData objects with merged crdt
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
				crdt, found = findStored(cd, db)
			}

			if found {
				crdt.Merge(cd.CRDT.Delta(crdt.Version()))
			} else {
				crdt = cd.CRDT
			}

			err := updateCache(cd, crdt, cache)
			if err != nil {
				errc <- err
				return
			}
			cd.CRDT = crdt
			cd.Updated = true

			select {
			case out <- cd: // pass the data package on to the next stage
			case <-ctx.Done(): // listen for pipeline shutdown
				return
			}

		}
	}()

	return out, errc, nil

}
