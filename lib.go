package main

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// Get stock - get a stock asset from ledger
func get_stock(stub shim.ChaincodeStubInterface, id string) (Stock, error) {
	var stock Stock
	stockAsBytes, err := stub.GetState(id)                  	//getState retreives a key/value from the ledger
	if err != nil {                                          	//this seems to always succeed, even if key didn't exist
		return stock, errors.New("Failed to find stock - " + id)
	}
	json.Unmarshal(stockAsBytes, &stock)                   		//un stringify it aka JSON.parse()

	if stock.Id != id {                                     //test if stock is actually here or just nil
		return stock, errors.New("Stock does not exist - " + id)
	}

	return stock, nil
}

// Get User - get the user asset from ledger
func get_user(stub shim.ChaincodeStubInterface, id string) (User, error) {
	var user User
	userAsBytes, err := stub.GetState(id)                     //getState retreives a key/value from the ledger
	if err != nil {                                            //this seems to always succeed, even if key didn't exist
		return user, errors.New("Failed to get user - " + id)
	}
	json.Unmarshal(userAsBytes, &user)                       //un stringify it aka JSON.parse()

	if len(user.Name) == 0 {                              //test if user is actually here or just nil
		return user, errors.New("User does not exist - " + id + ", '" + user.Name)
	}
	
	return user, nil
}

func get_transaction(stub shim.ChaincodeStubInterface, id string) (Trade, error) {
	var tran Trade
	tranAsBytes, err := stub.GetState(id)                     //getState retreives a key/value from the ledger
	if err != nil {                                            //this seems to always succeed, even if key didn't exist
		return tran, errors.New("Failed to get transaction - " + id)
	}
	json.Unmarshal(tranAsBytes, &tran)                       //un stringify it aka JSON.parse()
	
	if tran.Id != id {                                     //test if stock is actually here or just nil
		return tran, errors.New("Transaction does not exist - " + id)
	}

	return tran, nil
}

// ========================================================
// Input Sanitation - dumb input checking, look for empty strings
// ========================================================
func sanitize_arguments(strs []string) error{
	for i, val:= range strs {
		if len(val) <= 0 {
			return errors.New("Argument " + strconv.Itoa(i) + " must be a non-empty string")
		}
		if len(val) > 32 {
			return errors.New("Argument " + strconv.Itoa(i) + " must be <= 32 characters")
		}
	}
	return nil
}
