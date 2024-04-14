package generic_operator

import (
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/pkg/errors"
)

func getInputsAbiByMethod(avsName string, avsAbi *abi.ABI, method string) (abi.Arguments, error) {
	methodAbi, ok := avsAbi.Methods[method]
	if !ok {
		return nil, fmt.Errorf("not found the method %s in avs %s", method, avsName)
	}

	return methodAbi.Inputs, nil

}

func packCallParams(avsName string, avsAbi *abi.ABI, method string, params []interface{}) ([]byte, error) {
	inputs, err := getInputsAbiByMethod(avsName, avsAbi, method)
	if err != nil {
		return nil, errors.Wrap(err, "get inputs abi failed")
	}

	if len(inputs) == 0 {
		return nil, fmt.Errorf("the %s inputs %s len is zero, cannot use bls sign", avsName, method)
	}

	if len(params) != 0 {
		if len(inputs) == 1 {
			return nil, fmt.Errorf(
				"the %s inputs %s params not nil, but the inputs len is 1, cannot use bls sign",
				avsName, method,
			)
		} else {
			// no inputs need
			return []byte{}, nil
		}
	}

	inputsWithoutBlsSign := inputs[:len(inputs)-1]

	raw, err := inputsWithoutBlsSign.Pack(params...)
	if err != nil {
		return nil, errors.Wrapf(err, "pack params %s method %s failed", avsName, method)
	}

	return raw, nil
}
