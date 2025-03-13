package db

import "fmt"

type Account struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
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
	row := db.QueryRow(`select Id, Email from "User" where Id = ?`, id)

	err := row.Scan(&acct.Id, &acct.Email)
	if err != nil {
		return nil, err
	}

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
