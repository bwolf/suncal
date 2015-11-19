package main

import (
	"fmt"
	"github.com/bwolf/suncal"
	"time"
)

func main() {
	coords := suncal.GeoCoordinates{48.137222, 11.575556} // Munich, Germany
	date := time.Now()
	fmt.Printf("Coordinates %s\nDetailed date used %s\n", coords, date)

	for n := 0; n < 100; n++ {
		date = date.AddDate(0, 0, 1)
		sun := suncal.SunCal(coords, date)
		delta := sun.Set.Sub(sun.Rise)
		y, m, d := date.Date()
		fmt.Printf("Date %d-%02d-%02d  \u2191 %02d:%02d  \u2193 %02d:%02d  \u0394 %s\n", y, m, d,
			sun.Rise.Hour(), sun.Rise.Minute(),
			sun.Set.Hour(), sun.Set.Minute(),
			delta)
	}
}
