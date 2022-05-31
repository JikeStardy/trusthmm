package main

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"trust/common"
	"trust/common/entity/env"
	"trust/common/entity/trust"
)

type TrustContract interface {
	InitTrust(shim.ChaincodeStubInterface, []string) pb.Response
	QueryTrust(shim.ChaincodeStubInterface, []string) pb.Response
	AddTrust(shim.ChaincodeStubInterface, []string) pb.Response
	UpdateTrust(shim.ChaincodeStubInterface, []string) pb.Response
}

var (
	contract    = "TrustContract"
	initTrust   = "InitTrust"
	queryTrust  = "QueryTrust"
	addTrust    = "AddTrust"
	updateTrust = "UpdateTrust"
)

type TrustChainCode struct{}

func (tc *TrustChainCode) Init(shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(common.OK)
}

func (tc *TrustChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Printf("invoke %s -> %s\n", contract, function)

	if len(args) == 0 {
		return shim.Error(fmt.Sprintf("Incorrect Arguments for %s", contract))
	}

	t, err := trust.FromJSON(args[0])
	if err != nil {
		return shim.Error(fmt.Sprintf("Broken Data: %s", err))
	}

	// TODO 根据论文，对 TrustContract 各个方法的输入皆为 Pseudonym，则先 InvokeChainCode 将 Pseudonym 转为 CarId
	resp := stub.InvokeChaincode(env.PseudonymContract,
		common.CreateInvokeArgs("QueryCarIdByPseudonym", t.Id), env.Channel)
	if resp.Status == shim.ERROR {
		return resp
	}

	// TODO 此时将 Pseudonym 转为 CarId
	t.Id = string(resp.GetPayload())

	switch function {
	case initTrust:
		return tc.InitTrust(stub, t)
	case queryTrust:
		return tc.QueryTrust(stub, t)
	case addTrust:
		return tc.AddTrust(stub, t)
	case updateTrust:
		return tc.UpdateTrust(stub, t)
	default:
		return shim.Error(common.NoSuchFunction(contract, function))
	}
}

func (tc *TrustChainCode) getVehicleTrustState(
	stub shim.ChaincodeStubInterface, carId string) (*trust.VehicleTrust, error) {
	vtAsBytes, err := stub.GetState(carId)
	vtState, err := trust.FromJSON(string(vtAsBytes))
	if err != nil {
		return nil, err
	}
	return &vtState, nil
}

func (tc *TrustChainCode) InitTrust(stub shim.ChaincodeStubInterface, vt trust.VehicleTrust) pb.Response {
	// 将 vt 初始值设为 env.InitTrustValue = 100
	vt.VehicleTrust = env.InitialTrustValue
	// 将完整对象存入
	err := stub.PutState(vt.Id, vt.ToJsonBytes())
	if err != nil {
		return shim.Error(fmt.Sprintf("Init Trust Value Fail: %s", err))
	}

	return shim.Success(common.OK)
}

func (tc *TrustChainCode) QueryTrust(stub shim.ChaincodeStubInterface, vt trust.VehicleTrust) pb.Response {
	vtState, err := tc.getVehicleTrustState(stub, vt.Id)
	if err != nil {
		return shim.Error(fmt.Sprintf("Broken Data: %s", err))
	}

	return shim.Success(vtState.ToJsonBytes())
}

func (tc *TrustChainCode) AddTrust(stub shim.ChaincodeStubInterface, vt trust.VehicleTrust) pb.Response {
	vtState, err := tc.getVehicleTrustState(stub, vt.Id)
	if err != nil {
		return shim.Error(fmt.Sprintf("Broken Data: %s", err))
	}

	// 增加 vt 值
	vtState.VehicleTrust += vt.VehicleTrust

	err = stub.PutState(vtState.Id, vtState.ToJsonBytes())
	if err != nil {
		return shim.Error(fmt.Sprintf("%s Fail: %s", addTrust, err))
	}

	return shim.Success(common.OK)
}

func (tc *TrustChainCode) UpdateTrust(stub shim.ChaincodeStubInterface, vt trust.VehicleTrust) pb.Response {
	// TODO 是否需要判定此 VT 是否已存在 State 中
	err := stub.PutState(vt.Id, vt.ToJsonBytes())
	if err != nil {
		return shim.Error(fmt.Sprintf("Update Vehicle Trust Fail: %s", err))
	}

	return shim.Success(common.OK)
}

func main() {
	err := shim.Start(new(TrustChainCode))
	if err != nil {
		fmt.Printf("Error starting IdContract: %s", err)
	}
}
