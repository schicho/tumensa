package tumensa

import (
	"fmt"
	"time"
)

func PrintDateAndDay(time time.Time) {
	fmt.Printf("\033[35;1m%s %s:\033[0m\n", time.Format("02.01.2006"), time.Weekday())
}

func PrettyPrintMenus(menus []Menu) {
	for _, menu := range menus {
		fmt.Printf("\033[36;1m%s:\033[0m\n", menu.Name)
		for _, dish := range menu.Dishes {
			fmt.Printf("    - %s : %s\n", dish.Price, dish.Name)
		}
	}
}
