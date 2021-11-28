package domain

type Table interface {
	Read(csv_table []byte) (n int, err error)
	ParseQuery(query Query) Table
}

type Query interface {
	ParseRawQuery(rawQuery string) []string // to postfix form
}