// Here we replace the item that has been unused for the longest time.

package main

import (
	"container/list"
	"fmt"
)

// T is for key, value types
type T interface{}

// Element has key and value
type Element struct {
	key T
	val T
}
type Map map[T]*list.Element

type LRU struct {
	cache Map
	link  *list.List
	size  int
}

// New returns new LRU cache of the given size
func New(size int) *LRU {
	return &LRU{
		cache: make(Map),
		link:  list.New(),
		size:  size,
	}

}

// Set sets the given key,value pair in the cache.
// Returns true if new value is set in the cache.
func (l *LRU) Set(key, val T) bool {
	el := Element{key, val}
	defer l.dump()

	if e, ok := l.cache[key]; ok {
		e.Value = el
		l.link.MoveToFront(e)
		return false
	}
	if len(l.cache) < l.size {
		l.cache[key] = l.link.PushFront(el)
		return true

	}
	e := l.link.Back()
	delete(l.cache, e.Value.(Element).key)
	l.link.Remove(e)
	l.cache[key] = l.link.PushFront(el)
	return true
}

// Get gets a value from the cache
func (l *LRU) Get(key T) (T, bool) {
	defer l.dump()
	val, ok := l.cache[key]
	if !ok {
		return nil, ok
	}
	l.link.MoveToFront(val)
	el := val.Value.(Element)
	return el.val, true
}

// dump dumps cache content for debugging
func (l *LRU) dump() {
	e := l.link.Back()
	fmt.Printf("|")
	for i := 0; i < l.size; i++ {
		var val T
		val = " "
		if e != nil {
			val = e.Value.(Element).val
			e = e.Prev()
		}
		fmt.Printf("  %v  |", val)
	}
	fmt.Println()
}

func main() {
	c := New(4)
	c.Set("A", "A")
	c.Set("B", "B")
	c.Set("C", "C")
	c.Set("D", "D")
	c.Set("E", "E")
	c.Get("D")
	c.Set("F", "F")
}
