package sensor

import "testing"

func TestMeasureRequest(t *testing.T) {
	MeasureRequest(SessionCollection[].conn.RemoteAddr().String(), func(meta interface{}, data []byte) {

	})
}