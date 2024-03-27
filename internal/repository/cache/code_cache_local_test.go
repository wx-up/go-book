package cache

import (
	"fmt"
	"testing"

	lru "github.com/hashicorp/golang-lru/v2"
)

type D struct {
	Cnt int
}

func TestLocalCodeCache_Verify(t *testing.T) {
	cache, _ := lru.New[string, *D](2)
	cache.Add("key1", &D{Cnt: 1})
	item, _ := cache.Get("key1")
	item.Cnt += 10
	item, _ = cache.Get("key1")
	fmt.Println(item.Cnt)

	cache.Add("key2", &D{Cnt: 2})
	item, _ = cache.Get("key2")
	fmt.Println(item.Cnt)

	cache.Add("key3", &D{Cnt: 3})
	fmt.Println(cache.Get("key1"))
}
