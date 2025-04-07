package pkg

import (
	"crypto/sha256"
	"errors"

	"github.com/cbergoon/merkletree"
)

/**
This file is now OBSOLETE

I first wanted to do proofs about Validator information,
but then went to do Randao
*/

type ValidatorData struct {
	Data ValidatorInfo
}

type ValidatorInfo struct {
	Index     string
	Balance   string
	Status    string
	Validator ValidatorDetails
}

type ValidatorStringContent struct {
	Field string
}

type ValidatorBoolContent struct {
	Field bool
}

func defaultHashFunc(b []byte) ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write(b); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func (s ValidatorStringContent) CalculateHash() ([]byte, error) {
	return defaultHashFunc([]byte(s.Field))
}

func (s ValidatorStringContent) Equals(other merkletree.Content) (bool, error) {
	if otherT, ok := other.(ValidatorStringContent); !ok {
		return false, errors.New("value is not of type ValidatorStringContent")
	} else {
		return s.Field == otherT.Field, nil
	}
}

func (b ValidatorBoolContent) CalculateHash() ([]byte, error) {
	h := sha256.New()
	val := []byte{0}
	if b.Field == true {
		val = []byte{1}
	}
	if _, err := h.Write(val); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func (b ValidatorBoolContent) Equals(other merkletree.Content) (bool, error) {
	if otherT, ok := other.(ValidatorBoolContent); !ok {
		return false, errors.New("value is not of type ValidatorBoolContent")
	} else {
		return b.Field == otherT.Field, nil
	}
}
