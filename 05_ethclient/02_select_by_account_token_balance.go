package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
)

func getERC20Balance(rpcURL, contractAddress, userAddress string) (*big.Int, error) {
	// 连接以太坊节点
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %v", err)
	}
	defer client.Close()

	// 将地址转换为common.Address类型
	contractAddr := common.HexToAddress(contractAddress)
	userAddr := common.HexToAddress(userAddress)

	// 构造balanceOf调用的数据
	// 方法签名: balanceOf(address)
	// 函数选择器是前4个字节的keccak256哈希: 0x70a08231
	// 参数是32字节填充的用户地址
	data := hexutil.MustDecode("0x70a08231") // balanceOf的函数选择器
	data = append(data, common.LeftPadBytes(userAddr.Bytes(), 32)...)

	// 调用合约
	msg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}
	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, fmt.Errorf("contract call failed: %v", err)
	}

	// 解析返回的余额(32字节的大整数)
	balance := new(big.Int).SetBytes(result)
	return balance, nil
}

func main() {
	rpcURL := "https://sepolia.infura.io/v3/bdb2ede84fe04e41a6fc9b2c9506d8c7" // 替换为你的Infura项目ID
	contractAddress := "0x24cD7079cA1329Bd083E436527C8C4B7c2942Fb7"           // USDT合约地址
	userAddress := "0x0E2Ea7deceA74b1dAc4cfaeC79d6e7f702078A41"               // 示例地址(币安热钱包)

	balance, err := getERC20Balance(rpcURL, contractAddress, userAddress)
	if err != nil {
		log.Fatalf("Failed to get balance: %v", err)
	}

	// 假设代币有6位小数(如USDT)
	decimalDivisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil)
	humanBalance := new(big.Float).Quo(
		new(big.Float).SetInt(balance),
		new(big.Float).SetInt(decimalDivisor),
	)

	fmt.Printf("Token balance: %s (raw: %s)\n", humanBalance.String(), balance.String())
}
