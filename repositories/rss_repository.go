package repositories

import (
	"github.com/coocood/freecache"
)

type RssRepository struct {
	Cache         *freecache.Cache
	ExpireSeconds int
}

func (r RssRepository) GetBy(screenName string) *string {
	got, err := r.Cache.Get([]byte(screenName))
	if err != nil {
		return nil
	}
	result := string(got)
	return &result
}

func (r RssRepository) SetBy(screenName string, value string) error {
	key := []byte(screenName)
	val := []byte(value)
	return r.Cache.Set(key, val, r.ExpireSeconds)
}
