package main

import (
	"bytes"
	"fmt"
	"trust/common"
	"trust/common/entity/car"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

var (
	contract        = "CarContract"
	existCar        = "ExistCar"
	registerCar     = "RegisterCar"
	queryCarById    = "QueryCarById"
	queryCarByRange = "QueryCarByRange"
	updateCar       = "UpdateCar"
)

type CarContract interface {
	ExistCar(shim.ChaincodeStubInterface, string) pb.Response
	RegisterCar(shim.ChaincodeStubInterface, string) pb.Response
	QueryCarById(shim.ChaincodeStubInterface, string) pb.Response
	QueryCarByRange(shim.ChaincodeStubInterface, string) pb.Response
	UpdateCar(shim.ChaincodeStubInterface, string) pb.Response
}

type CarChainCode struct{}

func (cc *CarChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(common.OK)
}

func (cc *CarChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Printf("invoke %s -> %s\n", contract, function)

	// 定位分支逻辑
	switch function {
	case existCar:
		return cc.ExistCar(stub, args[0])
	case registerCar:
		return cc.RegisterCar(stub, args[0])
	case queryCarById:
		return cc.QueryCarById(stub, args[0])
	case queryCarByRange:
		return cc.QueryCarByRange(stub, args[0])
	case updateCar:
		return cc.UpdateCar(stub, args[0])
	default:
		return shim.Error(common.NoSuchFunction(contract, function))
	}
}

// ExistCar 利用 Car 的 Id 来查询车辆是否已经注册
func (cc *CarChainCode) ExistCar(stub shim.ChaincodeStubInterface, carId string) pb.Response {
	// 对 Car 的 Id 进行判空
	if carId == "" {
		return shim.Error(fmt.Sprintf("Invalid CarId: %s", carId))
	}

	carAsJsonBytes, err := stub.GetState(carId)
	// 获取失败
	if err != nil {
		return shim.Error(fmt.Sprintf("Get Car Info Fail: %s", err))
	}
	// 检查 GetState 的两个返回值是否都为 nil
	if carAsJsonBytes != nil {
		return shim.Success(common.NOT_EXIST)
	}
	return shim.Success(common.EXIST)
}

// RegisterCar 登记车辆信息
func (cc *CarChainCode) RegisterCar(stub shim.ChaincodeStubInterface, arg string) pb.Response {
	// 从 JSON 反序列化出 car 对象
	carToRegister, err := car.FromJSON(arg)
	if err != nil {
		return shim.Error(fmt.Sprintf("Broken Data: %s", err))
	}

	// 执行车辆信息登记
	err = stub.PutState(carToRegister.CarId, carToRegister.ToJsonBytes())
	if err != nil {
		return shim.Error(fmt.Sprint("Register Fail: %s", err))
	}

	return shim.Success(common.OK)
}

// QueryCarById 根据Id查询车辆信息
func (cc *CarChainCode) QueryCarById(stub shim.ChaincodeStubInterface, carId string) pb.Response {
	// 对 Car 的 Id 进行判空
	if carId == "" {
		return shim.Error(fmt.Sprintf("Invalid CarId: %s", carId))
	}

	carAsJsonBytes, err := stub.GetState(carId)
	if err != nil {
		return shim.Error(fmt.Sprintf("Get Car Info Fail: %s", err))
	}

	c, err := car.FromJSON(string(carAsJsonBytes))
	if err != nil {
		return shim.Error(fmt.Sprintf("Broken Data: %s", err))
	}

	return shim.Success(c.ToJsonBytes())
}

// QueryCarByRange 根据给定范围执行范围查询 TODO
func (cc *CarChainCode) QueryCarByRange(stub shim.ChaincodeStubInterface, arg string) pb.Response {
	return shim.Error("Not yet implemented")
}

// UpdateCar 根据给定信息更新车辆信息
func (cc *CarChainCode) UpdateCar(stub shim.ChaincodeStubInterface, newCarInfo string) pb.Response {
	if newCarInfo == "" {
		return shim.Error("New Car Info Should Not Be Empty")
	}

	c, err := car.FromJSON(newCarInfo)
	if err != nil {
		return shim.Error(fmt.Sprintf("Broken Data: %s", err))
	}

	if bytes.Equal(cc.ExistCar(stub, c.CarId).Payload, common.EXIST) {
		err = stub.PutState(c.CarId, c.ToJsonBytes())
		if err != nil {
			return shim.Error(fmt.Sprintf("Update Car Info Fail: %s", err))
		}
		return shim.Success(common.OK)
	} else {
		return shim.Error("Please Register This Car First")
	}
}

func main() {
	err := shim.Start(new(CarChainCode))
	if err != nil {
		fmt.Printf("Error starting IdContract: %s", err)
	}
}
