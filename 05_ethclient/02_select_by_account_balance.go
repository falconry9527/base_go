package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math"
	"math/big"
)

func main() {
	// 查询账户余额
	// 1. 连接到以太坊节点 (可以使用Infura或本地节点)
	// 主网链
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/bdb2ede84fe04e41a6fc9b2c9506d8c7")
	if err != nil {
		fmt.Printf("未查询到余额信息")
		log.Fatal(err)
	}

	// 2. 要查询的账户地址
	account := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc454e4438f44e") // 示例地址

	// 3. 查询余额 (返回的是 wei 单位)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		fmt.Printf("未查询到余额信息")
		log.Fatal(err)
	}

	// 4. 将 wei 转换为 ether
	// 1 ether = 10^18 wei
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

	// 5. 输出结果
	fmt.Printf("账户 %s 的余额: %f ETH\n", account.Hex(), ethValue)
	fmt.Printf("原始余额(wei): %s\n", balance.String())
}
