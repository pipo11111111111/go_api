package handlres

import (
	dbfunctions "facturadirecta_api/db_functions"
	b2b "facturadirecta_api/facturadirecta"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func SendByEmail(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	enableCors(&w)
	var filtro dbfunctions.ProdutoFiltro

	// Decodificar o filtro da requisição
	if err := json.NewDecoder(r.Body).Decode(&filtro); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	// Buscar fatura do banco de dados
	id, err := dbfunctions.GetIdFaturaDireta(filtro)
	if err != nil {
		http.Error(w, "Error fetching fatura", http.StatusInternalServerError)
		log.Printf("Error fetching clients: %v", err)
		return
	}
	//Email producao
	/* email:= id[0].Email */
	//Email teste
	email :="filiperodrigues@ibrain.pt"
	 contact:= b2b.EmailRequest{
		To: []string{email},
	 }

	// Preparar as linhas da fatura para o formato esperado
	 
	log.Printf("mandar %s", contact)
	// Chamar a função CreateInvoice
	if err := b2b.SendByEmail(contact,id[0].IdFacturaDirecta ); err != nil {
		http.Error(w, "Failed to create invoice: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to create invoice %v: %v", contact, err)
 
		return
	}

	// Responder com sucesso
/* 	fmt.Fprintf(w, "Invoice sended successfully") */
	// Retornar resposta de sucesso
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": fmt.Sprintf("Fatura %s  enviada com Sucesso para %s ",id[0].IdFacturaDirecta, id[0].Email),
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}
