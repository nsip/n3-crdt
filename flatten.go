// flatten.go

package crdt

import "fmt"

//
// Flatten takes a map of a json file and returns a new one where nested maps are replaced
// by dot-delimited keys.
//
// Also returns a sorted list of keys if updating in same order is important
// as when setting versions in the crdt.
//
func Flatten(m map[string]interface{}) map[string]interface{} {
	o := make(map[string]interface{})
	for k, v := range m {
		switch child := v.(type) {
		case map[string]interface{}: // nested map (object)
			nm := Flatten(child)
			for nk, nv := range nm {
				o[k+"."+nk] = nv
			}
		case []interface{}: // array
			for i, sv := range child {
				switch member := sv.(type) {
				case map[string]interface{}: // object inside array
					nm := Flatten(member)
					for nk, nv := range nm {
						key := fmt.Sprintf("%s.%d.%s", k, i, nk)
						o[key] = nv
					}
				default:
					key := fmt.Sprintf("%s.%d", k, i)
					o[key] = sv
				}
			}
		default:
			o[k] = v
		}
	}

	return o
}
