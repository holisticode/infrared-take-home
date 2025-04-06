package pkg

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
	mt "github.com/txaty/go-merkletree"
)

type ByteArray struct {
	data []byte
}

func (bt *ByteArray) Serialize() ([]byte, error) {
	return bt.data, nil
}

type BeaconStateDataLeaf struct {
	Field string
	data  []byte
}

func (bt *BeaconStateDataLeaf) Serialize() ([]byte, error) {
	return bt.data, nil
}

func AssembleRandaoProof(
	clint *ethclient.Client,
	randaoIndex int,
) error {
	/*
		resp, err := http.Get(RANDAO_API)
		if err != nil {
			return err
		}

		var randaoData RandaoData
		if err := json.NewDecoder(resp.Body).Decode(&randaoData); err != nil {
			return err
		}
		resp.Body.Close()

		randao := randaoData.Data.Randao

					resp, err = http.Get(BEACON_STATE_API)
					if err != nil {
						return err
					}
				defer resp.Body.Close()
			var state BeaconStateData
			if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
				return fmt.Errorf("failed decoding json from API: %v", err)
			}
	*/
	jsonFile, err := os.Open("pkg/beacon_state_test_data.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		return fmt.Errorf("failed to open json file: %v", err)
	}
	log.Println("Successfully opened json file")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteJSON, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	var state BeaconStateData
	if err := json.Unmarshal(byteJSON, &state); err != nil {
		return fmt.Errorf("failed decoding json from file: %v", err)
	}

	beaconStateLeaves, err := calculateBeaconStateLeaves(&state.Data)
	if err != nil {
		return err
	}
	beaconStateTreeBlocks, err := getMerkleTreeDataBlocksFromStrings(beaconStateLeaves)
	if err != nil {
		return err
	}
	cf := getMerkleTreeconfig()
	beaconStateTree, err := mt.New(cf, beaconStateTreeBlocks)
	if err != nil {
		return err
	}

	beaconStateRoot := beaconStateTree.Root
	rdProof, err := beaconStateTree.Proof(beaconStateTreeBlocks[14])
	if err != nil {
		return err
	}
	ok, err := beaconStateTree.Verify(beaconStateTreeBlocks[14], rdProof)
	if err != nil {
		return err
	}
	if ok {
		log.Println("successfully verified randao proof in randao mixes")
	} else {
		log.Fatal("verification or randao proof in randao mixes failed!")
	}

	fmt.Printf("beaconStateRoot %x\n", beaconStateRoot)

	randaoMixes := state.Data.RandaoMixes
	randao := randaoMixes[randaoIndex]
	bt := make([]mt.DataBlock, len(randaoMixes))

	index := -1
	for i, r := range randaoMixes {
		if r == randao {
			index = i
		}
		b, err := hex.DecodeString(r[2:])
		if err != nil {
			return err
		}
		bt[i] = &ByteArray{data: b}
	}

	if index == -1 {
		return errors.New("the randao wasn't in the list, as was expected")
	}

	control, err := mt.New(cf, bt)
	if err != nil {
		return err
	}
	ctrlRoot := control.Root
	cproof, err := control.Proof(bt[randaoIndex])
	if err != nil {
		return err
	}
	ok, err = control.Verify(bt[randaoIndex], cproof)
	if err != nil {
		return err
	}
	if ok {
		log.Println("successfully verified randao proof in block root")
	} else {
		log.Fatal("verification or randao proof in block root failed!")
	}
	fmt.Printf("ctrlRoot %x\n", ctrlRoot)

	/*
		proof := mtree.GenerateProof(index)

			for _, p := range proof {
				fmt.Printf("%x\n", p)
			}
	*/

	return nil

}

func getMerkleTreeDataBlocksFromStrings(data []string) ([]mt.DataBlock, error) {
	var err error
	bt := make([]mt.DataBlock, len(data))
	for i, s := range data {
		data := s
		var arr []byte
		// if it's a hex string we should decode it
		if strings.HasPrefix(s, "0x") {
			data = s[2:]
			arr, err = hex.DecodeString(data)
			// TODO: This is probably INCORRECT due to the JSON encoding,
			// which returns strings when there are uints or other number formats
			// the correct implementation is probably hashing number values, not their string rep
		} else {
			arr = []byte(data)
		}
		if err != nil {
			return nil, err
		}
		b := &ByteArray{data: arr}
		bt[i] = b
	}
	return bt, nil
}

