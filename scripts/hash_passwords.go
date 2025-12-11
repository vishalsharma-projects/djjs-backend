package main

import (
	"flag"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := flag.String("password", "", "The password to hash")
	flag.Parse()

	if *password == "" {
		log.Fatal("Please provide a password using -password flag")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Error hashing password: %v", err)
	}

	fmt.Printf("Plain Password: %s\n", *password)
	fmt.Printf("Hashed Password: %s\n", string(hashedPassword))
}
