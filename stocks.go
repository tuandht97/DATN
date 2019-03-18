package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ----- Stock ----- //
type Stock struct {
	ObjectType 	string        	`json:"docType"` 	// field for couchdb
	Id       	string          `json:"id"`
	Code       	string          `json:"code"`      	// mã
	Count      	int        		`json:"count"`		// số lượng chứng chỉ tạo ra
	Price       int         	`json:"price"`    	// giá một chứng chỉ
	Creator     UserInfo 		`json:"creator"`		// người tạo
}

// ----- User ----- //
type User struct {
	ObjectType 	string 			`json:"docType"`    // field for couchdb
	Id        	string 			`json:"id"`			
	Name   		string 			`json:"name"`		// tên
	Wallet    	[]Asset 		`json:"wallet"`		// ví
}

// ----- UserInfo ----- //
type UserInfo struct {
	Id         	string 			`json:"id"`		
	Name   		string 			`json:"name"`   	
}

// ----- Asset ----- //
type Asset struct {
	Id 			string 			`json:"id"`
	Code        string 			`json:"code"`		// mã chứng chỉ quỹ
	Count   	int 			`json:"count"`  	// số lượng
}

// ----- Trade ----- //
type Trade struct {
	ObjectType 	string 			`json:"docType"`    // field for couchdb
	Id			string 			`json:"id"`	
	Stock 		Asset			`json:"stock"`		// mã
	Seller		UserInfo		`json:"seller"`		// thông tin người bán
	Buyer		UserInfo		`json:"buyer"`		// thông tin người mua
	Time 		string 			`json:"time"`		// thời gian giao dịch
}

// Main
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode - %s", err)
	}
}

// Init - initialize the chaincode  
// Returns - shim.Success or error
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println(" ")
	fmt.Println("starting invoke, for - " + function)

	// Handle different functions
	if function == "init" {                    					// khởi tạo trạng thái
		return t.Init(stub)
	} else if function == "init_stock" {      					// tạo chứng chỉ quỹ
		return init_stock(stub, args)
	} else if function == "update_price" {      				// cập nhật giá chứng chỉ quỹ
		return update_price(stub, args)
	} else if function == "get_list_stock"{    					// xem toàn bộ mã chứng chỉ quỹ
		return get_list_stock(stub)
	} else if function == "init_user" {      					// tạo người dùng
		return init_user(stub, args)
	} else if function == "init_transaction"{   				// tạo giao dịch
		return init_transaction(stub, args)
	} else if function == "get_list_user"{    					// xem toàn bộ mã chứng chỉ quỹ
		return get_list_user(stub)
	} else if function == "get_list_transaction"{   			// xem toàn bộ danh sách các giao dịch
		return get_list_transaction(stub)
	} else if function == "get_list_transaction_by_user"{    	// danh sách các giao dịch của người dùng với id nhập vào
		return get_list_transaction_by_user(stub, args)
	} else if function == "get_list_user_have_stock_by_id"{    	// danh sách thông tin người dùng và số lượng mã người dùng đó có với mã có id nhập vào
		return get_list_user_have_stock_by_id(stub, args)
	}

	// error out
	fmt.Println("Received unknown invoke function name - " + function)
	return shim.Error("Received unknown invoke function name - '" + function + "'")
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Error("Unknown supported call - Query()")
}
