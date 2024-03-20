// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contractERC20PresetFixedSupply

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ContractERC20PresetFixedSupplyMetaData contains all meta data concerning the ContractERC20PresetFixedSupply contract.
var ContractERC20PresetFixedSupplyMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"name\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"symbol\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"initialSupply\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"allowance\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"approve\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"balanceOf\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"burn\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"burnFrom\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decreaseAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"subtractedValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"increaseAllowance\",\"inputs\":[{\"name\":\"spender\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"addedValue\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"name\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"symbol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"totalSupply\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transfer\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferFrom\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"Approval\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"spender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Transfer\",\"inputs\":[{\"name\":\"from\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"to\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162000e6038038062000e608339810160408190526200003491620002dd565b8351849084906200004d9060039060208501906200016a565b508051620000639060049060208401906200016a565b5050506200007881836200008260201b60201c565b50505050620003d6565b6001600160a01b038216620000dd5760405162461bcd60e51b815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f206164647265737300604482015260640160405180910390fd5b8060026000828254620000f1919062000372565b90915550506001600160a01b038216600090815260208190526040812080548392906200012090849062000372565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a35050565b828054620001789062000399565b90600052602060002090601f0160209004810192826200019c5760008555620001e7565b82601f10620001b757805160ff1916838001178555620001e7565b82800160010185558215620001e7579182015b82811115620001e7578251825591602001919060010190620001ca565b50620001f5929150620001f9565b5090565b5b80821115620001f55760008155600101620001fa565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126200023857600080fd5b81516001600160401b038082111562000255576200025562000210565b604051601f8301601f19908116603f0116810190828211818310171562000280576200028062000210565b816040528381526020925086838588010111156200029d57600080fd5b600091505b83821015620002c15785820183015181830184015290820190620002a2565b83821115620002d35760008385830101525b9695505050505050565b60008060008060808587031215620002f457600080fd5b84516001600160401b03808211156200030c57600080fd5b6200031a8883890162000226565b955060208701519150808211156200033157600080fd5b50620003408782880162000226565b60408701516060880151919550935090506001600160a01b03811681146200036757600080fd5b939692955090935050565b600082198211156200039457634e487b7160e01b600052601160045260246000fd5b500190565b600181811c90821680620003ae57607f821691505b60208210811415620003d057634e487b7160e01b600052602260045260246000fd5b50919050565b610a7a80620003e66000396000f3fe608060405234801561001057600080fd5b50600436106100cf5760003560e01c806342966c681161008c57806395d89b411161006657806395d89b41146101ad578063a457c2d7146101b5578063a9059cbb146101c8578063dd62ed3e146101db57600080fd5b806342966c681461015c57806370a082311461017157806379cc67901461019a57600080fd5b806306fdde03146100d4578063095ea7b3146100f257806318160ddd1461011557806323b872dd14610127578063313ce5671461013a5780633950935114610149575b600080fd5b6100dc6101ee565b6040516100e9919061087f565b60405180910390f35b6101056101003660046108f0565b610280565b60405190151581526020016100e9565b6002545b6040519081526020016100e9565b61010561013536600461091a565b610298565b604051601281526020016100e9565b6101056101573660046108f0565b6102bc565b61016f61016a366004610956565b6102de565b005b61011961017f36600461096f565b6001600160a01b031660009081526020819052604090205490565b61016f6101a83660046108f0565b6102eb565b6100dc610304565b6101056101c33660046108f0565b610313565b6101056101d63660046108f0565b610393565b6101196101e9366004610991565b6103a1565b6060600380546101fd906109c4565b80601f0160208091040260200160405190810160405280929190818152602001828054610229906109c4565b80156102765780601f1061024b57610100808354040283529160200191610276565b820191906000526020600020905b81548152906001019060200180831161025957829003601f168201915b5050505050905090565b60003361028e8185856103cc565b5060019392505050565b6000336102a68582856104f1565b6102b185858561056b565b506001949350505050565b60003361028e8185856102cf83836103a1565b6102d99190610a15565b6103cc565b6102e83382610739565b50565b6102f68233836104f1565b6103008282610739565b5050565b6060600480546101fd906109c4565b6000338161032182866103a1565b9050838110156103865760405162461bcd60e51b815260206004820152602560248201527f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f77604482015264207a65726f60d81b60648201526084015b60405180910390fd5b6102b182868684036103cc565b60003361028e81858561056b565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b6001600160a01b03831661042e5760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b606482015260840161037d565b6001600160a01b03821661048f5760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b606482015260840161037d565b6001600160a01b0383811660008181526001602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92591015b60405180910390a3505050565b60006104fd84846103a1565b9050600019811461056557818110156105585760405162461bcd60e51b815260206004820152601d60248201527f45524332303a20696e73756666696369656e7420616c6c6f77616e6365000000604482015260640161037d565b61056584848484036103cc565b50505050565b6001600160a01b0383166105cf5760405162461bcd60e51b815260206004820152602560248201527f45524332303a207472616e736665722066726f6d20746865207a65726f206164604482015264647265737360d81b606482015260840161037d565b6001600160a01b0382166106315760405162461bcd60e51b815260206004820152602360248201527f45524332303a207472616e7366657220746f20746865207a65726f206164647260448201526265737360e81b606482015260840161037d565b6001600160a01b038316600090815260208190526040902054818110156106a95760405162461bcd60e51b815260206004820152602660248201527f45524332303a207472616e7366657220616d6f756e7420657863656564732062604482015265616c616e636560d01b606482015260840161037d565b6001600160a01b038085166000908152602081905260408082208585039055918516815290812080548492906106e0908490610a15565b92505081905550826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161072c91815260200190565b60405180910390a3610565565b6001600160a01b0382166107995760405162461bcd60e51b815260206004820152602160248201527f45524332303a206275726e2066726f6d20746865207a65726f206164647265736044820152607360f81b606482015260840161037d565b6001600160a01b0382166000908152602081905260409020548181101561080d5760405162461bcd60e51b815260206004820152602260248201527f45524332303a206275726e20616d6f756e7420657863656564732062616c616e604482015261636560f01b606482015260840161037d565b6001600160a01b038316600090815260208190526040812083830390556002805484929061083c908490610a2d565b90915550506040518281526000906001600160a01b038516907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef906020016104e4565b600060208083528351808285015260005b818110156108ac57858101830151858201604001528201610890565b818111156108be576000604083870101525b50601f01601f1916929092016040019392505050565b80356001600160a01b03811681146108eb57600080fd5b919050565b6000806040838503121561090357600080fd5b61090c836108d4565b946020939093013593505050565b60008060006060848603121561092f57600080fd5b610938846108d4565b9250610946602085016108d4565b9150604084013590509250925092565b60006020828403121561096857600080fd5b5035919050565b60006020828403121561098157600080fd5b61098a826108d4565b9392505050565b600080604083850312156109a457600080fd5b6109ad836108d4565b91506109bb602084016108d4565b90509250929050565b600181811c908216806109d857607f821691505b602082108114156109f957634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b60008219821115610a2857610a286109ff565b500190565b600082821015610a3f57610a3f6109ff565b50039056fea2646970667358221220e80e5ee5d6a68ccdd2764e0ef75c9475131f675f2e7d72b9951a981b9d30f0a464736f6c634300080c0033",
}

