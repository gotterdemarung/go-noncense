package noncense


// Async NONCEs adder
type NonceAdder struct {
	holder *NoncesHolder
	in chan chanNode
}

// Internal structure for channel communication
type chanNode struct {
	value HString
	out chan bool
}
// Internal method, listening for incoming requests from channel
func (a *NonceAdder) accept() {
	for {
		req := <- a.in
		req.out <- a.holder.Add(req.value) == nil
	}
}

// Constructor
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

// Adds new NONCE
// Async, returns boolean channel
func (a *NonceAdder) Add(value string) <- chan bool {
	out := make(chan bool)

	a.in <- chanNode{
		value: NewHString(value),
		out: out,
	}

	return out
}
