package models

import (
	"bufio"
	"encoding/csv"
	io_ "github.com/macococo/go-webbase/io"
	"github.com/macococo/go-webbase/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	BASE_URL    = "http://www.post.japanpost.jp/zipcode/dl/kogaki/zip"
	TEMP_DIR    = "./tmp"
	CSV_COLUMNS = 15
)

var (
	cache map[string]*Address
)

type Address struct {
	Zipcode  string `json:"zipcode"`
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	Address3 string `json:"address3"`
}

func ReloadAddressCache(f func()) {
	newCache := make(map[string]*Address)

	err := os.Mkdir(TEMP_DIR, 0777)
	utils.HandleError(err)

	dest := downloadZip("ken_all")
	files := io_.Unzip(dest, TEMP_DIR)

	log.Println("reloading address...")

	for _, path := range files {
		file, err := os.Open(path)
		utils.HandleError(err)

		defer file.Close()

		reader := bufio.NewReaderSize(file, 4096)
		for {
			line, _, err := reader.ReadLine()
			if err == io.EOF {
				break
			} else if err != nil {
				utils.HandleError(err)
			}

			value, err := utils.ShiftJISToUtf8(string(line))
			record := csvToStrings(value)

			if len(record) != CSV_COLUMNS {
				continue
			}

			address := NewAddressByRecord(record)
			newCache[address.Zipcode] = address
		}
	}

	log.Println("reloaded address.")

	cache = newCache

	f()
}

func NewAddressByRecord(record []string) *Address {
	address := Address{}

	address.Zipcode = record[2]
	address.Address1 = record[6]
	address.Address2 = record[7]
	address.Address3 = record[8]

	return &address
}

func GetAddress(zipCode string) *Address {
	if cache == nil {
		return nil
	}
	return cache[zipCode]
}

func downloadZip(name string) string {
	fileName := name + ".zip"
	url := BASE_URL + "/" + fileName

	response, err := http.Get(url)
	utils.HandleError(err)

	body, err := ioutil.ReadAll(response.Body)
	utils.HandleError(err)

	dest := TEMP_DIR + "/" + fileName
	file, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0666)
	utils.HandleError(err)

	defer file.Close()

	file.Write(body)

	return dest
}

func csvToStrings(str string) []string {
	reader := csv.NewReader(strings.NewReader(str))
	reader.LazyQuotes = true

	record, err := reader.Read()
	utils.HandleError(err)

	return record
}
