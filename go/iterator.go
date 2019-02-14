package lru

type Iterator struct {
	C    <-chan lruPair
	stop chan struct{}
}

func (i *Iterator) Stop() {
	defer func() {
		recover()
	}()

	close(i.stop)

	for range i.C {
	}
}

func newIterator() (*Iterator, chan<- lruPair, <-chan struct{}) {
	itemChan := make(chan lruPair)
	stopChan := make(chan struct{})
	return &Iterator{
		C:    itemChan,
		stop: stopChan,
	}, itemChan, stopChan
}
