package generic_operator

import (
	"fmt"

	"github.com/Layr-Labs/eigensdk-go/logging"
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

func PackCallParams(avsName string, avsAbi *abi.ABI, method string, params []interface{}) ([]byte, error) {
	inputs, err := getInputsAbiByMethod(avsName, avsAbi, method)
	if err != nil {
		return nil, errors.Wrap(err, "get inputs abi failed")
	}

	if len(inputs) == 0 {
		return nil, fmt.Errorf("the %s inputs %s len is zero, cannot use bls sign", avsName, method)
	}

	if len(inputs) == 1 {
		if len(params) == 0 {

			// no inputs need
			return []byte{}, nil
		} else {

			return nil, fmt.Errorf("the %s inputs not empty but no params raw to decode", avsName)
		}
	}

	inputsWithoutBlsSign := inputs[:len(inputs)-1]

	raw, err := inputsWithoutBlsSign.Pack(params...)
	if err != nil {
		return nil, errors.Wrapf(err, "pack params %s method %s failed", avsName, method)
	}

	return raw, nil
}

func UnpackCallParams(logger logging.Logger, avsName string, avsAbi *abi.ABI, method string, paramsRaw []byte) ([]interface{}, error) {
	inputs, err := getInputsAbiByMethod(avsName, avsAbi, method)
	if err != nil {
		return nil, errors.Wrap(err, "get inputs abi failed")
	}

	if len(inputs) == 0 {
		return nil, fmt.Errorf("the %s inputs %s len is zero, cannot use bls sign", avsName, method)
	}

	if len(inputs) == 1 {
		if len(paramsRaw) == 0 {
			// no inputs need
			return []interface{}{}, nil
		} else {
			return nil, fmt.Errorf("the %s inputs not empty but no params raw to decode", avsName)
		}

	}

	inputsWithoutBlsSign := inputs[:len(inputs)-1]

	logger.Debug("unpack by api input", "abi", inputsWithoutBlsSign, "bytes", paramsRaw)

	params, err := inputsWithoutBlsSign.Unpack(paramsRaw)
	if err != nil {
		return nil, errors.Wrapf(err, "unpack params %s method %s failed", avsName, method)
	}

	return params, nil
}
