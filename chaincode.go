package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
	"os"
	"time"
)

var logger = shim.NewLogger("fabric-boilerplate")
//==============================================================================================================================
//	 Structure Definitions
//==============================================================================================================================
//	SimpleChaincode - A blank struct for use with Shim (An IBM Blockchain included go file used for get/put state
//					  and other IBM Blockchain functions)
//==============================================================================================================================
type SimpleChaincode struct {
}

type ECertResponse struct {
	OK string `json:"OK"`
}



type User struct {
	UserId       string   `json:"userId"` //Same username as on certificate in CA
	Salt         string   `json:"salt"`
	Hash         string   `json:"hash"`
	FirstName    string   `json:"firstName"`
	LastName     string   `json:"lastName"`
	Things       []string `json:"things"` //Array of thing IDs
	Address      string   `json:"address"`
	PhoneNumber  string   `json:"phoneNumber"`
	EmailAddress string   `json:"emailAddress"`
}

type BrokerageRequest struct {
	RequestID      			string   `json:"RequestID"` //Unique ID
	Submitter         		string   `json:"Submitter"`
	Approver         		string   `json:"Approver"`
	Documents    			string `json:"Documents"`
	PersonalDetails    		string `json:"PersonalDetails"`
	KYCDetails       	  	string `json:"KYCDetails"` //Array of thing IDs
	Status      			string   `json:"Status"`
	DocValidationReport  	string   `json:"DocValidationReport"`
	FacialValidation 		string   `json:"FacialValidation"`
	Video					string	 `json:"Video"`
	TimeStamps				string `json:"TimeStamps"`
	Meeting 			    string `json:"Meeting"`
	UpdateType              string `json:"UpdateType"`
}

type KyckUser struct {
	UserId       string   `json:"userId"` //Same username as on certificate in CA
	FirstName    string   `json:"firstName"`
	LastName     string   `json:"lastName"`
	Address      string   `json:"address"`
	PhoneNumber  string   `json:"phoneNumber"`
	Documents    			[]byte 	`json:"Documents"`
	PersonalDetails    		[]byte 	`json:"PersonalDetails"`
	KYCDetails       	  	[]byte 	`json:"KYCDetails"` 
	DocValidationReport  	[]byte  `json:"DocValidationReport"`
	TimeStamp   			string
	userType				string
	Rights					[]byte
}
/**** Accessor can be Broker, Govt agency, Regulator, etc ****/
type KyckAccessor struct {
	AccessorId      string   `json:"AccessorId"` //Same username as on certificate in CA
	Name     		string   `json:"Name"`
	Address  		string   `json:"Address"`
	Email    		[]string `json:"Email"`
	Phone     		string   `json:"Phone"`
	userType		string
}


type BrokerageResponse struct {
	status string	 
}

type BrokerageRequestTimeStamp struct {
	Submit 					string 
	MeetingConfirmation		string
	FinalStatus				string

}
type Thing struct {
	Id          string `json:"id"`
	Description string `json:"description"`
}

//=================================================================================================================================
//  Index collections - In order to create new IDs dynamically and in progressive sorting
//  Example:
//    signaturesAsBytes, err := stub.GetState(signaturesIndexStr)
//    if err != nil { return nil, errors.New("Failed to get Signatures Index") }
//    fmt.Println("Signature index retrieved")
//
//    // Unmarshal the signatures index
//    var signaturesIndex []string
//    json.Unmarshal(signaturesAsBytes, &signaturesIndex)
//    fmt.Println("Signature index unmarshalled")
//
//    // Create new id for the signature
//    var newSignatureId string
//    newSignatureId = "sg" + strconv.Itoa(len(signaturesIndex) + 1)
//
//    // append the new signature to the index
//    signaturesIndex = append(signaturesIndex, newSignatureId)
//    jsonAsBytes, _ := json.Marshal(signaturesIndex)
//    err = stub.PutState(signaturesIndexStr, jsonAsBytes)
//    if err != nil { return nil, errors.New("Error storing new signaturesIndex into ledger") }
//=================================================================================================================================
var usersIndexStr = "_users"
var thingsIndexStr = "_things"
var applicationIndexStr ="_applications"

var indexes = []string{usersIndexStr, thingsIndexStr,applicationIndexStr}

//==============================================================================================================================
//	Invoke - Called on chaincode invoke. Takes a function name passed and calls that function. Passes the
//  		 initial arguments passed are passed on to the called function.
//==============================================================================================================================

