package suncal

import "testing"

func TestDefault(t *testing.T) {
	rise := &DumbTime{7, 18}
	dawn := &DumbTime{19, 00}
	r := Calc(GeoCoord{50.0, 10.0}, MySummerTime, 2005, 9, 30)

	if !r.rise.Eq(rise) {
		t.Errorf("SunCal failed for rise, got %s, want %s", r.rise.String(), rise.String())
	}

	if !r.dawn.Eq(dawn) {
		t.Errorf("SunCal failed for dawn, got %v, want %q", r.dawn.String(), dawn.String())
	}
}

func TestMunich(t *testing.T) {
	r := Calc(GeoCoord{49.0, 11.0}, MyNormalTime, 2015, 11, 18)

	expRise := DumbTime{7, 29}
	if !r.rise.Eq(&expRise) {
		t.Errorf("SunCal failed for rise, got %s, want %s", r.rise.String(), expRise.String())
	}

	expDawn := DumbTime{16, 31}
	if !r.dawn.Eq(&expDawn) {
		t.Errorf("SunCal failed for dawn, got %v, want %q", r.dawn.String(), expDawn.String())
	}
}
