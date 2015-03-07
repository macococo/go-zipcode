package tasks

import (
	"archive/zip"
	"bufio"
	"code.google.com/p/go.text/encoding/japanese"
	"code.google.com/p/go.text/transform"
	"encoding/csv"
	"fmt"
	"github.com/macococo/go-webbase/utils"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	BASE_URL    = "http://www.post.japanpost.jp/zipcode/dl/kogaki/zip"
	TEMP_DIR    = "./tmp"
	CSV_COLUMNS = 15
)

func ImportAll() {
	err := os.Mkdir(TEMP_DIR, 0777)
	utils.HandleError(err)

	dest := downloadZip("ken_all")
	files := unzip(dest)

	for _, path := range files {
		log.Println(path)

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

			value, err := sjisToUtf8(string(line))
			record := csvToStrings(value)

			if len(record) != CSV_COLUMNS {
				continue
			}

			for _, val := range record {
				fmt.Println(val)
			}
		}
	}
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

func unzip(src string) []string {
	files := []string{}

	r, err := zip.OpenReader(src)
	utils.HandleError(err)

	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		utils.HandleError(err)

		defer rc.Close()

		path := filepath.Join(TEMP_DIR, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			utils.HandleError(err)

			defer f.Close()

			_, err = io.Copy(f, rc)
			utils.HandleError(err)

			files = append(files, path)
		}
	}

	return files
}

func sjisToUtf8(str string) (string, error) {
	ret, err := ioutil.ReadAll(transform.NewReader(strings.NewReader(str), japanese.ShiftJIS.NewDecoder()))
	if err != nil {
		return "", err
	}
	return string(ret), err
}

func csvToStrings(str string) []string {
	reader := csv.NewReader(strings.NewReader(str))
	reader.LazyQuotes = true

	record, err := reader.Read()
	utils.HandleError(err)

	return record
}
