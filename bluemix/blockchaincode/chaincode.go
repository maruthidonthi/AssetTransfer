package main

import(
	"errors"
	"fmt"
	"strconv"
	"encoding/json"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var assetIndexStr = "_assetindex" // name for the key/value that will store all fruites
var openTradesStr = "_opentrades"  //name for the key/value that will store all open openTrades

type Asset struct{
	AssetId string `json:"assetId"`
	Price int `json:"price"`
	UserId string `json:"userId"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	var Aval int
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	// Write the state to the ledger
	err = stub.PutState("abc", []byte(strconv.Itoa(Aval)))				//making a test var "abc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}

	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an emtpy array of strings to clear the index
	err = stub.PutState(assetIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
// ============================================================================================================================
// Run - Our entry point for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func (t *SimpleChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "delete" {										//deletes an entity from its state
		return t.Delete(stub, args)
	} else if function == "write" {											//writes a value to the chaincode state
		return t.Write(stub, args)
	} else if function == "init_asset" {									//create a new asset
		return t.init_asset(stub, args)
	} else if function == "set_user" {										//change owner of a asset
		return t.set_user(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {													//read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query")
}

// ============================================================================================================================
// Read - read a variable from chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var assetId, resp string
	var err error

	if len(args)!=1 {
		return nil, errors.New("incorrect number of arguments")
	}

	assetId = args[0]

	resAsbytes, err := stub.GetState(assetId)

	if err != nil {
		resp = "{\"Error\":\" Failed to get state for " + assetId + "\"}"
		return nil, errors.New(resp)
	}
	return resAsbytes, nil
}
// ============================================================================================================================
// Delete - remove a key/value pair from state
// ============================================================================================================================
func (t *SimpleChaincode) Delete(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	if len(args) !=1 {
		return nil, errors.New("incorrect arguments Expecting 1 argument")
	}

	assetId := args[0]
	//remove the asset from chaincode state
	err := stub.DelState(assetId)

	if err != nil{
		return nil, errors.New("failed to delete state")
	}

	// get the asset Index
	assetAsBytes, err := stub.GetState(assetIndexStr)

	if err != nil {
		return nil, errors.New("Failed to get asset index")
	}

	var assetIndex []string

	//un stringify it aka JSON.parse()
	//remove marble from index
	json.Unmarshal(assetAsBytes, &assetIndex)
	for i, val := range assetIndex{
			fmt.Println(strconv.Itoa(i)+ " - looking at "+ val +" for "+ assetId)
			if val == assetId{
				fmt.Println("found asset")
				assetIndex = append(assetIndex[:i], assetIndex[i+1:]...)
				for x:= range assetIndex{
					fmt.Println(string(x)+ "-" + assetIndex[x])
				}
				break
			}
		}
		jsonAsBytes, _ := json.Marshal(assetIndex)
		err = stub.PutState(assetIndexStr, jsonAsBytes)
		return nil, nil
}

// ============================================================================================================================
// Write - write variable into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) Write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var assetid, value string
	var err error

	if len(args) !=2 {
		return nil, errors.New("incorrect number of arguments Expecting 2")
	}

	assetid = args[0]
	value = args[1]

	err = stub.PutState(assetid, []byte(value))

	if err !=nil{
		return nil, err;
	}

	return nil, nil
}
// ============================================================================================================================
// Init Marble - create a new marble, store into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) init_asset(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error

	//		0						1								2
	//"assetId-aa"	"price-10"	"userid- donthi"

	if len(args) !=3 {
		return nil, errors.New("Incorrect number of arguments expecting 3")
	}
	 fmt.Println("- start intiating assets")

	 if len(args[0]) <= 0{
		 return nil, errors.New("1st argument must be a non-empty string")
	 }
	 if len(args[1]) <= 0{
		 return nil, errors.New("2nd argument must be a non-empty string")
	 }
	 if len(args[2]) <= 0{
		 return nil, errors.New("3rd argument must be a non-empty string")
	 }

	 price, err := strconv.Atoi(args[1])
	 if err != nil{
		 return nil, errors.New("2nd argument must be integer")
	 }

	 assetId := strings.ToLower(args[0])
	 userId := strings.ToLower(args[2])

	 //make json for asset struct0
	 str := `{"assetId": "` + assetId + `", "price": ` + strconv.Itoa(price) + `, "userId": "` + userId + `"}`
	 err = stub.PutState(args[0], []byte(str))
	 if err != nil{
		 return nil, err
	 }
	 // get the asset index
	 assetAsbytes, err := stub.GetState(assetIndexStr)
	 if err != nil{
		 return nil, errors.New("failed to get asset index")
	 }

	 var assetIndex []string
	 json.Unmarshal(assetAsbytes, &assetIndex)
	 // add asset to the index list
	 assetIndex = append(assetIndex, args[0])
	 fmt.Println("asset Index:", assetIndex)
	 jsonAsBytes, _ := json.Marshal(assetIndex)
	 err = stub.PutState(assetIndexStr, jsonAsBytes)
	  fmt.Println("end init asset method")

		return nil, nil
}

// ============================================================================================================================
// Set User Permission on Marble
// ============================================================================================================================
func (t *SimpleChaincode) set_user(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var err error

	// 			0    						1
	// "assetId"    "userId - bob"
	if len(args) < 2 {
		return nil, errors.New("Incorrect number of argumenrs, Expecting 2")
	}
	fmt.Println("- start set user")
	fmt.Println(args[0] + " - " + args[1])
	assetAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get thing")
	}
	res := Asset{}
	json.Unmarshal(assetAsBytes, &res)										//un stringify it aka JSON.parse()
	res.UserId = args[1]														//change the user

	jsonAsBytes, _ := json.Marshal(res)
	err = stub.PutState(args[0], jsonAsBytes)								//rewrite the asset with id as key
	if err != nil {
		return nil, err
	}

	fmt.Println("- end set user")
	return nil, nil
}
