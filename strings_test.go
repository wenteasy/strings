package strings_test

import (
	"testing"

	"github.com/dustin/go-humanize"
	"github.com/wenteasy/strings"
)

func TestComma(t *testing.T) {

	tpls := []struct {
		v    int
		want string
	}{
		{1234567890, "1,234,567,890"},
		{12345678900, "12,345,678,900"},
		{1234, "1,234"},
		{123, "123"},
		{-1234, "-1,234"},
		{-123, "-123"},
		{-1234567890, "-1,234,567,890"},
		{-123456789, "-123,456,789"},
		{0, "0"},
	}

	for idx, tpl := range tpls {
		got := strings.Comma(tpl.v)
		if got != tpl.want {
			t.Errorf("[%d]%d want %s got %s", idx, tpl.v, tpl.want, got)
		}
	}
}

func BenchmarkComma(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strings.Comma(1234567890)
	}
}

func BenchmarkHumanizeComma(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		humanize.Comma(1234567890)
	}
}
