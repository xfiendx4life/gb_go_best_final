package csvreader_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	cs "github.com/xfiendx4life/gb_go_best_final/pkg/csvreader"
)

func TestReadHeaders(t *testing.T) {
	r := cs.NewRawData()
	rdr := strings.NewReader("iso_code,continent,location,date,total_cases")
	err := r.ReadLine(rdr)
	assert.Nil(t, err)
	assert.Equal(t, "iso_code", r.GetHeaders()[0])
}
