package str

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	iutil "github.com/luyingjie/utils/internal/util"

	"github.com/luyingjie/utils/conv"

	"github.com/luyingjie/utils/generates/rand"
)

func Replace(origin, search, replace string, count ...int) string {
	n := -1
	if len(count) > 0 {
		n = count[0]
	}
	return strings.Replace(origin, search, replace, n)
}

func ReplaceI(origin, search, replace string, count ...int) string {
	n := -1
	if len(count) > 0 {
		n = count[0]
	}
	if n == 0 {
		return origin
	}
	length := len(search)
	searchLower := strings.ToLower(search)
	for {
		originLower := strings.ToLower(origin)
		if pos := strings.Index(originLower, searchLower); pos != -1 {
			origin = origin[:pos] + replace + origin[pos+length:]
			if n--; n == 0 {
				break
			}
		} else {
			break
		}
	}
	return origin
}

func Count(s, substr string) int {
	return strings.Count(s, substr)
}

func CountI(s, substr string) int {
	return strings.Count(ToLower(s), ToLower(substr))
}

func ReplaceByArray(origin string, array []string) string {
	for i := 0; i < len(array); i += 2 {
		if i+1 >= len(array) {
			break
		}
		origin = Replace(origin, array[i], array[i+1])
	}
	return origin
}

func ReplaceIByArray(origin string, array []string) string {
	for i := 0; i < len(array); i += 2 {
		if i+1 >= len(array) {
			break
		}
		origin = ReplaceI(origin, array[i], array[i+1])
	}
	return origin
}

func ReplaceByMap(origin string, replaces map[string]string) string {
	return iutil.ReplaceByMap(origin, replaces)
}

func ReplaceIByMap(origin string, replaces map[string]string) string {
	for k, v := range replaces {
		origin = ReplaceI(origin, k, v)
	}
	return origin
}

func ToLower(s string) string {
	return strings.ToLower(s)
}

func ToUpper(s string) string {
	return strings.ToUpper(s)
}

func UcFirst(s string) string {
	return iutil.UcFirst(s)
}

func LcFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	if IsLetterUpper(s[0]) {
		return string(s[0]+32) + s[1:]
	}
	return s
}

func UcWords(str string) string {
	return strings.Title(str)
}

func IsLetterLower(b byte) bool {
	return iutil.IsLetterLower(b)
}

func IsLetterUpper(b byte) bool {
	return iutil.IsLetterUpper(b)
}

func IsNumeric(s string) bool {
	return iutil.IsNumeric(s)
}

func SubStr(str string, start int, length ...int) (substr string) {
	lth := len(str)

	// Simple border checks.
	if start < 0 {
		start = 0
	}
	if start >= lth {
		start = lth
	}
	end := lth
	if len(length) > 0 {
		end = start + length[0]
		if end < start {
			end = lth
		}
	}
	if end > lth {
		end = lth
	}
	return str[start:end]
}

func SubStrRune(str string, start int, length ...int) (substr string) {
	rs := []rune(str)
	lth := len(rs)

	if start < 0 {
		start = 0
	}
	if start >= lth {
		start = lth
	}
	end := lth
	if len(length) > 0 {
		end = start + length[0]
		if end < start {
			end = lth
		}
	}
	if end > lth {
		end = lth
	}
	return string(rs[start:end])
}

func StrLimit(str string, length int, suffix ...string) string {
	if len(str) < length {
		return str
	}
	addStr := "..."
	if len(suffix) > 0 {
		addStr = suffix[0]
	}
	return str[0:length] + addStr
}

func StrLimitRune(str string, length int, suffix ...string) string {
	rs := []rune(str)
	if len(rs) < length {
		return str
	}
	addStr := "..."
	if len(suffix) > 0 {
		addStr = suffix[0]
	}
	return string(rs[0:length]) + addStr
}

