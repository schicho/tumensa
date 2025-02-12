package tumensa

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"regexp"
	"strconv"
	"time"
)

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
	menuJson, err := extractMenuJson(resp)
	if err != nil {
		return nil, err
	}
	return parseMenuJson(bytes.NewReader(menuJson), weekday)
}

// extractMenuJson extracts the embedded menu json from the HTTP response body.
// The respone body contains a JSON GraphQL response, which in itself contains the actual menu JSON
// embedded under the key "data.nodeByUri.menuplanCurrentWeek".
func extractMenuJson(resp io.Reader) ([]byte, error) {
	var raw map[string]map[string]any
	if err := json.NewDecoder(resp).Decode(&raw); err != nil {
		return nil, err
	}

	menu, ok := raw["data"]["nodeByUri"].(map[string]any)["menuplanCurrentWeek"].(string)
	if !ok {
		return nil, errors.New("failed to extract embedded menu json")
	}

	return []byte(menu), nil
}

// parseMenuJson parses the menu JSON and returns a slice of menu structs.
func parseMenuJson(menuJson io.Reader, weekday time.Weekday) ([]Menu, error) {
	var raw map[string]any
	if err := json.NewDecoder(menuJson).Decode(&raw); err != nil {
		return nil, err
	}

	// "menus" here is the menu types, e.g. "Mittagsmenü", "Tagesgericht", etc.
	menusList, ok := raw["menus"].([]any)
	if !ok {
		return nil, errors.New("failed to parse menu list")
	}
	// Go requires an element wise conversion of the interface{} slice to the desired type.
	menus := make([]map[string]any, 0, len(menusList))
	for _, m := range menusList {
		menu, ok := m.(map[string]any)
		if !ok {
			return nil, errors.New("failed to parse menus")
		}
		menus = append(menus, menu)
	}

	var parsedMenus []Menu = make([]Menu, 0, len(menus))

	for _, m := range menus {
		parsedMenuName, ok := m["name"].(string)
		if !ok {
			return nil, errors.New("failed to parse menu name")
		}

		// "menus" here are in fact dishes.
		// The data structure is odd, as only at this point the weekday is used to index into the map.
		// Subsequently, the dishes are stored in a slice of maps, where each map contains the dish name and price.
		dishesDay, ok := m["menus"].(map[string]any)
		if !ok {
			return nil, errors.New("failed to parse dishes by day")
		}
		dishListSpecifiedDay, ok := dishesDay[strconv.Itoa(int(weekday))]
		if !ok {
			return nil, errors.New("no key for dishes of specified day")
		}

		dishesList, ok := dishListSpecifiedDay.([]any)
		if !ok {
			return nil, errors.New("failed to parse dishes list for specified day")
		}
		dishes := make([]map[string]any, 0, len(dishesList))
		for _, d := range dishesList {
			dish, ok := d.(map[string]any)
			if !ok {
				return nil, errors.New("failed to parse dish")
			}
			dishes = append(dishes, dish)
		}

		parsedDishes := make([]Dish, 0, len(dishes))

		for _, d := range dishes {
			parsedDishName, ok := d["title_de"].(string)
			if !ok {
				return nil, errors.New("failed to parse dish name")
			}

			parsedDishPrice, ok := d["price"].(string)
			if !ok {
				return nil, errors.New("failed to parse dish price")
			}

			parsedDishes = append(parsedDishes, Dish{
				Name:  cleanDishName(parsedDishName),
				Price: parsedDishPrice,
			})
		}
		parsedMenus = append(parsedMenus, Menu{
			Name:   cleanMenuName(parsedMenuName),
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
