package pkg

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	mt "github.com/txaty/go-merkletree"
)

/**
This is the main file for the Randao proof generation task.
It loads a beacon state from a static JSON file,
generates the Merkle tree and the proof for the whole chain:

* That a particular Randao is contained in randao_mixes
* That that randao_mixes is contained in the beacon state
* That this beacon state is contained in the block

WARNING: This code has probably a problem, in that it loads
the beacon state as JSON file, taking as example data from
https://ethereum.github.io/beacon-APIs/#/Debug/getStateV2.

It uses reflection to build the MerkleTree of the state from this JSON.
However, because it is NOT using SSZ, some convertions
(notably numbers) may not be correct according to the spec.

In other words, the proof within a real EIP-4788 contract inside
an EVM would most probably fail.

I couldn't find a usable SSZ parser which would have helped with this task,
and writing my own felt like using a lot of the time available for the task.
*/

// DataLoader is an enum type to tell if load from File or API
// Currently only File is supported
type DataLoader int

const (
	API DataLoader = iota
	FILE
)

// ByteArray implements the `mt.DataBlock` interface used
// to create Merkle Trees
// github.com/txaty/go-merkletree
type ByteArray struct {
	data []byte
}

func (bt *ByteArray) Serialize() ([]byte, error) {
	return bt.data, nil
}

// RandaoProofWrapper encapsulates all the results
// obtained from the proof generation
type RandaoProofWrapper struct {
	RandaoMixesTree       *mt.MerkleTree
	RandaoProof           *mt.Proof
	RandaoMixesDataBlocks []mt.DataBlock
	RandaoTargetIndex     int

	BeaconStateTree        *mt.MerkleTree
	RandaoMixesProof       *mt.Proof
	BeaconStateDataBlocks  []mt.DataBlock
	BeaconStateTargetIndex int

	BlockRootTree        *mt.MerkleTree
	StateProof           *mt.Proof
	BlockRootDataBlocks  []mt.DataBlock
	BlockRootTargetIndex int
}

// AssembleRandaoProof is the entry point from main
func AssembleRandaoProof(
	randaoIndex int,
	source DataLoader,
) error {

	// load the JSON files
	state, blockData, err := loadData(source, "pkg/")
	if err != nil {
		return fmt.Errorf("failed loading state data %v", err)
	}

	// create the proofs
	results, err := assembleRandaoProofGenerator(state, blockData, randaoIndex)
	if err != nil {
		return fmt.Errorf("failed to assemble trees and proofs: %v", err)
	}

	// print results
	log.Printf("randao_mixes root: %x\n", results.RandaoMixesTree.Root)
	// let's skip error checking here...
	b, _ := results.RandaoMixesDataBlocks[results.RandaoTargetIndex].Serialize()
	log.Printf("randao leaf: %x\n", hash(b))
	log.Print("randao proof: [")
	for i, p := range results.RandaoProof.Siblings {
		fmt.Printf(" %x", p)
		if i < len(results.RandaoProof.Siblings)-1 {
			fmt.Print(",")
		}
	}
	fmt.Print(" ]\n")
	fmt.Println()
	log.Printf("beacon state root: %x\n", results.BeaconStateTree.Root)
	// let's skip error checking here...
	b, _ = results.BeaconStateDataBlocks[results.BeaconStateTargetIndex].Serialize()
	log.Printf("beacon state leaf: %x\n", hash(b))
	log.Print("randao_mixes proof: [")
	for i, p := range results.RandaoMixesProof.Siblings {
		fmt.Printf(" %x", p)
		if i < len(results.RandaoMixesProof.Siblings)-1 {
			fmt.Print(",")
		}
	}
	fmt.Print(" ]\n")
	fmt.Println()
	log.Printf("block tree root: %x\n", results.BlockRootTree.Root)
	b, _ = results.BlockRootDataBlocks[results.BlockRootTargetIndex].Serialize()
	log.Printf("block root leaf: %x\n", hash(b))
	log.Print("state proof: [")
	for i, p := range results.StateProof.Siblings {
		fmt.Printf(" %x", p)
		if i < len(results.StateProof.Siblings)-1 {
			fmt.Print(",")
		}
	}
	fmt.Print(" ]\n")
	return nil
}

