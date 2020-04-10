package cacheman

import (
	"strings"
	"sync"

	"github.com/VictoriaMetrics/fastcache"
)

// CacheManager - manage fastcache. This extends the github.com/VictoriaMetrics/fastcache to have a Del that supports wildcards
type CacheManager struct {
	cache     *fastcache.Cache
	keys      []string
	MaxLength int
}

// New - creates a new CacheManager
func New(max int) *CacheManager {
	if max == 0 {
		max = 25165824
	}
	return &CacheManager{MaxLength: max}
}

// Set the cache
func (cm *CacheManager) Set(key string, value []byte) {

	if cm.keys == nil {
		cm.keys = make([]string, 1)
		cm.keys[0] = key
	} else {
		cm.keys = append(cm.keys, key)
	}

	if cm.cache == nil {
		cm.cache = fastcache.New(cm.MaxLength)
	}

	cm.cache.SetBig([]byte(key), value)
}

// Get - get the cache content
func (cm *CacheManager) Get(dst []byte, key string) []byte {

	if key == "" {
		return []byte{}
	}

	return cm.cache.GetBig(dst, []byte(key))
}

// Del - delete an item in the cache
func (cm *CacheManager) Del(keyPattern string) {

	if keyPattern == "" {
		return
	}

	hassufx := strings.HasSuffix(keyPattern, "*")

	if !hassufx {
		cm.cache.Del([]byte(keyPattern))
		return
	}

	// If the cache key has an asterisk at the end,
	// we will search through the keys stored
	if hassufx {

		func() {
			// We create a mutex to block changes to the keys
			mutex := &sync.Mutex{}

			// remove the star character
			pfx := keyPattern[0 : len(keyPattern)-1]

			newkeys := make([]string, 0)

			mutex.Lock()         // block changes to the keys while looping
			defer mutex.Unlock() // unlock when the function returns
			for _, v := range cm.keys {
				if strings.HasPrefix(v, pfx) {
					cm.cache.Del([]byte(v))
					continue // skip adding to new keys
				}
				newkeys = append(newkeys, v)
			}

			cm.keys = make([]string, 0)           // reset size
			cm.keys = append(cm.keys, newkeys...) // add the new keys
			return
		}()

		return
	}

	return

}

// Has - check if a cache item is present in the cache
func (cm *CacheManager) Has(key string) bool {
	return cm.cache.Has([]byte(key))
}

// Reset - reset the cache content
func (cm *CacheManager) Reset() {
	cm.cache.Reset()
}

// ListKeys - gets the list of keys
func (cm *CacheManager) ListKeys() []string {
	return cm.keys
}
