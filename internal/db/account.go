package db

import (
	"database/sql"
	"fmt"
	"sync"
)

type Account struct {
	Id       string    `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password,omitempty"`
	FoodShop *FoodShop `json:"foodShop,omitempty"`
}

func (a *Account) Insert() error {

	_, err := db.Exec(`insert into "User" (Id, Email, Password) values (?,?,?)`, a.Id, a.Email, a.Password)
	if err != nil {
		return err
	}

	return nil
}

func GetAccountById(id string) (*Account, error) {

	var acct Account
	var foodShop *FoodShop
	var wg sync.WaitGroup
	var acctErr error
	var foodShopErr error

	wg.Add(2)
	go func() {
		defer wg.Done()
		row := db.QueryRow(`select Id, Email from "User" u where Id = ?`, id)

		acctErr = row.Scan(&acct.Id, &acct.Email)
		if acctErr != nil {
			fmt.Println(acctErr)
			// return nil, err
		}

	}()

	go func() {
		defer wg.Done()
		foodShop, foodShopErr = GetFoodShopByUserId(id)
		fmt.Println(foodShopErr)
		if foodShopErr != nil && foodShopErr != sql.ErrNoRows {
			fmt.Println(foodShopErr)
			// return nil, err
		}

	}()

	wg.Wait()

	if acctErr != nil {
		return nil, acctErr
	}

	if foodShopErr != nil && foodShopErr != sql.ErrNoRows {
		return nil, foodShopErr
	}
	fmt.Println(foodShop)
	acct.FoodShop = foodShop

	return &acct, nil
}

func GetAccountByEmail(email string) (*Account, error) {
	fmt.Println("getting account")
	var acct Account
	row := db.QueryRow(`select Id, Email from "User" where Email = ?`, email)

	err := row.Scan(&acct.Id, &acct.Email)
	fmt.Println("getting account")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(acct)

	return &acct, nil
}

func GetUserWPasswordByEmail(email string) (*Account, error) {

	var acctInfo Account
	err := db.QueryRow(`select Id, Email, Password from "User" where Email = ?`, email).Scan(&acctInfo.Id, &acctInfo.Email, &acctInfo.Password)

	return &acctInfo, err
}
