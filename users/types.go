package users

import (
	"github.com/gramework/threadsafe/cache"
)

type (
	System struct {
		store *Store
		cache *cache.Instance
	}
)
