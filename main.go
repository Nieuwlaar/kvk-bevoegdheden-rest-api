package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
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
	r := chi.NewRouter()
	rend := render.New()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(corsMiddleware) // Add the CORS middleware

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
		handleBevoegdheid(w, r, rend)
	})

	r.Post("/api/signatory-rights/{kvkNummer}", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Signatory Rights Request")
		handleSignatoryRight(w, r, rend)
	})

	r.Post("/api/company-certificate/{kvkNummer}", func(w http.ResponseWriter, r *http.Request) {
		handleCompanyCertificate(w, r, rend)
	})

	r.Post("/api/lpid/{kvkNummer}", func(w http.ResponseWriter, r *http.Request) {
		handleLPID(w, r, rend)
	})

	http.ListenAndServe(":3333", r)
}

func handleBevoegdheid(w http.ResponseWriter, r *http.Request, rend *render.Render) {
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
}

func handleSignatoryRight(w http.ResponseWriter, r *http.Request, rend *render.Render) {
	kvkNummer := chi.URLParam(r, "kvkNummer")
	var person struct {
		Geslachtsnaam            string `json:"geslachtsnaam"`
		Voornamen                string `json:"voornamen"`
		Geboortedatum            string `json:"geboortedatum"`
		VoorvoegselGeslachtsnaam string `json:"voorvoegselGeslachtsnaam"`
	}

	// Read the request body
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		rend.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}
	r.Body.Close()

	// Decode the request body into the person struct
	if err := json.Unmarshal(bodyBytes, &person); err != nil {
		log.Println("Error decoding person:", err)
		rend.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid person data"})
		return
	}

	// Decode the request body into the identityNP struct
	identityNP := models.IdentityNP{}
	if err := json.Unmarshal(bodyBytes, &identityNP); err != nil {
		log.Println("Error decoding identityNP:", err)
		rend.JSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid identityNP data"})
		return
	}

	// Log the incoming person details
	log.Println("Received signatory rights request for:", person.Voornamen, person.VoorvoegselGeslachtsnaam, person.Geslachtsnaam, person.Geboortedatum)

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

	// Extracting necessary fields from the bevoegdheidResponse
	bevoegdheidUittreksel := bevoegdheidResponse.BevoegdheidUittreksel
	isAuthorized := false
	fullNameRequest := strings.TrimSpace(person.Voornamen + " " + person.Geslachtsnaam)
	if person.VoorvoegselGeslachtsnaam != "" {
		fullNameRequest = strings.TrimSpace(person.Voornamen + " " + person.VoorvoegselGeslachtsnaam + " " + person.Geslachtsnaam)
	}
	log.Println("Fullnames Request: " + fullNameRequest)
	matchedFullName := ""
	for _, personEntry := range bevoegdheidUittreksel.AlleFunctionarissen {
		fullName := strings.Join(strings.Fields(strings.Join([]string{personEntry.Voornamen, personEntry.VoorvoegselGeslachtsnaam, personEntry.Geslachtsnaam}, " ")), " ")
		log.Println("Fullnames XML: " + fullName)
		log.Println("Geboortedatum XML: " + personEntry.Geboortedatum)
		log.Println("Geboortedatum Request: " + person.Geboortedatum)
		if personEntry.Geboortedatum == person.Geboortedatum && fullName == strings.TrimSpace(fullNameRequest) && (personEntry.Interpretatie.IsBevoegd == "Ja" || personEntry.Interpretatie.IsBevoegd == "Yes") {
			isAuthorized = true
			matchedFullName = fullName
			break
		}
	}

	signatoryRightResponse := map[string]interface{}{
		"data": map[string]interface{}{
			"full_name":         matchedFullName,
			"date_of_birth":     person.Geboortedatum,
			"is_authorized":     isAuthorized,
			"id":                "NLNHR." + bevoegdheidUittreksel.KvkNummer, // Conversion to EUID
			"legal_person_name": bevoegdheidUittreksel.Naam,
			"legal_form":        bevoegdheidUittreksel.PersoonRechtsvorm,
		},
		"metadata": generateMetadata(),
	}
	rend.JSON(w, http.StatusOK, signatoryRightResponse)
}

