package handlres

import (
	"facturadirecta_api/db"
	dbfunctions "facturadirecta_api/db_functions"
	b2b "facturadirecta_api/facturadirecta"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func InserOrcamento(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	enableCors(&w)
	var filtro dbfunctions.ProdutoFiltro

	// Decodificar o filtro da requisição
	if err := json.NewDecoder(r.Body).Decode(&filtro); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	// Buscar orcamento do banco de dados
	orcamento, err := dbfunctions.GetOrcamento(filtro)
	if err != nil {
		http.Error(w, "Error fetching Orçamento", http.StatusInternalServerError)
		log.Printf("Error fetching clients: %v", err)
		return
	}

	// Buscar as linhas do orcamento
	linhas, err := dbfunctions.GetLinhasOrcamento(orcamento[0].Id)
	if err != nil {
		http.Error(w, "Error fetching linhas", http.StatusInternalServerError)
		log.Printf("Error fetching linhas: %v", err)
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

	// Verificar se o orcamento já existe no banco de dados
	var existingClientID string
	query := `SELECT idfacturadirecta FROM facturadirecta_orc WHERE idpai = $1`
	err = conn.QueryRow(query, orcamento[0].Id).Scan(&existingClientID)
	if err == nil {
		// Se o orcamento já existir
		http.Error(w, fmt.Sprintf("Orçamento com código %d já existe na Facturadirecta.", orcamento[0].Id), http.StatusConflict)
		log.Printf("Orçamento já existe na Facturadirecta: %d", orcamento[0].Id)
		return
	} else if err.Error() != "sql: no rows in result set" {
		// Se ocorrer outro erro ao verificar no banco
		http.Error(w, "Error checking if client exists", http.StatusInternalServerError)
		log.Printf("Error checking if client exists: %v", err)
		return
	}
	// Preparar as linhas da fatura para o formato esperado
	invoiceLines := []b2b.EstimatesLines{}
	for _, linha := range linhas {
		var codart = linha.Desart.String

		if linha.Codartfacturadirecta.Valid {
			invoiceLine := b2b.EstimatesLines{
				Document:  linha.Codartfacturadirecta.String,
				Account:   "700000",
				Text:      codart,
				Quantity:  linha.Qtdart,
				UnitPrice: linha.Preunit,
				Tax: []string{
					"S_IVA_21"},
			}
			invoiceLines = append(invoiceLines, invoiceLine)
		} else {
			invoiceLine := b2b.EstimatesLines{
				Text:      codart,
				Quantity:  linha.Qtdart,
				UnitPrice: linha.Preunit,
				Tax: []string{
					"S_IVA_21"},
			}
			invoiceLines = append(invoiceLines, invoiceLine)
		}
	}
	// Criar o Payload para enviar
	content := b2b.EstimatesContent{

		Type: "estimate",
		Main: b2b.EstimatesMain{
			DocNumber: b2b.EstimatesDocNumber{
				Series: "P",
			},
			BaseState: "pending",
			Contact:   orcamento[0].CodcliB2b,
			Currency:  "EUR",
			Lines:     invoiceLines,
		},
	}

	// Chamar a função CreateInvoice
	if err := b2b.CreateEstimate(content, orcamento[0].Ultimoid, orcamento[0].Id); err != nil {
		http.Error(w, "Failed to create orçamento: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to create orçamento %v: %v", content, err)

		return
	}
	// Retornar resposta de sucesso
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": fmt.Sprintf("Orçamento %s inserido com sucesso na Facturadirecta!", strconv.Itoa(orcamento[0].Id)),
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
	// Responder com sucesso
	/*
		fmt.Fprintf(w, "Invoice created successfully") */
}
