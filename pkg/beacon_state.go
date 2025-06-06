package pkg

type Root string
type Slot string
type Epoch string
type Gwei string
type Bitlist string
type ValidatorIndex string
type Version string
type Hash32 [32]byte
type Bytes32 []string

/**
This file contains mappings to the data structures for beacon state
as defined in https://github.com/ethereum/consensus-specs/blob/dev/specs/phase0/beacon-chain.md

NOTE: As I was using JSON, many data types do NOT match the actual spec, although
the same spec clearly defines which data types are mapped to string for JSON
(which this implementation uses)
*/

// BeaconStateData is the wrapper needed to fetch
// JSON data for the beacon state
type BeaconStateData struct {
	Data BeaconStateSimplified
}

// RandaoData is the wrapper needed to fetch
// JSON data for the randao value
// In this implementation not needed
type RandaoData struct {
	Data Randao
}

// Randao is the container for the randao value
// when unpacked from API
type Randao struct {
	Randao string
}

// Fork as per spec
type Fork struct {
	PreviousVersion string
	CurrentVersion  string
	Epoch           string
}

// BeaconBlockHeader
type BeaconBlockHeader struct {
	Slot          string
	ProposerIndex string
	ParentRoot    string
	StateRoot     string
	BodyRoot      string
}

// Eth1Data
type Eth1Data struct {
	DepositRoot  string `json:"deposit_root"`
	DepositCount string `json:"deposit_count"`
	BlockHash    string `json:"block_hash"`
}

// ValidatorDetails holds validator data
type ValidatorDetails struct {
	Pubkey                     string
	WithdrawalCredentials      string
	EffectiveBalance           string
	Slashed                    bool
	ActivationEligibilityEpoch string
	ActivationEpoch            string
	ExitEpoch                  string
	WithdrawableEpoch          string
}

// Checkpoint
type Checkpoint struct {
	Epoch string
	Root  string
}

// AttestationData
type AttestationData struct {
	Slot            string
	Index           string
	BeaconBlockRoot string `json:"beacon_block_root"`
	Source          Checkpoint
	Target          Checkpoint
}

// PendingAttestation
type PendingAttestation struct {
	AggregationBits string
	Data            AttestationData
	InclusionDelay  string
	ProposerIndex   string
}

// BeaconStateSimplified is actually the container for
// BeaconState. Called it Simplified because it uses string
// data types for the JSON fields (Understood only later
// that the spec explicitly specifies this for JSON endpoints)
type BeaconStateSimplified struct {
	GenesisTime           string `json:"genesis_time"`
	GenesisValidatorsRoot string `json:"genesis_validators_root"`
	Slot                  string
	Fork                  Fork
	//History
	LatestBlockHeader BeaconBlockHeader `json:"latest_block_header"`
	BlockRoots        []string          `json:"block_roots"`
	StateRoots        []string          `json:"state_roots"`
	HistoricalRoots   []string          `json:"historical_roots"`

	//Eth1
	Eth1Data         Eth1Data   `json:"eth1_data"`
	Eth1DataVotes    []Eth1Data `json:"eth1_data_votes"`
	Eth1DepositIndex string     `json:"eth1_deposit_index"`
	//Registry
	Validators []ValidatorDetails
	Balances   []string
	//Randomness
	RandaoMixes []string `json:"randao_mixes"`
	//Slashings
	Slashings []string
	//Attestations
	// NOTE: For some reason the examples from my beacon node were returning
	// previous_epoch_participations and current_epoch_participations
	// So this item seems to not correspond to the actual spec!
	//PreviousEpochAttestations []PendingAttestation `json:"previous_epoch_attestations"`
	//CurrentEpochAttestations  []PendingAttestation `json:"current_epoch_attestations"`
	PreviousEpochAttestations []string `json:"previous_epoch_attestations"`
	CurrentEpochAttestations  []string `json:"current_epoch_attestations"`
	//Finality
	JustificationBits           string     `json:"justification_bits"`
	PreviousJustifiedCheckpoint Checkpoint `json:"previous_justified_checkpoint"`
	CurrentJustifiedCheckpoint  Checkpoint `json:"current_justified_checkpoint"`
	FinalizedCheckpoint         Checkpoint `json:"finalized_checkpoint"`
}
