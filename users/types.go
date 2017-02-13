package users

import (
	"github.com/gramework/threadsafe/cache"
	"github.com/gramework/threadsafe/store"
)

type (
	System struct {
		store *store.Store
		cache *cache.Instance
	}
)
