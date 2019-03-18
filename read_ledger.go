package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Get list stock 
func get_list_stock(stub shim.ChaincodeStubInterface) pb.Response {
	type ListStock struct {
		Stocks   []Stock   `json:"stocks"`
	}
	var listStock ListStock

	// ---- Get All Stock --- //
	resultsIterator, err := stub.GetStateByRange("s0", "s9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()
	
	for resultsIterator.HasNext() {
		aKeyValue, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryKeyAsStr := aKeyValue.Key
		queryValAsBytes := aKeyValue.Value
		fmt.Println("on stock id - ", queryKeyAsStr)
		var stock Stock
		json.Unmarshal(queryValAsBytes, &stock)                  	
		fmt.Println(stock)
		listStock.Stocks = append(listStock.Stocks, stock)   		

	}
	fmt.Println("stock array - ", listStock.Stocks)

	//change to array of bytes
	listStockAsBytes, _ := json.Marshal(listStock)              	
	return shim.Success(listStockAsBytes)
}

func get_list_user(stub shim.ChaincodeStubInterface) pb.Response {
	type ListUser struct {
		Users   []User   `json:"users"`
	}
	var listUser ListUser

	// ---- Get All user --- //
	usersIterator, err := stub.GetStateByRange("u0", "u9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer usersIterator.Close()
	
	for usersIterator.HasNext() {
		aKeyValue, err := usersIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryKeyAsStr := aKeyValue.Key
		queryValAsBytes := aKeyValue.Value
		fmt.Println("on user id - ", queryKeyAsStr)
		var user User
		json.Unmarshal(queryValAsBytes, &user)                  	
		listUser.Users = append(listUser.Users, user)   		
	}
	fmt.Println("user array - ", listUser.Users)

	//change to array of bytes
	listUserAsBytes, _ := json.Marshal(listUser)              	
	return shim.Success(listUserAsBytes)
}

func get_list_transaction(stub shim.ChaincodeStubInterface) pb.Response {
	type ListTrade struct {
		Trans   []Trade   `json:"transactions"`
	}
	var listTran ListTrade

	// ---- Get All user --- //
	tranIterator, err := stub.GetStateByRange("t0", "t9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer tranIterator.Close()
	
	for tranIterator.HasNext() {
		aKeyValue, err := tranIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryKeyAsStr := aKeyValue.Key
		queryValAsBytes := aKeyValue.Value
		fmt.Println("on transaction id - ", queryKeyAsStr)
		var transaction Trade
		json.Unmarshal(queryValAsBytes, &transaction)                  	
		listTran.Trans = append(listTran.Trans, transaction)   		
	}
	fmt.Println("transaction array - ", listTran.Trans)

	//change to array of bytes
	listTranAsBytes, _ := json.Marshal(listTran)              	
	return shim.Success(listTranAsBytes)
}

func get_list_transaction_by_user(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type ListTrade struct {
		Trans   []Trade   `json:"transactions"`
	}
	
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	user_id := args[0]

	_, err := get_user(stub, user_id)
	if err != nil {
		return shim.Error("This user does not exist - " + user_id)
	}

	var listTran ListTrade

	// ---- Get All transaction --- //
	tranIterator, err := stub.GetStateByRange("t0", "t9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer tranIterator.Close()
	
	for tranIterator.HasNext() {
		aKeyValue, err := tranIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryValAsBytes := aKeyValue.Value
		var transaction Trade
		json.Unmarshal(queryValAsBytes, &transaction)
		if (transaction.Buyer.Id == user_id || transaction.Seller.Id == user_id) {         	
			listTran.Trans = append(listTran.Trans, transaction)   
		}	
	}
	fmt.Println("transaction array of user-id " + user_id, listTran.Trans)

	//change to array of bytes
	listTranAsBytes, _ := json.Marshal(listTran)              	
	return shim.Success(listTranAsBytes)
}

func get_list_user_have_stock_by_id(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	type UserHaveStock struct {
		Id         	string 			`json:"id"`		
		Name   		string 			`json:"name"`
		Count 		int 			`json:"count"`
	}

	type ListUser struct {
		Users   []UserHaveStock   `json:"users"`
	}

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	stock_id := args[0]

	_, err := get_stock(stub, stock_id)
	if err != nil {
		return shim.Error("This stock does not exist - " + stock_id)
	}

	var listUser ListUser

	// ---- Get All transaction --- //
	userIterator, err := stub.GetStateByRange("u0", "u9999999999999999999")
	if err != nil {
		return shim.Error(err.Error())
	}
	defer userIterator.Close()
	
	for userIterator.HasNext() {
		aKeyValue, err := userIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		queryValAsBytes := aKeyValue.Value
		var user User
		var userHaveStock UserHaveStock
		json.Unmarshal(queryValAsBytes, &user)
		for _, asset := range user.Wallet {
			if(asset.Id == stock_id) {
				userHaveStock.Id = user.Id
				userHaveStock.Name = user.Name
				userHaveStock.Count = asset.Count
				listUser.Users = append(listUser.Users, userHaveStock)  
			}
		}
	}
	fmt.Println("list users have stock_id " + stock_id, listUser.Users)

	//change to array of bytes
	listTranAsBytes, _ := json.Marshal(listUser)              	
	return shim.Success(listTranAsBytes)
}