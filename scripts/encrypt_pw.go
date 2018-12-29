package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/cozy/cozy-stack/pkg/crypto"
)

func main() {
	var pass []byte
	// check passphrase environment variable or prompt user for passphrase
	if passstr, found := os.LookupEnv("passphrase"); found {
		pass = []byte(passstr)
	} else {
		reader := bufio.NewScanner(os.Stdin)
		fmt.Print("enter passphrase: ")
		reader.Scan()
		pass = reader.Bytes()
		if err := reader.Err(); err != nil {
			log.Fatal(err)
		}
	}
	hash, err := crypto.GenerateFromPassphrase(pass)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(hash))
}
