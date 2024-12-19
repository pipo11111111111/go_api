package handlres

import (
	"facturadirecta_api/db"
	dbfunctions "facturadirecta_api/db_functions"
	b2b "facturadirecta_api/facturadirecta"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func InserClient(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	enableCors(&w)
	var filtro dbfunctions.ProdutoFiltro

	if err := json.NewDecoder(r.Body).Decode(&filtro); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}
	clientes, err := dbfunctions.GetClient(filtro)
	if err != nil {
		http.Error(w, "Error fetching clients", http.StatusInternalServerError)
		log.Printf("Error fetching clients: %v", err)
		return
	}
	// Verificar se o cliente já existe na tabela facturadirecta_clientes
	conn, err := db.OpenConnection()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		log.Printf("Error connecting to database: %v", err)
		return
	}
	defer conn.Close()

	// Verificar se o cliente já existe no banco de dados
	var existingClientID string
	query := `SELECT codclifacturadirecta FROM facturadirecta_clientes WHERE codcli = $1`
	err = conn.QueryRow(query, clientes[0].Codcli).Scan(&existingClientID)
	if err == nil {
		// Se o cliente já existir
		http.Error(w, fmt.Sprintf("Cliente com código %d já existe na Facturadirecta.", clientes[0].Codcli), http.StatusConflict)
		log.Printf("Cliente já existe na Facturadirecta: %d", clientes[0].Codcli)
		return
	} else if err.Error() != "sql: no rows in result set" {
		// Se ocorrer outro erro ao verificar no banco
		http.Error(w, "Error checking if client exists", http.StatusInternalServerError)
		log.Printf("Error checking if client exists: %v", err)
		return
	}

	contact := b2b.Content{
		Type: "contact",
		Main: b2b.Main{
			FiscalID: clientes[0].Nif.String,
			Name:     clientes[0].Nomcli.String,
			Email:    clientes[0].Email.String,
			Phone:    clientes[0].Telefone.String,
			Address:  clientes[0].Morcli.String,
			ZipCode:  clientes[0].Codpostal.String,
			City:     clientes[0].Localidade.String,
			Country:  clientes[0].Codpai.String,
			Currency: "EUR", // Valor fixo ou proveniente de outra fonte
			Accounts: b2b.Account{
				Client: "430000", // Valor fixo ou derivado do contexto
			},
		},
	}
	log.Printf("mandar %s", contact)

	if err := b2b.CreateContact(contact, clientes[0].Codcli); err != nil {
		http.Error(w, "Failed to create contact: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to create contact %v: %v", contact, err)
		return
	}
	// Retornar resposta de sucesso
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": fmt.Sprintf("Cliente %s inserido com sucesso na Facturadirecta!", clientes[0].Nomcli.String),
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}
