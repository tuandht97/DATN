package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// Init Stock - create a new stock, store into chaincode state
func init_stock(stub shim.ChaincodeStubInterface, args []string) (pb.Response) {
	var err error
		fmt.Println("starting init_stock")

		if len(args) != 5 {
			return shim.Error("Incorrect number of arguments. Expecting 5")
		}

		err = sanitize_arguments(args)
		if err != nil {
			return shim.Error(err.Error())
		}

		id := args[0]
		code := args[1]
		count, err := strconv.Atoi(args[2])
		if err != nil {
			return shim.Error("2rd argument must be a numeric string")
		}
		price, err := strconv.Atoi(args[3])
		if err != nil {
			return shim.Error("3rd argument must be a numeric string")
		}
		user_id := args[4]
		
		// check user
		user, err := get_user(stub, user_id)
		if err != nil {
			fmt.Println("Failed to find user -" + user_id)
			return shim.Error(err.Error())
		}

		// check stock 
		cstock, err := get_stock(stub, id)
		if err == nil {
			return shim.Error("This stock already exists - " + id)
		}

		if cstock.Code == code{
			return shim.Error("This stock already exists - " + code)
		}

		var stock Stock
		stock.ObjectType = "stock"
		stock.Id = id
		stock.Code = code
		stock.Count = count
		stock.Price = price
		stock.Creator.Id = user.Id
		stock.Creator.Name = user.Name

		stockAsBytes, _ := json.Marshal(stock)                         
		err = stub.PutState(stock.Id, stockAsBytes)                    
		if err != nil {
			fmt.Println("Could not store stock")
			return shim.Error(err.Error())
		}

		var asset Asset
		asset.Id = stock.Id
		asset.Code = stock.Code
		asset.Count = stock.Count
		user.Wallet = append(user.Wallet, asset)

		userAsBytes, _ := json.Marshal(user)                         
		err = stub.PutState(user.Id, userAsBytes)                    
		if err != nil {
			fmt.Println("Could not store user")
			return shim.Error(err.Error())
		}

		fmt.Println("- end init_stock")
		return shim.Success(nil)
}

// Init User - create a new user, store into chaincode state
func init_user(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
		fmt.Println("starting init_user")
	
		if len(args) != 2 {
			return shim.Error("Incorrect number of arguments. Expecting 2")
		}
	
		//input sanitation
		err = sanitize_arguments(args)
		if err != nil {
			return shim.Error(err.Error())
		}
	
		var user User
		user.ObjectType = "user"
		user.Id =  args[0]
		user.Name = args[1]
		user.Wallet = nil
		fmt.Println(user)
	
		//check if user already exists
		_, err = get_user(stub, user.Id)
		if err == nil {
			fmt.Println("This user already exists - " + user.Id)
			return shim.Error("This user already exists - " + user.Id)
		}
	
		//store user
		userAsBytes, _ := json.Marshal(user)                         //convert to array of bytes
		err = stub.PutState(user.Id, userAsBytes)                    //store owner by its Id
		if err != nil {
			fmt.Println("Could not store user")
			return shim.Error(err.Error())
		}
	
		fmt.Println("- end init_user")
		return shim.Success(nil)
}

// Update price of stock
func update_price(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting update_price")

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	var id = args[0]
	new_price, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("2rd argument must be a numeric string")
	}

	stockAsBytes, err := stub.GetState(id)
	if err != nil {
		return shim.Error("Failed to get stock")
	}
	res := Stock{}
	json.Unmarshal(stockAsBytes, &res)           

	res.Price = new_price
	jsonAsBytes, _ := json.Marshal(res)           //convert to array of bytes
	err = stub.PutState(args[0], jsonAsBytes)     //rewrite the stock with id as key
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println(res.Code + "->" + strconv.Itoa(new_price))
	fmt.Println("- end update price")
	return shim.Success(nil)
}

