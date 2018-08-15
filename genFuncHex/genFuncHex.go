package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/crypto"
)

// calc hashId for function
func keccak(args []string) {
	for _, arg := range args {
		hash := crypto.Keccak256Hash([]byte(arg))
		fmt.Printf("arg:%s,hash:%s\n", arg, hash.String())
	}
}

func main() {
	fmt.Printf("usage: %s transfer(address,uint256)\n\n", os.Args[0])
	fmt.Printf("args:%v\n", os.Args[1:])
	keccak(os.Args[1:])
}
