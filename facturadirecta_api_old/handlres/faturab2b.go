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

func InserFatura(w http.ResponseWriter, r *http.Request) {
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
	fatura, err := dbfunctions.GetFatura(filtro)
	if err != nil {
		http.Error(w, "Error fetching fatura", http.StatusInternalServerError)
		log.Printf("Error fetching clients: %v", err)
		return
	}

	// Buscar as linhas da fatura
	linhas, err := dbfunctions.GetLinhasFatura(fatura[0].Id)
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

	// Verificar se o cliente já existe no banco de dados
	var existingClientID string
	query := `SELECT idfacturadirecta FROM facturadirecta_fvc WHERE idpai = $1`
	err = conn.QueryRow(query, fatura[0].Id).Scan(&existingClientID)
	if err == nil {
		// Se o cliente já existir
		http.Error(w, fmt.Sprintf("Fatura com código %d já existe na Facturadirecta.", fatura[0].Id), http.StatusConflict)
		log.Printf("Fatura já existe na Facturadirecta: %d", fatura[0].Id)
		return
	} else if err.Error() != "sql: no rows in result set" {
		// Se ocorrer outro erro ao verificar no banco
		http.Error(w, "Error checking if fatura exists", http.StatusInternalServerError)
		log.Printf("Error checking if fatura exists: %v", err)
		return
	}
	// Preparar as linhas da fatura para o formato esperado
	invoiceLines := []b2b.Lines{}
	for _, linha := range linhas {
		var codart = linha.Desart.String

		if linha.Codartfacturadirecta.Valid {
			invoiceLine := b2b.Lines{
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

			invoiceLine := b2b.Lines{

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
	content := b2b.InvoiceContent{

		Type: "invoice",
		Main: b2b.InvoiceMain{
			DocNumber: b2b.DocNumber{
				Series: "F",
			},
			Contact:  fatura[0].CodcliB2b,
			Currency: "EUR",
			Lines:    invoiceLines,
		},
	}
	fmt.Print(content)
	// Chamar a função CreateInvoice
	if err := b2b.CreateInvoice(content, fatura[0].Ultimoid, fatura[0].Id); err != nil {
		http.Error(w, "Failed to create invoice: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to create invoice %v: %v", content, err)

		return
	}

	// Responder com sucesso
	/* fmt.Fprintf(w, "Invoice created successfully") */
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": fmt.Sprintf("Fatura %s inserida com sucesso na Facturadirecta!", strconv.Itoa(fatura[0].Id)),
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}




func InserRectificativa(w http.ResponseWriter, r *http.Request) {
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
	fatura, err := dbfunctions.GetFaturaRectificativa(filtro)
	if err != nil {
		http.Error(w, "Error fetching fatura", http.StatusInternalServerError)
		log.Printf("Error fetching clients: %v", err)
		return
	}

	// Buscar as linhas da fatura
	linhas, err := dbfunctions.GetLinhasFatura(fatura[0].Id)
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

	// Verificar se o cliente já existe no banco de dados
	var existingClientID string
	query := `SELECT idfacturadirecta FROM facturadirecta_fvc_rectificativa WHERE idpai = $1`
	err = conn.QueryRow(query, fatura[0].Id).Scan(&existingClientID)
	if err == nil {
		// Se o cliente já existir
		http.Error(w, fmt.Sprintf("Fatura com código %d já existe na Facturadirecta.", fatura[0].Id), http.StatusConflict)
		log.Printf("Fatura já existe na Facturadirecta: %d", fatura[0].Id)
		return
	} else if err.Error() != "sql: no rows in result set" {
		// Se ocorrer outro erro ao verificar no banco
		http.Error(w, "Error checking if fatura exists", http.StatusInternalServerError)
		log.Printf("Error checking if fatura exists: %v", err)
		return
	}
	// Preparar as linhas da fatura para o formato esperado
	invoiceLines := []b2b.Lines{}
	for _, linha := range linhas {
		var codart = linha.Desart.String

		if linha.Codartfacturadirecta.Valid {
			invoiceLine := b2b.Lines{
				Document:  linha.Codartfacturadirecta.String,
				Account:   "700000",
				Text:      codart,
				Quantity:  -linha.Qtdart,
				UnitPrice: linha.Preunit,
				Tax: []string{
					"S_IVA_21"},
			}
			invoiceLines = append(invoiceLines, invoiceLine)
		} else {

			invoiceLine := b2b.Lines{

				Text:      codart,
				Quantity:  -linha.Qtdart,
				UnitPrice: linha.Preunit,
				Tax: []string{
					"S_IVA_21"},
			}
			invoiceLines = append(invoiceLines, invoiceLine)
		}
		/* // Verifica se a descrição (linha.Desart) está vazia
		if linha.Desart.String == "" {
			// Se Desart estiver vazio, usa o valor de Codart
			codart = linha.Codart.String
		} else {
			// Caso contrário, usa o valor de Desart
			codart = linha.Desart.String
		} */

	}

	// Criar o Payload para enviar
	content := b2b.InvoiceContent{

		Type: "invoice",
		Main: b2b.InvoiceMain{
			DocNumber: b2b.DocNumber{
				Series: "R",
			},
			Contact:  fatura[0].CodcliB2b,
			Currency: "EUR",
			Lines:    invoiceLines,
			CorrectedInvoice: fatura[0].Idfacturadirecta,
			Notes: "Fatura Rectificada referente " +fatura[0].Nomefactura,
		},
	}
	fmt.Print(content)
	// Chamar a função CreateInvoice
	if err := b2b.CreateRectificativa(content, fatura[0].Ultimoid, fatura[0].Id); err != nil {
		http.Error(w, "Failed to create invoice: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to create invoice %v: %v", content, err)

		return
	}

	// Responder com sucesso
	/* fmt.Fprintf(w, "Invoice created successfully") */
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": fmt.Sprintf("Fatura %s inserida com sucesso na Facturadirecta!", strconv.Itoa(fatura[0].Id)),
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}