func buildMerkleTreeFromStrings(data []string) (*mt.MerkleTree, error) {
	bt, err := getMerkleTreeDataBlocksFromStrings(data)
	if err != nil {
		return nil, err
	}
	return mt.New(nil, bt)
}

/*
func hashStringArray(data []string) (RootHash, error) {
	var list []merkletree.Content

	for _, s := range data {
		data := s
		if strings.HasPrefix(s, "0x") {
			data = s[2:]
		}
		list = append(list, ValidatorStringContent{Field: data})
	}

	t, err := merkletree.NewTree(list)
	if err != nil {
		return RootHash{0}, err
	}

	//Get the Merkle Root of the tree
	mr := t.MerkleRoot()
	return RootHash(mr), nil
}
*/

func calculateBeaconStateLeaves(bs *BeaconStateSimplified) ([]string, error) {
	var list []string
	state := reflect.ValueOf(*bs)

	for i := 0; i < state.NumField(); i++ {
		field := state.Field(i).Interface()
		fieldType := reflect.TypeOf(field)
		if fieldType == reflect.TypeOf("") {
			// this is a direct leaf, hash it and encode to hex
			strVal := field.(string)
			if strings.HasPrefix(strVal, "0x") {
				list = append(list, strVal)
			} else {
				// NOTE: same hack, if it's a number this is probably incorrect
				fieldHash := hash([]byte(strVal))
				list = append(list, hex.EncodeToString(fieldHash))
			}
		} else if fieldType.Kind() == reflect.Slice && fieldType.Elem() == reflect.TypeOf("") {
			subTree, err := buildMerkleTreeFromStrings(field.([]string))
			if err != nil {
				return nil, err
			}
			subRoot := subTree.Root
			subRootStr := hex.EncodeToString(subRoot)
			list = append(list, subRootStr)
		} else if fieldType.Kind() == reflect.Slice && fieldType.Elem().Kind() == reflect.Struct {
			var subList []string
			for i := 0; i < reflect.ValueOf(field).Len(); i++ {
				strlist := iterateFieldsOfStruct(reflect.ValueOf(field).Index(i))
				subTree, err := buildMerkleTreeFromStrings(strlist)
				if err != nil {
					return nil, err
				}
				subRoot := subTree.Root
				subRootStr := hex.EncodeToString(subRoot)
				subList = append(subList, subRootStr)
			}
			subTree, err := buildMerkleTreeFromStrings(subList)
			if err != nil {
				return nil, err
			}
			subRoot := subTree.Root
			subRootStr := hex.EncodeToString(subRoot)
			list = append(list, subRootStr)
		} else if fieldType.Kind() == reflect.Struct {
			strlist := iterateFieldsOfStruct(reflect.ValueOf(field))
			subTree, err := buildMerkleTreeFromStrings(strlist)
			if err != nil {
				return nil, err
			}
			subRoot := subTree.Root
			subRootStr := hex.EncodeToString(subRoot)
			list = append(list, subRootStr)
		} else {
			return nil, fmt.Errorf("unexpected validator field type: %v", fieldType)
		}
	}

	return list, nil
}

func iterateFieldsOfStruct(field reflect.Value) []string {
	//structType := reflect.ValueOf(field)
	var strlist []string
	for k := 0; k < field.NumField(); k++ {
		f := field.Field(k).Interface()
		if reflect.TypeOf(f) == reflect.TypeOf(true) {
			if reflect.ValueOf(f).Bool() == false {
				strlist = append(strlist, string([]byte{0}))
			} else {
				strlist = append(strlist, string([]byte{1}))
			}
		} else {
			strlist = append(strlist, f.(string))
		}
	}
	return strlist
}

func getMerkleTreeconfig() *mt.Config {
	return &mt.Config{
		// Customizable hash function used for tree generation.
		HashFunc: defaultHashFunc,
		// Number of goroutines run in parallel.
		// If RunInParallel is true and NumRoutine is set to 0, use number of CPU as the number of goroutines.
		NumRoutines: 0,
		// Mode of the Merkle Tree generation.
		Mode: mt.ModeProofGenAndTreeBuild,
		// If RunInParallel is true, the generation runs in parallel, otherwise runs without parallelization.
		// This increase the performance for the calculation of large number of data blocks, e.g. over 10,000 blocks.
		RunInParallel: true,
		// SortSiblingPairs is the parameter for OpenZeppelin compatibility.
		// If set to `true`, the hashing sibling pairs are sorted.
		SortSiblingPairs: false,
		// If true, the leaf nodes are NOT hashed before being added to the Merkle Tree.
		DisableLeafHashing: false,
	}
}
