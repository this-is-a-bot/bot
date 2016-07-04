package steam

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
)

const discountTableName = "steam_discount_game"
const featuredTableName = "steam_featured_game"

var (
	queryAllDiscounts string
	queryAllFeatured  string
)

// Prepare queries.
func init() {
	fields := []string{
		"name", "link", "img_src", "review", "price_before", "price_now",
		"discount"}
	queryAllDiscounts = fmt.Sprintf(
		"SELECT %s FROM %s", strings.Join(fields, ", "), discountTableName)

	fields_featured := []string{
		"name", "link", "img_src", "headline", "price_before", "price_now",
		"discount"}
	queryAllFeatured = fmt.Sprintf(
		"SELECT %s FROM %s", strings.Join(fields_featured, ","), featuredTableName)
}

// Corresponds to rows in `steam_discount_game` table.
type SteamGame struct {
	Name        string  `json:"name"`
	URL         string  `json:"url"`
	ImgSrc      string  `json:"imgSrc"`
	Review      string  `json:"review"`
	PriceBefore float32 `json:"priceBefore"`
	PriceNow    float32 `json:"priceNow"`
	Discount    string  `json:"discount"`
}

// Get all discounts in current discount table.
func GetDiscounts(db *sql.DB) ([]SteamGame, error) {
	rows, err := db.Query(queryAllDiscounts)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := make([]SteamGame, 0)
	for rows.Next() {
		var game SteamGame

		err = rows.Scan(
			&game.Name, &game.URL, &game.ImgSrc, &game.Review, &game.PriceBefore,
			&game.PriceNow, &game.Discount)
		if err != nil {
			return nil, err
		}

		res = append(res, game)
	}
	return res, nil
}

// Get all featured games in current featured table.
func GetFeatured(db *sql.DB, feature string) ([]SteamGame, error) {
	if !IsValidFeature(feature) {
		feature = "win"
	}
	queryOneFeature := fmt.Sprintf(
		"%s where feature_type='featured_%s'", queryAllFeatured, feature)
	rows, err := db.Query(queryOneFeature)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	res := make([]SteamGame, 0)
	for rows.Next() {
		var game SteamGame

		err = rows.Scan(
			&game.Name, &game.URL, &game.ImgSrc, &game.Review, &game.PriceBefore,
			&game.PriceNow, &game.Discount)
		if err != nil {
			return nil, err
		}
		res = append(res, game)
	}
	return res, nil
}

func IsValidFeature(feature string) bool {
	switch feature {
	case
		"win",
		"linux",
		"mac":
		return true
	}
	return false
}
