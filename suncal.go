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
	lat, lon float64
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

// TODO was ist DK?
// TODO was ist T?
// TODO Funktionsname
// In German: Zeitgleichung
func BerechneZeitgleichung(T float64) (float64, float64) {
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
	DK := Asin(Sin(e) * Sin(L))

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
	return dRA, DK
}

type MyTimeZone float64

const (
	MyWorldTime  MyTimeZone = 0.0
	MyNormalTime MyTimeZone = 1.0
	MySummerTime MyTimeZone = 2.0
)

func floatTimeToDumbTime(time float64) *DumbTime {
	var minutes int = int(60.0*(time-float64(int(time))) + 0.5)

	var hours int = int(time)
	if minutes >= 60.0 {
		minutes -= 60.0
		hours++
	} else if minutes < 0 {
		minutes += 60.0
		hours--
		if hours < 0.0 {
			hours += 24.0
		}
	}

	return &DumbTime{hours, minutes}
}

// Apply timezone to world time
func applyTimezone(worldTime, timezone float64) float64 {
	var t float64 = worldTime + timezone // in hours
	if t < 0.0 {
		t += 24.0
	} else if t >= 24.0 {
		t -= 24.0
	}
	return t
}

// TODO Zeitzone als Parameter
func Calc(coord GeoCoord, timezone MyTimeZone, year, month, day int) SunCal {
	const JD2000 float64 = 2451545.0
	jd := MkJulianDate(year, month, day)

	var t float64 = (jd - JD2000) / 36525.0
	const h float64 = -50.0 / 60.0 * RAD // TODO buergerlich, astronomisch oder militaerisch
	var lat float64 = coord.lat * RAD
	var lon float64 = coord.lon

	// Zeitzone := float64(timezone) // TODO TZ
	tz := float64(timezone) // TODO TZ

	var timeEqu, DK float64 = BerechneZeitgleichung(t) // TODO CalcTimeEquation
	var timeDiff float64 = 12.0 * Acos((Sin(h)-Sin(lat)*Sin(DK))/(Cos(lat)*Cos(DK))) / Pi
	var zoneTimeRise float64 = 12.0 - timeDiff - timeEqu
	var zoneTimeDawn float64 = 12.0 + timeDiff - timeEqu
	var worldTimeRise float64 = zoneTimeRise - lon/15.0
	var worldTimeDawn float64 = zoneTimeDawn - lon/15.0
	rise := applyTimezone(worldTimeRise, tz)
	dawn := applyTimezone(worldTimeDawn, tz)
	dtRise := floatTimeToDumbTime(rise)
	dtDawn := floatTimeToDumbTime(dawn)

	// TODO Ausgabe in finaler Version nicht erforderlich
	fmt.Printf("Aufgang %d:%02d\n", dtRise.hour, dtRise.minute)
	fmt.Printf("Untergang %d:%02d\n", dtDawn.hour, dtDawn.minute)

	return SunCal{*dtRise, *dtDawn}

	// Vergleich mit CalSky.com
	// Aufgang        :  7h18.4m Untergang      : 19h00.6m
}

func main() {
	var c = GeoCoord{50.0, 10.0}
	var r = Calc(c, MyNormalTime, 2005, 9, 30)
	fmt.Printf("Aufgang %d:%02d\n", r.rise.hour, r.rise.minute)
	fmt.Printf("Untergang %d:%02d\n", r.dawn.hour, r.dawn.minute)
}
