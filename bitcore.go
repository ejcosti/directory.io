package main

import (
	"fmt"
	"os"
	"bufio"
	"math/big"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
)

const ResultsPerPage = 500000

const KeyTemplate = `<span id="%s"><a href="/warning:understand-how-this-works!/%s">+</a> <span title="%s">%s </span> <a href="https://blockchain.info/address/%s">%34s</a> <a href="https://blockchain.info/address/%s">%34s</a></span>
`
const KeyTemplate2 = `%s, %s, %s, %s %s`

var (
	// Total bitcoins
	total = new(big.Int).SetBytes([]byte{
		0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFE,
		0xBA, 0xAE, 0xDC, 0xE6, 0xAF, 0x48, 0xA0, 0x3B, 0xBF, 0xD2, 0x5E, 0x8C, 0xD0, 0x36, 0x41, 0x40,
	})

	// One
	one = big.NewInt(1)
)

type Key struct {
	private      string
	number       string
	compressed   string
	uncompressed string
}

func compute(count *big.Int) (keys [ResultsPerPage]Key, length int) {
	var padded [32]byte

	var i int
	for i = 0; i < ResultsPerPage; i++ {
		// Increment our counter
		count.Add(count, one)

		// Check to make sure we're not out of range
		if count.Cmp(total) > 0 {
			break
		}

		// Copy count value's bytes to padded slice
		copy(padded[32-len(count.Bytes()):], count.Bytes())

		// Get private and public keys
		privKey, public := btcec.PrivKeyFromBytes(btcec.S256(), padded[:])

		// Get compressed and uncompressed addresses for public key
		caddr, _ := btcutil.NewAddressPubKey(public.SerializeCompressed(), &chaincfg.MainNetParams)
		uaddr, _ := btcutil.NewAddressPubKey(public.SerializeUncompressed(), &chaincfg.MainNetParams)

		// Encode addresses
		wif, _ := btcutil.NewWIF(privKey, &chaincfg.MainNetParams, false)
		keys[i].private = wif.String()
		keys[i].number = count.String()
		keys[i].compressed = caddr.EncodeAddress()
		keys[i].uncompressed = uaddr.EncodeAddress()
	}
	return keys, i
}


func main() {
	start_val := new(big.Int)

	start_val, sok := start_val.SetString("74968309822361279513706563288209140390982877136853926828208222868511747", 10)
  if !sok {
      	fmt.Println("SetString: error")
        return
  }

  for i := 1; i < 1000; i++ {
  	my_val := new(big.Int)
  	temp_val := ResultsPerPage * i
  	my_val, mok := my_val.SetString(fmt.Sprint(temp_val), 10)
  	if !mok {
      	fmt.Println("SetString: error")
        return
  	}

  	val := new(big.Int)
  	val.Add(my_val, start_val)
  	fmt.Println(val)

  	filename := fmt.Sprintf("%d.txt", val)

		f, _ := os.Create(filename)
  	w := bufio.NewWriter(f)

  	keys, length := compute(val)

		for i := 0; i < length; i++ {
			key := keys[i]
			fmt.Fprintf(w, KeyTemplate2, key.number, key.private, key.uncompressed, key.compressed, "\r\n")
		}

		w.Flush()
  }

	


}
