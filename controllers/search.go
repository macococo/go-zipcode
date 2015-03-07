package controllers

import (
	"fmt"
	"github.com/macococo/go-webbase/utils"
	"net/http"
)

func SearchController(w http.ResponseWriter, r *http.Request) {
	zipcode := utils.GetParam(r, "zipcode", "")
	fmt.Println("zipcode: " + zipcode)
}
