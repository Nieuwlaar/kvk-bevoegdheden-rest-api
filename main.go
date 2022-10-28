package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	kvkBevoegdheden "github.com/kvk-innovatie/kvk-bevoegdheden"
	"github.com/kvk-innovatie/kvk-bevoegdheden/models"
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

	r.Get("/api/test-inschrijvingen", func(w http.ResponseWriter, r *http.Request) {
		files, err := ioutil.ReadDir("./cache-inschrijvingen/")
		if err != nil {
			rend.JSON(w, http.StatusNotFound, err)
		}
		fileNames := []string{}
		for _, file := range files {
			if !file.IsDir() {
				fn := strings.TrimSuffix(file.Name(), ".xml")
				fileNames = append(fileNames, fn)
			}
		}
		rend.JSON(w, http.StatusOK, fileNames)
	})

	r.Post("/api/bevoegdheid/{kvkNummer}", func(w http.ResponseWriter, r *http.Request) {
		kvkNummer := chi.URLParam(r, "kvkNummer")
		identityNP := models.IdentityNP{}
		json.NewDecoder(r.Body).Decode(&identityNP)

		bevoegdheidResponse, err := kvkBevoegdheden.GetBevoegdheid(kvkNummer, identityNP, os.Getenv("CERTIFICATE_KVK"), os.Getenv("PRIVATE_KEY_KVK"), true, "preprd")

		if err == kvkBevoegdheden.ErrInschrijvingNotFound {
			rend.JSON(w, http.StatusNotFound, err)
			return
		} else if err == kvkBevoegdheden.ErrInvalidInput {
			rend.JSON(w, http.StatusBadRequest, err)
			return
		} else if err != nil {
			rend.JSON(w, http.StatusInternalServerError, err)
			return
		}

		rend.JSON(w, http.StatusOK, bevoegdheidResponse)
	})

	http.ListenAndServe(":3333", r)
}
