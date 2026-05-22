package main

import (
	"log"
	"os"
	"strings"

	cipherenv "cipherenv"

	"golang.org/x/term"
)

func askCred() string {
	print("Enter Password: ")
	bytePassword, err := term.ReadPassword(0)
	if err != nil {
		log.Fatal(err)
	}
	// println("\nPassword typed: " + string(bytePassword))
	password := string(bytePassword)

	return strings.TrimSpace(password)
}

func create(envFilename, orgFilename, secretKey string) {
	plainData, err := os.ReadFile(orgFilename)
	if err != nil {
		println(err.Error())
		return
	}
	println("start: " + envFilename)
	err = cipherenv.Create(envFilename, secretKey, plainData)
	if err != nil {
		println(err.Error())
		return
	}
}

func read(envFilename, secretKey string) {
	env, err := cipherenv.LoadEnv(envFilename, secretKey)
	if err != nil {
		log.Fatal(err)
	}
	for _, key := range env.Keys() {
		value := env.Get(key)
		println(key + " = " + value)
	}
}

func usage() {
	envFilename := os.Getenv("ENV_FILE")
	secretKey := "secret:" + os.Getenv("ENV_KEY")
	env, err := cipherenv.LoadEnv(envFilename, secretKey)
	if err != nil {
		return
	}
	for _, key := range env.Keys() {
		value := env.Get(key)
		println(key + " = " + value)
	}
}

func main() {
	if len(os.Args) < 2 {
		println("create: cipherenv [cipherfilename] [orgenvfile]")
		println("read:   cipherenv [filename]")
		return
	}
	envFilename := os.Args[1]
	secretKey := askCred()
	if len(os.Args) == 3 {
		orgFilename := os.Args[2]
		create(envFilename, orgFilename, secretKey)
		return
	}
	read(envFilename, secretKey)
}
