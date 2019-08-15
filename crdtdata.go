// crdtdata.go

package crdt

import "github.com/nsip/vvmap"

//
// data type passed through all stages of the
// send pipeline
//
type CRDTData struct {
	//
	// unique id determined for this object
	// will be taken from the object if it has a
	// declared unique id (in config)
	// If no unique id present in the object one will be
	// assigned.
	//
	N3id string
	//
	// The data model this object associated with through
	// classification
	//
	DataModel string
	//
	// Type of the object, derived from classifier
	//
	Type string
	//
	// map containing the original json
	//
	RawData map[string]interface{}
	//
	// the crdt to hold the data
	//
	CRDT *vvmap.Map
	//
	// encoded binary of the crdt
	//
	EncodedCRDT []byte
	//
	// user id to identify the owner of the
	// changes to the data
	//
	UserId string
	//
	// streaming topic to pubish to
	//
	TopicName string
	//
	// Flag whether any new data was added
	//
	Updated bool
}
