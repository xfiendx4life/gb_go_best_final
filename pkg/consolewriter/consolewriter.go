package consolewriter

import (
	"fmt"
	"io"
	"strings"

	"github.com/xfiendx4life/gb_go_best_final/pkg/csvreader"
)

type Consolewriter interface {
	Write(io.Writer, csvreader.Table)
}

type CWriter struct{}

func NewConsoleWriter() Consolewriter {
	return &CWriter{}
}

func (cw *CWriter) Write(out io.Writer, dat csvreader.Table) {
	tab := dat.GetTable()
	maxLen := getLongestWordLen(tab) + 1
	keys := make([]string, 0, len(tab))
	for k := range tab {
		keys = append(keys, k)
	}
	var l int
	for k := range tab {
		l = len(tab[k])
		break
	}
	line := printLine(keys, maxLen)
	fmt.Fprintln(out, line)
	for i := 0; i < l; i++ {
		vals := make([]string, 0)
		for _, k := range keys {
			vals = append(vals, tab[k][i])
		}
		line = printLine(vals, maxLen)
		fmt.Fprintln(out, line)
	}

}

func printLine(data []string, size int) string {
	res := ""
	for _, item := range data {
		res += item + strings.Repeat(" ", size-len(item))
	}
	return res
}

func getLongestWordLen(data map[string][]string) int {
	m := 0
	for k, v := range data {
		if len(k) > m {
			m = len(k)
		}
		for _, item := range v {
			if len(item) > m {
				m = len(item)
			}
		}
	}
	return m
}
