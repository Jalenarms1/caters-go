package db

import (
	"database/sql"
	"fmt"
	"time"
)

type FoodShopSchedule struct {
	Id         string `json:"id"`
	FoodShopID string `json:"foodShopId"`
	OpenDt     string `json:"openDt"`
	CloseDt    string `json:"closeDt"` // -- 20250101
	OpenTm     string `json:"openTm"`  // -- 0930, 1330
	CloseTm    string `json:"closeTm"` //
	CreatedAt  string `json:"createdAt"`
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
	`

	rows, err := db.Query(query, foodShopId)
	var items []FoodShopSchedule
	for rows.Next() {
		var fi FoodShopSchedule
		var openDate, closeDate, createdAt interface{}
		var openTime, closeTime sql.NullString
		err := rows.Scan(
			&fi.Id,
			&fi.FoodShopID,
			&openDate,
			&closeDate,
			&openTime,
			&closeTime,
			&createdAt,
		)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(openDate.(time.Time).Format("2006-01-02"))
		fmt.Println(openDate, closeDate, createdAt)
		fmt.Println(openTime.String[:5], closeTime.String[:5])

		fi.OpenDt = openDate.(time.Time).Format("2006-01-02")
		fi.CloseDt = closeDate.(time.Time).Format("2006-01-02")
		fi.CreatedAt = createdAt.(time.Time).Format("2006-01-02")
		fi.OpenTm = openTime.String[:5]
		fi.CloseTm = closeTime.String[:5]

		fmt.Println(fi)

		items = append(items, fi)

	}

	return &items, err
}
