package bigbang

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/shopspring/decimal"
)

const (
	testNodeAPI = "http://ip:port"
)

func Test_getBlockHeight(t *testing.T) {
	token := ""
	c := NewClient(testNodeAPI, token, true)

	r, err := c.getBlockHeight()

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("height:", r)
	}

}

func Test_gettxpool(t *testing.T) {
	token := ""
	c := NewClient(testNodeAPI, token, true)
	r, err := c.getUTXOsInPool()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
}

func Test_getBlockByHeight(t *testing.T) {
	token := ""
	c := NewClient(testNodeAPI, token, true)
	r, err := c.getBlockByHeight(0)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
}

func Test_getBlockHash(t *testing.T) {
	token := ""
	c := NewClient(testNodeAPI, token, true)

	height := uint64(3036)

	r, err := c.getBlockHash(height)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
}

func Test_getBalance(t *testing.T) {
	token := ""
	c := NewClient(testNodeAPI, token, true)

	address := "1s98h31v48qdkjnxsxp1p9z7k2xcajh8jvmy7wvdnsxfc3fww9hdawjam"

	r, err := c.getBalance(address, "0000000083c3a4ac4f18efaadb6452002c922be7f7204a669e9a6b279aa3a4a2")

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
}

func Test_getTransaction(t *testing.T) {
	token := ""
	c := NewClient(testNodeAPI, token, true)
	txid := "5dcbc00095199a79a9f3eec4afc281d03ab52b0834dbe1959709df65c361d904" //"9KBoALfTjvZLJ6CAuJCGyzRA1aWduiNFMvbqTchfBVpF"

	r, err := c.getTransaction(txid)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r)
	}
}

func Test_convert(t *testing.T) {

	amount := uint64(5000000001)

	amountStr := fmt.Sprintf("%d", amount)

	fmt.Println(amountStr)

	d, _ := decimal.NewFromString(amountStr)

	w, _ := decimal.NewFromString("100000000")

	d = d.Div(w)

	fmt.Println(d.String())

	d = d.Mul(w)

	fmt.Println(d.String())

	r, _ := strconv.ParseInt(d.String(), 10, 64)

	fmt.Println(r)

	fmt.Println(time.Now().UnixNano())
}

func Test_getTransactionByAddresses(t *testing.T) {
	addrs := "ARAA8AnUYa4kWwWkiZTTyztG5C6S9MFTx11"

	token := ""
	c := NewClient(testNodeAPI, token, true)
	result, err := c.getMultiAddrTransactions(0, -1, addrs)

	if err != nil {
		t.Error("get transactions failed!")
	} else {
		for _, tx := range result {
			fmt.Println(tx.TxID)
		}
	}
}

func Test_getContractAccountInfo(t *testing.T) {
	regid := "3291379-2" //"1549609-1"
	address := "WPhr838tCoAMu22qvLg7JL6y6c8WESFchQ"

	token := ""
	c := NewClient(testNodeAPI, token, true)

	r, err := c.getContractAccountBalence(regid, address)
	fmt.Println(err)
	fmt.Println(r)
}
