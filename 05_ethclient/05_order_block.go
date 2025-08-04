package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

// 订阅区块
func main() {
	client, err := ethclient.Dial("wss://sepolia.infura.io/ws/v3/bdb2ede84fe04e41a6fc9b2c9506d8c7")
	if err != nil {
		log.Fatal(err)
	}
	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			fmt.Println("-----header -------------")
			fmt.Println(header.Number.Uint64())     // 5671744
			fmt.Println(header.Time)                // 1527211625
			fmt.Println(header.Difficulty.Uint64()) // 3217000136609065
			fmt.Println(header.Hash())              // 0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9

			//block, err := client.BlockByHash(context.Background(), header.Hash())
			//if err != nil {
			//	log.Fatal(err)
			//}
			//fmt.Println("-----block-------------")
			//fmt.Println(block.Number().Uint64())     // 5671744
			//fmt.Println(block.Time())                // 1712798400
			//fmt.Println(block.Difficulty().Uint64()) // 0
			//fmt.Println(block.Hash().Hex())          // 0x9e8751ebb5069389b855bba72d94902cc385042661498a415979b7b6ee9ba4b9
			//fmt.Println(len(block.Transactions()))   // 70
		}
	}
}
