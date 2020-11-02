package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"cosmogony.com/sales/internal/models"
	dal "cosmogony.com/sales/internal/storage"
	"github.com/gorilla/mux"
)

func CreateCotizacion(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	cotizacion := &models.Cotizacion{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(cotizacion); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	// Guardar en persistent storage
	result, err := dal.CreateCotizacion(cotizacion)
	if err != nil {
		// w.WriteHeader(http.StatusInternalServerError)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	w.WriteHeader(http.StatusCreated)
	result = strings.ReplaceAll(result, `"`, `'`)
	w.Write([]byte(`{"message": "` + result + `"}`))
}

func ReadCotizacion(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Not an integer"}`))
		return
	}

	// Guardar en persistent storage
	cotizacion, err := dal.ReadCotizacion(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	v := r.URL.Query()
	fmt.Println(v)

	encoder := json.NewEncoder(w)
	if err = encoder.Encode(*cotizacion); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "` + err.Error() + `"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
}
