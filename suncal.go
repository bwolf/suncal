/* Package suncal calculates the rise (sunrise) and the set (sunset)
   of the sun with about minute accuracy.

   The calculations are based on the ideas at
   http://lexikon.astronomie.info/zeitgleichung/ (german).

   Author: M. Geiger
   License: Apache License Version 2.0, January 2004

   Repo: http://github.com/bwolf/suncal.git
*/

package suncal

import (
	"fmt"
	. "math"
	"time"
)

const (
	// Private constants
	pi2    = 2 * Pi
	jd2000 = 2451545.0
)

type SunInfo struct {
	Rise time.Time
	Set  time.Time
}

type GeoCoordinates struct {
	Latitude  float64
	Longitude float64
}

func (c GeoCoordinates) String() string {
	return fmt.Sprintf("Latitude: %v Longitude %v", c.Latitude, c.Longitude)
}

func mkJulianDate(date time.Time) float64 {
	const hour = 12.0
	const minute = 0.0
	const second = 0.0

	year, month, day := date.Date()

	if month <= 2 {
		month = month + 12
		year = year - 1
	}

	// Gregorian calendar
	var gregor int = (year / 400) - (year / 100) + (year / 4)

	return 2400000.5 + 365.0*float64(year) - 679004.0 + float64(gregor) +
		float64(int64(30.6001*(float64(month)+1))) + float64(day) + hour/24.0 +
		minute/1440.0 + second/86400.0

}

func inPi(x float64) float64 {
	var n int = (int)(x / pi2)
	x = x - float64(n)*pi2
	if x < 0 {
		x += pi2
	}
	return x
}

// Tilt of the earth axis
func eps(t float64) float64 {
	return Pi / 180 * (23.43929111 + (-46.8150*t-0.00059*t*t+0.001813*t*t*t)/3600.0)
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
	dk := Asin(Sin(e) * Sin(l)) // declination

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

func timeFromFloatTime(date time.Time, ftime float64) time.Time {
	var minute int = int(60.0*(ftime-float64(int(ftime))) + 0.5)
	var hour int = int(ftime)

	if minute >= 60.0 {
		minute -= 60.0
		hour++
	} else if minute < 0 {
		minute += 60.0
		hour--
		if hour < 0.0 {
			hour += 24.0
		}
	}

	return time.Date(date.Year(), date.Month(), date.Day(),
		hour, minute, 0, 0,
		date.Location())
}

// Apply timezone to world time
func applyTimezone(worldTime float64, date time.Time) float64 {
	_, tzoffsetMinutes := date.Zone()
	offset := float64(tzoffsetMinutes / 3600)

	if t := worldTime + offset; t < 0.0 {
		return t + 24.0
	} else if t >= 24.0 {
		return t - 24.0
	} else {
		return t
	}
}

// Calculate sunrise and sunset for given coordinates and date for the default twilight.
func SunCal(coords GeoCoordinates, date time.Time) SunInfo {
	jd := mkJulianDate(date)

	t := (jd - jd2000) / 36525.0
	// default -50, civil -6, astronomic -18, nautic -12 arcs
	// h := float64(twilightKind) / 60.0 * Pi / 180
	h := -50.0 / 60.0 * Pi / 180
	lat := coords.Latitude * Pi / 180
	lon := coords.Longitude

	timeEqu, dk := calcTimeEquation(t)
	timeDiff := 12.0 * Acos((Sin(h)-Sin(lat)*Sin(dk))/(Cos(lat)*Cos(dk))) / Pi
	zoneTimeRise := 12.0 - timeDiff - timeEqu
	zoneTimeDawn := 12.0 + timeDiff - timeEqu
	worldTimeRise := zoneTimeRise - lon/15.0
	worldTimeDawn := zoneTimeDawn - lon/15.0
	var rise float64 = applyTimezone(worldTimeRise, date)
	var dawn float64 = applyTimezone(worldTimeDawn, date)

	dtRise := timeFromFloatTime(date, rise)
	dtDawn := timeFromFloatTime(date, dawn)

	return SunInfo{dtRise, dtDawn}

	// Compared with CalSky.com
	// Rise:  7h18.4m
	// Set:  19h00.6m
}

// EOF
