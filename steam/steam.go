package steam

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
)

const tableName = "steam_discount_game"

var (
	queryAllDiscounts string
)

// Prepare queries.
func init() {
	fields := []string{
		"name", "link", "img_src", "review", "price_before", "price_now",
		"discount"}
	queryAllDiscounts = fmt.Sprintf(
		"SELECT %s FROM %s", strings.Join(fields, ", "), tableName)
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
	defer rows.Close()

	if err != nil {
		return nil, err
	}

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
