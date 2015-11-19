package suncal

import (
	"fmt"
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	coord := GeoCoord{50.0, 10.0}
	fmt.Printf("Coordinates %+v\n", coord)
	rise := time.Date(2005, 9, 30, 7, 18, 0, 0, time.FixedZone("CET", 7200))
	set := time.Date(2005, 9, 30, 19, 00, 0, 0, time.FixedZone("CET", 7200))
	date := time.Date(2005, 9, 30, 0, 0, 0, 0, time.FixedZone("CET", 7200))

	sun := SunCal(coord, date)

	if !sun.Rise.Equal(rise) {
		t.Errorf("SunCal failed for rise, got %s, want %s", sun.Rise.String(), rise.String())
	}

	if !sun.Set.Equal(set) {
		t.Errorf("SunCal failed for dawn, got %v, want %q", sun.Set.String(), set.String())
	}

	fmt.Printf("Date %d-%02d-%02d\n", date.Year(), date.Month(), date.Day())
	fmt.Printf("Rise %02d:%02d\n", sun.Rise.Hour(), sun.Rise.Minute())
	fmt.Printf("Set  %02d:%02d\n", sun.Set.Hour(), sun.Set.Minute())
}

func TestMunich(t *testing.T) {
	coord := GeoCoord{49.0, 11.0}
	fmt.Printf("Coordinates %+v\n", coord)
	rise := time.Date(2005, 11, 18, 7, 29, 0, 0, time.FixedZone("CET", 3600))
	set := time.Date(2005, 11, 18, 16, 31, 0, 0, time.FixedZone("CET", 3600))
	date := time.Date(2005, 11, 18, 0, 0, 0, 0, time.FixedZone("CET", 3600))
	sun := SunCal(coord, date)

	if !sun.Rise.Equal(rise) {
		t.Errorf("SunCal failed for rise, got %s, want %s", sun.Rise.String(), rise.String())
	}

	if !sun.Set.Equal(set) {
		t.Errorf("SunCal failed for dawn, got %v, want %q", sun.Set.String(), set.String())
	}

	fmt.Printf("Date %d-%02d-%02d\n", date.Year(), date.Month(), date.Day())
	fmt.Printf("Rise %02d:%02d\n", sun.Rise.Hour(), sun.Rise.Minute())
	fmt.Printf("Set  %02d:%02d\n", sun.Set.Hour(), sun.Set.Minute())
}