// ContractERC20PresetFixedSupplyABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractERC20PresetFixedSupplyMetaData.ABI instead.
var ContractERC20PresetFixedSupplyABI = ContractERC20PresetFixedSupplyMetaData.ABI

// ContractERC20PresetFixedSupplyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ContractERC20PresetFixedSupplyMetaData.Bin instead.
var ContractERC20PresetFixedSupplyBin = ContractERC20PresetFixedSupplyMetaData.Bin

// DeployContractERC20PresetFixedSupply deploys a new Ethereum contract, binding an instance of ContractERC20PresetFixedSupply to it.
func DeployContractERC20PresetFixedSupply(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string, initialSupply *big.Int, owner common.Address) (common.Address, *types.Transaction, *ContractERC20PresetFixedSupply, error) {
	parsed, err := ContractERC20PresetFixedSupplyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ContractERC20PresetFixedSupplyBin), backend, name, symbol, initialSupply, owner)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ContractERC20PresetFixedSupply{ContractERC20PresetFixedSupplyCaller: ContractERC20PresetFixedSupplyCaller{contract: contract}, ContractERC20PresetFixedSupplyTransactor: ContractERC20PresetFixedSupplyTransactor{contract: contract}, ContractERC20PresetFixedSupplyFilterer: ContractERC20PresetFixedSupplyFilterer{contract: contract}}, nil
}

