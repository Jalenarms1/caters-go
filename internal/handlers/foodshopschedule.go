package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Jalenarms1/caters-go/internal/db"
	"github.com/Jalenarms1/caters-go/internal/types"
)

func HandlerNewFoodShopSchedule(w http.ResponseWriter, r *http.Request) *types.Error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return &types.Error{
			Err:        err,
			ReturnCode: http.StatusInternalServerError,
		}
	}
	defer r.Body.Close()

	var foodShopSchedule db.FoodShopSchedule
	err = json.Unmarshal(body, &foodShopSchedule)
	if err != nil {
		return &types.Error{
			Err:        err,
			ReturnCode: http.StatusInternalServerError,
		}
	}

	fmt.Println(foodShopSchedule)

	err = foodShopSchedule.Insert()
	if err != nil {
		return &types.Error{
			Err:        err,
			ReturnCode: http.StatusInternalServerError,
		}
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func HandleGetFoodShopSchedule(w http.ResponseWriter, r *http.Request) *types.Error {
	foodShopId := r.PathValue("foodShopId")

	items, err := db.GetFoodShopSchedule(foodShopId)
	if err != nil && err != sql.ErrNoRows {
		return &types.Error{
			Err:        err,
			ReturnCode: http.StatusInternalServerError,
		}
	}

	fmt.Println("items again", items)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(items)
	if err != nil {
		return &types.Error{
			Err:        err,
			ReturnCode: http.StatusInternalServerError,
		}
	}

	return nil
}

func HandleDeleteScheduleSlot(w http.ResponseWriter, r *http.Request) *types.Error {

	slotId := r.PathValue("slotId")

	err := db.DeleteScheduleSlot(slotId)

	if err != nil {
		return &types.Error{
			Err:        err,
			ReturnCode: http.StatusInternalServerError,
		}
	}

	return nil
}