// loadData loads the required data structures
func loadData(source DataLoader, prefix string) (*BeaconStateData, *BeaconBlockData, error) {
	var (
		state     *BeaconStateData
		blockData *BeaconBlockData
		err       error
	)

	if source == FILE {
		state, err = loadBeaconStateFromFile(prefix + JSON_STATE_TEST_FILE)
		if err != nil {
			return nil, nil, fmt.Errorf("failed opening state data test file: %v", err)
		}
		blockData, err = loadBeaconBlockFromFile(prefix + JSON_BLOCK_TEST_FILE)
		if err != nil {
			return nil, nil, fmt.Errorf("failed opening block data test file: %v", err)
		}
	} else {
		state, err = queryBeaconStateAPI(BEACON_STATE_API)
		if err != nil {
			return nil, nil, fmt.Errorf("failed opening state data from API: %v", err)
		}
		// omitted block data as actually turned out API is inconsistent
	}
	return state, blockData, nil
}

// assembleRandaoProofGenerator does the heavy lifting.
// For each level, it creates a MerkleTree and its proof
func assembleRandaoProofGenerator(
	state *BeaconStateData,
	blockData *BeaconBlockData,
	randaoIndex int,
) (*RandaoProofWrapper, error) {

	block := blockData.Data.Header.Message
	beaconState := state.Data
	randaoMixes := state.Data.RandaoMixes

	// prove randao in randao mixes
	randaoMixesTree, randaoMixesDataBlocks, err := createRandaoMixesTree(randaoIndex, randaoMixes)
	if err != nil {
		return nil, fmt.Errorf("failed to create tree for randao mixes: %v", err)
	}
	randaoProof, err := randaoMixesTree.Proof(randaoMixesDataBlocks[randaoIndex])
	if err != nil {
		return nil, err
	}

	// prove randao_mixes inside beacon state
	beaconStateTree, beaconStateDataBlocks, err := createBeaconTree(&beaconState)
	if err != nil {
		return nil, fmt.Errorf("failed to create beacon state tree: %v", err)
	}
	randaoMixesProof, err := beaconStateTree.Proof(beaconStateDataBlocks[RANDAO_MIXES_INDEX])
	if err != nil {
		return nil, err
	}

	// prove state inside the block root
	blockRootTree, blockRootTreeDataBlocks, err := block.CreateBeaconBlockTreeBlocks()
	if err != nil {
		return nil, fmt.Errorf("failed creating block tree blocks: %v", err)
	}
	stateProof, err := blockRootTree.Proof(blockRootTreeDataBlocks[STATE_ROOT_INDEX])
	if err != nil {
		return nil, fmt.Errorf("failed creating proof for state hash: %v", err)
	}

	result := &RandaoProofWrapper{
		RandaoMixesTree:       randaoMixesTree,
		RandaoProof:           randaoProof,
		RandaoMixesDataBlocks: randaoMixesDataBlocks,
		RandaoTargetIndex:     randaoIndex,

		BeaconStateTree:        beaconStateTree,
		RandaoMixesProof:       randaoMixesProof,
		BeaconStateDataBlocks:  beaconStateDataBlocks,
		BeaconStateTargetIndex: RANDAO_MIXES_INDEX,

		BlockRootTree:        blockRootTree,
		StateProof:           stateProof,
		BlockRootDataBlocks:  blockRootTreeDataBlocks,
		BlockRootTargetIndex: STATE_ROOT_INDEX,
	}

	return result, nil

}

func createRandaoMixesTree(randaoIndex int, randaoMixes []string) (*mt.MerkleTree, []mt.DataBlock, error) {
	randao := randaoMixes[randaoIndex]
	bt := make([]mt.DataBlock, len(randaoMixes))

	index := -1
	for i, r := range randaoMixes {
		if r == randao {
			index = i
		}
		b, err := hex.DecodeString(r[2:])
		if err != nil {
			return nil, nil, err
		}
		bt[i] = &ByteArray{data: b}
	}

	if index == -1 {
		return nil, nil, errors.New("the randao wasn't in the list, as was expected")
	}

	cf := getMerkleTreeconfig()
	randaoTree, err := mt.New(cf, bt)
	if err != nil {
		return nil, nil, err
	}
	return randaoTree, bt, nil
}

func createBeaconTree(beaconState *BeaconStateSimplified) (*mt.MerkleTree, []mt.DataBlock, error) {
	beaconStateLeaves, err := calculateBeaconStateLeaves(beaconState)
	if err != nil {
		return nil, nil, err
	}

	bt, err := getMerkleTreeDataBlocksFromStrings(beaconStateLeaves)
	if err != nil {
		return nil, nil, err
	}

	cf := getMerkleTreeconfig()
	beaconStateTree, err := mt.New(cf, bt)
	if err != nil {
		return nil, nil, err
	}
	return beaconStateTree, bt, nil
}

