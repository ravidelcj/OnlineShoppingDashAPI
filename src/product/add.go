package product

import (
	"databaseHelper"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"model"
	"net/http"
)

var jsonMap = make(map[string]interface{})

//A global json object which is sent in case error occurs
//it takes the status and message of the error
//status 0 indicates error
func initJson(status int, message string) {
	jsonMap["status"] = status
	jsonMap["message"] = message
}

//add product
//product have a product-id product-name tag price company stock
func AddProd(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", "application/json")
	if req.Method != "POST" {
		fmt.Println("Not a POST request")
		initJson(0, "Not a POST request")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		fmt.Println(err)
		initJson(0, "request body error")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	var prod model.Product

	err = json.Unmarshal(body, &prod)

	if err != nil {
		fmt.Println(err)
		initJson(0, "Incorrect Keys in json")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	isFieldsCorrect := checkAllFields(prod)

	if isFieldsCorrect == false {
		fmt.Println("Not every key value pair present")
		initJson(0, "Not every key value pair present")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	//Various checking can be applied on the product sent to check whether
	//Same products are not being entered again and again
	//Here only product id is being used to check redundant for th
	//sake of convenience
	isProdExist := databaseHelper.CheckProdExist(prod)

	if isProdExist {
		fmt.Println("Product Already exist")
		initJson(0, "Product already exist")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}

	err = databaseHelper.AddProduct(prod)

	if err != nil {
		fmt.Println("Error in Adding to database")
		initJson(0, "Internal database error")
		json.NewEncoder(res).Encode(jsonMap)
		return
	}
	fmt.Println("Product added")
	initJson(1, "Product added")
	json.NewEncoder(res).Encode(jsonMap)
}

//Check whether every field is present or not
func checkAllFields(prod model.Product) bool {
	if prod.ProductId == "" || prod.ProductName == "" || prod.Price == 0 || prod.Tag == "" || prod.Company == "" || prod.Stock == 0 {
		return false
	}
	return true
}
