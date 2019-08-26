// objectclassifier.go

package crdt

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/nats-io/nuid"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

//
// Identifies & classifies the object passed in from the
// upstream reader.
//
// Uses the config in ./config/datatype.toml for deriving the
// data model, unique id etc.
//
//
func objectClassifier(ctx context.Context, userid string, topicName string, in <-chan map[string]interface{}) (
	<-chan CRDTData, // emits CRDTData objects with classification elements
	<-chan error, // emits errors encountered to the pipeline manager
	error) { // any error encountered when creating this component

	out := make(chan CRDTData)
	errc := make(chan error, 1)

	// load the classifier definitions;
	// each data-model type characterised by properties of the
	// json data.
	//
	type classifier struct {
		Data_model     string
		Required_paths []string
		N3id           string
		Links          []string
		Unique         []string
	}
	type classifiers struct {
		Classifier []classifier
	}
	var c classifiers
	classifierFile := fmt.Sprintf("./contexts/%s/%s/crdt/config/datatypes.toml", userid, topicName)
	if _, err := toml.DecodeFile(classifierFile, &c); err != nil {
		return nil, nil, err
	}

	go func() {
		defer close(out)
		defer close(errc)
		for jsonMap := range in { // read json object (map) from upstream source

			rawJson, err := json.Marshal(jsonMap) // we need json bytes for use with gjson
			if err != nil {
				errc <- errors.Wrap(err, "json marshal error")
				return
			}

			cd := CRDTData{}
			classified := false
			var dataModel, objectType, n3id string
			//
			// check the data by comparing with the known
			// classificaiton attributes from the config
			//
			for _, classifier := range c.Classifier {

				// now apply classification
				results := gjson.GetManyBytes(rawJson, classifier.Required_paths...)
				found := 0
				for _, r := range results {
					if r.Exists() {
						found++
					}
				}
				if len(classifier.Required_paths) == found {
					classified = true
				}
				if classified {
					// find the unique identifier for this object
					// if no id available use a nuid
					result := gjson.GetBytes(rawJson, classifier.N3id)
					if result.Exists() {
						n3id = result.String()
					} else {
						n3id = nuid.Next()
					}
					dataModel = classifier.Data_model
					break
				}
			}

			// default if model isn't classified
			if !classified {
				dataModel = "JSON"
			}

			// set the object type
			// if only 1 top level key, derive object type from it (SIF)
			// otherwise default to the datamodel as type (eg. xAPI)
			keys := []string{}
			for k := range jsonMap {
				keys = append(keys, k)
			}
			if len(keys) == 1 {
				objectType = keys[0]
			} else {
				objectType = dataModel
			}

			//
			// store all metadata for use in other stages
			//
			cd.DataModel = dataModel
			cd.Type = objectType
			cd.N3id = n3id
			cd.RawData = jsonMap
			cd.UserId = userid
			cd.TopicName = topicName

			select {
			case out <- cd: // pass the data package on to the next stage
			case <-ctx.Done(): // listen for pipeline shutdown
				return
			}

		}
	}()

	return out, errc, nil

}
