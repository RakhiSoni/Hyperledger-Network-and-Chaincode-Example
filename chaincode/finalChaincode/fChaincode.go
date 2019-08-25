package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/chaincode/shim/ext/cid"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type productDetails struct {
	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Name       string `json:"name"`    //the fieldtags are needed to keep case from bouncing around
	Price      int    `json:"price"`
}

type product struct {
	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Name       string `json:"name"`    //the fieldtags are needed to keep case from bouncing around
	Color      string `json:"color"`
	Owner      string `json:"owner"`
	Price      int    `json:"price"`
}

type org struct {
	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Name       string `json:"name"`    //the fieldtags are needed to keep case from bouncing around
	Desc       string `json:"desc"`
	Size       int    `json:"size"`
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	// fmt.Println("invoke is running " + function)

	// Handle different functions
	switch function {
	case "addProduct":
		//create a new marble
		return t.addProduct(stub, args)
	case "readProduct":
		//read a marble
		return t.readProduct(stub, args)
	case "addOrgDetails":
		//read a marble
		return t.addOrgDetails(stub, args)
	case "readPrivateDetails":
		//read a marble
		return t.readPrivateDetails(stub, args)
	default:
		//error
		fmt.Println("invoke did not find func: " + function)
		return shim.Error("Received unknown function invocation")
	}
}

// ============================================================
// initMarble - create a new marble, store into chaincode state
// ============================================================
func (t *SimpleChaincode) addProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	type productTransientInput struct {
		ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
		Name       string `json:"name"`    //the fieldtags are needed to keep case from bouncing around
		Color      string `json:"color"`
		Owner      string `json:"owner"`
		Price      int    `json:"price"`
	}

	// ==== Input sanitation ====
	fmt.Println("- start init product")
	fmt.Println("- Arguments --> ", args)

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments.")
	}

	transMap, err := stub.GetTransient()
	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}
	fmt.Println("- transMap   --> ", transMap)

	fmt.Println("---> TransMap ERROR Reason   --> ")
	fmt.Println(transMap["product"])
	fmt.Println("------- END   --> ")

	if _, ok := transMap["product"]; !ok {
		return shim.Error("product must be a key in the transient map")
	}

	fmt.Println("- Length transMap    --> " + string(len(transMap["product"])))

	if len(transMap["product"]) == 0 {
		return shim.Error("product value in the transient map must be a non-empty JSON string")
	}

	var productInput productTransientInput
	err = json.Unmarshal(transMap["product"], &productInput)

	fmt.Println("-json.Unmarshal    --> ", err)
	fmt.Println("---------")
	fmt.Println(productInput)
	fmt.Println("---------")

	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(transMap["product"]))
	}

	if len(productInput.Name) == 0 {
		return shim.Error("name field must be a non-empty string")
	}
	if len(productInput.Color) == 0 {
		return shim.Error("color field must be a non-empty string")
	}
	if len(productInput.Owner) == 0 {
		return shim.Error("owner field must be a non-empty string")
	}
	if productInput.Price <= 0 {
		return shim.Error("price field must be a positive integer")
	}

	// Get the client ID object
	id, err := cid.New(stub)
	if err != nil {
		return shim.Error(err.Error())
	}
	mspid, err := id.GetMSPID()
	if err != nil {
		return shim.Error(err.Error())
	}

	if mspid == "ManufacturerMSP" {
		// ==== Check if product already exists ====
		productAsBytes, err := stub.GetPrivateData("collectionProducts", productInput.Name)
		if err != nil {
			return shim.Error("Failed to get Product: " + err.Error())
		} else if productAsBytes != nil {
			fmt.Println("This Product already exists: " + productInput.Name)
			return shim.Error("This Product already exists: " + productInput.Name)
		}

		// ==== Create product object, marshal to JSON, and save to state ====
		product := &product{
			ObjectType: "product",
			Name:       productInput.Name,
			Color:      productInput.Color,
			Owner:      productInput.Owner,
			Price:      productInput.Price,
		}
		productJSONasBytes, err := json.Marshal(product)
		if err != nil {
			return shim.Error(err.Error())
		}
		err = stub.PutPrivateData("collectionProducts", productInput.Name, productJSONasBytes)

		if err != nil {
			return shim.Error(err.Error())
		}
	} else {
		fmt.Println("- IN Else ManufacturerMSP")
		shim.Error("Wrong MSP")
	}

	fmt.Println("mspid: " + mspid)

	// === Save product to state ===
	fmt.Println("- end Add Product")
	return shim.Success(nil)
}

