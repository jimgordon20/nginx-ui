package analytics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExtension(t *testing.T) {
	assert.Equal(t, ".webp", fileExtension("/assets/cover.webp"))
	assert.Equal(t, ".png", fileExtension("/img/a.PNG?v=1"))
	assert.Equal(t, ".json", fileExtension("/api/data.json#frag"))
	assert.Equal(t, "(none)", fileExtension("/"))
	assert.Equal(t, "(none)", fileExtension("/comics/issue-12"))
	assert.Equal(t, "(none)", fileExtension(""))
}

func TestAggregateExtensions(t *testing.T) {
	paths := []KeyValue{
		{Key: "/a.webp", Value: 50},
		{Key: "/b.webp", Value: 30},
		{Key: "/c.png", Value: 20},
		{Key: "/page", Value: 10},
	}

	result := AggregateExtensions(paths, 10)
	assert.Len(t, result, 3)
	assert.Equal(t, ".webp", result[0].Extension)
	assert.Equal(t, 80, result[0].Count)
	assert.Equal(t, ".png", result[1].Extension)
	assert.Equal(t, 20, result[1].Count)
	assert.Equal(t, "(none)", result[2].Extension)
	assert.Equal(t, 10, result[2].Count)
}

func TestAggregateExtensionsLimit(t *testing.T) {
	paths := []KeyValue{
		{Key: "/a.webp", Value: 5},
		{Key: "/b.png", Value: 4},
		{Key: "/c.jpg", Value: 3},
		{Key: "/d.gif", Value: 2},
	}
	result := AggregateExtensions(paths, 2)
	assert.Len(t, result, 2)
	assert.Equal(t, ".webp", result[0].Extension)
	assert.Equal(t, ".png", result[1].Extension)
}
