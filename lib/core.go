package noncense

import (
	"hash/crc32"
	"errors"
)

/**
 * Htring - hashed string structure
 */
type Htring struct {
	Value string
	HashCode uint32
}

func NewHString(value string) Htring {

	h := crc32.NewIEEE()
	h.Write([]byte(value))
	return Htring{Value: value, HashCode: h.Sum32()}
}

func (s *Htring) trim(mod uint32) uint32 {
	return s.HashCode % mod;
}

/**
 * Represents hash map node
 */
type mapNode struct {
	prev *mapNode
	next *mapNode
	value *Htring
}

/**
 * Represents linked list node
 */
type listNode struct {
	next *listNode
	mapItem *mapNode
}




// -------------------------




/**
 * Main nonces holder
 */
type NoncesHolder struct {
	count uint32
	served uint32
	sizeMap uint32
	sizeList uint32

	hashMap []*mapNode
	load []uint32

	first *listNode
	last *listNode
}

func NewNoncesHolder(sizeMap uint32, sizeList uint32) *NoncesHolder {
	nh := NoncesHolder{
		count: 0,
		sizeList: sizeList,
		sizeMap: sizeMap,
		hashMap: make([]*mapNode, sizeMap),
		load: make([]uint32, sizeMap),
	}

	return &nh
}

func (h *NoncesHolder) GetLoadFactor() uint32 {

	if h.count == 0 {
		return 0;
	} else if h.count == 1 {
		return 1;
	}

	// Calculating max load
	var max uint32
	var i uint32
	for i = 0; i < h.sizeMap; i++ {
		if h.load[i] > max {
			max = h.load[i]
		}
	}

	return max
}

func (h *NoncesHolder) Add(value Htring) error {

	if h.Has(value) {
		return errors.New("Value already presents")
	}

	// Building map node
	mn := mapNode{value: &value}

	// Placing map node into hashmap
	i := value.trim(h.sizeMap)
	h.load[i]++
	if (h.hashMap[i] == nil) {
		h.hashMap[i] = &mn
	} else {
		mn.next = h.hashMap[i]
		h.hashMap[i].prev = &mn
		h.hashMap[i] = &mn
	}

	// Building list node
	ln := listNode{mapItem: &mn}

	h.count++
	h.served++
	if h.last == nil {
		// Empty nonces holder
		h.first = &ln
		h.last = &ln
	} else {
		// Moving last to new one
		h.last.next = &ln
		h.last = &ln

		// Truncate check
		if h.count > h.sizeList {

			// Truncate requested
			h.count--

			// Removing from map
			t := h.first.mapItem
			j := t.value.trim(h.sizeMap)
			h.load[j]--
			if t.next != nil && t.prev != nil {
				// Entry inside list
				t.prev.next = t.next
				t.next.prev = t.prev
			} else if t.prev != nil {
				// Entry is last
				t.prev.next = nil
			} else if t.next != nil {
				// Entry is first
				t.next.prev = nil
				h.hashMap[j] = t.next
 			} else {
				// Entry is the only one
				h.hashMap[j] = nil
			}

			// Removing from list
			h.first = h.first.next
		}
	}

	return nil
}

func (h *NoncesHolder) Has(value Htring) bool {
	i := value.trim(h.sizeMap)

	var mn *mapNode;

	mn = h.hashMap[i];

	for mn != nil {
		if mn.value.HashCode == value.HashCode && mn.value.Value == value.Value {
			return true
		}

		mn = mn.next
	}

	return false
}





// -------------------------



/**
 * Async nonces added
 */
type NonceAdder struct {
	holder *NoncesHolder
	in chan chanNode
}

type chanNode struct {
	value Htring
	out chan bool
}

func (a *NonceAdder) accept() {
	for {
		req := <- a.in
		req.out <- a.holder.Add(req.value) == nil
	}
}

func NewNoncesAdder(capacity uint32) *NonceAdder {

	// Constructing
	a := NonceAdder{
		holder: NewNoncesHolder(capacity / 2, capacity),
	}

	// Making incoming queue
	a.in = make(chan chanNode)

	// Starting goroutine
	go a.accept()

	return &a;
}

func (a *NonceAdder) Add(value string) <- chan bool {
	out := make(chan bool)

	a.in <- chanNode{
		value: NewHString(value),
		out: out,
	}

	return out
}
