package main

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const apiKey = "https://mainnet.infura.io/v3/c5a6f56c3082466bbf869fd2896bf024"

func main() {
	file := "./keys.txt"
	keys := Readfile(file)
	for _, key := range keys {
		addr := AddrGenerator(key)
		balance := GetBalance(addr)

		fmt.Printf("Private key: %s\n", key)
		fmt.Printf("Address: 0x%x\n", addr)
		fmt.Println("Balance: ", balance)
		fmt.Println("-------------------------------------------------")
	}
}

func Readfile(path string) []string {
	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func AddrGenerator(keyHex string) common.Address {
	key, err := crypto.HexToECDSA(keyHex)
	if err != nil {
		panic(err)
	}

	publicKey := key.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("invalid key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	return fromAddress
}

func GetBalance(addr common.Address) float64 {
	client, err := ethclient.Dial(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	balance, err := client.BalanceAt(context.Background(), addr, nil)
	if err != nil {
		log.Fatal(err)
	}

	balancefloat := new(big.Float).SetInt(balance)
	i := new(big.Float)
	i.SetString("1000000000000000000")
	f := new(big.Float).Quo(balancefloat, i)
	b, _ := f.Float64()

	return b
}

