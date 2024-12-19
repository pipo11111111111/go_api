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

func InserProduct(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	enableCors(&w)
	var filtro dbfunctions.ProdutoFiltro

	if err := json.NewDecoder(r.Body).Decode(&filtro); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}
	produtos, err := dbfunctions.GetProdutos(filtro)
	if err != nil {
		http.Error(w, "Error fetching products", http.StatusInternalServerError)
		log.Printf("Error fetching products: %v", err)
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
	// Arrays para armazenar sucessos e erros
	var insertedProducts []string
	var failedProducts []string

	for _, produto := range produtos {
		// Verificar se o produto já existe no banco de dados
		var existingClientID string
		query := `SELECT codartfacturadirecta FROM facturadirecta_produtos WHERE codart = $1`
		err = conn.QueryRow(query, produto.Codart).Scan(&existingClientID)
		if err == nil {
			// Se o produto já existir
			http.Error(w, fmt.Sprintf("Produto com código %s já existe na Facturadirecta.", produtos[0].Codart), http.StatusConflict)
			failedProducts = append(failedProducts, produto.Codart)
			log.Printf("Produto já existe na Facturadirecta: %s", produtos[0].Codart)
			return
		} else if err.Error() != "sql: no rows in result set" {
			// Se ocorrer outro erro ao verificar no banco
			http.Error(w, "Error checking if product exists", http.StatusInternalServerError)
			failedProducts = append(failedProducts, produto.Codart)
			log.Printf("Error checking if product exists: %v", err)
			return
		}

		contact := b2b.ProductContent{
			Type: "product",
			Main: b2b.ProductMain{
				Sku: produto.Codart,
				Name:     produto.Codart,
				Currency: "EUR",
				Sales: b2b.Sales{
					Price:       produto.Precotabela,
					Description: produto.Codartext.String,
					Tax:         []string{"S_IVA_21"},
					Account:     "700000",
				},
			},
		}

		if err := b2b.CreateProduct(contact, produto.Codart); err != nil {
			http.Error(w, "Failed to create contact: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Failed to create contact %v: %v", contact, err)
			return
		}
	}

	// Retornar resposta de sucesso
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message":           "Processamento concluído",
		"inserted_products": insertedProducts,
		"failed_products":   failedProducts,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		log.Printf("Error encoding response: %v", err)
	}
}
