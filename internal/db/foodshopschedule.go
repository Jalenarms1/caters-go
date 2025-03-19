package db

import (
	"fmt"
)

type FoodShopSchedule struct {
	Id         string `json:"id"`
	FoodShopID string `json:"foodShopId"`
	OpenDt     string `json:"openDt"`
	CloseDt    string `json:"closeDt"` // -- 20250101
	OpenTm     int64  `json:"openTm"`  // -- 09:30, 13:30
	CloseTm    int64  `json:"closeTm"` //
	CreatedAt  int64  `json:"createdAt"`
}

func (fi *FoodShopSchedule) Insert() error {
	query := `
		INSERT INTO FoodShopSchedule (
			Id, FoodShopId, OpenDt, CloseDt, OpenTm, CloseTm, CreatedAt
		) 
		VALUES (
			?, ?, ?, ?, ?, ?, ?
		)
	`

	_, err := db.Exec(query, fi.Id, fi.FoodShopID, fi.OpenDt, fi.CloseDt, fi.OpenTm, fi.CloseTm, fi.CreatedAt)

	return err
}

func GetFoodShopSchedule(foodShopId string) (*[]FoodShopSchedule, error) {
	query := `
	SELECT 
		Id, 
		FoodShopId, 
		OpenDt, 
		CloseDt, 
		OpenTm,
		CloseTm, 
		CreatedAt
	FROM FoodShopSchedule
	WHERE FoodShopId = ?
	ORDER BY OpenDt, OpenTm
	`

	rows, err := db.Query(query, foodShopId)
	var items []FoodShopSchedule
	for rows.Next() {
		var fi FoodShopSchedule
		err := rows.Scan(
			&fi.Id,
			&fi.FoodShopID,
			&fi.OpenDt,
			&fi.CloseDt,
			&fi.OpenTm,
			&fi.CloseTm,
			&fi.CreatedAt,
		)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("foodschedulde", fi)

		items = append(items, fi)

	}
	// fmt.Println("items", items)
	return &items, err
}

func DeleteScheduleSlot(id string) error {
	query := `
	DELETE FROM FoodShopSchedule WHERE Id = ?
	`
	_, err := db.Exec(query, id)

	return err
}