func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	logger.Infof("Invoke is running " + function)

	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "reset_indexes" {
		return t.reset_indexes(stub, args)
	} else if function == "add_user" {
		return t.add_user(stub, args)
	} else if function == "add_thing" {
		return t.add_thing(stub, args)
	}else if function == "add_resource"{
        return t.add_resource(stub, args)
    }else if function == "create_brokerage_request" {	//Create a new application
		return t.create_brokerage_request(stub, args[0])
	}else if function == "update_brokerage_application" {
		//updateType, jsonData, brokerageRequestId - input arguments
		return t.update_brokerage_application(stub, args[0])
	}

	return nil, errors.New("Received unknown invoke function name")
}

//=================================================================================================================================
//	Query - Called on chaincode query. Takes a function name passed and calls that function. Passes the
//  		initial arguments passed are passed on to the called function.
//=================================================================================================================================
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	logger.Infof("Query is running " + function)

	if function == "get_user" {
		return t.get_user(stub, args[1])
	} else if function == "get_thing" {
		return t.get_thing(stub, args)
	} else if function == "get_all_things" {
		return t.get_all_things(stub, args)
	} else if function == "authenticate" {
		return t.authenticate(stub, args)
	}else if function == "get_resource"{
        return t.get_resource(stub, args)
    }else if function == "get_brokerage_request"{
        return t.get_brokerage_request(stub, args[0])
    }else if function == "get_all_brokerage_requests"{
        return t.get_all_brokerage_requests(stub, args[0])
    }

	return nil, errors.New("Received unknown query function name")
}

//=================================================================================================================================
//  Main - main - Starts up the chaincode
//=================================================================================================================================

func main() {

	// LogDebug, LogInfo, LogNotice, LogWarning, LogError, LogCritical (Default: LogDebug)
	logger.SetLevel(shim.LogInfo)

	logLevel, _ := shim.LogLevel(os.Getenv("SHIM_LOGGING_LEVEL"))
	shim.SetLoggingLevel(logLevel)

	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting SimpleChaincode: %s", err)
	}
}

//==============================================================================================================================
//  Init Function - Called when the user deploys the chaincode
//==============================================================================================================================

func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	var empty []string
	jsonAsBytes, _ := json.Marshal(empty)								//marshal an empty array of strings to clear the index
	var err = stub.PutState(applicationIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	//Create a table to store all the Brokerage Applications submitted
	err = stub.CreateTable("BrokerageRequests", []*shim.ColumnDefinition{
			&shim.ColumnDefinition{Name: "RequestID"			, Type:shim.ColumnDefinition_STRING,	Key: true},
			&shim.ColumnDefinition{Name: "Submitter"			, Type:shim.ColumnDefinition_STRING,	Key:false},
			&shim.ColumnDefinition{Name: "Approver"				, Type:shim.ColumnDefinition_STRING, 	Key:false},
			&shim.ColumnDefinition{Name: "Documents"			, Type:shim.ColumnDefinition_BYTES, 	Key:false},
			&shim.ColumnDefinition{Name: "PersonalDetails"		, Type:shim.ColumnDefinition_BYTES, 	Key:false},
			&shim.ColumnDefinition{Name: "KYCDetails"			, Type:shim.ColumnDefinition_BYTES, 	Key:false},
			&shim.ColumnDefinition{Name: "Status"				, Type:shim.ColumnDefinition_STRING, 	Key:false},
			&shim.ColumnDefinition{Name: "DocValidationReport"	, Type:shim.ColumnDefinition_BYTES, 	Key:false},
			&shim.ColumnDefinition{Name: "FacialValidation"		, Type:shim.ColumnDefinition_BYTES, 	Key:false},
			&shim.ColumnDefinition{Name: "VideoRecording"		, Type:shim.ColumnDefinition_BYTES, 	Key:false},
			&shim.ColumnDefinition{Name: "TimeStamps"		    , Type:shim.ColumnDefinition_BYTES, 	Key:false},
			&shim.ColumnDefinition{Name: "Meeting"		        , Type:shim.ColumnDefinition_STRING, 	Key:false},
	})
	if err != nil{ return nil, errors.New( "Failed creating Brokerage Requests Table")}

	return nil, nil
}

//==============================================================================================================================
//  Utility Functions
//==============================================================================================================================

// "create":  true -> create new ID, false -> append the id
func append_id(stub *shim.ChaincodeStub, indexStr string, id string, create bool) ([]byte, error) {

	indexAsBytes, err := stub.GetState(indexStr)
	if err != nil {
		return nil, errors.New("Failed to get " + indexStr)
	}

	// Unmarshal the index
	var tmpIndex []string
	json.Unmarshal(indexAsBytes, &tmpIndex)

	// Create new id
	var newId = id
	if create {
		newId += strconv.Itoa(len(tmpIndex) + 1)
	}

	// append the new id to the index
	tmpIndex = append(tmpIndex, newId)

	jsonAsBytes, _ := json.Marshal(tmpIndex)
	err = stub.PutState(indexStr, jsonAsBytes)
	if err != nil {
		return nil, errors.New("Error storing new " + indexStr + " into ledger")
	}

	return []byte(newId), nil

}

