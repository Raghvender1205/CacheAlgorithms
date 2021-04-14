// Combined lRU and LFU
/*
We can modify the LFU implementation a bit and add a global sequential id to each item. When accessing/adding an item we increment 
the id and store it in the item. So the items that are recently accessed will have a higher sequential ids than the older ones. 
Now we can compare the sequential ids when two items have the same hit counts.
*/

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
	tick  uint64
}

func (e1 *E) Less(e2 *E) bool {
	if e1.hits < e2.hits {
		return true
	}
	if e1.hits == e2.hits {
		return e1.tick < e2.tick
	}
	return false
}

type PriorityQueue []*E

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Less(pq[j])
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

func New(size int) *LRFU {
	c := &LRFU{
		size:  size,
		cache: make(map[T]*E),
	}
	return c
}

type LRFU struct {
	pq    PriorityQueue
	cache map[T]*E
	size  int
	tick  uint64
}

func (c *LRFU) Get(key T) (T, bool) {
	defer c.dump()
	e, ok := c.cache[key]
	if !ok {
		return nil, false
	}
	c.tick++
	e.hits++
	e.tick = c.tick
	c.pq.update(e)
	return e.value, true
}

func (c *LRFU) Set(key, value T) bool {
	defer c.dump()
	c.tick++
	e, ok := c.cache[key]
	if ok {
		e.hits++
		e.tick = c.tick
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
		e.tick = c.tick
		c.cache[key] = e
		return true
	}

	e = &E{key: key, value: value, hits: 1, tick: c.tick}
	heap.Push(&c.pq, e)
	c.cache[key] = e
	return true

}

// dump dumps cache content for debugging
func (c *LRFU) dump() {
	fmt.Printf("|")
	for i := 0; i < c.size; i++ {
		if i < c.pq.Len() {
			e := c.pq[i]
			fmt.Printf("  %v(%d,%d)  |", e.value, e.hits, e.tick)
			continue
		}
		fmt.Printf("          |")
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
