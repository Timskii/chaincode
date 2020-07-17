package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type Datas struct {
	ID                    string `json:"id"`
	Hash                  string
	Date                  string  `json:"date"`
	Type                  string  `json:"type"`
	Status                string  `json:"status"`
}

type Data struct {
	JsonMessage    string
	HashMessage    string
	OldJsonMessage string
	OldHashMessage string
}

type DatasChaincode struct {
	Stub shim.ChaincodeStubInterface
}

func (i *DatasChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	log.Println("Init")
	return shim.Success(nil)
}

var response *Response

func (i *DatasChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	i.Stub = stub
	function, args := i.Stub.GetFunctionAndParameters()
	domain := args[0]
	id := args[1]
	response = newResponse(i.Stub.GetTxID(), domain, id)
	log.Printf("function=%s\n", function)
	if function == "delete" {
		return i.delete(args)
	} else if function == "get" {
		return i.get(args)
	} else if function == "register" {
		return i.register(args)
	}
	return response.error("Invalid invoke function name. Expecting \"register\" \"delete\" \"get\"")
}

func (i *DatasChaincode) register(args []string) pb.Response {
	log.Printf("args = %v", args)

	jsonMessage := []byte(args[2])
	hash := args[3]

	data := &Datas{}
	err := json.Unmarshal(jsonMessage, &data)
	log.Println("register: Unmarshal message")

	if err != nil {
		log.Printf("register: err Unmarshal %v", err)
		return response.error(err.Error())
	}

	if len(args) > 4 {
		oldHash := args[4]
		log.Printf("register: oldHash = <%#v>", oldHash)
		if oldHash != "" {
			err := i.checkHash(data.ID, oldHash)
			if err != nil {
				log.Printf("err %v", err)
				return response.error(err.Error())
			}
		}
	}

	data.Hash = hash
	response.Hash = data.Hash
	keyData := data.ID
	assetInfo := &AssetInfo{Key: keyData, Status: "Ok"}
	log.Println("register: get datasBytes")
	datasBytes, err := json.Marshal(data)
	if err != nil {
		assetInfo.Status = "Error"
		log.Printf("err = %v", err)
		return response.error("register: error marshal data = " + err.Error())
	}
	err = i.Stub.PutState(keyData, datasBytes)
	if err != nil {
		assetInfo.Status = "Error"
		log.Printf("register: err PutState = %v", err)
		return response.error(err.Error())
	}
	response.AssetInfo = *assetInfo

	return response.success("REGISTER")
}

func (i *DatasChaincode) checkHash(id, hash string) error {

	datasBytes, err := i.Stub.GetState(id)
	if err != nil {
		log.Printf("checkHash: GetState by key %v err =%v", id, err)
		return err
	}
	data := &Datas{}
	//m := make(map[string]interface{})
	err = json.Unmarshal(datasBytes, &data)

	if err != nil {
		log.Printf("checkHash: err Unmarshal %v", err)
		return err
	}
	if data.Hash == hash {
		return nil
	} else {
		return errors.New("Hash is invalid!")
	}
}

func (i *DatasChaincode) delete(args []string) pb.Response {

	ID := args[2]
	err := i.Stub.DelState(ID)
	if err != nil {
		return response.error("Failed to delete state:" + err.Error())
	}

	return response.success("DELETE")
}

func (i *DatasChaincode) get(args []string) pb.Response {

	var buffer bytes.Buffer

	ID := args[2]
	datasBytes, err := i.Stub.GetState(ID)

	if err != nil {
		return response.error("Failed to get data:" + err.Error())
	} else if datasBytes == nil {
		return response.error("Data does not exist")
	}
	buffer.Write(datasBytes)
	log.Printf("DatasBytes = <%v>", buffer.String())
	return response.success(buffer.String()) //TODO: return []byte
}

func main() {
	err := shim.Start(new(DatasChaincode))
	if err != nil {
		log.Printf("Error starting DatasChaincode: %s", err)
	}
}