func init_transaction(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	fmt.Println("starting init_transaction")

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	// input sanitation
	err = sanitize_arguments(args)
	if err != nil {
		return shim.Error(err.Error())
	}

	trade_id := args[0]
	stock_id := args[1]
	stock_count, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("2rd argument must be a numeric string")
	}
	seller_id := args[3]
	buyer_id := args[4]
	time := args[5]

	// check stock 
	ctransaction, err := get_transaction(stub, trade_id)
	if err == nil {
		return shim.Error("This transaction already exists - " + trade_id)
	}
	_ = ctransaction

	// check if seller already exists
	seller, err := get_user(stub, seller_id)
	if err != nil {
		return shim.Error("This seller does not exist - " + seller_id)
	}

	// check if buyer already exists
	buyer, err := get_user(stub, buyer_id)
	if err != nil {
		return shim.Error("This buyer does not exist - " + buyer_id)
	}

	stock, err := get_stock(stub, stock_id)
	if err != nil {
		return shim.Error("This stock does not exist - " + stock_id)
	}

	fmt.Println(buyer.Id + " - " + buyer.Name + " buy " + args[2] + " code " + stock.Code + " from " + seller.Id + " - " + seller.Name)

	// get seller
	sellerAsBytes, err := stub.GetState(seller_id)
	if err != nil {
		return shim.Error("Failed to get seller info")
	}
	res := User{}
	json.Unmarshal(sellerAsBytes, &res)           //un stringify it aka JSON.parse()

	// check wallet seller
	if len(res.Wallet) > 0 {
		for i := len(res.Wallet) - 1; i >= 0; i-- {
			if res.Wallet[i].Id == stock.Id {
				if res.Wallet[i].Count >= stock_count {
					break
				} else {
					return shim.Error("The amount in the wallet is not enough")
				}
			}
		}
	} else {
		return shim.Error("The amount in the wallet is not enough")
	}

	update_wallet(stub, seller, stock.Id, stock.Code, stock_count, 1)
	update_wallet(stub, buyer, stock.Id, stock.Code, stock_count, 0)

	var transaction Trade
	transaction.ObjectType = "trade"
	transaction.Id = trade_id
	transaction.Stock.Id = stock_id
	transaction.Stock.Code = stock.Code
	transaction.Stock.Count = stock_count
	transaction.Seller.Id = seller.Id
	transaction.Seller.Name = seller.Name
	transaction.Buyer.Id = buyer.Id
	transaction.Buyer.Name = buyer.Name
	transaction.Time = time

	tradeAsBytes, _ := json.Marshal(transaction)                         
	err = stub.PutState(transaction.Id, tradeAsBytes)                    
	if err != nil {
		fmt.Println("Could not store transaction")
		return shim.Error(err.Error())
	}

	fmt.Println("- end init_transaction")
	return shim.Success(nil)
}

func update_wallet(stub shim.ChaincodeStubInterface, user User, stock_id string, stock_code string, count int, operation int) pb.Response {
	var err error
	fmt.Println("starting update_walllet")

	if len(user.Wallet) > 0 {
		var check = 0
		for i := len(user.Wallet) - 1; i >= 0; i-- {
			if user.Wallet[i].Code == stock_code {
				check = 1
				if operation == 0 {
					user.Wallet[i].Count += count
				} else {
					user.Wallet[i].Count -= count
					if user.Wallet[i].Count <= 0 {
						user.Wallet = append(user.Wallet[:i], user.Wallet[i+1:]...)
					}
				}
			}
		}
		if (check == 0) {
			var asset Asset
			asset.Id = stock_id
			asset.Code = stock_code
			asset.Count = count
			user.Wallet = append(user.Wallet, asset)
		}
	} else {
		var asset Asset
		asset.Id = stock_id
		asset.Code = stock_code
		asset.Count = count
		user.Wallet = append(user.Wallet, asset)
	}

	jsonAsBytes, _ := json.Marshal(user)          
	err = stub.PutState(user.Id, jsonAsBytes)     
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println(user)
	fmt.Println("- end update price")
	return shim.Success(nil)
}