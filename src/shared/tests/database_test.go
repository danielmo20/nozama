package main

import (
	nozama "nozama/src/shared"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrimaryKey(t *testing.T) {
	primaryKey := nozama.GeneratePrimaryKey()
	assert.NotEmpty(t, primaryKey)

}
