package cipherenv

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"errors"
	"os"
	"strings"
)

func encryptByGCM(key, plainText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize()) // Unique nonce is required(NonceSize 12byte)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	cipherText := gcm.Seal(nil, nonce, plainText, nil)
	cipherText = append(nonce, cipherText...)

	return cipherText, nil
}

func decryptByGCM(key, cipherText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := cipherText[:gcm.NonceSize()]
	plainByte, err := gcm.Open(nil, nonce, cipherText[gcm.NonceSize():], nil)
	if err != nil {
		return nil, err
	}

	return plainByte, nil
}

func save(filename string, password, plainData []byte) error {
	cipherData, err := encryptByGCM(password, plainData)
	if err != nil {
		return err
	}
	// 書き込み
	err = os.WriteFile(filename, cipherData, 0644)
	if err != nil {
		return err
	}
	return nil
}

func load(filename string, password []byte) ([]byte, error) {
	cipherData, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	plainData, err := decryptByGCM(password, cipherData)
	if err != nil {
		return nil, err
	}
	return plainData, nil
}

func calKey(secretKey string) []byte {
	key := []byte(secretKey)
	md5key := md5.Sum(key)
	// md5hex := fmt.Sprintf("%x", md5key)
	// println(md5hex)
	return md5key[:]
}

var envCache map[string]string = map[string]string{}

func parseEnv(envData []byte) error {
	envCache = map[string]string{}
	for _, line := range strings.Split(string(envData), "\n") {
		// println(line)
		line = strings.TrimPrefix(line, "\a")
		line = strings.TrimSuffix(line, "\a")
		if strings.HasPrefix(line, "#") {
			continue
		}
		if strings.TrimSpace(line) == "" {
			continue
		}
		index := strings.Index(line, "=")
		if index < 0 {
			return errors.New("syntax error: " + line)
		}
		key := strings.TrimSpace(line[:index])
		// println(key)
		if key == "" {
			continue
		}
		value := strings.TrimSpace(line[(index + 1):])
		if value[0] == '"' && value[len(value)-1] == '"' {
			value = value[1 : len(value)-1]
		}
		envCache[key] = value
	}
	return nil
}

func LoadEnv(envFilename, secretKey string) error {
	password := calKey(secretKey)
	envData, err := load(envFilename, password)
	if err != nil {
		return err
	}
	err = parseEnv(envData)
	if err != nil {
		return err
	}
	return nil
}

func Keys() []string {
	keys := make([]string, 0)
	for k := range envCache {
		keys = append(keys, k)
	}
	return keys
}

func Get(key string) string {
	return envCache[key]
}

func Create(envFilename, secretKey string,
	plainData []byte) error {
	password := calKey(secretKey)
	err := save(envFilename, password, plainData)
	if err != nil {
		return err
	}
	return nil
}