//==============================================================================================================================
//  Invoke Functions
//==============================================================================================================================
func (t *SimpleChaincode) reset_indexes(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	for _, i := range indexes {
		// Marshal the index
		var emptyIndex []string

		empty, err := json.Marshal(emptyIndex)
		if err != nil {
			return nil, errors.New("Error marshalling")
		}
		err = stub.PutState(i, empty);

		if err != nil {
			return nil, errors.New("Error deleting index")
		}
		logger.Infof("Delete with success from ledger: " + i)
	}
	return nil, nil
}

func (t *SimpleChaincode) add_user(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	//Args
	//			0				1
	//		  index		user JSON object (as string)

	id, err := append_id(stub, usersIndexStr, args[0], false)
	if err != nil {
		return nil, errors.New("Error creating new id for user " + args[0])
	}

	err = stub.PutState(string(id), []byte(args[1]))
	if err != nil {
		return nil, errors.New("Error putting user data on ledger")
	}

	return nil, nil
}

func (t *SimpleChaincode) add_thing(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	// args
	// 		0			1
	//	   index	   thing JSON object (as string)

	id, err := append_id(stub, thingsIndexStr, args[0], false)
	if err != nil {
		return nil, errors.New("Error creating new id for thing " + args[0])
	}

	err = stub.PutState(string(id), []byte(args[1]))
	if err != nil {
		return nil, errors.New("Error putting thing data on ledger")
	}

	return nil, nil

}

func (t *SimpleChaincode) create_brokerage_request(stub *shim.ChaincodeStub, jsonData string) ([]byte, error) {

	fmt.Println("Input request object :: " + jsonData)

	/**** Convert the incoming arguments from json to bytearray ****/
	var bytesArray = []byte(jsonData)

	/**** Copy the incoming json data to a struct b ****/
	var b BrokerageRequest;
	json.Unmarshal(bytesArray, &b)
	fmt.Println("B value :: " + b.RequestID)

	/**** Create an object for inserting TimeStamps ****/
	timeStampJson := []byte(t.get_current_time())

	/****  Insert the details of the Brokerage application into a new row in the Table structure ****/
	fmt.Println("Inserting row now")
	tx, err := stub.InsertRow( "BrokerageRequests" , shim.Row{
				Columns: []*shim.Column{
					&shim.Column{Value: &shim.Column_String_{String_: b.RequestID}},
					&shim.Column{Value: &shim.Column_String_{String_: b.Submitter}},
					&shim.Column{Value: &shim.Column_String_{String_: b.Approver}},
					&shim.Column{Value: &shim.Column_Bytes	{Bytes: []byte(b.Documents)}},
					&shim.Column{Value: &shim.Column_Bytes	{Bytes: []byte(b.PersonalDetails)}},
					&shim.Column{Value: &shim.Column_Bytes	{Bytes: []byte(b.KYCDetails)}},
					&shim.Column{Value: &shim.Column_String_{String_: b.Status}},
					&shim.Column{Value: &shim.Column_Bytes{Bytes: []byte(b.DocValidationReport)}},
					&shim.Column{Value: &shim.Column_Bytes{Bytes: []byte(b.FacialValidation)}},
					&shim.Column{Value: &shim.Column_Bytes{Bytes: []byte(b.Video)}},
					&shim.Column{Value: &shim.Column_Bytes{Bytes: []byte(timeStampJson)}},
					&shim.Column{Value: &shim.Column_String_{String_: b.Meeting}},
				},
			})
	fmt.Println("ERROR :: " , err)

	if err != nil {
		fmt.Println("Error while updating record :: ", b.RequestID)
	}else{
		fmt.Println("Response from Insert Row ::", tx)
	}
	
	/*Returning nil response for now*/
	return nil, nil
}

