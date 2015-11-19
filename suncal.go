package suncal

import (
	"fmt"
	. "math"
)

const (
	pi2    = 2 * Pi
	jd2000 = 2451545.0
	rad    = 0.017453292519943295769236907684886 // TODO replace with to_radians'!
)

type DumbTime struct {
	hour, minute int
}

func NewDumbTime(time float64) DumbTime {
	return floatTimeToDumbTime(time)
}

func (lhs *DumbTime) Eq(rhs *DumbTime) bool {
	return lhs.hour == rhs.hour && lhs.minute == rhs.minute
}

func (dt *DumbTime) String() string {
	return fmt.Sprintf("%02d:%02d", dt.hour, dt.minute) // TODO 00:00
}

type SunInfo struct {
	rise DumbTime
	set  DumbTime
}

type GeoCoord struct {
	lat float64 // Latitude, geographische Breite
	lon float64 // Longitude, geographische Laenge
}

// TODO replace with go's own timezone if possible or change type to float64 and name it timezone offset
type MyTimeZone float64

const (
	MyWorldTime  MyTimeZone = 0.0
	MyNormalTime MyTimeZone = 1.0
	MySummerTime MyTimeZone = 2.0
)

// TODO use go date
// TODO function name
// TODO function name indicating float64 result, as special case of julian date in go
func mkJulianDate(year, month, day int) float64 {
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
func eps(t float64) float64 {
	return rad * (23.43929111 + (-46.8150*t-0.00059*t*t+0.001813*t*t*t)/3600.0)
}

// https://en.wikipedia.org/wiki/Equation_of_time
func calcTimeEquation(t float64) (float64, float64) {
	raMid := 18.71506921 + 2400.0513369*t + (2.5862e-5-1.72e-9*t)*t*t
	m := inPi(pi2 * (0.993133 + 99.997361*t))
	l := inPi(pi2 * (0.7859453 + m/pi2 + (6893.0*Sin(m)+72.0*Sin(2.0*m)+6191.2*t)/1296.0e3))
	e := eps(t)
	ra := Atan(Tan(l) * Cos(e))

	if ra < 0.0 {
		ra += Pi
	}
	if l > Pi {
		ra += Pi
	}

	ra = 24.0 * ra / pi2
	dk := Asin(Sin(e) * Sin(l))

	// Ensure 0 <= raMid < 24
	raMid = 24.0 * inPi(pi2*raMid/24.0) / pi2

	dRA := raMid - ra
	if dRA < -12.0 {
		dRA += 24.0
	} else if dRA > 12.0 {
		dRA -= 24.0
	}

	dRA = dRA * 1.0027379
	return dRA, dk
}

func floatTimeToDumbTime(time float64) DumbTime {
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

	return DumbTime{hours, minutes}
}

// Apply timezone to world time
func applyTimezone(worldTime, timezone float64) float64 {
	if t := worldTime + timezone; t < 0.0 {
		return t + 24.0
	} else if t >= 24.0 {
		return t - 24.0
	} else {
		return t
	}
}

// TODO function name
// TODO use go date
// TODO use go timezone?
func SunCalc(coord GeoCoord, timezone MyTimeZone, year, month, day int) SunInfo {
	jd := mkJulianDate(year, month, day)

	t := (jd - jd2000) / 36525.0
	const h = -50.0 / 60.0 * rad // TODO buergerlich, astronomisch oder militaerisch
	lat := coord.lat * rad
	lon := coord.lon

	// Zeitzone := float64(timezone) // TODO TZ
	tz := float64(timezone) // TODO TZ

	timeEqu, dk := calcTimeEquation(t)
	timeDiff := 12.0 * Acos((Sin(h)-Sin(lat)*Sin(dk))/(Cos(lat)*Cos(dk))) / Pi
	zoneTimeRise := 12.0 - timeDiff - timeEqu
	zoneTimeDawn := 12.0 + timeDiff - timeEqu
	worldTimeRise := zoneTimeRise - lon/15.0
	worldTimeDawn := zoneTimeDawn - lon/15.0
	rise := applyTimezone(worldTimeRise, tz)
	dawn := applyTimezone(worldTimeDawn, tz)
	dtRise := NewDumbTime(rise)
	dtDawn := NewDumbTime(dawn)

	// TODO Ausgabe in finaler Version nicht erforderlich
	fmt.Printf("Aufgang %d:%02d\n", dtRise.hour, dtRise.minute)
	fmt.Printf("Untergang %d:%02d\n", dtDawn.hour, dtDawn.minute)

	return SunInfo{dtRise, dtDawn}

	// Vergleich mit CalSky.com
	// Aufgang        :  7h18.4m Untergang      : 19h00.6m
}

func main() {
	coord := GeoCoord{50.0, 10.0}
	sun := SunCalc(coord, MyNormalTime, 2005, 9, 30)
	fmt.Printf("Aufgang %d:%02d\n", sun.rise.hour, sun.rise.minute)
	fmt.Printf("Untergang %d:%02d\n", sun.set.hour, sun.set.minute)
}
