package goutil

type FloatQueue struct {
	adder   func(acc, val float64) float64
	baseVal float64
	size    int
	start   int
	qlen    int
	items   []float64
	sum     float64
}

func NewFloatQueue(size int, baseValue float64, initItems []float64, adder func(float64, float64) float64) *FloatQueue {
	if adder == nil {
		adder = floatAdder
	}
	ret := &FloatQueue{
		adder:   adder,
		baseVal: baseValue,
		items:   make([]float64, size),
		size:    size,
	}
	ret.Clear()
	for _, v := range initItems {
		ret.Enqueue(v)
	}
	return ret
}

func floatAdder(base, val float64) float64 { return base + val }

func (k *FloatQueue) Size() int            { return k.size }
func (k *FloatQueue) Length() int          { return k.qlen }
func (k *FloatQueue) Sum() float64         { return k.sum }
func (k *FloatQueue) Item(idx int) float64 { return k.items[(k.start+idx)%k.size] }

func (k *FloatQueue) Items() []float64 {
	ret := make([]float64, k.qlen)
	for i := range ret {
		ret = append(ret, k.Item(i))
	}
	return ret
}

func (k *FloatQueue) Clear() {
	k.start = 0
	k.qlen = 0
	k.sum = k.baseVal
}

func (k *FloatQueue) Enqueue(val float64) {
	var idx int
	if k.qlen == k.size {
		k.sum = k.adder(k.sum, -k.items[k.start])
		idx = k.start
		k.start = (idx + 1) % k.size
	} else {
		idx = (k.start + k.qlen) % k.size
		k.qlen++
	}
	k.items[idx] = val
	k.sum = k.adder(k.sum, val)
}

type QueueReducer interface {
	Add(base, val interface{}) interface{}
	Sub(base, val interface{}) interface{}
}

type CircularQueue struct {
	reducer QueueReducer
	baseVal interface{}
	size    int
	start   int
	qlen    int
	items   []interface{}
	sum     interface{}
}

func NewCircularQueue(size int, baseValue interface{}, initItems []interface{}, reducer QueueReducer) *CircularQueue {
	ret := &CircularQueue{
		reducer: reducer,
		baseVal: baseValue,
		items:   make([]interface{}, size),
		size:    size,
	}
	ret.Clear()
	for _, v := range initItems {
		ret.Enqueue(v)
	}
	return ret
}

func (k *CircularQueue) Size() int                { return k.size }
func (k *CircularQueue) Length() int              { return k.qlen }
func (k *CircularQueue) Sum() interface{}         { return k.sum }
func (k *CircularQueue) Item(idx int) interface{} { return k.items[(k.start+idx)%k.size] }

func (k *CircularQueue) Clear() {
	k.start = 0
	k.qlen = 0
	k.sum = k.baseVal
}

func (k *CircularQueue) Enqueue(val interface{}) {
	var idx int
	if k.qlen == k.size {
		k.sum = k.reducer.Sub(k.sum, k.items[k.start])
		idx = k.start
		k.start = (idx + 1) % k.size
	} else {
		idx = (k.start + k.qlen) % k.size
		k.qlen++
	}
	k.items[idx] = val
	k.sum = k.reducer.Add(k.sum, val)
}
