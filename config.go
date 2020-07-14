// config.go

package crdt

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
)

//
// check to see if config files exist
// if not creates useable defaults
//
func createDefaultConfig(filePath string) error {

	// check we can access/write configs
	configPath := fmt.Sprintf("%s/config", filePath)
	err := os.MkdirAll(configPath, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "cannot create config folder")
	}

	// now create classifier config
	classifierFile := fmt.Sprintf("%s/datatypes.toml", configPath)
	if fileExists(classifierFile) {
		return nil
	}

	log.Printf("%s not found, creating default classifier config...\n", classifierFile)
	err = writeDefaultClassifierConfig(classifierFile)
	if err != nil {
		return errors.Wrap(err, "cannot create default classifier config")
	}
	log.Println("...default classifier config created.")

	// any other configs here.

	return nil

}

//
// create the default data classifier config
//
func writeDefaultClassifierConfig(fname string) error {
	f, err := os.Create(fname)
	defer f.Close()
	_, err = f.WriteString(classifierConfigText)
	if err != nil {
		return err
	}
	f.Sync()
	return nil
}

//
// check whether file is already there
//
func fileExists(fname string) bool {

	_, err := os.Open(fname)
	if err != nil {
		return false
	}
	log.Println("found existing config file: ", fname)
	return true
}

var classifierConfigText = `
# 
# Classify section does a basic file/object identification
# based on fields in the object
# 
# required paths are those needed to identify an object as belonging
# to a particular data_model
# 
# n3id is the field to be used for unique identification of an
# object in the data store, if empty n3 will assign a unique id.
# 
# links are the features of the object to connect to the 
# overall data graph
# 
# unique fields are those used to construct a unique linking key
# for the object if no suitable single property is available
# 
# 
[[classifier]]
data_model = "SIF"
required_paths = ["*.RefId"]
n3id = "*.RefId"
links = ["RefId","LocalId"]

[[classifier]]
data_model = "XAPI"
required_paths = ["actor.name", "actor.mbox", "object.id", "verb.id"]
n3id = "id"
links = ["actor.mbox","actor.name","object.id","object.definition.name"]

[[classifier]]
data_model = "Syllabus"
required_paths = ["learning_area", "subject", "stage"]
n3id = "id"
links = ["learning_area", "subject", "stage"]
unique = ["subject","stage"]

[[classifier]]
data_model = "Subject"
required_paths = ["Subject.subject", "Subject.synonyms"]
n3id = "id"
links = ["Subject.learning_area", "Subject.subject", "Subject.stage", "Subject.synonyms"]
unique = ["Subject.subject","Subject.stage"]

[[classifier]]
data_model = "Lesson"
required_paths = ["Lesson.learning_area", "Lesson.lesson_id"]
n3id = "id"
links = ["Lesson.learning_area", "Lesson.subject", "Lesson.stage"]
unique = ["Lesson.subject","Lesson.stage"]

[[classifier]]
data_model = "LessonSequence"
required_paths = ["thearea", "thecourse", "thesubject", "thestage"]
n3id = "lessonId"
links = ["thearea", "thesubject", "thestage"]
unique = ["thesubject","thestage"]

[[classifier]]
data_model = "LessonSchedule"
required_paths = ["thecolor", "thecourse"]
n3id = "scheduleId"
links = ["userId"]

[[classifier]]
data_model = "OtfProviderItem"
required_paths = ["providerNodeId", "externalReference"]
n3id = "providerNodeId"
links = ["externalReference"]

[[classifier]]
data_model = "OtfNLPLink"
required_paths = ["nlpNodeId", "nlpReference"]
n3id = "nlpNodeId"
links = ["linkReference"]

[[classifier]]
data_model = "OtfScale"
required_paths = ["progressionLevel", "partiallyAchieved"]
n3id = "scaleItemId"
links = ["progressionLevel"]

[[classifier]]
data_model = "OtfMappedScore"
required_paths = ["scoreId", "mappedScore"]
n3id = "scoreId"
links = ["assessment"]
`
