package future

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testFuture(args ...interface{}) interface{} {
	a := args[0].(int)
	b := args[1].(int)
	return a + b
}

func TestFuture(t *testing.T) {
	data := []int{1, 3}
	except := 4
	fu := GentFuture(testFuture, data[0], data[1])
	fu.Invoke()
	i := fu.Get()
	assert.True(t, i.(int) == except)
}
