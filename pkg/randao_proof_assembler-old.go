package pkg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/ethclient"
	ssz "github.com/prysmaticlabs/go-ssz"
	//"github.com/protolambda/zssz"
)

/**
	This file is now OBSOLETE.

	It contains first attempts navigating through the complexity
	of getting this done.

	Specifically, I tried several ways of using ssz parsers,
	but they all kinda failed on me...

	Which was the reason why I moved to use JSON
**/

func AssembleRandaoProofBig(
	clint *ethclient.Client,
) error {
	resp, err := http.Get(RANDAO_API)
	if err != nil {
		return err
	}

	var randao RandaoData
	if err := json.NewDecoder(resp.Body).Decode(&randao); err != nil {
		return err
	}
	//fmt.Println(randao.Data.Randao)
	resp.Body.Close()

	resp, err = http.Get(BEACON_STATE_API)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	/*
		databytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		statestr := string(databytes)
		fmt.Println(statestr)
	*/
	//state := statestr[9 : len(statestr)-1]
	//fmt.Println(state)

	var data BeaconStateData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return err
	}

	allRoots := make([][32]byte, 0)

	genesisRoot, err := ssz.HashTreeRoot(data.Data.GenesisTime)
	if err != nil {
		return err
	}
	allRoots = append(allRoots, genesisRoot)
	genesisValidatorRoot, err := ssz.HashTreeRoot(data.Data.GenesisValidatorsRoot)
	if err != nil {
		return err
	}
	allRoots = append(allRoots, genesisValidatorRoot)
	slotRoot, err := ssz.HashTreeRoot(data.Data.Slot)
	if err != nil {
		return err
	}
	allRoots = append(allRoots, slotRoot)
	forkRoot, err := ssz.HashTreeRoot(data.Data.Fork)
	if err != nil {
		return err
	}
	allRoots = append(allRoots, forkRoot)
	latestBlockHeaderRoot, err := ssz.HashTreeRoot(data.Data.LatestBlockHeader)
	if err != nil {
		return err
	}
	allRoots = append(allRoots, latestBlockHeaderRoot)
	blockRoots, err := ssz.HashTreeRoot(data.Data.BlockRoots)
	if err != nil {
		return err
	}
	allRoots = append(allRoots, blockRoots)

	stateRoots, err := ssz.HashTreeRoot(data.Data.StateRoots)
	if err != nil {
		return err
	}
	allRoots = append(allRoots, stateRoots)

	historicalRoots, err := ssz.HashTreeRoot(data.Data.HistoricalRoots)
	if err != nil {
		return err
	}
	allRoots = append(allRoots, historicalRoots)
	eth1DataRoots, err := ssz.HashTreeRoot(data.Data.Eth1Data)
	if err != nil {
		return err
	}
	allRoots = append(allRoots, eth1DataRoots)
	eth1DataVotesRoots, err := ssz.HashTreeRoot(data.Data.Eth1DataVotes)
	if err != nil {
		return err
	}
	allRoots = append(allRoots, eth1DataVotesRoots)
	eth1DepIndexRoots, err := ssz.HashTreeRoot(data.Data.Eth1DepositIndex)
	if err != nil {
		return err
	}
	allRoots = append(allRoots, eth1DepIndexRoots)
	validatorRoots, err := ssz.HashTreeRoot(data.Data.Validators)
	if err != nil {
		return err
	}
	allRoots = append(allRoots, validatorRoots)
	balancesRoots, err := ssz.HashTreeRoot(data.Data.Balances)
	if err != nil {
		return err
	}
	allRoots = append(allRoots, balancesRoots)
	randaoRoots, err := ssz.HashTreeRoot(data.Data.RandaoMixes)
	if err != nil {
		return err
	}
	allRoots = append(allRoots, randaoRoots)
	slashingRoots, err := ssz.HashTreeRoot(data.Data.Slashings)
	if err != nil {
		return err
	}
	allRoots = append(allRoots, slashingRoots)
	prevEpochAttRoots, err := ssz.HashTreeRoot(data.Data.PreviousEpochAttestations)
	if err != nil {
		return err
	}
	allRoots = append(allRoots, prevEpochAttRoots)
	currEpochAttRoots, err := ssz.HashTreeRoot(data.Data.CurrentEpochAttestations)
	if err != nil {
		return err
	}
	allRoots = append(allRoots, currEpochAttRoots)
	justBitsRoots, err := ssz.HashTreeRoot(data.Data.JustificationBits)
	if err != nil {
		return err
	}
	allRoots = append(allRoots, justBitsRoots)
	prevJustRoots, err := ssz.HashTreeRoot(data.Data.PreviousJustifiedCheckpoint)
	if err != nil {
		return err
	}
	allRoots = append(allRoots, prevJustRoots)
	currJustRoots, err := ssz.HashTreeRoot(data.Data.CurrentJustifiedCheckpoint)
	if err != nil {
		return err
	}
	allRoots = append(allRoots, currJustRoots)
	finChkRoots, err := ssz.HashTreeRoot(data.Data.FinalizedCheckpoint)
	if err != nil {
		return err
	}
	allRoots = append(allRoots, finChkRoots)

	beaconStateRoot, err := ssz.HashTreeRoot(allRoots)
	if err != nil {
		return err
	}

	//state := fmt.Sprintf("%v", data.Data)
	fmt.Printf("%x\n", beaconStateRoot)

	/*
		root, err := ssz.HashTreeRoot(data.Data)
		if err != nil {
			return err
		}
	*/
	/*
		randaoMixes := data.Data.RandaoMixes

		for _, r := range randaoMixes {
			if r == randao.Data.Randao {
				fmt.Println(r)
			}
		}

		mixesRoot, err := hashStringArray(randaoMixes)
		if err != nil {
			return err
		}
		fmt.Printf("%x\n", mixesRoot)
	*/
	/*
		valRoot, err := hashValidatorRoot(&data.Data.Validator)
		if err != nil {
			return err
		}
		fmt.Printf("%x\n", valRoot)
	*/
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
