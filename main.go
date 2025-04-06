package main

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/holisticode/infrared-take-home/pkg"
)

func main() {
	log.Println("connecting to client...")
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatalf("failed to connect to ethclient: %v", err)
	}
	log.Println("generating proof...")
	//eth2endpoint := "http://localhost:9596/eth/v2/beacon/blocks/2"
	//	eth2endpoint := "http://localhost:9596/eth/v2/beacon/blocks/2"
	if err := pkg.AssembleRandaoProof(client, 7); err != nil {
		log.Printf("failed generating proof: %v", err)
	}
}
