package suncal

import (
	"fmt"
	. "math"
)

const (
	pi2    = 2 * Pi
	JD2000 = 2451545.0
	RAD    = 0.017453292519943295769236907684886 // TODO replace with to_radians'!
)

type DumbTime struct {
	hour, minute int
}

func (dtl *DumbTime) Eq(dtr *DumbTime) bool {
	return dtl.hour == dtr.hour && dtl.minute == dtr.minute
}

func (dt *DumbTime) String() string {
	return fmt.Sprintf("%d:%02d", dt.hour, dt.minute) // TODO 00:00
}

type SunCal struct {
	rise DumbTime
	dawn DumbTime
}

// Latitude = Geographische Breite
// Longitude = Geographische Laenge

type GeoCoord struct {
	lat, long float64
}

func MkJulianDate(year, month, day int) float64 {
	const hour = 12.0
	const minute = 0.0
	const second = 0.0

	if month <= 2 {
		month = month + 12
		year = year - 1
	}

	var gregor int = (year / 400) - (year / 100) + (year / 4) // Gregorianischer Kalender
	return 2400000.5 + 365.0*float64(year) - 679004.0 + float64(gregor) +
		float64(int64(30.6001*(float64(month)+1))) + float64(day) + hour/24.0 +
		minute/1440.0 + second/86400.0

}

// TODO whats this? can be replaced by builtin?
func inPi(x float64) float64 {
	var n int = (int)(x / pi2)
	x = x - float64(n)*pi2
	if x < 0 {
		x += pi2
	}
	return x
}

// Neigung der Erdachse
func Eps(T float64) float64 {
	return RAD * (23.43929111 + (-46.8150*T-0.00059*T*T+0.001813*T*T*T)/3600.0)
}

// TODO gib mehrere Werte zurueck anstelle des pointers
// TODO was ist DK
// TODO was ist T
// TODO Funktionsname
func BerechneZeitgleichung(DK *float64, T float64) float64 {
	var RA_Mittel float64 = 18.71506921 + 2400.0513369*T + (2.5862e-5-1.72e-9*T)*T*T

	var M float64 = inPi(pi2 * (0.993133 + 99.997361*T))
	var L float64 = inPi(pi2 * (0.7859453 + M/pi2 +
		(6893.0*Sin(M)+72.0*Sin(2.0*M)+6191.2*T)/1296.0e3))

	var e float64 = Eps(T)
	var RA float64 = Atan(Tan(L) * Cos(e))

	if RA < 0.0 {
		RA += Pi
	}
	if L > Pi {
		RA += Pi
	}

	RA = 24.0 * RA / pi2
	*DK = Asin(Sin(e) * Sin(L))

	// Damit 0 <= RA_Mittel < 24
	RA_Mittel = 24.0 * inPi(pi2*RA_Mittel/24.0) / pi2

	var dRA float64 = RA_Mittel - RA
	if dRA < -12.0 {
		dRA += 24.0
	}
	if dRA > 12.0 {
		dRA -= 24.0
	}

	dRA = dRA * 1.0027379
	return dRA
}

type MyTimeZone float64

const (
	MyWorldTime  MyTimeZone = 0.0
	MyNormalTime MyTimeZone = 1.0
	MySummerTime MyTimeZone = 2.0
)

// TODO Zeitzone als Parameter
func Calc(coord GeoCoord, timezone MyTimeZone, year, month, day int) SunCal {
	const JD2000 float64 = 2451545.0
	var JD float64 = MkJulianDate(year, month, day)

	var T float64 = (JD - JD2000) / 36525.0
	const h float64 = -50.0 / 60.0 * RAD
	var B float64 = coord.lat * RAD // geographische Breite
	var GeographischeLaenge float64 = coord.long

	// const Zeitzone float64 = 0 // Weltzeit
	// const Zeitzone float64 = 1 // Winterzeit
	// const Zeitzone float64 = 2.0 // Sommerzeit
	Zeitzone := float64(timezone)

	var DK float64 = 0
	var Zeitgleichung float64 = BerechneZeitgleichung(&DK, T)

	var Zeitdifferenz float64 = 12.0 * Acos((Sin(h)-Sin(B)*Sin(DK))/(Cos(B)*Cos(DK))) / Pi
	var AufgangOrtszeit float64 = 12.0 - Zeitdifferenz - Zeitgleichung
	var UntergangOrtszeit float64 = 12.0 + Zeitdifferenz - Zeitgleichung
	var AufgangWeltzeit float64 = AufgangOrtszeit - GeographischeLaenge/15.0
	var UntergangWeltzeit float64 = UntergangOrtszeit - GeographischeLaenge/15.0

	var Aufgang float64 = AufgangWeltzeit + Zeitzone // In Stunden
	if Aufgang < 0.0 {
		Aufgang += 24.0
	} else if Aufgang >= 24.0 {
		Aufgang -= 24.0
	}

	var Untergang float64 = UntergangWeltzeit + Zeitzone
	if Untergang < 0.0 {
		Untergang += 24.0
	} else if Untergang >= 24.0 {
		Untergang -= 24.0
	}

	var AufgangMinuten int = int(60.0*(Aufgang-float64(int(Aufgang))) + 0.5)
	var AufgangStunden int = int(Aufgang)
	if AufgangMinuten >= 60.0 {
		AufgangMinuten -= 60.0
		AufgangStunden++
	} else if AufgangMinuten < 0.0 {
		AufgangMinuten += 60.0
		AufgangStunden--
		if AufgangStunden < 0.0 {
			AufgangStunden += 24.0
		}
	}

	var UntergangMinuten int = int(60.0*(Untergang-float64(int(Untergang))) + 0.5)
	var UntergangStunden int = int(Untergang)
	if UntergangMinuten >= 60.0 {
		UntergangMinuten -= 60.0
		UntergangStunden++
	} else if UntergangMinuten < 0 {
		UntergangMinuten += 60.0
		UntergangStunden--
		if UntergangStunden < 0.0 {
			UntergangStunden += 24.0
		}
	}

	// TODO Ausgabe in finaler Version nicht erforderlich
	fmt.Printf("Aufgang %d:%02d\n", AufgangStunden, AufgangMinuten)
	fmt.Printf("Untergang %d:%02d\n", UntergangStunden, UntergangMinuten)

	return SunCal{
		DumbTime{AufgangStunden, AufgangMinuten},
		DumbTime{UntergangStunden, UntergangMinuten},
	}

	// Vergleich mit CalSky.com
	// Aufgang        :  7h18.4m Untergang      : 19h00.6m
}

func main() {
	var c = GeoCoord{50.0, 10.0}
	var r = Calc(c, MyNormalTime, 2005, 9, 30)
	fmt.Printf("Aufgang %d:%02d\n", r.rise.hour, r.rise.minute)
	fmt.Printf("Untergang %d:%02d\n", r.dawn.hour, r.dawn.minute)
}
