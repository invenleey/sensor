package count

import (
	"testing"
	"time"
)

func TestAddErrorLog(t *testing.T) {
	AddErrorOperation("test01")
	time.Sleep(time.Second * 70)

	AddErrorOperation("test01")
	time.Sleep(time.Second * 130)

	AddErrorOperation("test01")
	time.Sleep(time.Second * 190)

	time.Sleep(time.Second * 5)
}
