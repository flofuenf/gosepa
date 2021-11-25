package lib

import (
	"strconv"
	"strings"
)

// DecimalsNumber returns the number of decimals in a float
func DecimalsNumber(f float64) int {
	s := strconv.FormatFloat(f, 'f', -1, 64)
	p := strings.Split(s, ".")
	if len(p) < 2 {
		return 0
	}
	return len(p[1])
}

// ToCents returns the cents representation in int64
func ToCents(f float64) (int64, error) {
	s := strconv.FormatFloat(f, 'f', 2, 64)
	sc := strings.Replace(s, ".", "", 1)
	return strconv.ParseInt(sc, 10, 64)
}

// ToEuro returns the euro representation in float64
func ToEuro(i int64) (float64, error) {
	d := strconv.FormatInt(i, 10)
	df := d[:len(d)-2] + "." + d[len(d)-2:]
	return strconv.ParseFloat(df, 64)
}
