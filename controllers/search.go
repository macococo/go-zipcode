package controllers

import (
	http_ "github.com/macococo/go-webbase/http"
	"github.com/macococo/go-zipcode/models"
	"net/http"
	"strings"
)

func SearchController(w http.ResponseWriter, r *http.Request) {
	zipcode := http_.GetRequestParam(r, "zipcode", "")
	zipcode = strings.Replace(zipcode, "-", "", -1)

	address := models.GetAddress(zipcode)

	http_.WriteJsonResponse(w, address)
}
