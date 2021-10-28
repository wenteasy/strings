package strings

import (
	"fmt"
	"strings"
)

const (
	CommaMinusSign = rune('-')
	CommaSeparator = ","
)

func Comma(v int) string {

	var b strings.Builder
	buf := fmt.Sprintf("%d", v)
	l := len(buf)

	b.Grow(l + l/3 + 1)

	//1234567890 -> 1,234,567,890

	for idx, n := range buf {
		b.WriteRune(n)
		p := l - idx - 1
		if (p != 0) && (p%3 == 0) && n != CommaMinusSign {
			b.WriteString(CommaSeparator)
		}
	}

	return b.String()
}
