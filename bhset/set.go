package bhset

import (
	"github.com/emirpasic/gods/v2/sets/hashset"
)

func IsUnique[T comparable](keys ...T) bool {
	set := hashset.New[T]()
	for _, key := range keys {
		if set.Contains(key) {
			return false
		}

		set.Add(key)
	}

	return true
}
