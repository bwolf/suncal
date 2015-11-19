# Suncal

A simple package to calculate the rise and the set of the sun with about minute accuracy.

The calculations are based on the ideas presented at http://lexikon.astronomie.info/zeitgleichung/.

Samples of the results have been successfully compared to http://skycal.com/.

```go	
import fmt
import suncal

coords := suncal.GeoCoordinates{48.137222, 11.575556}
date := time.Now()
sun := suncal.SunCal(coords, date)
fmt.Printf("Rise %s\nSet %s\n", sun.Rise, sun.Set)
```
EOF
