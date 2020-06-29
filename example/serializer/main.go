package main

import (
	"encoding/json"
	"fmt"
	"github.com/elliotmr/plug/example"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

type serializer struct{}

func (s serializer) Marshal(data map[string]*anypb.Any) ([]byte, error) {
	var err error
	resp := make(map[string]interface{})
	for k, v := range data {
		// TODO: use generics to serialize the values without the "value" key
		resp[k], err = anypb.UnmarshalNew(v, proto.UnmarshalOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal any: %w", err)
		}
	}
	return json.Marshal(resp)
}

func main() {
	err := example.RunSerializer(&serializer{})
	if err != nil {
		fmt.Println(err.Error(), ", exiting...")
	}
}
