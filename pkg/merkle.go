package pkg

/**
	This file contains manual MerkleTree construction code I copied from elsewhere.
	I was mainly interested in the VerifyProof function, so that I would be able to
	verify the merkle proof without using the github.com/txaty/go-merkletree
	which uses its own custom interface with `DataBlock`s.

	Basically a fall back solution
**/

import (
	"encoding/hex"
)

// Structure for the Merkle Tree
type MerkleTree struct {
	Leaves [][]byte
	Root   []byte
}

// Create a new Merkle Tree given a list of leaves
func NewMerkleTree(leaves [][]byte) *MerkleTree {
	if len(leaves) == 0 {
		return nil
	}

	tree := &MerkleTree{Leaves: leaves}
	tree.Root = tree.buildTree(leaves)
	return tree
}

// Build the tree and return the root hash
func (mt *MerkleTree) buildTree(leaves [][]byte) []byte {
	hashes := leaves

	for len(hashes) > 1 {
		level := [][]byte{}
		for i := 0; i < len(hashes); i += 2 {
			if i+1 == len(hashes) {
				level = append(level, hashes[i])
			} else {
				level = append(level, hash(append(hashes[i], hashes[i+1]...)))
			}
		}
		hashes = level
	}
	return hashes[0]
}

// Generate Merkle Proof for an element
func (mt *MerkleTree) GenerateProof(index int) [][]byte {
	if index < 0 || index >= len(mt.Leaves) {
		return nil
	}

	proof := [][]byte{}
	hashes := mt.Leaves[:]
	for len(hashes) > 1 {
		if index%2 == 0 {
			// Index is even, add the next sibling hash to the proof
			if index+1 < len(hashes) {
				proof = append(proof, hashes[index+1])
			}
		} else {
			// Index is odd, add the previous sibling hash to the proof
			proof = append(proof, hashes[index-1])
		}
		hashes = hashes[:(len(hashes)+1)/2]
		index /= 2
	}
	return proof
}

// needed to introduce this to avoid name clash with the
// hash variable in VerifyProof
var calchash = hash

func VerifyProof(leaf []byte, proof [][]byte, root []byte, leafIndex int) bool {
	hash := leaf
	for _, sibling := range proof {
		if leafIndex%2 == 0 {
			hash = calchash(append(hash, sibling...))
		} else {
			hash = calchash(append(sibling, hash...))
		}
		leafIndex /= 2
	}
	return hex.EncodeToString(hash) == hex.EncodeToString(root)
}
