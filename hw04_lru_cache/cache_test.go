package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		c.Clear()
		_, ok := c.Get("aaa")
		require.False(t, ok)
	})

	t.Run("rewrite the same elemenet in the cache", func(t *testing.T) {
		c := NewCache(5)

		for _, v := range [7]int{0, 1, 2, 3, 4, 5, 6} {
			c.Set("a", v)
		}

		val, ok := c.Get("a")
		require.True(t, ok)
		require.Equal(t, val, 6)
	})

	t.Run("check capacity of the cache", func(t *testing.T) {
		c := NewCache(3)

		c.Set("a", 0) // [a: 0]
		c.Set("b", 1) // [b: 1, a: 0]
		c.Set("c", 2) // [c: 2, b: 1, a: 0]
		c.Set("d", 3) // [d: 3, c: 2, b: 1]
		c.Set("c", 4) // [c: 4, d: 3, b: 1]

		_, ok := c.Get("a")
		require.False(t, ok)

		val, ok := c.Get("c")
		require.True(t, ok)
		require.Equal(t, val, 4)
		_, ok = c.Get("a")
		require.False(t, ok)

		c.Clear()
		_, ok = c.Get("c")
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
