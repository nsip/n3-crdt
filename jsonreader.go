// jsonreader.go

package crdt

import (
	"context"
	"encoding/json"
	"io"

	"github.com/pkg/errors"
)

//
// Iterator for json objects presented through a reader; file,
// http request, stdin etc.
//
// Reads json data from an array in a stream such as a file
// and feeds each parsed object into the ingest pipeline.
//
// ctx - required context for pipeline management
// r - reader accessing json data
//
func jsonReaderSource(ctx context.Context, r io.Reader) (
	<-chan map[string]interface{}, // source emits json objects read from file as map
	<-chan error, // emits any errors encountered to the pipeline
	error) { // any error when creating the source stage itself

	out := make(chan map[string]interface{})
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer close(errc)

		d := json.NewDecoder(r)

		// read opening brace "["
		_, err := d.Token()
		if err != nil {
			errc <- errors.Wrap(err, "unexpected token; json file should be json array")
			return
		}

		// read json objects one by one
		for d.More() {

			var m map[string]interface{}
			err := d.Decode(&m)
			if err != nil {
				errc <- errors.Wrap(err, "unable to decode json object.")
				return
			}

			select {
			case out <- m: // pass the map onto the next stage
			case <-ctx.Done(): // listen for pipeline shutdown
				return
			}

		}
	}()

	return out, errc, nil
}
