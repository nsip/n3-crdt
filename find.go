// find.go

package crdt

import (
	"time"

	"github.com/dgraph-io/badger"
	"github.com/nsip/vvmap"
	"github.com/pkg/errors"
	"github.com/tidwall/buntdb"
)

//
// helpers for accessing crdts in long-term storage & caches
//

//
// look for this object in the local cache
//
func findCached(cd CRDTData, cache *buntdb.DB) (*vvmap.Map, bool) {

	crdt := vvmap.New(vvmap.ID(cd.UserId), DefaultResolver)

	var crdtstring string
	err := cache.View(func(tx *buntdb.Tx) error {
		var txErr error
		crdtstring, txErr = tx.Get(cd.N3id)
		if txErr != nil {
			return txErr
		}
		return nil
	})
	if err != nil {
		return crdt, false
	}

	decoded, err := DecodeCRDT([]byte(crdtstring))
	if err != nil {
		return crdt, false
	}

	return decoded, true

}

//
// look for this object in the long-term storage
//
func findStored(cd CRDTData, db *badger.DB) (*vvmap.Map, bool) {

	crdt := vvmap.New(vvmap.ID(cd.UserId), DefaultResolver)

	var valCopy []byte
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(cd.N3id))
		if err != nil {
			return err
		}
		valCopy, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return crdt, false
	}

	decoded, err := DecodeCRDT(valCopy)
	if err != nil {
		return crdt, false
	}

	return decoded, true

}

//
// add the crdt to the local cache, with TTL expiry
//
func updateCache(cd CRDTData, crdt *vvmap.Map, cache *buntdb.DB) error {

	encoded, err := EncodeCRDT(crdt)
	if err != nil {
		return errors.Wrap(err, "unable to encode crdt updateCache():")
	}
	err = cache.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(cd.N3id, string(encoded), &buntdb.SetOptions{Expires: true, TTL: time.Hour})
		return err
	})
	if err != nil {
		return errors.Wrap(err, "unable to updateCache():")
	}

	return nil
}