func (t *SimpleChaincode) addOrgDetails(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var collectionName string

	type orgTransientInput struct {
		ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
		Name       string `json:"name"`    //the fieldtags are needed to keep case from bouncing around
		Desc       string `json:"desc"`
		Size       int    `json:"size"`
	}

	// ==== Input sanitation ====
	fmt.Println("- start init product")
	fmt.Println("- Arguments --> ", args)

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments.")
	}

	transMap, err := stub.GetTransient()

	fmt.Println("- transMap --> ", transMap)

	if err != nil {
		return shim.Error("Error getting transient: " + err.Error())
	}

	fmt.Println("- transMap[] --> ", transMap["orgDetails"])

	if _, ok := transMap["orgDetails"]; !ok {
		return shim.Error("orgDetails must be a key in the transient map")
	}

	if len(transMap["orgDetails"]) == 0 {
		return shim.Error("orgDetails value in the transient map must be a non-empty JSON string")
	}

	var orgInput orgTransientInput
	err = json.Unmarshal(transMap["orgDetails"], &orgInput)

	fmt.Println("-json.Unmarshal ----> --> ", err)

	if err != nil {
		return shim.Error("Failed to decode JSON of: " + string(transMap["orgDetails"]))
	}

	if len(orgInput.Name) == 0 {
		return shim.Error("name field must be a non-empty string")
	}
	if len(orgInput.Desc) == 0 {
		return shim.Error("desc field must be a non-empty string")
	}

	if orgInput.Size <= 0 {
		return shim.Error("size field must be a positive integer")
	}

	//  Get ABAC Attribute
	value, found, err := cid.GetAttributeValue(stub, "pID")

	if err != nil {
		return shim.Error("could not do abac")
	}

	fmt.Println("value is ", value)
	fmt.Println("found is ", found)

	if value == "GRP1" {
		collectionName = "collection1PrivateProducts"
		fmt.Println("found is ", found)
	} else if value == "GRP1" {
		collectionName = "collection2PrivateProducts"
		fmt.Println("found is ", found)
	} else {
		return shim.Error("Not a valid user")
	}

	org := &org{
		ObjectType: "org",
		Name:       orgInput.Name,
		Desc:       orgInput.Desc,
		Size:       orgInput.Size,
	}
	orgJSONasBytes, err := json.Marshal(org)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutPrivateData(collectionName, orgInput.Name, orgJSONasBytes)

	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save product to state ===
	fmt.Println("- end Add Details")
	return shim.Success(nil)
}

// ===============================================
// readMarble - read a marble from chaincode state
// ===============================================
func (t *SimpleChaincode) readProduct(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the product to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetPrivateData("collectionProducts", name) //get the product from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Marble does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

func (t *SimpleChaincode) readPrivateDetails(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error
	var collectionName string

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the Orgamization to query")
	}
	name = args[0]

	value, found, err := cid.GetAttributeValue(stub, "pID")

	if value == "GRP1" {
		collectionName = "collection1PrivateProducts"
		fmt.Println("found is ", found)
	} else if value == "GRP1" {
		collectionName = "collection2PrivateProducts"
		fmt.Println("found is ", found)
	} else {
		return shim.Error("Not a valid user")
	}

	valAsbytes, err := stub.GetPrivateData(collectionName, name) //get the marble private details from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get private details for " + name + ": " + err.Error() + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Marble private details does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}
