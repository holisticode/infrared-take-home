package pkg

// file names
const JSON_STATE_TEST_FILE = "beacon_state_test_data.json"
const JSON_BLOCK_TEST_FILE = "block_root_test_data.json"

// leaf indexes
const STATE_ROOT_INDEX = 3
const RANDAO_MIXES_INDEX = 14

// APIs
const BEACON_CONTRACT_ADDRESS = "0x000F3df6D732807Ef1319fB7B8bB8522d0Beac02"
const VALIDATOR_INDEX_API = "http://localhost:9596/eth/v1/beacon/states/head/validators/%d"
const VALIDATORS_ALL_API = "http://localhost:9596/eth/v1/beacon/states/head/validators/"
const BEACON_BLOCK_API = "http://localhost:9596/eth/v2/beacon/blocks/head"
const RANDAO_API = "http://localhost:9596/eth/v1/beacon/states/head/randao"
const BEACON_STATE_API = "http://localhost:9596/eth/v2/debug/beacon/states/head"
const BEACON_STATE_ROOT_API = "http://localhost:9596/eth/v1/beacon/states/head/root"
