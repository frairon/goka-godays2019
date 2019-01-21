package utils

import (
	"fmt"
	"math/rand"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/storage"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomName() string {
	b := make([]byte, 5)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func RandomStoragePath() goka.ProcessorOption {
	return goka.WithStorageBuilder(storage.DefaultBuilder(fmt.Sprintf("/tmp/goka-%s", randomName())))
}
func RandomStorageViewPath() goka.ViewOption {
	return goka.WithViewStorageBuilder(storage.DefaultBuilder(fmt.Sprintf("/tmp/goka-%s", randomName())))
}
