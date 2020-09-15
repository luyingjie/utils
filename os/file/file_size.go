package file

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Size(path string) int64 {
	s, e := os.Stat(path)
	if e != nil {
		return 0
	}
	return s.Size()
}

func ReadableSize(path string) string {
	return FormatSize(Size(path))
}

func StrToSize(sizeStr string) int64 {
	i := 0
	for ; i < len(sizeStr); i++ {
		if sizeStr[i] == '.' || (sizeStr[i] >= '0' && sizeStr[i] <= '9') {
			continue
		} else {
			break
		}
	}
	var (
		unit      = sizeStr[i:]
		number, _ = strconv.ParseFloat(sizeStr[:i], 64)
	)
	if unit == "" {
		return int64(number)
	}
	switch strings.ToLower(unit) {
	case "b", "bytes":
		return int64(number)
	case "k", "kb", "ki", "kib", "kilobyte":
		return int64(number * 1024)
	case "m", "mb", "mi", "mib", "mebibyte":
		return int64(number * 1024 * 1024)
	case "g", "gb", "gi", "gib", "gigabyte":
		return int64(number * 1024 * 1024 * 1024)
	case "t", "tb", "ti", "tib", "terabyte":
		return int64(number * 1024 * 1024 * 1024 * 1024)
	case "p", "pb", "pi", "pib", "petabyte":
		return int64(number * 1024 * 1024 * 1024 * 1024 * 1024)
	case "e", "eb", "ei", "eib", "exabyte":
		return int64(number * 1024 * 1024 * 1024 * 1024 * 1024 * 1024)
	case "z", "zb", "zi", "zib", "zettabyte":
		return int64(number * 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024)
	case "y", "yb", "yi", "yib", "yottabyte":
		return int64(number * 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024)
	case "bb", "brontobyte":
		return int64(number * 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024)
	}
	return -1
}

func FormatSize(raw int64) string {
	var r float64 = float64(raw)
	var t float64 = 1024
	var d float64 = 1
	if r < t {
		return fmt.Sprintf("%.2fB", r/d)
	}
	d *= 1024
	t *= 1024
	if r < t {
		return fmt.Sprintf("%.2fK", r/d)
	}
	d *= 1024
	t *= 1024
	if r < t {
		return fmt.Sprintf("%.2fM", r/d)
	}
	d *= 1024
	t *= 1024
	if r < t {
		return fmt.Sprintf("%.2fG", r/d)
	}
	d *= 1024
	t *= 1024
	if r < t {
		return fmt.Sprintf("%.2fT", r/d)
	}
	d *= 1024
	t *= 1024
	if r < t {
		return fmt.Sprintf("%.2fP", r/d)
	}
	d *= 1024
	t *= 1024
	if r < t {
		return fmt.Sprintf("%.2fE", r/d)
	}
	d *= 1024
	t *= 1024
	if r < t {
		return fmt.Sprintf("%.2fZ", r/d)
	}
	d *= 1024
	t *= 1024
	if r < t {
		return fmt.Sprintf("%.2fY", r/d)
	}
	d *= 1024
	t *= 1024
	if r < t {
		return fmt.Sprintf("%.2fBB", r/d)
	}
	return "TooLarge"
}
