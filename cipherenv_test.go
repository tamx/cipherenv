package cipherenv

import (
	"os"
	"testing"

	"github.com/kuking/seof/crypto"
	"github.com/stretchr/testify/assert"
)

const (
	tempFilename string = ".env.tmp"
)

func TestMain(m *testing.M) {
	// ここにテストの初期化処理

	code := m.Run()
	// ここでテストのお片づけ
	os.Exit(code)
}

func TestCipher(t *testing.T) {
	password := "this is a very long password nobody should know about."
	plainData := crypto.RandBytes(100)
	key := calKey(password)

	err := save(tempFilename, key, plainData)
	assert.NoError(t, err)

	decryptedData, err := load(tempFilename, key)
	assert.NoError(t, err)
	assert.Equal(t, plainData, decryptedData)
}
