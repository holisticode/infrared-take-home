package pkg

// Verify is just a wrapper for the existing verification func
func Verify(leaf []byte, proof [][]byte, root []byte, index int) bool {
	return VerifyProof(leaf, proof, root, index)
}
