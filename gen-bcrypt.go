package main

import (
	"flag"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := flag.String("p", "", "The Password to HASH")
	flag.Parse()

	if len(*password) < 1 {
		fmt.Println("Please provide something to hash with -p")
		return
	}

	tohash := []byte(*password)
	hash, _ := bcrypt.GenerateFromPassword(tohash, 14)

	fmt.Println("Hash: ", string(hash))
}
