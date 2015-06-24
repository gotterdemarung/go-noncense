package noncense

import (
	"errors"
	"sync"
)

type Holder struct {
	sync.Mutex

	served      uint
	currentSize uint
	maxSize     uint
	last        *string
	current     *string
	data        map[string]*string
}

func NewHolder(size uint) (*Holder, error) {
	if size == 0 {
		return nil, errors.New("Holder size must be more than 0")
	}

	return &Holder{maxSize: size, data: make(map[string]*string, size)}, nil
}

func (holder *Holder) Add(value string) error {
	if holder.Has(value) {
		return errors.New("Value already presents")
	}

	copyValue := value

	holder.Lock()

	holder.currentSize++

	if holder.currentSize == 1 {
		holder.last = &copyValue
		holder.data[copyValue] = nil
	} else {
		holder.data[*holder.current] = &copyValue
	}

	holder.current = &copyValue

	if holder.currentSize > holder.maxSize {
		holder.pop()
	}

	holder.served++

	holder.Unlock()

	return nil
}

func (holder *Holder) AddAsync(value string) <-chan error {
	out := make(chan error)

	go func(outChan chan<- error) {
		outChan <- holder.Add(value)
	}(out)

	return out
}

func (holder *Holder) Has(value string) bool {
	_, ok := holder.data[value]

	return ok
}

func (holder *Holder) GetServedCount() uint {
	return holder.served
}

func (holder *Holder) pop() {

	previous := holder.data[*holder.last]

	delete(holder.data, *holder.last)

	holder.last = previous

	holder.currentSize--
}
