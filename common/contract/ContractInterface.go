package contract

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

type ContractInterface interface {
	Init(shim.ChaincodeStubInterface) pb.Response
	Invoke(shim.ChaincodeStubInterface) pb.Response
}
