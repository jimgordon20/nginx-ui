package analytics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSummarizeStatusRates(t *testing.T) {
	success, clientErr, serverErr, total := summarizeStatusRates(map[string]int{
		"200": 80,
		"404": 15,
		"500": 5,
	})
	assert.Equal(t, 100, total)
	assert.InDelta(t, 80.0, success, 0.01)
	assert.InDelta(t, 15.0, clientErr, 0.01)
	assert.InDelta(t, 5.0, serverErr, 0.01)
}

func TestSummarizeStatusRatesEmpty(t *testing.T) {
	success, clientErr, serverErr, total := summarizeStatusRates(map[string]int{})
	assert.Equal(t, 0, total)
	assert.Equal(t, 0.0, success)
	assert.Equal(t, 0.0, clientErr)
	assert.Equal(t, 0.0, serverErr)
}

func TestStringAndNumericFieldHelpers(t *testing.T) {
	fields := map[string]interface{}{
		"ip":         "1.2.3.4",
		"status":     float64(200),
		"bytes_sent": float64(1024),
		"timestamp":  int64(1700000000),
	}
	assert.Equal(t, "1.2.3.4", stringField(fields, "ip"))
	assert.Equal(t, "", stringField(fields, "missing"))
	assert.Equal(t, 200, intField(fields, "status"))
	assert.Equal(t, int64(1024), int64Field(fields, "bytes_sent"))
	assert.Equal(t, int64(1700000000), int64Field(fields, "timestamp"))
}
