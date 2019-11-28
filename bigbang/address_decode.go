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
	"encoding/base32"
	"encoding/binary"
	"encoding/hex"
	"errors"
)

type addressDecoder struct {
	wm *WalletManager //钱包管理者
}

//NewAddressDecoder 地址解析器
func NewAddressDecoder(wm *WalletManager) *addressDecoder {
	decoder := addressDecoder{}
	decoder.wm = wm
	return &decoder
}

//PrivateKeyToWIF 私钥转WIF
func (decoder *addressDecoder) PrivateKeyToWIF(priv []byte, isTestnet bool) (string, error) {
	return "", nil
}

func inverseBytes(data []byte) []byte {
	ret := make([]byte, len(data))
	for i := 0; i < len(data); i++ {
		ret[i] = data[len(data)-1-i]
	}
	return ret
}

//PublicKeyToAddress 公钥转地址
func (decoder *addressDecoder) PublicKeyToAddress(pub []byte, isTestnet bool) (string, error) {
	pubBigEndian := make([]byte, 0)
	chksum := crc24q(pub)
	tmp := [4]byte{}
	binary.BigEndian.PutUint32(tmp[:], chksum)

	pubBigEndian = append(pub, tmp[1:]...)

	address, err := decoder.wm.Client.importPubkey(hex.EncodeToString(inverseBytes(pub)))
	if err != nil {
		return "", errors.New("Failed to import public key into node!")
	}

	if address != "1" + base32.NewEncoding("0123456789abcdefghjkmnpqrstvwxyz").EncodeToString(pubBigEndian) {
		return "", errors.New("create address faile!")
	}
	return address, nil
}

//RedeemScriptToAddress 多重签名赎回脚本转地址
func (decoder *addressDecoder) RedeemScriptToAddress(pubs [][]byte, required uint64, isTestnet bool) (string, error) {
	return "", nil
}

//WIFToPrivateKey WIF转私钥
func (decoder *addressDecoder) WIFToPrivateKey(wif string, isTestnet bool) ([]byte, error) {
	return nil, nil

}
