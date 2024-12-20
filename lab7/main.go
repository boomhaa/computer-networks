package main

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"google.golang.org/api/option"
	"log"
	"math/big"
	"time"
)

func senddata(block *types.Block) {
	opt := option.WithCredentialsFile("./lab7.json")
	config := &firebase.Config{
		DatabaseURL: "https://lab7-6e4bd-default-rtdb.firebaseio.com",
	}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}
	client, err := app.Database(context.Background())
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	ref := client.NewRef("/" + block.Number().String())
	send_data := map[string]interface{}{
		"number":     block.Number().String(),
		"time":       block.Time(),
		"difficulty": block.Difficulty().String(),
		"hash":       block.Hash().String(),
	}

	if err := ref.Set(context.Background(), send_data); err != nil {
		log.Fatalf("error sending data: %v\n", err)
	}

	for _, tx := range block.Transactions() {
		if tx.To() == nil {
			continue
		}
		ref = client.NewRef("/" + block.Number().String() + "/transactions/" + tx.Hash().String())
		transaction_data := map[string]interface{}{
			"chainID":   tx.ChainId(),
			"hash":      tx.Hash().String(),
			"to":        tx.To().Hex(),
			"value":     tx.Value().String(),
			"cost":      tx.Cost().String(),
			"gas":       tx.Gas(),
			"gas price": tx.GasPrice().String()}
		if err := ref.Set(context.Background(), transaction_data); err != nil {
			log.Fatalf("error sending data: %v\n", err)
		}
	}

	fmt.Println("sent data block number: " + block.Number().String())
}

func main() {
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/1cb3c9bac5ad452eb9e994c328660611")
	if err != nil {
		log.Fatalln(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	header, err := client.HeaderByNumber(context.Background(), nil)
	prev := header.Number.String()
	for {
		header, err := client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}
		if prev != header.Number.String() {
			prev = header.Number.String()
			fmt.Println("New block")
			fmt.Println(header.Number.String()) // The lastes block in blockchain because nil pointer in header

			blockNumber := big.NewInt(header.Number.Int64())
			block, err := client.BlockByNumber(context.Background(), blockNumber) //get block with this number
			if err != nil {
				log.Fatal(err)
			}
			// all info about block
			fmt.Println(block.Number().String())
			fmt.Println(block.Time())
			fmt.Println(block.Difficulty().String())
			fmt.Println(block.Hash().String())
			fmt.Println(len(block.Transactions()))
			go senddata(block)
			fmt.Println()
		}
		time.Sleep(1000 * time.Millisecond)
	}
}
