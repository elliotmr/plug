package main

import (
	"fmt"
	"github.com/elliotmr/plug"
	"github.com/elliotmr/plug/example"
)

type kvImpl struct {
	m map[string][]byte
}

func (k *kvImpl) Get(key string) ([]byte, error) {
	v, exists := k.m[key]
	if !exists {
		return nil, fmt.Errorf("key [%s] not in database", key)
	}
	return v, nil
}

func (k *kvImpl) Put(key string, value []byte) error {
	k.m[key] = value
	return nil
}

func main() {
	s := &example.KVPlugin{Impl: &kvImpl{m: make(map[string][]byte)}}
	err := plug.Run(s)
	if err != nil {
		fmt.Println(err.Error(), ", exiting...")
	}
}
