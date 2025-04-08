package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	//"github.com/ethereum/go-ethereum/ethclient"

	"github.com/holisticode/infrared-take-home/pkg"
)

const DEFAULT_TARGET = 7

func main() {
	/*
		  // Initially planned to connect to some real node
			// setup (was thinking ganache or hardhat) to query nodes live
			// but I do not have experience with consensus client development nodes,
			// so finally I lef this unused

			log.Println("connecting to client...")
			client, err := ethclient.Dial("http://localhost:8545")
			if err != nil {
				log.Fatalf("failed to connect to ethclient: %v", err)
			}
	*/
	if len(os.Args) < 2 {
		generateProof(DEFAULT_TARGET)
		return
	}

	cmd := os.Args[1]
	if cmd == "generate" {
		var target int
		if len(os.Args) > 2 && os.Args[2] != "" {
			tget, err := strconv.ParseInt(os.Args[2], 10, 64)
			if err != nil {
				panic(err)
			}
			target = int(tget)
		} else {
			target = DEFAULT_TARGET
		}
		generateProof(target)
	} else if cmd == "verify" {
		// poor man's arg control here...
		if len(os.Args) != 6 {
			fmt.Println("usage: verify 'leaf' 'proof' 'root' 'target'")
			return
		}

		leaf := os.Args[2]
		proof := os.Args[3]
		root := os.Args[4]
		target := os.Args[5]
		// skipping sanity checks for this exercise
		err := verify(leaf, proof, root, target)
		if err != nil {
			log.Fatal("failed to process verification command")
		}
	} else {
		log.Fatal("command not understood")
	}
}

func generateProof(target int) {
	log.Println("generating proof...")
	if err := pkg.AssembleRandaoProof(target, pkg.FILE); err != nil {
		log.Printf("failed generating proof: %v", err)
	}
}

func verify(l, p, r, t string) error {
	leaf, err := hex.DecodeString(l)
	if err != nil {
		return err
	}
	proofs := strings.Split(p, ",")
	proof := make([][]byte, len(proofs))
	for i, p := range proofs {
		proof[i], err = hex.DecodeString(p)
		if err != nil {
			return err
		}
	}

	root, err := hex.DecodeString(r)
	if err != nil {
		return err
	}

	target, err := strconv.Atoi(t)
	if err != nil {
		return err
	}

	ok := pkg.Verify(leaf, proof, root, target)
	if ok {
		log.Println("proof validation successful")
	} else {
		log.Fatal("proof validation FAILED")
	}
	return nil
}