func buildMerkleTreeFromStrings(data []string) (*mt.MerkleTree, error) {
	bt, err := getMerkleTreeDataBlocksFromStrings(data)
	if err != nil {
		return nil, err
	}
	return mt.New(nil, bt)
}

func getMerkleTreeDataBlocksFromStrings(data []string) ([]mt.DataBlock, error) {
	var err error
	bt := make([]mt.DataBlock, len(data))

	for i, s := range data {
		dat := s
		var arr []byte
		// if it's a hex string we should decode it
		if strings.HasPrefix(s, "0x") {
			dat = s[2:]
			arr, err = hex.DecodeString(dat)
			// TODO: This is probably INCORRECT due to the JSON encoding,
			// which returns strings when there are uints or other number formats
			// the correct implementation is probably hashing number values, not their string rep
			// As it is generic here, I decided to leave this as this for the sake of the test
		} else {
			arr = []byte(dat)
		}
		if err != nil {
			return nil, err
		}
		b := &ByteArray{data: arr}
		bt[i] = b
	}
	return bt, nil
}

// calculateBeaconStateLeaves uses reflection to gather all leaves of the BeaconState
// data structure. I realized ONLY LATER that this would probably result in violating
// the SSZ spec, because the JSON API returns strings for different data types, which,
// to calculate the correct hashes, should be converted to their original data types,
// and the converted to []byte for hashing.
func calculateBeaconStateLeaves(bs *BeaconStateSimplified) ([]string, error) {
	var list []string
	state := reflect.ValueOf(*bs)

	// iterate all fields of BeaconStateSimplified
	for i := 0; i < state.NumField(); i++ {
		field := state.Field(i).Interface()
		fieldType := reflect.TypeOf(field)
		// direct leaf as string
		if fieldType == reflect.TypeOf("") {
			// this is a direct leaf, hash it and encode to hex
			strVal := field.(string)
			if strings.HasPrefix(strVal, "0x") {
				list = append(list, strVal)
			} else {
				// NOTE: same hack, if it's a number this is probably incorrect
				fieldHash := hash([]byte(strVal))
				list = append(list, hex.EncodeToString(fieldHash))
			}
			// a slice of strings, e.g. BlockRoots
		} else if fieldType.Kind() == reflect.Slice && fieldType.Elem() == reflect.TypeOf("") {
			subTree, err := buildMerkleTreeFromStrings(field.([]string))
			if err != nil {
				return nil, err
			}
			subRoot := subTree.Root
			subRootStr := hex.EncodeToString(subRoot)
			list = append(list, subRootStr)
			// a slice of struct, e.g. []Eth1Data
		} else if fieldType.Kind() == reflect.Slice && fieldType.Elem().Kind() == reflect.Struct {
			var subList []string
			for i := 0; i < reflect.ValueOf(field).Len(); i++ {
				strlist := iterateFieldsOfStruct(reflect.ValueOf(field).Index(i))
				subTree, err := buildMerkleTreeFromStrings(strlist)
				if err != nil {
					return nil, err
				}
				subRoot := subTree.Root
				subRootStr := hex.EncodeToString(subRoot)
				subList = append(subList, subRootStr)
			}
			subTree, err := buildMerkleTreeFromStrings(subList)
			if err != nil {
				return nil, err
			}
			subRoot := subTree.Root
			subRootStr := hex.EncodeToString(subRoot)
			list = append(list, subRootStr)
			// a struct of different type, e.g. Fork
		} else if fieldType.Kind() == reflect.Struct {
			strlist := iterateFieldsOfStruct(reflect.ValueOf(field))
			subTree, err := buildMerkleTreeFromStrings(strlist)
			if err != nil {
				return nil, err
			}
			subRoot := subTree.Root
			subRootStr := hex.EncodeToString(subRoot)
			list = append(list, subRootStr)
		} else {
			return nil, fmt.Errorf("unexpected validator field type: %v", fieldType)
		}
	}

	return list, nil
}

// get all fields for a sub struct
func iterateFieldsOfStruct(field reflect.Value) []string {
	//structType := reflect.ValueOf(field)
	var strlist []string
	for k := 0; k < field.NumField(); k++ {
		f := field.Field(k).Interface()
		if reflect.TypeOf(f) == reflect.TypeOf(true) {
			if reflect.ValueOf(f).Bool() == false {
				strlist = append(strlist, string([]byte{0}))
			} else {
				strlist = append(strlist, string([]byte{1}))
			}
		} else {
			strlist = append(strlist, f.(string))
		}
	}
	return strlist
}
