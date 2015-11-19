package suncal

import "testing"

func TestDefault(t *testing.T) {
	rise := &DumbTime{7, 18}
	set := &DumbTime{19, 00}
	sun := SunCalc(GeoCoord{50.0, 10.0}, MySummerTime, 2005, 9, 30)

	if !sun.rise.Eq(rise) {
		t.Errorf("SunCal failed for rise, got %s, want %s", sun.rise.String(), rise.String())
	}

	if !sun.set.Eq(set) {
		t.Errorf("SunCal failed for dawn, got %v, want %q", sun.set.String(), set.String())
	}
}

func TestMunich(t *testing.T) {
	sun := SunCalc(GeoCoord{49.0, 11.0}, MyNormalTime, 2015, 11, 18)

	expectedRise := DumbTime{7, 29}
	if !sun.rise.Eq(&expectedRise) {
		t.Errorf("SunCal failed for rise, got %s, want %s", sun.rise.String(), expectedRise.String())
	}

	expectedSet := DumbTime{16, 31}
	if !sun.set.Eq(&expectedSet) {
		t.Errorf("SunCal failed for dawn, got %v, want %q", sun.set.String(), expectedSet.String())
	}
}