// ContractERC20PresetFixedSupply is an auto generated Go binding around an Ethereum contract.
type ContractERC20PresetFixedSupply struct {
	ContractERC20PresetFixedSupplyCaller     // Read-only binding to the contract
	ContractERC20PresetFixedSupplyTransactor // Write-only binding to the contract
	ContractERC20PresetFixedSupplyFilterer   // Log filterer for contract events
}

// ContractERC20PresetFixedSupplyCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractERC20PresetFixedSupplyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractERC20PresetFixedSupplyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractERC20PresetFixedSupplyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractERC20PresetFixedSupplyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractERC20PresetFixedSupplyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractERC20PresetFixedSupplySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractERC20PresetFixedSupplySession struct {
	Contract     *ContractERC20PresetFixedSupply // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                   // Call options to use throughout this session
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// ContractERC20PresetFixedSupplyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractERC20PresetFixedSupplyCallerSession struct {
	Contract *ContractERC20PresetFixedSupplyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                         // Call options to use throughout this session
}

// ContractERC20PresetFixedSupplyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractERC20PresetFixedSupplyTransactorSession struct {
	Contract     *ContractERC20PresetFixedSupplyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                         // Transaction auth options to use throughout this session
}

// ContractERC20PresetFixedSupplyRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractERC20PresetFixedSupplyRaw struct {
	Contract *ContractERC20PresetFixedSupply // Generic contract binding to access the raw methods on
}

// ContractERC20PresetFixedSupplyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractERC20PresetFixedSupplyCallerRaw struct {
	Contract *ContractERC20PresetFixedSupplyCaller // Generic read-only contract binding to access the raw methods on
}

// ContractERC20PresetFixedSupplyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractERC20PresetFixedSupplyTransactorRaw struct {
	Contract *ContractERC20PresetFixedSupplyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContractERC20PresetFixedSupply creates a new instance of ContractERC20PresetFixedSupply, bound to a specific deployed contract.
