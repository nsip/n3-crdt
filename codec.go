// codec.go

package crdt

import (
	"bytes"
	"encoding/gob"

	"github.com/nsip/vvmap"
	"github.com/pkg/errors"
)

//
// binary encding for messages going to internal datastore.
//
func EncodeCRDT(vvm *vvmap.Map) ([]byte, error) {

	encBuf := new(bytes.Buffer)
	encoder := gob.NewEncoder(encBuf)
	err := encoder.Encode(vvm)
	if err != nil {
		return nil, errors.Wrap(err, "Encoder unable to binary encode crdt: ")
	}
	return encBuf.Bytes(), nil

}

//
// binary decoding for messages coming from datastore.
//
func DecodeCRDT(encoded []byte) (*vvmap.Map, error) {

	decBuf := bytes.NewBuffer(encoded)
	decoder := gob.NewDecoder(decBuf)
	crdtOut := vvmap.New("", DefaultResolver)
	err := decoder.Decode(&crdtOut)
	if err != nil {
		return nil, errors.Wrap(err, "Error decoding crdt: ")
	}
	return crdtOut, nil
}
