package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
	"sync"
	"time"
)

var transactionChan chan []Transaction

func GetTxsByBlockNumber(blockNumber *big.Int) []*MorphTransaction {
	start := time.Now()
	var wg sync.WaitGroup

	transactionChan = make(chan []Transaction, 60)
	wg.Add(len(tables))
	var txs []Transaction

	for _, table := range tables {
		go getTxsByBlockNumber(table, blockNumber, &wg)
	}
	var k = 0
	for {
		if k == 60 {
			close(transactionChan)
		}
		k += 1
		val, ok := <-transactionChan
		if !ok {
			break
		}
		if len(val) == 0 {
			continue
		}
		txs = append(txs, val...)
	}

	wg.Wait()
	transactionChan = nil
	mtxs := dealTxs(&txs) // 处理一个区块的事务
	fmt.Println(time.Now(), "Get the transactions of block number("+blockNumber.String()+") spend "+time.Since(start).String())
	return mtxs
}

func GetTxByHash(hash common.Hash) []*MorphTransaction {
	start := time.Now()
	var wg sync.WaitGroup

	transactionChan = make(chan []Transaction, 60)
	wg.Add(len(tables))
	var txs []Transaction

	for _, table := range tables {
		go getTxByHash(table, hash, &wg)
	}
	var k = 0
	for {
		if k == 60 {
			close(transactionChan)
		}
		k += 1
		val, ok := <-transactionChan
		if !ok {
			break
		}
		if len(val) == 0 {
			continue
		}
		txs = append(txs, val...)
	}

	wg.Wait()
	transactionChan = nil
	mtxs := dealTxs(&txs) // 处理一个区块的事务
	fmt.Println(time.Now(), "Get the transactions by hash(\""+hash.Hex()+"\") spend "+time.Since(start).String())
	return mtxs
}

// dealTxs deal transactions of a block
func dealTxs(txs *[]Transaction) []*MorphTransaction {
	txSet := make(map[common.Hash]bool, 500)
	var morphTxs []*MorphTransaction
	dealNil := false
	for _, tx := range *txs {
		var morphTx *MorphTransaction = nil
		if tx.Hash == nil && !dealNil { // 处理挖矿奖励
			morphTx = dealRewardTx(&tx)
			dealNil = true
		} else if _, ok := txSet[*tx.Hash]; !ok { // 交易尚未处理过,防止之前存了重复的交易
			txSet[*tx.Hash] = true
			morphTx = dealGeneralTx(&tx)
		}

		if morphTx != nil {
			morphTxs = append(morphTxs, morphTx)
		}
	}
	return morphTxs
}

func dealGeneralTx(tx *Transaction) *MorphTransaction {
	morphTx := NewMorphTransaction()
	morphTx.Hash = tx.Hash.Hex()
	morphTx.BlockNumber = (*big.Int)(tx.BlockNumber) // 区块奖励交易记录, 只有区块号和区块Hash
	morphTx.BlockHash = tx.BlockHash.Hex()           // 区块奖励交易记录, 只有区块号和区块Hash
	morphTx.From = tx.From.Hex()
	if tx.To != nil { // create contract
		morphTx.To = tx.To.Hex()
	}
	for _, tr := range tx.Transfer {
		morphTr := dealTransfer(tr)
		morphTx.Transfer = append(morphTx.Transfer, morphTr)
	}
	return morphTx
}

func dealRewardTx(tx *Transaction) *MorphTransaction {
	morphTx := NewMorphTransaction()
	morphTx.Hash = tx.BlockNumber.String()
	morphTx.BlockNumber = (*big.Int)(tx.BlockNumber) // 区块奖励交易记录, 只有区块号和区块Hash
	morphTx.BlockHash = tx.BlockHash.Hex()           // 区块奖励交易记录, 只有区块号和区块Hash
	for _, tr := range tx.Transfer {
		morphTr := dealTransfer(tr)
		morphTx.Transfer = append(morphTx.Transfer, morphTr)
	}
	return morphTx
}

func dealTransfer(tr *Transfer) *MorphTransfer {
	morphTx := NewMorphTransfer()
	morphTx.From = tr.From.Address
	morphTx.To = tr.To.Address
	morphTx.Type = tr.Type
	morphTx.Index = tr.Nonce
	return morphTx
}

// getTxsByBlockNumber 获取某一区块号的所有交易
func getTxsByBlockNumber(table string, blockNumber *big.Int, wg *sync.WaitGroup) {
	defer wg.Done()
	sqlServer := openSqlServer()
	defer closeSqlServer(sqlServer)
	if sqlServer == nil {
		fmt.Println("sqlServer is nil")
		return
	}
	rows, err := sqlServer.Query("SELECT info FROM " + table + " WHERE blockNumber=\"" + (*hexutil.Big)(blockNumber).String() + "\";")
	defer rows.Close() // 非常重要：关闭rows释放持有的数据库链接
	if err != nil {
		fmt.Println("Query failed", err)
		return
	}
	// 循环读取结果集中的数据
	var txs []Transaction
	for rows.Next() {
		var (
			info string
			tx   Transaction
		)
		err = rows.Scan(&info)
		if err != nil {
			fmt.Println("Scan failed", err)
			return
		}
		err = json.Unmarshal([]byte(info), &tx)
		if err != nil {
			fmt.Println("Unmarshal failed", err)
			return
		}
		txs = append(txs, tx)
	}
	transactionChan <- txs
}

// getTxsByHash 获取某一交易
func getTxByHash(table string, hash common.Hash, wg *sync.WaitGroup) {
	defer wg.Done()
	sqlServer := openSqlServer()
	defer closeSqlServer(sqlServer)
	if sqlServer == nil {
		fmt.Println("sqlServer is nil")
		return
	}
	rows, err := sqlServer.Query("SELECT info FROM " + table + " WHERE hash=\"" + hash.Hex() + "\";")
	defer rows.Close() // 非常重要：关闭rows释放持有的数据库链接
	if err != nil {
		fmt.Println("Query failed", err)
		return
	}
	// 循环读取结果集中的数据
	var txs []Transaction
	for rows.Next() {
		var (
			info string
			tx   Transaction
		)
		err = rows.Scan(&info)
		if err != nil {
			fmt.Println("Scan failed", err)
			return
		}
		err = json.Unmarshal([]byte(info), &tx)
		if err != nil {
			fmt.Println("Unmarshal failed", err)
			return
		}
		txs = append(txs, tx)
	}
	transactionChan <- txs
}
