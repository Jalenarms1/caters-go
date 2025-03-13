package db

type FoodShopItem struct {
	Id              string  `json:"id"`
	FoodShopID      string  `json:"foodShopId"`
	Label           string  `json:"label"`
	Description     *string `json:"description"`
	Category        string  `json:"category"`
	Images          string  `json:"images"`
	Price           int     `json:"price"`
	IsInStock       bool    `json:"isInStock"`
	DefaultToppings *string `json:"defaultToppings"`
}

func (fi *FoodShopItem) Insert() error {
	query := `
		INSERT INTO FoodShopItem (
			Id, FoodShopId, Label, Description, Category, Images, Price, IsInStock, DefaultToppings
		) 
		VALUES (
			?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`

	_, err := db.Exec(query, fi.Id, fi.FoodShopID, fi.Label, fi.Description, fi.Category, fi.Images, fi.Price, fi.IsInStock, fi.DefaultToppings)

	return err
}

func GetFoodItemCategories() (*[]string, error) {
	rows, err := db.Query("select e.enumlabel from pg_enum e join pg_type t on t.oid = e.enumtypid where t.typname = 'fooditemcategory'")
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

func GetFoodShopItems(foodShopId string) (*[]FoodShopItem, error) {
	query := `
	SELECT 
		Id, 
		FoodShopId, 
		Label, 
		Description, 
		Category, 
		Images, 
		Price, 
		IsInStock, 
		DefaultToppings
	FROM FoodShopItem
	WHERE FoodShopId = ?
	`

	rows, err := db.Query(query, foodShopId)
	var items []FoodShopItem
	for rows.Next() {
		var fi FoodShopItem
		_ = rows.Scan(
			&fi.Id,
			&fi.FoodShopID,
			&fi.Label,
			&fi.Description,
			&fi.Category,
			&fi.Images,
			&fi.Price,
			&fi.IsInStock,
			&fi.DefaultToppings,
		)

		items = append(items, fi)

	}

	return &items, err
}
