package cacheman

import (
	"errors"
	"strings"
	"sync"

	"github.com/VictoriaMetrics/fastcache"
)

// CacheManager - manage fastcache. This extends the github.com/VictoriaMetrics/fastcache to have a Del that supports wildcards
type CacheManager struct {
	cache     *fastcache.Cache
	keys      []string
	MaxLength int
	mutex     *sync.Mutex
}

var (
	ErrKeyDoesNotExist  error = errors.New(`key does not exist`)
	ErrKeyNotSet        error = errors.New(`key not set`)
	ErrKeyPatternNotSet error = errors.New(`key pattern not set`)
)

// New - creates a new CacheManager
func New(max int) *CacheManager {
	if max == 0 {
		max = 25165824
	}
	return &CacheManager{
		MaxLength: max,
		cache:     fastcache.New(max),
		keys:      make([]string, 0),
	}
}

// Set the cache
func (cm *CacheManager) Set(key string, value []byte) error {
	if key == "" {
		return ErrKeyNotSet
	}

	if cm.keys == nil {
		cm.keys = make([]string, 0)
	}

	if cm.cache == nil {
		cm.cache = fastcache.New(cm.MaxLength)
	}

	cm.cache.SetBig([]byte(key), value)
	cm.keys = append(cm.keys, key)

	return nil
}

// Get - get the cache content
func (cm *CacheManager) Get(dst []byte, key string) []byte {
	if key == "" {
		return []byte{}
	}

	return cm.cache.GetBig(dst, []byte(key))
}

// GetWithErr - get the cache content with error
func (cm *CacheManager) GetWithErr(key string) ([]byte, error) {
	if key == "" {
		return []byte{}, ErrKeyNotSet
	}

	return cm.cache.GetBig(nil, []byte(key)), nil
}

// Del - delete an item in the cache
func (cm *CacheManager) Del(keyPattern string) error {

	if keyPattern == "" {
		return ErrKeyPatternNotSet
	}

	if hassufx := strings.HasSuffix(keyPattern, "*"); !hassufx {
		cm.cache.Del([]byte(keyPattern))
		return nil
	}

	// We create a mutex to block changes to the keys
	if cm.mutex == nil {
		cm.mutex = &sync.Mutex{}
	}

	// If the cache key has an asterisk at the end,
	// we will search through the keys stored
	func() {

		// remove the star character
		pfx := keyPattern[0 : len(keyPattern)-1]
		newkeys := make([]string, 0)
		cm.mutex.Lock()         // block changes to the keys while looping
		defer cm.mutex.Unlock() // unlock when the function returns

		for _, v := range cm.keys {
			if strings.HasPrefix(v, pfx) {
				cm.cache.Del([]byte(v))
				continue // skip adding to new keys
			}
			newkeys = append(newkeys, v)
		}

		cm.keys = make([]string, 0)           // reset size
		cm.keys = append(cm.keys, newkeys...) // add the new keys
	}()

	return nil
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