func (t *SimpleChaincode) update_brokerage_application(stub *shim.ChaincodeStub, jsonData string) ([]byte, error) {

	fmt.Println("Input request object :: " + jsonData)

	/****New Data to be written****/
	var bytesArray = []byte(jsonData)

	var inputBrokerageRequest BrokerageRequest

	fmt.Println("Unmarshalling data to BrokerageRequest")
	json.Unmarshal(bytesArray, &inputBrokerageRequest)

	/****First get the data stored****/
	fmt.Println("Fetching data from table BrokerageRequests")
	brokerageRequestRow := t.fetch_from_table(stub, inputBrokerageRequest.RequestID)

	/****Convert to local Struct here****/
	fmt.Println("Converting Row to Struct data from table BrokerageRequests")
	brokerageRequest := t.getStructFromRow(brokerageRequestRow)	
	
	var timeStampJson []byte

	updateType := inputBrokerageRequest.UpdateType

	fmt.Println("UpdateType is :: " + updateType)

	if updateType == "MEETING" {
		brokerageRequest.Meeting = jsonData
		timeStampJson = []byte(t.get_current_time())
	}else if updateType == "VIDEO" {
		brokerageRequest.Video = jsonData
	}else if updateType == "STATUS"{
		brokerageRequest.Status = jsonData
	}
	
	fmt.Println("Replacing row now...")
	/**** Store the data ****/
		tx,err := stub.ReplaceRow( "BrokerageRequests" , shim.Row{
				Columns: []*shim.Column{
					&shim.Column{Value: &shim.Column_String_{String_: brokerageRequest.RequestID}},
					&shim.Column{Value: &shim.Column_String_{String_: brokerageRequest.Submitter}},
					&shim.Column{Value: &shim.Column_String_{String_: brokerageRequest.Approver}},
					&shim.Column{Value: &shim.Column_Bytes	{Bytes: []byte(brokerageRequest.Documents)}},
					&shim.Column{Value: &shim.Column_Bytes	{Bytes: []byte(brokerageRequest.PersonalDetails)}},
					&shim.Column{Value: &shim.Column_Bytes	{Bytes: []byte(brokerageRequest.KYCDetails)}},
					&shim.Column{Value: &shim.Column_String_{String_: brokerageRequest.Status}},
					&shim.Column{Value: &shim.Column_Bytes{Bytes: []byte(brokerageRequest.DocValidationReport)}},
					&shim.Column{Value: &shim.Column_Bytes{Bytes: []byte(brokerageRequest.FacialValidation)}},
					&shim.Column{Value: &shim.Column_Bytes{Bytes: []byte(brokerageRequest.Video)}},
					&shim.Column{Value: &shim.Column_Bytes{Bytes: []byte(timeStampJson)}},
					&shim.Column{Value: &shim.Column_String_{String_: brokerageRequest.Meeting}},
				},
		})
		fmt.Println("Error ::: ",err)

		if err != nil {
			fmt.Println("Error while updating record :: ", brokerageRequest.RequestID)
		}else{
			fmt.Println("Response from Insert Row ::", tx)
		}
	
	return timeStampJson, nil
}

func (t *SimpleChaincode) toJson()(string){
	var returnValue string;
	
	return returnValue
}

func (t *SimpleChaincode) get_current_time() ([]byte) {
	var timeStampObject BrokerageRequestTimeStamp
	timenow := time.Now()
	fmt.Printf("Time.now %s", timenow)
	timeStampObject.Submit = timenow.UTC().Format(time.UnixDate);
    timeStampJson, _ := json.Marshal(timeStampObject)

	return timeStampJson
}



/*
	Return the struct from the table.
*/
func(t *SimpleChaincode) getStructFromRow(row shim.Row)(BrokerageRequest){
	
	var brokerageRequest BrokerageRequest
	for index := range row.Columns {
		column := row.Columns[index]
		if index == 0 {
			brokerageRequest.RequestID = column.GetString_()
		}else if index == 1 {
			brokerageRequest.Submitter = column.GetString_()
		}else if index == 2 {
			brokerageRequest.Approver = column.GetString_()
		}else if index == 3 {
			brokerageRequest.Documents = column.GetString_()
		}else if index == 4 {
			brokerageRequest.PersonalDetails = column.GetString_()
		}else if index == 5 {
			brokerageRequest.KYCDetails = column.GetString_()
		}else if index == 6 {
			brokerageRequest.Status = column.GetString_()
		}else if index == 7 {
			brokerageRequest.DocValidationReport = column.GetString_()
		}else if index == 8 {
			brokerageRequest.FacialValidation = column.GetString_()
		}else if index == 9 {
			brokerageRequest.Video = column.GetString_()
		}else if index == 10 {
			brokerageRequest.TimeStamps = column.GetString_()
		}else if index == 11 {
			brokerageRequest.Meeting = column.GetString_()
		}
		index ++
	}
	return brokerageRequest
}

