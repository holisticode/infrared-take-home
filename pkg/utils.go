package pkg

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	mt "github.com/txaty/go-merkletree"
)

// hash is the default hash function used basically everywhere
// using SHA256
func hash(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}

// loadBeaconStateFromFile loads the beacon state from JSON file
func loadBeaconStateFromFile(filename string) (*BeaconStateData, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open json file: %v", err)
	}
	log.Println("Successfully opened beacon state json file")
	defer jsonFile.Close()

	byteJSON, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	var state BeaconStateData
	if err := json.Unmarshal(byteJSON, &state); err != nil {
		return nil, fmt.Errorf("failed decoding json from file: %v", err)
	}

	return &state, nil
}

// loadBeaconBlockFromFile loads the JSON file for the beacon block
func loadBeaconBlockFromFile(filename string) (*BeaconBlockData, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open block root json file: %v", err)
	}
	log.Println("Successfully opened json file")
	defer jsonFile.Close()

	byteJSON, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	var block BeaconBlockData
	if err := json.Unmarshal(byteJSON, &block); err != nil {
		return nil, fmt.Errorf("failed decoding json from file: %v", err)
	}

	return &block, nil
}

// Initially I was using an actual lodestar node to connect to Sepolia
// However, that node was not returning the state conforming to the spec at
// // https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md#
// Specifically the previous_epoch_attestations were called previous_epoch_participations,
// same for current_epoch_attestations
//
// So I decided to use static files for testing
func queryBeaconStateAPI(apiURL string) (*BeaconStateData, error) {
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var state BeaconStateData
	if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
		return nil, fmt.Errorf("failed decoding json from API: %v", err)
	}

	return &state, nil
}

// queryRandaoAPI was used during the initial phase to get the actual
// randao, but it is not used at the moment.
func queryRandaoAPI(apiURL string) (*RandaoData, error) {
	resp, err := http.Get(RANDAO_API)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var randaoData RandaoData
	if err := json.NewDecoder(resp.Body).Decode(&randaoData); err != nil {
		return nil, err
	}

	//randao := randaoData.Data.Randao
	return &randaoData, nil
}

// getMerkleTreeconfig returns the config for the MerkleTree library
func getMerkleTreeconfig() *mt.Config {
	return &mt.Config{
		// Customizable hash function used for tree generation.
		HashFunc: defaultHashFunc,
		// Number of goroutines run in parallel.
		// If RunInParallel is true and NumRoutine is set to 0, use number of CPU as the number of goroutines.
		NumRoutines: 0,
		// Mode of the Merkle Tree generation.
		Mode: mt.ModeProofGenAndTreeBuild,
		// If RunInParallel is true, the generation runs in parallel, otherwise runs without parallelization.
		// This increase the performance for the calculation of large number of data blocks, e.g. over 10,000 blocks.
		RunInParallel: true,
		// SortSiblingPairs is the parameter for OpenZeppelin compatibility.
		// If set to `true`, the hashing sibling pairs are sorted.
		SortSiblingPairs: false,
		// If true, the leaf nodes are NOT hashed before being added to the Merkle Tree.
		DisableLeafHashing: false,
	}
}
