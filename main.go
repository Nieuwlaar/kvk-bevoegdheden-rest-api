package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	kvkExtract "github.com/privacybydesign/kvk-extract"
	"github.com/privacybydesign/kvk-extract/models"
	"github.com/unrolled/render"
)

func main() {
	// runAll()
	r := chi.NewRouter()
	rend := render.New()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Post("/api/bevoegdheid", func(w http.ResponseWriter, r *http.Request) {
		bevoegdheidExtract := &models.BevoegdheidExtract{}
		json.NewDecoder(r.Body).Decode(&bevoegdheidExtract)

		err := kvkExtract.GetBevoegdheidExtract(bevoegdheidExtract, os.Getenv("CERTIFICATE_KVK"), os.Getenv("PRIVATE_KEY_KVK"), true, "preprd")

		if err == kvkExtract.ErrPersonNotOnExtract {
			rend.JSON(w, http.StatusNotFound, err)
			return
		} else if err != nil {
			panic(err)
		}

		rend.JSON(w, http.StatusOK, bevoegdheidExtract)
	})

	http.ListenAndServe(":3333", r)
}
