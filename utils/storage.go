package utils

import (
	"fmt"
	"math/rand"

	"github.com/lovoo/goka"
	"github.com/lovoo/goka/storage"
)

// RandomStoragePath is a helper option that uses a random name for storage to avoid clashes
func RandomStoragePath() goka.ProcessorOption {
	return goka.WithStorageBuilder(storage.DefaultBuilder(fmt.Sprintf("/tmp/goka-%x", rand.Int())))
}

// RandomStorageViewPath is a helper option that uses a random name for storage to avoid clashes
func RandomStorageViewPath() goka.ViewOption {
	return goka.WithViewStorageBuilder(storage.DefaultBuilder(fmt.Sprintf("/tmp/goka-%x", rand.Int())))
}
