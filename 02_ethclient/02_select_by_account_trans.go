package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	// client, err := ethclient.Dial("https://mainnet.infura.io/v3/bdb2ede84fe04e41a6fc9b2c9506d8c7")
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/2BDYF5KGJGH5VAIIPXVFB7K6KWDVQ36M1G")

	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	account := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e")
	ctx := context.Background()

	// 获取最新区块号
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	endBlock := header.Number.Int64()
	startBlock := endBlock - 10000 // 查询最近10000个区块

	// 限制查询的交易数量
	maxTransactions := 20
	transactionsFound := 0

	startTime := time.Now()
	fmt.Printf("开始查询区块范围 %d 到 %d...\n", startBlock, endBlock)

	for i := endBlock; i >= startBlock && transactionsFound < maxTransactions; i-- {
		block, err := client.BlockByNumber(ctx, big.NewInt(i))
		if err != nil {
			log.Printf("获取区块 %d 失败: %v", i, err)
			continue
		}

		for _, tx := range block.Transactions() {
			if transactionsFound >= maxTransactions {
				break
			}

			// 获取发送方地址
			from, err := client.TransactionSender(ctx, tx, block.Hash(), 0)
			if err != nil {
				continue
			}

			// 检查是否是该账户的交易
			if from == account || (tx.To() != nil && *tx.To() == account) {
				transactionsFound++
				printTransaction(tx, block.Number(), from)
			}
		}

		// 每查询100个区块打印一次进度
		if i%100 == 0 {
			fmt.Printf("已扫描到区块 %d，找到 %d 笔交易...\n", i, transactionsFound)
		}
	}

	fmt.Printf("查询完成，耗时 %v\n", time.Since(startTime))
}

func printTransaction(tx *types.Transaction, blockNumber *big.Int, from common.Address) {
	fmt.Println("\n==================================")
	fmt.Printf("交易哈希: %s\n", tx.Hash().Hex())
	fmt.Printf("区块高度: %d\n", blockNumber)
	fmt.Printf("发送方: %s\n", from.Hex())
	if tx.To() != nil {
		fmt.Printf("接收方: %s\n", tx.To().Hex())
	} else {
		fmt.Println("接收方: 合约创建")
	}
	fmt.Printf("金额: %s wei\n", tx.Value().String())
	fmt.Printf("Gas 价格: %s wei\n", tx.GasPrice().String())
	fmt.Printf("Nonce: %d\n", tx.Nonce())
	fmt.Println("==================================")
}
