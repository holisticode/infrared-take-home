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

type Fork struct {
	PreviousVersion string
	CurrentVersion  string
	Epoch           string
}

type BeaconBlockHeader struct {
	Slot          string
	ProposerIndex string
	ParentRoot    string
	StateRoot     string
	BodyRoot      string
}
type Eth1Data struct {
	DepositRoot  string `json:"deposit_root"`
	DepositCount string `json:"deposit_count"`
	BlockHash    string `json:"block_hash"`
}

type Checkpoint struct {
	Epoch string
	Root  string
}
type AttestationData struct {
	Slot            string
	Index           string
	BeaconBlockRoot string `json:"beacon_block_root"`
	Source          Checkpoint
	Target          Checkpoint
}

type PendingAttestation struct {
	AggregationBits string
	Data            AttestationData
	InclusionDelay  string
	ProposerIndex   string
}

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

type BeaconState struct {
	GenesisTime           string
	GenesisValidatorsRoot Root
	Slot                  Slot
	Fork                  Fork
	//History
	LatestBlockHeader BeaconBlockHeader
	BlockRoots        []Root
	StateRoots        []Root
	HistoricalRoots   []Root
	//Eth1
	Eth1Data         Eth1Data
	Eth1DataVotes    []Eth1Data
	Eth1DepositIndex string
	//Registry
	Validators []ValidatorDetails
	Balances   []Gwei
	//Randomness
	RandaoMixes []Bytes32
	//Slashings
	Slashings []Gwei
	//Attestations
	PreviousEpochAttestations []PendingAttestation
	CurrentEpochAttestations  []PendingAttestation
	//Finality
	//JustificationBits           []Bitvector
	JustificationBits           any
	PreviousJustifiedCheckpoint Checkpoint
	CurrentJustifiedCheckpoint  Checkpoint
	FinalizedCheckpoint         Checkpoint
}
