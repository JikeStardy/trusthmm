package main

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"trust/common"
)

var (
	contract               = "PseudonymContract"
	queryCarIdByPseudonym  = "QueryIdByPseudonym"
	bindPseudonymWithCarId = "BindPseudonymWithCarId"
)

type PseudonymContract interface {
	QueryCarIdByPseudonym(shim.ChaincodeStubInterface, []string) pb.Response
	BindPseudonymWithCarId(shim.ChaincodeStubInterface, []string) pb.Response
	TransferPseudonym(shim.ChaincodeStubInterface, []string) pb.Response
	UnbindPseudonymWithCarId(shim.ChaincodeStubInterface, []string) pb.Response
}

type PseudonymChainCode struct{}

func (pc *PseudonymChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(common.OK)
}

func (pc *PseudonymChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Printf("invoke %s -> %s\n", contract, function)

	switch function {
	case queryCarIdByPseudonym:
		return pc.QueryCarIdByPseudonym(stub, args)
	case bindPseudonymWithCarId:
		return pc.BindPseudonymWithCarId(stub, args)
	default:
		return shim.Error(common.NoSuchFunction(contract, function))
	}
}

// QueryIdByPseudonym 通过假名/昵称查询真实 ID
func (pc *PseudonymChainCode) QueryCarIdByPseudonym(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var pseudonym string
	err := common.ParseArgs(args, &pseudonym)
	if err != nil {
		return shim.Error(err.Error())
	}
	if pseudonym == "" {
		return shim.Error("Invalid Pseudonym")
	}

	carId, err := stub.GetState(pseudonym)
	if err != nil {
		return shim.Error(fmt.Sprintf("Get CarId By Pseudonym Fail: %s", err))
	}

	if carId != nil {
		return shim.Success(carId)
	} else {
		return shim.Error("Please Bind Pseudonym to A CarId First")
	}
}

// BindPseudonymWithCarId 将假名/昵称与 CarId 进行绑定
func (pc *PseudonymChainCode) BindPseudonymWithCarId(
	stub shim.ChaincodeStubInterface,
	args []string) pb.Response {

	var pseudonym, carId string
	err := common.ParseArgs(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	if pseudonym == "" || carId == "" {
		return shim.Error("Invalid Pseudonym or CarId")
	}

	err = stub.PutState(pseudonym, []byte(carId))
	if err != nil {
		return shim.Error(fmt.Sprintf("Bind Pseudonym with CarId Fail: %s", err))
	}

	return shim.Success(common.OK)
}

// TransferPseudonym 将假名/昵称从 A 的 CarId 转移到 B 的 CarId TODO
func (pc *PseudonymChainCode) TransferPseudonym(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return shim.Error("Not yet implemented")
}

// UnbindPseudonymWithCarId 将假名/昵称与 CarId 解绑 TODO
func (pc *PseudonymChainCode) UnbindPseudonymWithCarId(shim.ChaincodeStubInterface, []string) pb.Response {
	return shim.Error("Not yet implemented")
}

func main() {
	err := shim.Start(new(PseudonymChainCode))
	if err != nil {
		fmt.Printf("Error starting IdContract: %s", err)
	}
}
