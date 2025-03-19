package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/Jalenarms1/caters-go/internal/utils"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5"
)

type FoodShop struct {
	Id                  string             `json:"id"`
	UserId              string             `json:"userId"`
	AccountID           string             `json:"accountId,omitempty"`
	UrlSlug             string             `json:"urlSlug"`
	Label               string             `json:"label"`
	Bio                 *string            `json:"bio"`
	FoodCategory        string             `json:"foodCategory"`
	Logo                string             `json:"logo"`
	IsDeliveryAvailable bool               `json:"isDeliveryAvailable"`
	Address             string             `json:"address"`
	City                string             `json:"city"`
	State               string             `json:"state"`
	ZipCode             string             `json:"zipCode"`
	Country             string             `json:"country"`
	Latitude            *float64           `json:"latitude"`
	Longitude           *float64           `json:"longitude"`
	MaxDeliveryRadius   *float64           `json:"maxDeliveryRadius"`
	DeliveryFee         *float64           `json:"deliveryFee"`
	CreatedAt           int64              `json:"createdAt"`
	FoodShopItems       []FoodShopItem     `json:"foodShopItems,omitempty"`
	FoodShopSchedule    []FoodShopSchedule `json:"foodShopSchedule,omitempty"`
}

func (f *FoodShop) Insert() error {
	// 	f.UrlSlug = fmt.Sprintf("%s-%s", f.Label, utils.GenerateRandomUrlSlug())
	uid, _ := uuid.NewV4()
	f.Id = uid.String()

	f.UrlSlug = fmt.Sprintf("%s-%s", strings.ToLower(strings.Replace(f.Label, " ", "-", -1)), utils.GenerateRandomUrlSlug())
	f.CreatedAt = time.Now().Unix()

	if !GetIsSlugAvailable(f.UrlSlug) {
		return errors.New("url slug already taken")
	}

	_, err := db.Exec("insert into FoodShop (Id, UserId, Label, UrlSlug, Bio, FoodCategory, Logo, IsDeliveryAvailable, Address, City, State, ZipCode, Country, Latitude, Longitude, MaxDeliveryRadius, DeliveryFee, CreatedAt) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", f.Id, f.UserId, f.Label, f.UrlSlug, f.Bio, f.FoodCategory, f.Logo, f.IsDeliveryAvailable, f.Address, f.City, f.State, f.ZipCode, f.Country, f.Latitude, f.Longitude, f.MaxDeliveryRadius, f.DeliveryFee, f.CreatedAt)

	if os.Getenv("IS_DEV") == "true" {
		cmd := exec.Command("sh", "-c", fmt.Sprintf("echo '%s' | sudo -S sh -c 'echo 127.0.0.1 %s.localhost >> /etc/hosts'", os.Getenv("SH_PASS"), f.UrlSlug))
		err := cmd.Run()
		if err != nil {
			fmt.Println("error adding subdomain to hosts")
		}
	}

	return err
}

func GetFoodShopCategories() (*[]string, error) {
	rows, err := db.Query("select e.enumlabel from pg_enum e join pg_type t on t.oid = e.enumtypid where t.typname = 'foodcategory' order by e.enumlabel")
	if err != nil {
		return nil, err
	}

	var categories []string
	for rows.Next() {
		var cat string
		_ = rows.Scan(&cat)

		categories = append(categories, cat)
	}

	return &categories, nil
}

func GetIsSlugAvailable(urlSlug string) bool {
	row := db.QueryRow("select Id from FoodShop where UrlSlug = ?", urlSlug)

	var id string
	err := row.Scan(&id)
	fmt.Println(err)
	fmt.Println(err == pgx.ErrNoRows)
	return err == sql.ErrNoRows
}

func GetFoodShop(urlSlug string) (*FoodShop, error) {
	query := `
	SELECT 
		Id, 
		UrlSlug,
		Label,
		Bio,
		FoodCategory,
		Logo,
		IsDeliveryAvailable,
		Address,
		City,
		State,
		ZipCode,
		Country,
		Latitude,
		Longitude,
		MaxDeliveryRadius,
		DeliveryFee,
		CreatedAt
	FROM FoodShop
	WHERE UrlSlug = ?
	`

	row := db.QueryRow(query, urlSlug)

	var foodShop FoodShop
	err := row.Scan(
		&foodShop.Id,
		&foodShop.UrlSlug,
		&foodShop.Label,
		&foodShop.Bio,
		&foodShop.FoodCategory,
		&foodShop.Logo,
		&foodShop.IsDeliveryAvailable,
		&foodShop.Address,
		&foodShop.City,
		&foodShop.State,
		&foodShop.ZipCode,
		&foodShop.Country,
		&foodShop.Latitude,
		&foodShop.Longitude,
		&foodShop.MaxDeliveryRadius,
		&foodShop.DeliveryFee,
		&foodShop.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		fooditems, _ := GetFoodShopItems(foodShop.Id)

		foodShop.FoodShopItems = *fooditems
	}()

	go func() {
		defer wg.Done()
		shopSchedule, _ := GetFoodShopSchedule(foodShop.Id)

		foodShop.FoodShopSchedule = *shopSchedule
	}()

	wg.Wait()

	return &foodShop, nil
}

func GetFoodShopByUserId(userId string) (*FoodShop, error) {

	query := `
	SELECT 
		s.Id, 
		s.UrlSlug,
		s.Label,
		s.Bio,
		s.FoodCategory,
		s.Logo,
		s.IsDeliveryAvailable,
		s.Address,
		s.City,
		s.State,
		s.ZipCode,
		s.Country,
		s.Latitude,
		s.Longitude,
		s.MaxDeliveryRadius,
		s.DeliveryFee,
		s.CreatedAt 
	FROM "User" u 
	JOIN FoodShop s on s.UserId = u.Id
	WHERE u.Id = ?
	`

	row := db.QueryRow(query, userId)

	var foodShop FoodShop
	err := row.Scan(
		&foodShop.Id,
		&foodShop.UrlSlug,
		&foodShop.Label,
		&foodShop.Bio,
		&foodShop.FoodCategory,
		&foodShop.Logo,
		&foodShop.IsDeliveryAvailable,
		&foodShop.Address,
		&foodShop.City,
		&foodShop.State,
		&foodShop.ZipCode,
		&foodShop.Country,
		&foodShop.Latitude,
		&foodShop.Longitude,
		&foodShop.MaxDeliveryRadius,
		&foodShop.DeliveryFee,
		&foodShop.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		fooditems, _ := GetFoodShopItems(foodShop.Id)

		foodShop.FoodShopItems = *fooditems
	}()

	go func() {
		defer wg.Done()
		shopSchedule, _ := GetFoodShopSchedule(foodShop.Id)

		foodShop.FoodShopSchedule = *shopSchedule
	}()

	wg.Wait()

	return &foodShop, nil
}