func NewContractERC20PresetFixedSupply(address common.Address, backend bind.ContractBackend) (*ContractERC20PresetFixedSupply, error) {
	contract, err := bindContractERC20PresetFixedSupply(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ContractERC20PresetFixedSupply{ContractERC20PresetFixedSupplyCaller: ContractERC20PresetFixedSupplyCaller{contract: contract}, ContractERC20PresetFixedSupplyTransactor: ContractERC20PresetFixedSupplyTransactor{contract: contract}, ContractERC20PresetFixedSupplyFilterer: ContractERC20PresetFixedSupplyFilterer{contract: contract}}, nil
}

// NewContractERC20PresetFixedSupplyCaller creates a new read-only instance of ContractERC20PresetFixedSupply, bound to a specific deployed contract.
func NewContractERC20PresetFixedSupplyCaller(address common.Address, caller bind.ContractCaller) (*ContractERC20PresetFixedSupplyCaller, error) {
	contract, err := bindContractERC20PresetFixedSupply(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractERC20PresetFixedSupplyCaller{contract: contract}, nil
}

// NewContractERC20PresetFixedSupplyTransactor creates a new write-only instance of ContractERC20PresetFixedSupply, bound to a specific deployed contract.
func NewContractERC20PresetFixedSupplyTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractERC20PresetFixedSupplyTransactor, error) {
	contract, err := bindContractERC20PresetFixedSupply(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractERC20PresetFixedSupplyTransactor{contract: contract}, nil
}

// NewContractERC20PresetFixedSupplyFilterer creates a new log filterer instance of ContractERC20PresetFixedSupply, bound to a specific deployed contract.
func NewContractERC20PresetFixedSupplyFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractERC20PresetFixedSupplyFilterer, error) {
	contract, err := bindContractERC20PresetFixedSupply(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractERC20PresetFixedSupplyFilterer{contract: contract}, nil
}

// bindContractERC20PresetFixedSupply binds a generic wrapper to an already deployed contract.
func bindContractERC20PresetFixedSupply(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ContractERC20PresetFixedSupplyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractERC20PresetFixedSupply.Contract.ContractERC20PresetFixedSupplyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.Contract.ContractERC20PresetFixedSupplyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.Contract.ContractERC20PresetFixedSupplyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ContractERC20PresetFixedSupply.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ContractERC20PresetFixedSupply.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplySession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ContractERC20PresetFixedSupply.Contract.Allowance(&_ContractERC20PresetFixedSupply.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _ContractERC20PresetFixedSupply.Contract.Allowance(&_ContractERC20PresetFixedSupply.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ContractERC20PresetFixedSupply.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplySession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ContractERC20PresetFixedSupply.Contract.BalanceOf(&_ContractERC20PresetFixedSupply.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _ContractERC20PresetFixedSupply.Contract.BalanceOf(&_ContractERC20PresetFixedSupply.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ContractERC20PresetFixedSupply.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplySession) Decimals() (uint8, error) {
	return _ContractERC20PresetFixedSupply.Contract.Decimals(&_ContractERC20PresetFixedSupply.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyCallerSession) Decimals() (uint8, error) {
	return _ContractERC20PresetFixedSupply.Contract.Decimals(&_ContractERC20PresetFixedSupply.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ContractERC20PresetFixedSupply.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplySession) Name() (string, error) {
	return _ContractERC20PresetFixedSupply.Contract.Name(&_ContractERC20PresetFixedSupply.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyCallerSession) Name() (string, error) {
	return _ContractERC20PresetFixedSupply.Contract.Name(&_ContractERC20PresetFixedSupply.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ContractERC20PresetFixedSupply.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplySession) Symbol() (string, error) {
	return _ContractERC20PresetFixedSupply.Contract.Symbol(&_ContractERC20PresetFixedSupply.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyCallerSession) Symbol() (string, error) {
	return _ContractERC20PresetFixedSupply.Contract.Symbol(&_ContractERC20PresetFixedSupply.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ContractERC20PresetFixedSupply.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplySession) TotalSupply() (*big.Int, error) {
	return _ContractERC20PresetFixedSupply.Contract.TotalSupply(&_ContractERC20PresetFixedSupply.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyCallerSession) TotalSupply() (*big.Int, error) {
	return _ContractERC20PresetFixedSupply.Contract.TotalSupply(&_ContractERC20PresetFixedSupply.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplySession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.Contract.Approve(&_ContractERC20PresetFixedSupply.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.Contract.Approve(&_ContractERC20PresetFixedSupply.TransactOpts, spender, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyTransactor) Burn(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.contract.Transact(opts, "burn", amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplySession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.Contract.Burn(&_ContractERC20PresetFixedSupply.TransactOpts, amount)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 amount) returns()
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyTransactorSession) Burn(amount *big.Int) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.Contract.Burn(&_ContractERC20PresetFixedSupply.TransactOpts, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyTransactor) BurnFrom(opts *bind.TransactOpts, account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.contract.Transact(opts, "burnFrom", account, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplySession) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.Contract.BurnFrom(&_ContractERC20PresetFixedSupply.TransactOpts, account, amount)
}

// BurnFrom is a paid mutator transaction binding the contract method 0x79cc6790.
//
// Solidity: function burnFrom(address account, uint256 amount) returns()
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyTransactorSession) BurnFrom(account common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.Contract.BurnFrom(&_ContractERC20PresetFixedSupply.TransactOpts, account, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplySession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.Contract.DecreaseAllowance(&_ContractERC20PresetFixedSupply.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.Contract.DecreaseAllowance(&_ContractERC20PresetFixedSupply.TransactOpts, spender, subtractedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplySession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.Contract.IncreaseAllowance(&_ContractERC20PresetFixedSupply.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.Contract.IncreaseAllowance(&_ContractERC20PresetFixedSupply.TransactOpts, spender, addedValue)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyTransactor) Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.contract.Transact(opts, "transfer", to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplySession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.Contract.Transfer(&_ContractERC20PresetFixedSupply.TransactOpts, to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyTransactorSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.Contract.Transfer(&_ContractERC20PresetFixedSupply.TransactOpts, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.contract.Transact(opts, "transferFrom", from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplySession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.Contract.TransferFrom(&_ContractERC20PresetFixedSupply.TransactOpts, from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyTransactorSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ContractERC20PresetFixedSupply.Contract.TransferFrom(&_ContractERC20PresetFixedSupply.TransactOpts, from, to, amount)
}

// ContractERC20PresetFixedSupplyApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the ContractERC20PresetFixedSupply contract.
type ContractERC20PresetFixedSupplyApprovalIterator struct {
	Event *ContractERC20PresetFixedSupplyApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractERC20PresetFixedSupplyApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractERC20PresetFixedSupplyApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractERC20PresetFixedSupplyApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractERC20PresetFixedSupplyApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractERC20PresetFixedSupplyApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractERC20PresetFixedSupplyApproval represents a Approval event raised by the ContractERC20PresetFixedSupply contract.
type ContractERC20PresetFixedSupplyApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*ContractERC20PresetFixedSupplyApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ContractERC20PresetFixedSupply.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &ContractERC20PresetFixedSupplyApprovalIterator{contract: _ContractERC20PresetFixedSupply.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *ContractERC20PresetFixedSupplyApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _ContractERC20PresetFixedSupply.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractERC20PresetFixedSupplyApproval)
				if err := _ContractERC20PresetFixedSupply.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyFilterer) ParseApproval(log types.Log) (*ContractERC20PresetFixedSupplyApproval, error) {
	event := new(ContractERC20PresetFixedSupplyApproval)
	if err := _ContractERC20PresetFixedSupply.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractERC20PresetFixedSupplyTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ContractERC20PresetFixedSupply contract.
type ContractERC20PresetFixedSupplyTransferIterator struct {
	Event *ContractERC20PresetFixedSupplyTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *ContractERC20PresetFixedSupplyTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractERC20PresetFixedSupplyTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(ContractERC20PresetFixedSupplyTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *ContractERC20PresetFixedSupplyTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractERC20PresetFixedSupplyTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractERC20PresetFixedSupplyTransfer represents a Transfer event raised by the ContractERC20PresetFixedSupply contract.
type ContractERC20PresetFixedSupplyTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*ContractERC20PresetFixedSupplyTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ContractERC20PresetFixedSupply.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &ContractERC20PresetFixedSupplyTransferIterator{contract: _ContractERC20PresetFixedSupply.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ContractERC20PresetFixedSupplyTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _ContractERC20PresetFixedSupply.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractERC20PresetFixedSupplyTransfer)
				if err := _ContractERC20PresetFixedSupply.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_ContractERC20PresetFixedSupply *ContractERC20PresetFixedSupplyFilterer) ParseTransfer(log types.Log) (*ContractERC20PresetFixedSupplyTransfer, error) {
	event := new(ContractERC20PresetFixedSupplyTransfer)
	if err := _ContractERC20PresetFixedSupply.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
