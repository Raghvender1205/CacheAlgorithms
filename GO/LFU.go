/// I used a map and a minimum priority queue to implement it. Priority queue sorts the items in a binary tree based on their hit counts. 
// Item at the root of the heap has the least hit count. Instead of removing the root, we update it’s values and reset hit count and Fix the tree. 
// calling heap.Fix() is equivalent to, but less expensive than, calling heap.Remove() followed by a heap.Push()of the new value.

package main

import (
	"container/heap"
	"fmt"
)

// T is for key, value types
type T interface{}

type E struct {
	value T
	key   T
	hits  int
	index int
}

type PriorityQueue []*E

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].hits < pq[j].hits
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := pq.Len()
	item := x.(*E)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := old.Len()
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) update(e *E) {
	heap.Fix(pq, e.index)
}

func New(size int) *LFU {
	c := &LFU{
		size:  size,
		cache: make(map[T]*E),
	}
	return c
}

type LFU struct {
	pq    PriorityQueue
	cache map[T]*E
	size  int
}

func (c *LFU) Get(key T) (T, bool) {
	defer c.dump()
	e, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	e.hits++
	c.pq.update(e)
	return e.value, true
}

func (c *LFU) Set(key, value T) bool {
	defer c.dump()
	e, ok := c.cache[key]
	if ok {
		e.hits = 1
		e.key = key
		e.value = value
		c.pq.update(e)
		return false
	}

	if c.pq.Len() == c.size {
		e = c.pq[0]
		delete(c.cache, e.key)
		e.key = key
		e.value = value
		e.hits = 1
		c.cache[key] = e
		return true
	}

	e = &E{key: key, value: value, hits: 1}
	heap.Push(&c.pq, e)
	c.cache[key] = e
	return true

}

// dump dumps cache content for debugging
func (c *LFU) dump() {
	fmt.Printf("|")
	for i := 0; i < c.size; i++ {
		if i < c.pq.Len() {
			e := c.pq[i]
			fmt.Printf("  %v(%d)  |", e.value, e.hits)
			continue
		}
		fmt.Printf("        |")
	}
	fmt.Println()
}

func main() {
	c := New(2)
	c.Set("A", "A")
	c.Set("B", "B")
	c.Get("A")
	c.Set("C", "C")
	c.Get("B")
	c.Get("C")
	c.Set("D", "D")
	c.Get("A")
	c.Get("C")
	c.Get("D")
}
