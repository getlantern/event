package event

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	numListeners = 100
	numMessages  = 1000
)

func TestBlocking(t *testing.T) {
	actual := int64(0)
	expected := runTest(true, func(val int64) {
		atomic.AddInt64(&actual, val)
	})
	assert.EqualValues(t, expected, atomic.LoadInt64(&actual))
}

func TestNonBlocking(t *testing.T) {
	actual := int64(0)
	expected := runTest(false, func(val int64) {
		atomic.AddInt64(&actual, val)
		time.Sleep(1 * time.Millisecond)
	})
	assert.True(t, expected > atomic.LoadInt64(&actual))
	assert.NotZero(t, atomic.LoadInt64(&actual))
}

func runTest(blocking bool, onMsg func(val int64)) int64 {
	expected := int64(0)

	d := NewDispatcher(blocking, 0)
	for i := 0; i < numListeners; i++ {
		d.AddListener(func(msg interface{}) {
			onMsg(msg.(int64))
		})
	}

	for i := 0; i < numMessages; i++ {
		atomic.AddInt64(&expected, numListeners)
		d.Dispatch(int64(1))
	}
	d.Close()

	time.Sleep(250 * time.Millisecond)
	return atomic.LoadInt64(&expected)
}
