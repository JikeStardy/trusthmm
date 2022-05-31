package main

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"trust/common"
	"trust/common/entity/event"
	"trust/common/entity/message"
)

var (
	contract      = "EventContract"
	submitMessage = "SubmitMessage"
)

type EventContract interface {
	SubmitMessage(shim.ChaincodeStubInterface, []string) pb.Response
}

type EventChainCode struct{}

func (mc *EventChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(common.OK)
}

func (mc *EventChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Printf("invoke %s -> %s\n", contract, function)

	switch function {
	case submitMessage:
		return mc.SubmitMessage(stub, args)
	default:
		return shim.Error(common.NoSuchFunction(contract, function))
	}
}

func (mc *EventChainCode) SubmitMessage(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var eventAsJson string
	err := common.ParseArgs(args, &eventAsJson)
	if err != nil {
		return shim.Error(fmt.Sprintf("Parse Event Fail: %s", err))
	}

	e, err := event.FromJSON(eventAsJson)
	if err != nil {
		return shim.Error(fmt.Sprintf("Broken Data: %s", err))
	}

	// TODO 05.30 17:47 CarContract 与 TrustContract 的 Key 冲突，考虑废除 CarContract
	// 确认 CarId 对应的车辆是否存在
	//resp := stub.InvokeChaincode(env.CarContract, common.CreateInvokeArgs("ExistCar", e.CarId), env.Channel)
	//if resp.Status == shim.ERROR || bytes.Equal(resp.GetPayload(), common.NOT_EXIST) {
	//	return resp
	//}

	// 05.30 12:43 一定要上复合键 05.30 18:49
	carMessageKey, err := shim.CreateCompositeKey("CarId~MessageId", []string{e.CarId, message.NewMessageId()})
	if err != nil {
		return shim.Error(fmt.Sprintf("Create Composite Key Fail: %s", err))
	}

	err = stub.PutState(carMessageKey, e.Message.ToJsonBytes())
	if err != nil {
		return shim.Error(fmt.Sprintf("Submit Message Fail: %s", err))
	}

	return shim.Success(common.OK)
}

func main() {
	err := shim.Start(new(EventChainCode))
	if err != nil {
		fmt.Printf("Error starting IdContract: %s", err)
	}
}
