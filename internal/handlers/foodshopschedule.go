package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Jalenarms1/caters-go/internal/db"
)

func HandlerNewFoodShopSchedule(w http.ResponseWriter, r *http.Request) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	var foodShopSchedule db.FoodShopSchedule
	err = json.Unmarshal(body, &foodShopSchedule)
	if err != nil {
		return err
	}

	fmt.Println(foodShopSchedule)

	err = foodShopSchedule.Insert()
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func HandleGetFoodShopSchedule(w http.ResponseWriter, r *http.Request) error {
	foodShopId := r.PathValue("foodShopId")

	items, err := db.GetFoodShopSchedule(foodShopId)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(items)
}