func Reverse(str string) string {
	runes := []rune(str)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func NumberFormat(number float64, decimals int, decPoint, thousandsSep string) string {
	neg := false
	if number < 0 {
		number = -number
		neg = true
	}
	// Will round off
	str := fmt.Sprintf("%."+strconv.Itoa(decimals)+"F", number)
	prefix, suffix := "", ""
	if decimals > 0 {
		prefix = str[:len(str)-(decimals+1)]
		suffix = str[len(str)-decimals:]
	} else {
		prefix = str
	}
	sep := []byte(thousandsSep)
	n, l1, l2 := 0, len(prefix), len(sep)
	c := (l1 - 1) / 3
	tmp := make([]byte, l2*c+l1)
	pos := len(tmp) - 1
	for i := l1 - 1; i >= 0; i, n, pos = i-1, n+1, pos-1 {
		if l2 > 0 && n > 0 && n%3 == 0 {
			for j := range sep {
				tmp[pos] = sep[l2-j-1]
				pos--
			}
		}
		tmp[pos] = prefix[i]
	}
	s := string(tmp)
	if decimals > 0 {
		s += decPoint + suffix
	}
	if neg {
		s = "-" + s
	}

	return s
}

func ChunkSplit(body string, chunkLen int, end string) string {
	if end == "" {
		end = "\r\n"
	}
	runes, endRunes := []rune(body), []rune(end)
	l := len(runes)
	if l <= 1 || l < chunkLen {
		return body + end
	}
	ns := make([]rune, 0, len(runes)+len(endRunes))
	for i := 0; i < l; i += chunkLen {
		if i+chunkLen > l {
			ns = append(ns, runes[i:]...)
		} else {
			ns = append(ns, runes[i:i+chunkLen]...)
		}
		ns = append(ns, endRunes...)
	}
	return string(ns)
}

func Compare(a, b string) int {
	return strings.Compare(a, b)
}

func Equal(a, b string) bool {
	return strings.EqualFold(a, b)
}

func Fields(str string) []string {
	return strings.Fields(str)
}

func HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

func CountWords(str string) map[string]int {
	m := make(map[string]int)
	buffer := bytes.NewBuffer(nil)
	for _, r := range []rune(str) {
		if unicode.IsSpace(r) {
			if buffer.Len() > 0 {
				m[buffer.String()]++
				buffer.Reset()
			}
		} else {
			buffer.WriteRune(r)
		}
	}
	if buffer.Len() > 0 {
		m[buffer.String()]++
	}
	return m
}

func CountChars(str string, noSpace ...bool) map[string]int {
	m := make(map[string]int)
	countSpace := true
	if len(noSpace) > 0 && noSpace[0] {
		countSpace = false
	}
	for _, r := range []rune(str) {
		if !countSpace && unicode.IsSpace(r) {
			continue
		}
		m[string(r)]++
	}
	return m
}

func WordWrap(str string, width int, br string) string {
	if br == "" {
		br = "\n"
	}
	var (
		current           int
		wordBuf, spaceBuf bytes.Buffer
		init              = make([]byte, 0, len(str))
		buf               = bytes.NewBuffer(init)
	)
	for _, char := range []rune(str) {
		if char == '\n' {
			if wordBuf.Len() == 0 {
				if current+spaceBuf.Len() > width {
					current = 0
				} else {
					current += spaceBuf.Len()
					spaceBuf.WriteTo(buf)
				}
				spaceBuf.Reset()
			} else {
				current += spaceBuf.Len() + wordBuf.Len()
				spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
			}
			buf.WriteRune(char)
			current = 0
		} else if unicode.IsSpace(char) {
			if spaceBuf.Len() == 0 || wordBuf.Len() > 0 {
				current += spaceBuf.Len() + wordBuf.Len()
				spaceBuf.WriteTo(buf)
				spaceBuf.Reset()
				wordBuf.WriteTo(buf)
				wordBuf.Reset()
			}
			spaceBuf.WriteRune(char)
		} else {
			wordBuf.WriteRune(char)
			if current+spaceBuf.Len()+wordBuf.Len() > width && wordBuf.Len() < width {
				buf.WriteString(br)
				current = 0
				spaceBuf.Reset()
			}
		}
	}

	if wordBuf.Len() == 0 {
		if current+spaceBuf.Len() <= width {
			spaceBuf.WriteTo(buf)
		}
	} else {
		spaceBuf.WriteTo(buf)
		wordBuf.WriteTo(buf)
	}
	return buf.String()
}

func RuneLen(str string) int {
	return LenRune(str)
}

func LenRune(str string) int {
	return utf8.RuneCountInString(str)
}

func Repeat(input string, multiplier int) string {
	return strings.Repeat(input, multiplier)
}

func Str(haystack string, needle string) string {
	if needle == "" {
		return ""
	}
	idx := strings.Index(haystack, needle)
	if idx == -1 {
		return ""
	}
	return haystack[idx+len([]byte(needle))-1:]
}

func Shuffle(str string) string {
	runes := []rune(str)
	s := make([]rune, len(runes))
	for i, v := range rand.Perm(len(runes)) {
		s[i] = runes[v]
	}
	return string(s)
}

func Split(str, delimiter string) []string {
	return strings.Split(str, delimiter)
}

func SplitAndTrim(str, delimiter string, characterMask ...string) []string {
	array := make([]string, 0)
	for _, v := range strings.Split(str, delimiter) {
		v = Trim(v, characterMask...)
		if v != "" {
			array = append(array, v)
		}
	}
	return array
}

func SplitAndTrimSpace(str, delimiter string) []string {
	array := make([]string, 0)
	for _, v := range strings.Split(str, delimiter) {
		v = strings.TrimSpace(v)
		if v != "" {
			array = append(array, v)
		}
	}
	return array
}

func Join(array []string, sep string) string {
	return strings.Join(array, sep)
}

func JoinAny(array interface{}, sep string) string {
	return strings.Join(conv.Strings(array), sep)
}

func Explode(delimiter, str string) []string {
	return Split(str, delimiter)
}

func Implode(glue string, pieces []string) string {
	return strings.Join(pieces, glue)
}

func Chr(ascii int) string {
	return string([]byte{byte(ascii % 256)})
}

func Ord(char string) int {
	return int(char[0])
}

func HideStr(str string, percent int, hide string) string {
	array := strings.Split(str, "@")
	if len(array) > 1 {
		str = array[0]
	}
	var (
		rs       = []rune(str)
		length   = len(rs)
		mid      = math.Floor(float64(length / 2))
		hideLen  = int(math.Floor(float64(length) * (float64(percent) / 100)))
		start    = int(mid - math.Floor(float64(hideLen)/2))
		hideStr  = []rune("")
		hideRune = []rune(hide)
	)
	for i := 0; i < hideLen; i++ {
		hideStr = append(hideStr, hideRune...)
	}
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(string(rs[0:start]))
	buffer.WriteString(string(hideStr))
	buffer.WriteString(string(rs[start+hideLen:]))
	if len(array) > 1 {
		buffer.WriteString("@" + array[1])
	}
	return buffer.String()
}

func Nl2Br(str string, isXhtml ...bool) string {
	r, n, runes := '\r', '\n', []rune(str)
	var br []byte
	if len(isXhtml) > 0 && isXhtml[0] {
		br = []byte("<br />")
	} else {
		br = []byte("<br>")
	}
	skip := false
	length := len(runes)
	var buf bytes.Buffer
	for i, v := range runes {
		if skip {
			skip = false
			continue
		}
		switch v {
		case n, r:
			if (i+1 < length) && (v == r && runes[i+1] == n) || (v == n && runes[i+1] == r) {
				buf.Write(br)
				skip = true
				continue
			}
			buf.Write(br)
		default:
			buf.WriteRune(v)
		}
	}
	return buf.String()
}

func AddSlashes(str string) string {
	var buf bytes.Buffer
	for _, char := range str {
		switch char {
		case '\'', '"', '\\':
			buf.WriteRune('\\')
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

func StripSlashes(str string) string {
	var buf bytes.Buffer
	l, skip := len(str), false
	for i, char := range str {
		if skip {
			skip = false
		} else if char == '\\' {
			if i+1 < l && str[i+1] == '\\' {
				skip = true
			}
			continue
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

func QuoteMeta(str string, chars ...string) string {
	var buf bytes.Buffer
	for _, char := range str {
		if len(chars) > 0 {
			for _, c := range chars[0] {
				if c == char {
					buf.WriteRune('\\')
					break
				}
			}
		} else {
			switch char {
			case '.', '+', '\\', '(', '$', ')', '[', '^', ']', '*', '?':
				buf.WriteRune('\\')
			}
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

func SearchArray(a []string, s string) int {
	for i, v := range a {
		if s == v {
			return i
		}
	}
	return -1
}

func InArray(a []string, s string) bool {
	return SearchArray(a, s) != -1
}
