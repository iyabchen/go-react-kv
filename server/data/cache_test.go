package data

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotExist(t *testing.T) {
	ds, err := NewMem()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	_, err = ds.GetOne(ctx, "random-id")
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "not exist"))

	err = ds.DeleteOne(ctx, "random-id")
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "not exist"))

	err = ds.Update(ctx, "random-id", "k", "v")
	assert.NotNil(t, err)
	assert.True(t, strings.Contains(err.Error(), "not exist"))

}
