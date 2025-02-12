package tumensa

import (
	"os"
	"testing"
	"time"
)

func TestParseMenuJSON(t *testing.T) {
	file, err := os.Open("testResources/menuPlanCurrentWeek.json")
	if err != nil {
		t.Fatal("Failed to read test file:", err)
	}
	defer file.Close()

	menus, err := parseMenuJson(file, time.Wednesday)
	if err != nil {
		t.Error("Failed to parse menu JSON:", err)
	}

	expectedMenusWednesday := []Menu{
		{
			Name: "Menü Veggie",
			Dishes: []Dish{
				{
					Name:  "Rösti Ratatouille Dip mit Kren und Kräutern",
					Price: "6.50",
				},
			},
		},
		{
			Name: "Menü Herzhaft",
			Dishes: []Dish{
				{
					Name:  "Koreanisches Bulgogi Eiernudeln",
					Price: "7.60",
				},
			},
		},
		{
			Name: "Tagesgerichte",
			Dishes: []Dish{
				{
					Name:  "Little Italy Burger mit BIO Rindfleisch und Pommes frites",
					Price: "11.90",
				},
				{
					Name:  "Melanzani Burger mit Pommes frites",
					Price: "9.90",
				},
			},
		},
	}

	if len(menus) != len(expectedMenusWednesday) {
		t.Errorf("Got %d menus, expected %d", len(menus), len(expectedMenusWednesday))
	}

	for i, menu := range menus {
		if menu.Name != expectedMenusWednesday[i].Name {
			t.Errorf("Menu %d: got %s, expected %s", i, menu.Name, expectedMenusWednesday[i])
		}
	}

	for i, menu := range menus {
		if len(menu.Dishes) != len(expectedMenusWednesday[i].Dishes) {
			t.Errorf("Menu %d: got %d dishes, expected %d", i, len(menu.Dishes), len(expectedMenusWednesday[i].Dishes))
		}

		for j, dish := range menu.Dishes {
			if dish != expectedMenusWednesday[i].Dishes[j] {
				t.Errorf("Menu %d, Dish %d: got %v, expected %v", i, j, dish, expectedMenusWednesday[i].Dishes[j])
			}
		}
	}
}
