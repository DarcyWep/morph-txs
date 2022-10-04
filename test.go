package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	_ "github.com/go-sql-driver/mysql" // 导入包但不使用, init()
	"math/big"
)

func test() {
	SetHost("202.114.6.20")
	//min, max := big.NewInt(7202102), big.NewInt(8852777)
	mtxs := GetTxsByBlockNumber(big.NewInt(7202102))
	for _, tx := range mtxs {
		fmt.Println(tx.Hash, tx.BlockNumber, tx.BlockHash, tx.From, tx.To, len(tx.Transfer))
	}
	mtxs = GetTxByHash(common.HexToHash("0x5d6a7e06db5644a7921829094752a36ed3498376ab4a9b5dc69e30661bfe6828"))
	for _, tx := range mtxs {
		fmt.Println(tx.Hash, tx.BlockNumber, tx.BlockHash, tx.From, tx.To, len(tx.Transfer))
	}
}
