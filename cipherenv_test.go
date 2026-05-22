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

func TestEnvParsing1(t *testing.T) {
	password := "test-password"
	envContent := []byte("KEY1=VALUE1\nKEY2=\"VALUE2\"\n#COMMENT=IGNORE\nEMPTY_VALUE=\n")

	err := Create(tempFilename, password, envContent)
	assert.NoError(t, err)
	defer os.Remove(tempFilename)

	env, err := LoadEnv(tempFilename, password)
	assert.NoError(t, err)

	assert.Equal(t, "VALUE1", env.Get("KEY1"))
	assert.Equal(t, "VALUE2", env.Get("KEY2"))
	assert.Equal(t, "", env.Get("EMPTY_VALUE")) // 空の値のテストを追加
	assert.Contains(t, env.Keys(), "KEY1")
	assert.Contains(t, env.Keys(), "KEY2")
	assert.Contains(t, env.Keys(), "EMPTY_VALUE")
	assert.NotContains(t, env.Keys(), "COMMENT") // コメントがキーに含まれないことを確認
}

func TestEnvParsing2(t *testing.T) {
	password := "test-password"
	envContent := []byte("KEY1=VALUE1\nKEY2=\"VALUE2\"\n#COMMENT=IGNORE")

	err := Create(tempFilename, password, envContent)
	assert.NoError(t, err)
	defer os.Remove(tempFilename)

	env, err := LoadEnv(tempFilename, password)
	assert.NoError(t, err)

	assert.Equal(t, "VALUE1", env.Get("KEY1"))
	assert.Equal(t, "VALUE2", env.Get("KEY2"))
	assert.Contains(t, env.Keys(), "KEY1")
}