func handleLPID(w http.ResponseWriter, r *http.Request, rend *render.Render) {
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

	// Extracting necessary fields from the bevoegdheidResponse
	bevoegdheidUittreksel := bevoegdheidResponse.BevoegdheidUittreksel

	lpidResponse := map[string]interface{}{
		"data": map[string]interface{}{
			"id":                "NLNHR." + bevoegdheidUittreksel.KvkNummer, // Conversion to EUID
			"legal_person_name": bevoegdheidUittreksel.Naam,
			"legal_form":        bevoegdheidUittreksel.PersoonRechtsvorm,
		},
		"metadata": generateMetadata(),
	}

	rend.JSON(w, http.StatusOK, lpidResponse)
}

func handleCompanyCertificate(w http.ResponseWriter, r *http.Request, rend *render.Render) {
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

	// Extracting necessary fields from the bevoegdheidResponse
	bevoegdheidUittreksel := bevoegdheidResponse.BevoegdheidUittreksel

	// Pruning authorized_persons data and combining names
	prunedAuthorizedPersons := []map[string]interface{}{}
	for _, person := range bevoegdheidUittreksel.AlleFunctionarissen {
		fullName := strings.Join(strings.Fields(strings.Join([]string{person.Voornamen, person.VoorvoegselGeslachtsnaam, person.Geslachtsnaam}, " ")), " ")
		var isAuthorized string
		switch person.Interpretatie.IsBevoegd {
		case "Ja":
			isAuthorized = "Yes"
		case "Nee":
			isAuthorized = "No"
		case "Niet vastgesteld":
			isAuthorized = "Not determined"
		default:
			isAuthorized = "Unknown"
		}

		prunedPerson := map[string]interface{}{
			"full_name":     fullName,
			"date_of_birth": person.Geboortedatum,
			"interpretatie": map[string]interface{}{
				"isAuthorized": isAuthorized,
			},
		}
		prunedAuthorizedPersons = append(prunedAuthorizedPersons, prunedPerson)
	}

	companyCertificate := map[string]interface{}{
		"id":                   "NLNHR." + bevoegdheidUittreksel.KvkNummer, // Conversion to EUID
		"legal_person_name":    bevoegdheidUittreksel.Naam,
		"legal_form":           bevoegdheidUittreksel.PersoonRechtsvorm,
		"registration_number":  bevoegdheidUittreksel.KvkNummer,
		"registered_country":   "NL", // Assuming the Member State is the Netherlands, modify if needed
		"registered_office":    bevoegdheidUittreksel.Adres,
		"postal_address":       bevoegdheidUittreksel.Adres, // Assuming postal address is the same as registered office, modify if needed
		"electronic_address":   bevoegdheidUittreksel.EmailAdres,
		"date_of_registration": bevoegdheidUittreksel.RegistratieAanvang,
		"capital_subscribed":   bevoegdheidUittreksel.BijzondereRechtstoestand, // Assuming this field maps correctly, modify if needed
		"status":               bevoegdheidUittreksel.BijzondereRechtstoestand, // Assuming this field maps correctly, modify if needed
		"authorized_persons":   prunedAuthorizedPersons,
		"object":               bevoegdheidUittreksel.SbiActiviteit, // Assuming this field maps correctly, modify if needed
	}

	companyCertificateResponse := map[string]interface{}{
		"data":     companyCertificate,
		"metadata": generateMetadata(),
	}

	rend.JSON(w, http.StatusOK, companyCertificateResponse)
}

func generateMetadata() map[string]interface{} {
	return map[string]interface{}{
		"issuing_authority_name": "Kamer van Koophandel",
		"issuer_id":              "NLNHR.59581883", // Kamer van Koophandel's EUID
		"issuing_country":        "NL",
		"issuance_date":          "2022-06-15T15:35:52.687Z",  // Example value, should be altered when the actual issuance is done
		"expiry_date":            "2025-06-15T15:35:52.687Z",  // Example value, should be altered when the actual issuance is done
		"schema":                 "http://schema.example.com", // Example value
		"revocation_information": "http://revoke.example.com", // Example value
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
