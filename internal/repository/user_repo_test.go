package repository

import (
	"testing"
)

func TestCacheUserRepository_FindById(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{
			name: "缓存未命中，但是查询命中",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
		})
	}
}
