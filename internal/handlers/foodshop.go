package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"

	"github.com/Jalenarms1/caters-go/internal/db"
	"github.com/Jalenarms1/caters-go/internal/types"
	"github.com/Jalenarms1/caters-go/internal/utils"
	"github.com/gofrs/uuid"
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

func HandleNewFoodShop(w http.ResponseWriter, r *http.Request) *types.Error {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return &types.Error{
			Err:        err,
			ReturnCode: http.StatusInternalServerError,
		}
	}
	defer r.Body.Close()

	var foodShop db.FoodShop
	err = json.Unmarshal(body, &foodShop)
	if err != nil {
		return &types.Error{
			Err:        err,
			ReturnCode: http.StatusBadRequest,
		}
	}

	fmt.Println(foodShop)

	foodShop.UserId = r.Context().Value(types.AuthKey).(string)

	uid, _ := uuid.NewV4()
	imagePath := path.Join("public/images", fmt.Sprintf("%s.png", uid))

	err = utils.SaveImage(imagePath, foodShop.Logo)
	if err != nil {
		return &types.Error{
			Err:        err,
			ReturnCode: http.StatusInternalServerError,
		}
	}

	foodShop.Logo = imagePath

	err = foodShop.Insert()
	if err != nil {
		fmt.Println(err)
		return &types.Error{
			Err:        errors.New("error occurred processing order"),
			ReturnCode: http.StatusInternalServerError,
		}
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
