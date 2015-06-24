package noncense

type nativeNode struct {
	next  *nativeNode
	prev  *nativeNode
	value string
}

type eventNode struct {
	value string
	out   chan bool
}

type NoncesAdderNative struct {
	values map[string]*nativeNode
	first  *nativeNode
	last   *nativeNode
	events chan eventNode
	count  int
}

func NewNoncesAdderNative(capacity int) *NoncesAdderNative {
	a := NoncesAdderNative{}
	a.events = make(chan eventNode)
	a.values = make(map[string]*nativeNode)
	a.count = capacity

	return &a
}

func (a *NoncesAdderNative) accept() {
	for {
		req := <-a.events
		req.out <- a.AddSync(req.value)
	}
}

func (a *NoncesAdderNative) Has(value string) bool {
	return a.values[value] != nil
}

func (a *NoncesAdderNative) AddSync(value string) bool {
	if a.Has(value) {
		return false
	}

	nn := nativeNode{value: value}
	a.values[value] = &nn

	if a.first == nil {
		a.first = &nn
		a.last = &nn
	} else {
		a.last.next = &nn
		nn.prev = a.last
		a.last = &nn
	}

	if len(a.values) > a.count {
		a.pop()
	}

	return true
}

func (a *NoncesAdderNative) pop() *nativeNode {
	if a.first == nil {
		return nil
	} else if a.first == a.last {
		node := a.first
		delete(a.values, node.value)
		a.first = nil
		a.last = nil

		return node
	} else {
		node := a.first
		delete(a.values, node.value)
		a.first = a.first.next
		a.first.prev = nil

		return node
	}
}

func (a *NoncesAdderNative) Add(value string) <-chan bool {
	out := make(chan bool)
	a.events <- eventNode{value: value, out: out}

	return out
}
