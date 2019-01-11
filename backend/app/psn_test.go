package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSearch(t *testing.T) {
	searched, err := SearchByName("Человек паук")
	assert.Nil(t, err)

	for _, s := range searched {
		t.Log(s)
		game, err := SearchByID(s.Id)
		assert.Nil(t, err)
		assert.Equal(t, s, game)
	}
}
