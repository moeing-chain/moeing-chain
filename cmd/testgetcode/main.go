package main

import (
	"bytes"
	"os"
	"strconv"
	"time"

	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/smartbch/smartbch/internal/bigutils"
	"github.com/smartbch/smartbch/internal/testutils"
)

func main() {
	runCount := 2000
	if len(os.Args) > 1 {
		n, err := strconv.ParseUint(os.Args[1], 10, 32)
		if err == nil {
			runCount = int(n)
		}
	}

	// see testdata/counter/contracts/Counter.sol
	creationBytecode := testutils.HexToBytes(`
608060405234801561001057600080fd5b5060cc8061001f6000396000f3fe60
80604052348015600f57600080fd5b506004361060325760003560e01c806361
bc221a1460375780636299a6ef146053575b600080fd5b603d607e565b604051
8082815260200191505060405180910390f35b607c6004803603602081101560
6757600080fd5b81019080803590602001909291905050506084565b005b6000
5481565b8060008082825401925050819055505056fea2646970667358221220
37865cfcfd438966956583c78d31220c05c0f1ebfd116aced883214fcb1096c6
64736f6c634300060c0033
`)
	deployedBytecode := testutils.HexToBytes(`
6080604052348015600f57600080fd5b506004361060325760003560e01c8063
61bc221a1460375780636299a6ef146053575b600080fd5b603d607e565b6040
518082815260200191505060405180910390f35b607c60048036036020811015
606757600080fd5b81019080803590602001909291905050506084565b005b60
005481565b8060008082825401925050819055505056fea26469706673582212
2037865cfcfd438966956583c78d31220c05c0f1ebfd116aced883214fcb1096
c664736f6c634300060c0033
`)

	startTime := time.Now()
	initBal := bigutils.NewU256(testutils.DefaultInitBalance)
	valPubKey := ed25519.GenPrivKey().PubKey()
	key, _ := testutils.GenKeyAndAddr()
	var stateRoot []byte

	println("total count:", runCount)
	for i := 0; i < runCount; i++ {
		_app := testutils.CreateTestApp0(startTime, initBal, valPubKey, key)
		//defer _app.Destroy()

		println("run #", i)
		tx, _, contractAddr := _app.DeployContractInBlock(key, creationBytecode)
		_app.EnsureTxSuccess(tx.Hash())
		if i == 0 {
			stateRoot = _app.StateRoot
		} else {
			if !bytes.Equal(stateRoot, _app.StateRoot) {
				println("stateRoot not equal!")
			}
		}

		if !bytes.Equal(deployedBytecode, _app.GetCode(contractAddr)) {
			panic("deployedBytecode not equal")
		}

		_app.Destroy()
	}
}
