package main

import (
	"fmt"
	"github.com/bwolf/suncal"
	"time"
)

func main() {
	coord := suncal.GeoCoord{50.0, 10.0}
	date := time.Now()

	doit := func(twilightName string, twilightKind float64) {
		sun := suncal.SunCal(coord, twilightKind, date)
		fmt.Printf("Sun rise and set (twilight %s) for date %s\n", twilightName, date)
		fmt.Printf("Rise %02d:%02d\n", sun.Rise.Hour(), sun.Rise.Minute())
		fmt.Printf("Set  %02d:%02d\n", sun.Set.Hour(), sun.Set.Minute())
	}

	doit("Normal", suncal.TwilightDefault)
}
