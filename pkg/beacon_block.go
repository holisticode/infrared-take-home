package pkg

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"

	mt "github.com/txaty/go-merkletree"
)

// BeaconBlockData is used for parsing from JSON
type BeaconBlockData struct {
	Data BlockPreamble
}

// BlockPreamble was needed because the example API call I extracted from
// https://ethereum.github.io/beacon-APIs/#/Beacon/getBlockHeader
// Is showing this structure
type BlockPreamble struct {
	Root      string
	Canonical bool
	Header    BlockMessageHeader
}

// BlockMessageHeader finally is the last wrapper for the actual Block data
type BlockMessageHeader struct {
	Message BeaconBlock
}

// BeaconBlock encapsulates data for the beacon blockTree
// The data types are not according to the SSZ spec because I am
// loading it via JSON, which returns strings
type BeaconBlock struct {
	Slot          string
	ProposerIndex string `json:"proposer_index"`
	ParentRoot    string `json:"parent_root"`
	StateRoot     string `json:"state_root"`
	BodyRoot      string `json:"body_root"`
}

// CreateBeaconBlockTreeBlocks creates leaves as `[]mt.DataBlock` for each leaf of the
// Block data type, and creates a MerkleTree from it
func (bb *BeaconBlock) CreateBeaconBlockTreeBlocks() (*mt.MerkleTree, []mt.DataBlock, error) {
	// represents the leaves
	bt := []mt.DataBlock{}

	// we know here this is an uint64, so let's do the right conversion
	slotUint, err := strconv.ParseUint(bb.Slot, 10, 64)
	if err != nil {
		return nil, nil, fmt.Errorf("failed parsing Slot value")
	}
	// same
	piUint, err := strconv.ParseUint(bb.ProposerIndex, 10, 64)
	if err != nil {
		return nil, nil, fmt.Errorf("failed parsing proposerIndex value")
	}
	// we convert these now to bytes before hashing
	var slotBytes = make([]byte, 8)
	binary.LittleEndian.PutUint64(slotBytes, slotUint)
	var pidxBytes = make([]byte, 8)
	binary.LittleEndian.PutUint64(pidxBytes, piUint)
	// create a leaf for each member
	bt = append(bt, &ByteArray{data: slotBytes})
	bt = append(bt, &ByteArray{data: pidxBytes})

	// to hash actual hashes, we decode these to byte arrays, otherwise we'd be
	// hashing their string representation
	pRoot, err := hex.DecodeString(bb.ParentRoot[2:])
	if err != nil {
		return nil, nil, fmt.Errorf("failed decoding parent root: %v", err)
	}
	bt = append(bt, &ByteArray{data: pRoot})
	sRoot, err := hex.DecodeString(bb.StateRoot[2:])
	if err != nil {
		return nil, nil, fmt.Errorf("failed decoding state root: %v", err)
	}
	bt = append(bt, &ByteArray{data: sRoot})
	bRoot, err := hex.DecodeString(bb.BodyRoot[2:])
	if err != nil {
		return nil, nil, fmt.Errorf("failed decoding body root: %v", err)
	}
	bt = append(bt, &ByteArray{data: bRoot})

	cf := getMerkleTreeconfig()
	// now we can create the tree
	blockTree, err := mt.New(cf, bt)
	if err != nil {
		return nil, nil, fmt.Errorf("failed creating block tree: %v", err)
	}
	return blockTree, bt, nil
}
