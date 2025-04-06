package pkg

import (
	"crypto/sha256"
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

// Hash function to hash inputs
func hash(data []byte) []byte {
	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}
