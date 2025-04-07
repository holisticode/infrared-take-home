package main

import (
	"log"

	//"github.com/ethereum/go-ethereum/ethclient"

	"github.com/holisticode/infrared-take-home/pkg"
)

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
	log.Println("generating proof...")
	if err := pkg.AssembleRandaoProof(7, pkg.FILE); err != nil {
		log.Printf("failed generating proof: %v", err)
	}
}
