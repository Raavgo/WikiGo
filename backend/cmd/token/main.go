package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func main() {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		panic(err)
	}
	fmt.Println(hex.EncodeToString(token))
}
