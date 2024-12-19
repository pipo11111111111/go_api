package handlres

import (
	dbfunctions "facturadirecta_api/db_functions"
	"encoding/json"
	"log"
	"net/http"
)

func ListIdFaturaDirecta(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	enableCors(&w)
	var filtro dbfunctions.ProdutoFiltro
	err := json.NewDecoder(r.Body).Decode(&filtro)
	if err != nil {
		// Se houver um erro no corpo da requisição, retorna um erro 400
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}
	fatura, err := dbfunctions.GetIdFaturaDireta(filtro)
	if err != nil {
		http.Error(w, "Error fetching fatura", http.StatusInternalServerError)
		log.Printf("Error fetching fatura: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(fatura); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}
 