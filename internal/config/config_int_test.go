package config_test

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_go_best_final/internal/config"
)

func createTestFile(data []byte) {
	f, err := os.Create("test.yaml")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Fatalf("can't close file")
		}
	}()
	f.Write(data)
}

func TestReadFromFile(t *testing.T) {
	data := []byte(`timeout: 2
loglevel: 5
logfile: access.txt
targetfile: target.csv`)
	createTestFile(data)
	res, err := config.ReadFromFile("test.yaml")
	os.Remove("test.yaml")
	assert.Nil(t, err)
	assert.Equal(t, data, res)
}

func TestReadFromFileSep(t *testing.T) {
	data := []byte(`timeout: 2
loglevel: 5
logfile: access.txt
targetfile: target.csv
separator: ","
`)
	createTestFile(data)
	res, err := config.ReadFromFile("test.yaml")
	assert.Nil(t, err)
	c := config.InitConfig()
	err = c.ReadConfig(res)
	os.Remove("test.yaml")
	assert.Nil(t, err)
	assert.Equal(t, ',', c.GetSeparator())
}
func TestReadFromFileError(t *testing.T) {
	_, err := config.ReadFromFile("test.yaml")
	os.Remove("test.yaml")
	assert.NotNil(t, err)
}
