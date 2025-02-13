package tumensa

import (
	"bytes"
	"encoding/json"
	"io"
	"regexp"
	"strconv"
	"time"
)

type gqlMensaRespBody struct {
	Data struct {
		NodeByUri struct {
			MenuplanCurrentWeek string `json:"menuplanCurrentWeek"`
		} `json:"nodeByUri"`
	} `json:"data"`
}

type mensaMenuPlan struct {
	Menus []struct {
		Name  string `json:"name"`
		Menus map[string][]struct {
			TitleDe string `json:"title_de"`
			Price   string `json:"price"`
		} `json:"menus"`
	} `json:"menus"`
}

type Menu struct {
	Name   string
	Dishes []Dish
}

type Dish struct {
	Name  string
	Price string
}

// ParseGQLResponse parses the JSON response from the GraphQL API and returns a slice of menu structs.
// The weekday parameter is used to select the menu for the current day.
func ParseGQLResponse(resp io.Reader, weekday time.Weekday) ([]Menu, error) {
	menuJson, err := parseGQLResponse(resp)
	if err != nil {
		return nil, err
	}
	return parseMenuJson(bytes.NewReader(menuJson), weekday)
}

// parseGQLResponse extracts the embedded menu json from the HTTP response body.
// The respone body contains a JSON GraphQL response, which in itself contains the actual menu JSON
// embedded under the key "data.nodeByUri.menuplanCurrentWeek".
func parseGQLResponse(resp io.Reader) ([]byte, error) {
	var gqlResponse gqlMensaRespBody
	if err := json.NewDecoder(resp).Decode(&gqlResponse); err != nil {
		return nil, err
	}
	menu := gqlResponse.Data.NodeByUri.MenuplanCurrentWeek
	return []byte(menu), nil
}

// parseMenuJson parses the retrieved menu plan json into a slice of our simplified menu structs.
func parseMenuJson(menuJson io.Reader, weekday time.Weekday) ([]Menu, error) {
	var menuPlan mensaMenuPlan
	if err := json.NewDecoder(menuJson).Decode(&menuPlan); err != nil {
		return nil, err
	}

	var parsedMenus []Menu = make([]Menu, 0, len(menuPlan.Menus))

	for _, m := range menuPlan.Menus {
		parsedMenuName := cleanMenuName(m.Name)
		// The data structure is odd here. Compare the provided test file "menuPlanCurrentWeek.json".
		// The data is first structured by the menu type like "Menü Veggie", "Menü Herzhaft", etc.
		// and only subsequently by the day of the week.
		// As such we need to use the weekday as a string key to access the correct menu.
		parsedDishes := make([]Dish, 0, len(m.Menus[strconv.Itoa(int(weekday))]))

		for _, d := range m.Menus[strconv.Itoa(int(weekday))] {
			parsedDishName := cleanDishName(d.TitleDe)
			parsedDishPrice := d.Price

			parsedDishes = append(parsedDishes, Dish{
				Name:  parsedDishName,
				Price: parsedDishPrice,
			})
		}
		parsedMenus = append(parsedMenus, Menu{
			Name:   parsedMenuName,
			Dishes: parsedDishes,
		})
	}
	return parsedMenus, nil
}

// cleanDishName removes <br /> tags followed by a newline from the dish name.
func cleanDishName(name string) string {
	re := regexp.MustCompile(`<br />\n`)
	return re.ReplaceAllString(name, " ")
}

// cleanMenuName removes the content in parentheses from the menu name.
func cleanMenuName(name string) string {
	re := regexp.MustCompile(` \(.*\)`)
	return re.ReplaceAllString(name, "")
}
