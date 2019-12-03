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
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"

	//"math/big"

	"github.com/blocktree/openwallet/log"
	"github.com/imroc/req"
	"github.com/tidwall/gjson"
)

type ClientInterface interface {
	Call(path string, request []interface{}) (*gjson.Result, error)
}

// A Client is a Bitcoin RPC client. It performs RPCs over HTTP using JSON
// request and responses. A Client must be configured with a secret token
// to authenticate with other Cores on the network.
type Client struct {
	BaseURL     string
	AccessToken string
	Debug       bool
	client      *req.Req
	//Client *req.Req
}

type Response struct {
	Code    int         `json:"code,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Message string      `json:"message,omitempty"`
	Id      string      `json:"id,omitempty"`
}


func NewClient(url, token string, debug bool) *Client {
	c := Client{
		BaseURL:     url,
		AccessToken: token,
		Debug:       debug,
	}

	api := req.New()
	//trans, _ := api.Client().Transport.(*http.Transport)
	//trans.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	c.client = api

	return &c
}

// Call calls a remote procedure on another node, specified by the path.
func (c *Client) Call(path string, request map[string]interface{}) (*gjson.Result, error) {

	var (
		body = make(map[string]interface{}, 0)
	)

	if c.client == nil {
		return nil, errors.New("API url is not setup. ")
	}

	authHeader := req.Header{
		"Accept":        "application/json",
		"Authorization": "Basic " + c.AccessToken,
	}

	//json-rpc
	body["jsonrpc"] = "2.0"
	body["id"] = "1"
	body["method"] = path
	body["params"] = request//req.BodyJSON(request)

	if c.Debug {
		log.Std.Info("Start Request API...")
	}

	r, err := c.client.Post(c.BaseURL, req.BodyJSON(&body), authHeader)

	if c.Debug {
		log.Std.Info("Request API Completed")
	}

	if c.Debug {
		log.Std.Info("%+v", r)
	}

	if err != nil {
		return nil, err
	}

	resp := gjson.ParseBytes(r.Bytes())
	err = isError(&resp)
	if err != nil {
		return nil, err
	}

	result := resp.Get("result")

	return &result, nil
}

// See 2 (end of page 4) http://www.ietf.org/rfc/rfc2617.txt
// "To receive authorization, the client sends the userid and password,
// separated by a single colon (":") character, within a base64
// encoded string in the credentials."
// It is not meant to be urlencoded.
func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

//isError 是否报错
func isError(result *gjson.Result) error {
	var (
		err error
	)

	/*
		//failed 返回错误
		{
			"result": null,
			"error": {
				"code": -8,
				"message": "Block height out of range"
			},
			"id": "foo"
		}
	*/

	if !result.Get("error").IsObject() {

		if !result.Get("result").Exists() {
			return errors.New("Response is empty! ")
		}

		return nil
	}



	errInfo := fmt.Sprintf("[%d]%s",
		result.Get("error.code").Int(),
		result.Get("error.message").String())
	err = errors.New(errInfo)

	return err
}

// 获取当前区块高度
func (c *Client) getBlockHeight() (uint64, error) {

	path := "getblockcount"
	request := map[string]interface{}{
	}
	resp, err := c.Call(path, request)

	if err != nil {
		return 0, err
	}

	return resp.Uint() - 1, nil
}

// 通过高度获取区块哈希
func (c *Client) getBlockHash(height uint64) (string, error) {

	path := "getblockhash"
	request := map[string]interface{}{
		"height":height,
	}

	resp, err := c.Call(path, request)

	if err != nil {
		return "", err
	}

	return resp.Array()[0].String(), nil
}

func (c *Client) importPubkey(pub string) (string, error) {
	path := "importpubkey"
	request := map[string]interface{}{
		"pubkey":pub,
	}
	resp, err := c.Call(path, request)

	if err != nil {
		return "", err
	}
	return resp.String(), nil
}

// 获取地址余额
func (c *Client) getBalance(address, anchor string) (*AddrBalance, error) {

	//var amount = uint64(0)
	//
	//path := "listunspent"
	//request := map[string]interface{}{
	//	"forkid": anchor,
	//	"address":address,
	//	"max": 0,
	//	"sum": true,
	//}
	//
	//resp, err := c.Call(path, request)
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//utxos := resp.Get("unspents").Array()
	//if len(utxos) == 0 {
	//	amount = 0
	//} else {
	//	amount = convertFromAmount(resp.Get("sum").String())
	//}
	//
	//return &AddrBalance{Address: address, Balance: big.NewInt(int64(amount))}, nil

	path := "getbalance"

	request := map[string]interface{}{
		"address":address,
	}

	resp, err := c.Call(path, request)

	if err != nil {
		return nil, err
	}

	if len(resp.Array()) == 0 {
		return &AddrBalance{
			Address:address,
			Balance:big.NewInt(0),
		}, nil
	}

	return &AddrBalance{
		Address:address,
		Balance:big.NewInt(int64(convertFromAmount(resp.Array()[0].Get("avail").String()))),
	}, nil
}

type UnSpent struct {
	TxID string
	Vout byte
	Amount uint64
}

func (c *Client) listUnnSpent (address, anchor string) ([]UnSpent, error) {

	ret := make([]UnSpent, 0)
	path := "listunspent"
	request := map[string]interface{}{
		"forkid": anchor,
		"address":address,
		"max": 0,
		"sum": true,
	}

	resp, err := c.Call(path, request)

	if err != nil {
		return nil, err
	}

	utxos := resp.Get("unspents").Array()

	if len(utxos) == 0 {
		return []UnSpent{}, nil
	} else {
		for _, utxo := range utxos {
			ret = append(ret, UnSpent{
				TxID:   utxo.Get("txid").String(),
				Vout:   byte(utxo.Get("out").Uint()),
				Amount: convertFromAmount(utxo.Get("amount").String()),
			})
		}
	}

	return ret, nil
}

type UTXOinPool struct {
	TxID string
	Vout byte
}

func (c *Client) getUTXOsInPool() ([]UTXOinPool, error) {
	ret := make([]UTXOinPool, 0)
	path := "gettxpool"

	request := map[string]interface{}{
		"detail":true,
	}

	resp, err := c.Call(path, request)
	if err != nil {
		return nil, err
	}

	if resp.Raw == "{}" {
		return []UTXOinPool{}, nil
	} else {
		txs := resp.Get("list").Array()
		for _, txid := range txs {
			path = "gettransaction"
			request = map[string]interface{}{
				"txid":txid.Get("hex").String(),
				"serialized": false,
			}
			resp, err = c.Call(path, request)
			if err != nil {
				return nil, err
			}
			vins := resp.Get("transaction").Get("vin").Array()
			for _, in := range vins {
				ret = append(ret, UTXOinPool{
					TxID: in.Get("txid").String(),
					Vout: byte(in.Get("vout").Uint()),
				})
			}
		}
	}

	return ret, nil
}

func isUnspentAlreadyInPool(utxos []UTXOinPool, utxo UnSpent) bool {
	if utxos == nil || len(utxos) == 0 {
		return false
	}

	for _, v := range utxos {
		if v.TxID == utxo.TxID && v.Vout == utxo.Vout {
			return true
		}
	}
	return false
}

// 获取区块信息
func (c *Client) getBlock(hash string) (*Block, error) {
	return nil, nil
}

func (c *Client) getBlockByHeight(height uint64) (*Block, error) {

	path := "getblockhash"
	request := map[string]interface{}{
		"height":height,
	}

	resp, err := c.Call(path, request)

	if err != nil {
		return nil, err
	}

	hash := resp.Array()[0].String()

	path = "getblock"
	request = map[string]interface{}{
		"block":hash,
	}

	resp, err = c.Call(path, request)

	if err != nil {
		return nil, err
	}

	return NewBlock(resp), nil
}

func (c *Client) getToFromTxID (txid string, vout byte) (string, error) {

	to := ""

	path := "gettransaction"
	request := map[string]interface{}{
		"txid":txid,
		"serialized": false,
	}

	resp, err := c.Call(path, request)

	if err != nil {
		return "", err
	}
	isChange := (vout == 1)
	for {
		if isChange {
			request = map[string]interface{}{
				"txid":resp.Get("transaction").Get("vin").Array()[0].Get("txid").String(),
				"serialized": false,
			}
			resp, err = c.Call(path, request)

			if err != nil {
				return "", err
			}

			if len(resp.Get("transaction").Get("vin").Array()) == 0 {
				to = resp.Get("transaction").Get("sendto").String()
				break
			}
			isChange = (resp.Get("transaction").Get("vin").Array()[0].Get("vout").Uint() == 1)

		} else {
			if len(resp.Get("transaction").Get("vin").Array()) == 0 {
				to = resp.Get("transaction").Get("sendto").String()
				break
			}
			request = map[string]interface{}{
				"txid":resp.Get("transaction").Get("vin").Array()[0].Get("txid").String(),
				"serialized": false,
			}
			resp, err = c.Call(path, request)

			if err != nil {
				return "", err
			}
			to = resp.Get("transaction").Get("sendto").String()
			break
		}
	}

	return to, nil
}

func (c *Client) getTransaction(txid string) (*Transaction, error) {
	path := "gettransaction"
	request := map[string]interface{}{
		"txid":txid,
		"serialized": false,
	}

	resp, err := c.Call(path, request)

	if err != nil {
		return nil, err
	}

	return c.NewTransaction(resp), nil
}


func (c *Client) getAnchor() (string, error) {
	path := "getblockhash"
	request := map[string]interface{}{
		"height":0,
	}

	resp, err := c.Call(path, request)

	if err != nil {
		return "", err
	}

	return  resp.Array()[0].String(), nil
}

func (c *Client) sendTransaction(rawTx string) (string, error) {
	path := "sendtransaction"

	request := map[string]interface{}{
		"txdata":rawTx,
	}

	resp, err := c.Call(path, request)

	if err != nil {
		return "", err
	}

	return resp.String(), nil

}

func (c *Client) getContractAccountBalence(regid, address string) (*AddrBalance, error) {
	return nil, errors.New("Contract is not supported!")
}
