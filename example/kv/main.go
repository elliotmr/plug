package main

import (
	"fmt"
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
	err := example.RunKV(&kvImpl{m: make(map[string][]byte)})
	if err != nil {
		fmt.Println(err.Error(), ", exiting...")
	}
}
