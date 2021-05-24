// config.go

package crdt

import (
	"fmt"
	"testing"
)

func TestAddTempClassifierConfig(t *testing.T) {
	type args struct {
		dataModel     string
		n3id          string
		requiredPaths []string
		links         []string
		unique        []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test",
			args: args{
				dataModel:     "TestModel",
				n3id:          "MyId",
				requiredPaths: []string{"PathA", "PathB", "PathC"},
				links:         []string{"LinkA", "LinkB", "LinkC"},
				unique:        []string{"UniqueA", "UniqueB", "UniqueC"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			temp := AddTempClassifierConfig(tt.args.dataModel, tt.args.n3id, tt.args.requiredPaths, tt.args.links, tt.args.unique)
			fmt.Println(temp)
			fmt.Println("----------------------")
			fmt.Println(classifierConfigTextTemp)
		})
	}
	fmt.Println("----------------------")
	fmt.Println(GetCurClassifierConfig())
	fmt.Println("----------------------")
	ClearTempClassifierConfig()
	fmt.Println(GetCurClassifierConfig())
}
