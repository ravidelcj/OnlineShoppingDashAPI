package databaseHelper

import (
	//"database/sql"
	"fmt"
	"model"
)

func AddProduct(prod model.Product) error {
	stmt, err := Db.Prepare("Insert Into product set product_id = ? , product_name = ? , tag = ? , price = ? , company = ? , stock = ? ")

	if err != nil {
		fmt.Println("Error in preparing statement ", err)
		return err
	}
	_, err = stmt.Exec(prod.ProductId, prod.ProductName, prod.Tag, prod.Price, prod.Company, prod.Stock)
	if err != nil {
		fmt.Println("Add Product : ", err)
		return err
	}
	return nil
}

func CheckProdExist(prod model.Product) bool {
	query := "select exists(select 1 from product where product_id = '" + prod.ProductId + "' );"

	var exist bool

	err := Db.QueryRow(query).Scan(&exist)

	if err != nil {
		fmt.Println(err)
		return false
	}

	return exist
}
