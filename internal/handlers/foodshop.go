package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/Jalenarms1/caters-go/internal/db"
	"github.com/Jalenarms1/caters-go/internal/types"
)

func HandlerGetFoodShopCategories(w http.ResponseWriter, r *http.Request) error {

	categories, err := db.GetFoodShopCategories()
	if err != nil {
		return err
	}

	fmt.Println(categories)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fmt.Sprintf(`{"categories": %v}`, categories))

	return nil
}

func HandleNewFoodShop(w http.ResponseWriter, r *http.Request) error {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	var foodShop db.FoodShop
	err = json.Unmarshal(body, &foodShop)
	if err != nil {
		return err
	}

	fmt.Println(foodShop)

	err = foodShop.Insert()
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)

	return nil
}

func HandleGetFoodShop(w http.ResponseWriter, r *http.Request) error {
	shopUrlSlug := r.Context().Value(types.UrlSlugKey)

	if shopUrlSlug == nil {
		return errors.New("invalid url slug")
	}

	foodShop, err := db.GetFoodShop(shopUrlSlug.(string))
	if err != nil {
		return err
	}

	if foodShop == nil {
		return errors.New("food shop not found")
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(foodShop)

	return err

}

func HandleGetMyFoodShop(w http.ResponseWriter, r *http.Request) error {
	uid := r.Context().Value(types.AuthKey).(string)

	foodShop, err := db.GetFoodShopByUserId(uid)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(foodShop)

}
