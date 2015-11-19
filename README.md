# Suncal

A simple package to calculate the rise (sunrise) and the set (sunset) of the sun with about minute accuracy.

The calculations are based on the ideas presented at http://lexikon.astronomie.info/zeitgleichung/ (german).

Samples of the results have been successfully compared to http://skycal.com/.

# License
See `LICENSE.txt` in the repository.


# Example
See examples/suncaldemo.go:
```go	
import fmt
import suncal

coords := suncal.GeoCoordinates{48.137222, 11.575556}
date := time.Now()
sun := suncal.SunCal(coords, date)
fmt.Printf("Rise %s\nSet %s\n", sun.Rise, sun.Set)
```

Example output:
```
go run examples/suncaldemo.go
Coordinates Latitude: 48.137222 Longitude 11.575556
Detailed date used 2015-11-19 18:19:46.965911373 +0100 CET
Date 2015-11-20  ↑ 07:27  ↓ 16:30  Δ 9h3m0s
Date 2015-11-21  ↑ 07:28  ↓ 16:29  Δ 9h1m0s
Date 2015-11-22  ↑ 07:29  ↓ 16:28  Δ 8h59m0s
Date 2015-11-23  ↑ 07:31  ↓ 16:27  Δ 8h56m0s
Date 2015-11-24  ↑ 07:32  ↓ 16:26  Δ 8h54m0s
Date 2015-11-25  ↑ 07:34  ↓ 16:25  Δ 8h51m0s
```

# License
See `LICENSE.txt` in the repository.

EOF
