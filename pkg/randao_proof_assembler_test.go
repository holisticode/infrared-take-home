package pkg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestRandaoProofGeneration tests that a merkle tree can be built from
// all randao_mixes, proves that a given randao is contained in the randao_mixes,
// and that the randao_mixes is contained in the whole beacon state,
// and that this beacon state is contained in the beacon block.

// It uses two different libraries for verification.
func TestRandaoProofGeneration(t *testing.T) {
	state, blockData, err := loadData(FILE, "")
	require.Empty(t, err)

	results, err := assembleRandaoProofGenerator(state, blockData, 7)
	require.Empty(t, err)

	tree := results.RandaoMixesTree
	blocks := results.RandaoMixesDataBlocks
	proof := results.RandaoProof
	target := results.RandaoTargetIndex
	ok, err := tree.Verify(blocks[target], proof)
	require.True(t, ok)
	require.Empty(t, err)
	// alternative check with homemade merkle tree...
	leaf, err := blocks[target].Serialize()
	require.Empty(t, err)
	leafhash := hash(leaf)
	ok = VerifyProof(leafhash, proof.Siblings, tree.Root, target)
	require.True(t, ok)

	tree = results.BeaconStateTree
	blocks = results.BeaconStateDataBlocks
	proof = results.RandaoMixesProof
	target = results.BeaconStateTargetIndex
	ok, err = tree.Verify(blocks[target], proof)
	require.True(t, ok)
	require.Empty(t, err)
	// alternative check with homemade merkle tree...
	leaf, err = blocks[target].Serialize()
	require.Empty(t, err)
	leafhash = hash(leaf)
	ok = VerifyProof(leafhash, proof.Siblings, tree.Root, target)
	require.True(t, ok)

	tree = results.BlockRootTree
	blocks = results.BlockRootDataBlocks
	proof = results.StateProof
	target = results.BlockRootTargetIndex
	ok, err = tree.Verify(blocks[target], proof)
	require.True(t, ok)
	require.Empty(t, err)
	// alternative check with homemade merkle tree...
	leaf, err = blocks[target].Serialize()
	require.Empty(t, err)
	leafhash = hash(leaf)
	ok = VerifyProof(leafhash, proof.Siblings, tree.Root, target)
	require.True(t, ok)
}
