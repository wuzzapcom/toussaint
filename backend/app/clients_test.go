package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetClientType(t *testing.T) {
	clientType, err := GetClientType("telegram")
	assert.Nil(t, err)
	assert.Equal(t, Telegram, clientType)

	clientType, err = GetClientType("smth",)
	assert.NotNil(t, err)
	assert.Equal(t, -1, int(clientType))
	//t.Log(clientType)
}
