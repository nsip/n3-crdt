//
// CRDT package is used to manage state of data objects
// being exchanged and stored over the N3 infrastrucure.
//
// When a client sends a group (array) of json objects, they are
// turned into crdts.
// On Send() the crdts merge with any known local state of the object
// On Recieve() objects can come from the orignal user or any other user
// of the system, so a merge with the known global state object occurs.
//
// Objects exit the crdt layer as the merged version of the object, as
// plain json.
//
// In N3 the Recieve() is typically from a streaming service such as
// NATS or kafka.
//
// If the user wants to just run locally config can set the
// send and receive to occur sequentially, otherwise send and receive
// are assumed to be independent operations.
//

package crdt

//
// placehoder so go get will downlaod package
// and location for godoc overview.
//