/*This function helps in getting the data stored from local database*/
func (t *SimpleChaincode) fetch_from_table(stub *shim.ChaincodeStub, requestId string)(shim.Row){
	fmt.Println("Fetching Row for RequestID :: " + requestId)
	var columns []shim.Column
	row := shim.Row{
		Columns: []*shim.Column{
					&shim.Column{Value: &shim.Column_String_{String_: requestId}},
					//&shim.Column{Value: &shim.Column_String_{String_: submitter}},
		},
	}
	queryCol := shim.Column{Value: &shim.Column_String_{String_: requestId}}
	columns = append(columns, queryCol)
	row.Columns[0].GetString_()
	//row.Columns[1].GetString_()
	row, _ = stub.GetRow("BrokerageRequests", columns)
	fmt.Println("ROW :: ", row)
	uid:= row.Columns[0].GetString_()
	submitter := row.Columns[1].GetString_()
	fmt.Println("UID is " + uid)
	fmt.Println("Submitter is" + submitter)
	return row
}



//==============================================================================================================================
//		Query Functions
//==============================================================================================================================

func (t *SimpleChaincode) get_user(stub *shim.ChaincodeStub, userID string) ([]byte, error) {

	bytes, err := stub.GetState(userID)

	if err != nil {
		return nil, errors.New("Could not retrieve information for this user")
	}

	return bytes, nil

}

func (t *SimpleChaincode) get_thing(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	//Args
	//			0
	//		thingID

	bytes, err := stub.GetState(args[0])

	if err != nil {
		return nil, errors.New("Error getting from ledger")
	}

	return bytes, nil

}

func (t *SimpleChaincode) get_all_things(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	indexAsBytes, err := stub.GetState(thingsIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get " + thingsIndexStr)
	}

	// Unmarshal the index
	var thingsIndex []string
	json.Unmarshal(indexAsBytes, &thingsIndex)

	var things []Thing
	for _, thing := range thingsIndex {

		bytes, err := stub.GetState(thing)
		if err != nil {
			return nil, errors.New("Unable to get thing with ID: " + thing)
		}

		var t Thing
		json.Unmarshal(bytes, &t)
		things = append(things, t)
	}

	thingsAsJsonBytes, _ := json.Marshal(things)
	if err != nil {
		return nil, errors.New("Could not convert things to JSON ")
	}

	return thingsAsJsonBytes, nil
}

func (t *SimpleChaincode) authenticate(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	// Args
	//	0		1
	//	userId	password

	var u User

	username := args[0]

	user, err := t.get_user(stub, username)

	// If user can not be found in ledgerstore, return authenticated false
	if err != nil {
		return []byte(`{ "authenticated": false }`), nil
	}

	//Check if the user is an employee, if not return error message
	err = json.Unmarshal(user, &u)
	if err != nil {
		return []byte(`{ "authenticated": false}`), nil
	}

	// Marshal the user object
	userAsBytes, err := json.Marshal(u)
	if err != nil {
		return []byte(`{ "authenticated": false}`), nil
	}

	// Return authenticated true, and include the user object
	str := `{ "authenticated": true, "user": ` + string(userAsBytes) + `  }`

	return []byte(str), nil
}


func (t *SimpleChaincode) add_resource(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

    fmt.Println("Chaincode running AddResource()")
	fmt.Println(args[0])
    id := args[0] + args[1]
	fmt.Println("KEY value for invoke - add resource -->> " +id)
    err := stub.PutState(string(id), []byte(args[2]))
	if err != nil {
		return nil, errors.New("Error putting resource data on ledger")
	}
    return nil, nil   

}

func (t *SimpleChaincode) get_resource(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

    fmt.Println("Chaincode running GetResource()")
	fmt.Println(args[0])
	fmt.Println(args[1])
    id := args[0] + args[1]
	fmt.Println("ID = " + id)
    path, err := stub.GetState(string(id))
	fmt.Println("PATH IS  = " + string(path))
	if err != nil {
		return nil, errors.New("Error getting resource data from ledger")
	}
    return path, nil

    
}

func (t *SimpleChaincode) get_brokerage_request(stub *shim.ChaincodeStub, requestId string) ([]byte, error) {
 	fmt.Println("Chaincode running get_brokerage_request()")
	 row := t.fetch_from_table(stub, requestId)
	 structure := t.getStructFromRow(row)
	 bytesArray,_ := json.Marshal(structure)
	 return bytesArray,nil
}

func (t *SimpleChaincode) get_all_brokerage_requests(stub *shim.ChaincodeStub, requestId string) ([]byte, error) {
 	fmt.Println("Chaincode running get_brokerage_request()")
	 row := t.fetch_from_table(stub, requestId)
	 structure := t.getStructFromRow(row)
	 bytesArray,_ := json.Marshal(structure)
	 return bytesArray,nil
}
