package media

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSaveUploadRejectsNonImage(t *testing.T) {
	err := SaveUpload("site", "", "notes.txt", strings.NewReader("hello"), 5)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "only image files")
}
