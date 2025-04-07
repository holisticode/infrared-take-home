package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/cbergoon/merkletree"
	"github.com/ethereum/go-ethereum/ethclient"
)

/**
This file is now OBSOLETE.

Also belongs to the initial attempts of proving validator information
*/

type RootHash [32]byte

func AssembleBeaconRootProofData(
	clint *ethclient.Client,
	validatorIndex uint64,
) error {
	// /eth/v2/beacon/blocks/number
	resp, err := http.Get(fmt.Sprintf(VALIDATOR_INDEX_API, validatorIndex))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var data ValidatorData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	valRoot, err := hashValidatorRoot(&data.Data.Validator)
	if err != nil {
		return err
	}
	fmt.Printf("%x\n", valRoot)

	/*
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		now := time.Now().Add(time.Duration(-20) * time.Second).Unix()
		timestamp := make([]byte, 32)
		binary.BigEndian.PutUint64(timestamp, uint64(now))

			ctrAddr := common.HexToAddress(BEACON_CONTRACT_ADDRESS)
			beaconContractCallMsg := eth.CallMsg{
				To:   &ctrAddr,
				Data: timestamp,
			}

			beaconContractResult, err := clint.CallContract(ctx, beaconContractCallMsg, nil)
			if err != nil {
				return err
			}
			fmt.Println(len(beaconContractResult))
			fmt.Printf("beacon contract result: %x\n", beaconContractResult)
	*/

	return nil

}

func hashValidatorRoot(valData *ValidatorDetails) (RootHash, error) {
	var list []merkletree.Content
	valDataStruct := reflect.ValueOf(*valData)

	for i := 0; i < valDataStruct.NumField(); i++ {
		field := valDataStruct.Field(i).Interface()
		if reflect.TypeOf(field) == reflect.TypeOf("") {
			list = append(list, ValidatorStringContent{Field: field.(string)})
		} else if reflect.TypeOf(field) == reflect.TypeOf(true) {
			list = append(list, ValidatorBoolContent{Field: field.(bool)})
		} else {
			return RootHash{0}, errors.New("unexpected validator field type, we only have string and bool")
		}
	}

	t, err := merkletree.NewTree(list)
	if err != nil {
		return RootHash{0}, err
	}

	//Get the Merkle Root of the tree
	mr := t.MerkleRoot()
	return RootHash(mr), nil
}
