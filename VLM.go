package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

/* Vehicle Declaration */
type Vehicle struct {
	Make          string `json:"make"`
	Model         string `json:"model"`
	Owner         string `json:"owner"`
	Color         string `json:"color"`
	ChasisNumber  int `json:"chasisnumber"`
  EngineNumber  string `json:"enginenumber"`
	Status        string `json:"status"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "getVehicleDetails" {
		return s.getVehicleDetails(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createVehicle" {
		return s.createVehicle(APIstub, args)
	} else if function == "getAllVehicles" {
		return s.getAllVehicles(APIstub)
	} else if function == "changeOwnerShip" {
		return s.changeOwnerShip(APIstub, args)
	} else if function == "getHistoryForVehicle" {
		  return s.getHistoryForVehicle(APIstub, args)
  }

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) getVehicleDetails(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	vehicleAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(vehicleAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	vehicles := []Vehicle{
		Vehicle{Make: "Honda", Model: "CBR150R", Owner: "Jeff", Color:"Blue", ChasisNumber: 152365111, EngineNumber: "TXG66053729"},
		Vehicle{Make: "Hero", Model: "Glamour",  Owner: "Smith", Color:"Orange", ChasisNumber: 5465956, EngineNumber: "OTXG6685512"},
		Vehicle{Make: "Mahindra", Model: "Rodeo", Owner: "Bob", Color:"Red", ChasisNumber: 168556546, EngineNumber: "BJYTG6855129"},
		Vehicle{Make: "TVS", Model: "Sport", Owner: "John", Color:"Blue", ChasisNumber: 15562251, EngineNumber: "KNLJ89523144"},
		Vehicle{Make: "Honda", Model: "Activa", Owner: "Sandhya", Color: "Red", ChasisNumber: 567612, EngineNumber: "TTYY5100C"},
		Vehicle{Make: "Vespa", Model: "Trendy", Owner: "Lakshmi", Color: "Brown", ChasisNumber: 3459667, EngineNumber: "XG66150CA"},
		Vehicle{Make: "Bajaj", Model: "Axv", Owner: "Satwik",Color: "Blue", ChasisNumber: 1358679, EngineNumber: "XG6605150CB"},
		Vehicle{Make: "HeroHonda", Model: "Infinit", Color: "Grey",ChasisNumber: 1456945,EngineNumber: "XKUK56556J", Owner: "Prabhu"},
		Vehicle{Make: "Yamaha", Model: "RX100", Color: "Metallic Blue",ChasisNumber: 1793568,EngineNumber: "UNNF96572U", Owner: "Syed"},
		Vehicle{Make: "TVS", Model: "Victor", Color: "Blue",ChasisNumber: 2366312,EngineNumber: "MKITH9895", Owner: "Anand"},
		Vehicle{Make: "Honda" , Model: "Unicorn", Color: "Blue",ChasisNumber: 9347521,EngineNumber: "POI89466Q", Owner: "Seshu"},
		Vehicle{Make: "Honda", Model: "CBShine", Color: "Moonsoon Gray",ChasisNumber: 12335245,EngineNumber: "TRFN6538I", Owner: "Vani"},
	}

	i := 0
	for i < len(vehicles) {
		fmt.Println("i is ", i)
    chasisnumber := strconv.Itoa(vehicles[i].ChasisNumber);
		vehicleAsBytes, _ := json.Marshal(vehicles[i])
		APIstub.PutState(chasisnumber, vehicleAsBytes)
		fmt.Println("Added", vehicles[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createVehicle(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

  chasisnumber,_ := strconv.Atoi(args[4]);
	vehicle := Vehicle{Make: args[0], Model: args[1], Owner: args[2], Color: args[3], ChasisNumber: chasisnumber, EngineNumber: args[5]}

	vehicleAsBytes, _ := json.Marshal(vehicle)
	APIstub.PutState(args[4], vehicleAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) getAllVehicles(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "0"
	endKey := "999999999999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getAllVehicles:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeOwnerShip(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	vehicleAsBytes, _ := APIstub.GetState(args[0])
	vehicle := Vehicle{}

	json.Unmarshal(vehicleAsBytes, &vehicle)
	vehicle.Owner = args[1]

	vehicleAsBytes, _ = json.Marshal(vehicle)
	APIstub.PutState(args[0], vehicleAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) getHistoryForVehicle(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	chasisnumber := args[0]

	fmt.Printf("- start getHistoryForVehicle: %s\n", chasisnumber)

	resultsIterator, err := stub.GetHistoryForKey(chasisnumber)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForVehicle returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
