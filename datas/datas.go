package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"crypto/sha256"
	"strings"
	"time"


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

type Response struct {
	Status    string
	Message   string
	Txn       string
	Hash      string
	Domain    string
	Id        string
	AssetInfo AssetInfo
}

type AssetInfo struct {
	Key         string
	Status      string
	Description string
}

const (
	OK             = 200
	ERRORTHRESHOLD = 400
	ERROR          = 500
)


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


func newResponse(txn string, domain string, id string) *Response {
	return &Response{Status: "Ok", Txn: txn, Domain: domain, Id: id}
}

func (r *Response) error(msg string) pb.Response {
	r.Status = "ERROR"
	r.Message = msg
	responseBytes, _ := json.Marshal(r)
	return pb.Response{
		Status:  ERROR,
		Message: msg,
		Payload: responseBytes,
	}
}

func (r *Response) success(msg string) pb.Response {
	responseBytes, _ := json.Marshal(r)
	return pb.Response{
		Status:  OK,
		Message: msg,
		Payload: responseBytes,
	}
}

const dateFormat = "2006-01-02"

type JSONTime struct {
	time.Time
}

func (jsonTime *JSONTime) UnmarshalJSON(p []byte) error {
	strInput := string(p)
	strInput = strings.Trim(strInput, `"`)
	t, err := time.Parse(dateFormat, strInput)
	if err != nil {
		return err
	}
	jsonTime.Time = t
	return nil
}

func generateHash(message []byte) []byte {
	hasher := sha256.New()
	hasher.Write(message)
	return hasher.Sum(nil)
}


func main() {
	err := shim.Start(new(DatasChaincode))
	if err != nil {
		log.Printf("Error starting DatasChaincode: %s", err)
	}
}
