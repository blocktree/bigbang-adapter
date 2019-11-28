/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package bigbang

import (
	"fmt"

	"github.com/blocktree/openwallet/crypto"
	"github.com/blocktree/openwallet/openwallet"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tidwall/gjson"
)

type Block struct {
	Hash                  string
	PrevBlockHash         string
	Fork                  string
	TransactionMerkleRoot string
	Timestamp             uint64
	Height                uint64
	Transactions          []string
}

type Transaction struct {
	TxID            string
	TimeStamp       uint64
	Type            string
	Anchor          string
	From            string
	Amount          uint64
	Fee             uint64
	To              string
	BlockHeight     uint64
	BlockHash       string
	Confirmations   uint64
	Memo            string
}

func (c *Client)NewTransaction(json *gjson.Result) *Transaction {
	obj := &Transaction{}

	obj.TxID = json.Get("transaction").Get("txid").String()
	obj.TimeStamp = json.Get("transaction").Get("time").Uint()
	obj.Type = json.Get("transaction").Get("type").String()
	obj.Anchor = json.Get("transaction").Get("anchor").String()
	from, _ := c.getToFromTxID(json.Get("transaction").Get("vin").Array()[0].Get("txid").String(), byte(json.Get("transaction").Get("vin").Array()[0].Get("vout").Uint()))
	obj.From = from
	obj.Amount = convertFromAmount(json.Get("transaction").Get("amount").String())
	obj.Fee = convertFromAmount(json.Get("transaction").Get("txfee").String())
	obj.To = json.Get("transaction").Get("sendto").String()
	obj.Confirmations = json.Get("transaction").Get("confirmations").Uint()
	obj.Memo = json.Get("transaction").Get("data").String()

	return obj
}

func NewBlock(json *gjson.Result) *Block {
	obj := &Block{}

	obj.Hash = json.Get("hash").String()
	obj.PrevBlockHash = json.Get("hashPrev").String()
	obj.Fork = json.Get("fork").String()
	obj.TransactionMerkleRoot = json.Get("txmint").String()
	obj.Timestamp = json.Get("time").Uint()
	obj.Height = json.Get("height").Uint()

	for _, tx := range json.Get("tx").Array() {
		obj.Transactions = append(obj.Transactions, tx.String())

	}

	return obj
}

//BlockHeader 区块链头
func (b *Block) BlockHeader() *openwallet.BlockHeader {

	obj := openwallet.BlockHeader{}
	//解析json
	obj.Hash = b.Hash
	//obj.Confirmations = b.Confirmations
	obj.Merkleroot = b.TransactionMerkleRoot
	obj.Previousblockhash = b.PrevBlockHash
	obj.Height = b.Height
	obj.Time = b.Timestamp
	obj.Symbol = Symbol

	return &obj
}

//UnscanRecords 扫描失败的区块及交易
type UnscanRecord struct {
	ID          string `storm:"id"` // primary key
	BlockHeight uint64
	TxID        string
	Reason      string
}

func NewUnscanRecord(height uint64, txID, reason string) *UnscanRecord {
	obj := UnscanRecord{}
	obj.BlockHeight = height
	obj.TxID = txID
	obj.Reason = reason
	obj.ID = common.Bytes2Hex(crypto.SHA256([]byte(fmt.Sprintf("%d_%s", height, txID))))
	return &obj
}
