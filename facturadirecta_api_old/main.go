package main

import (
	"fmt"
	"log"
	"net/http"

	"facturadirecta_api/configs"
	"facturadirecta_api/handlres"

	"github.com/go-chi/chi"
)

func main() {
	err := configs.Load()
	if err != nil {
		log.Fatalf("Error loading cinfig: %v", err)
	}

	//cria novo router
	r := chi.NewRouter()
	//Romafe

	//get cliente bd
	r.Post("/getCliente/", handlres.ListClient)
	//get faturas
	r.Post("/getFatura/", handlres.ListFatura)
	//get linhas fatura
	r.Post("/getLinhas/", handlres.ListLinhas)
	//get id facturasimples
	r.Post("/getIdFacturaDirecta/", handlres.ListIdFaturaDirecta)
	//get orcamento
	r.Post("/getOrcamento/", handlres.ListOrcamento)
	//get linhas orcamento
	r.Post("/getLinhasOrcamento/", handlres.ListLinhasOrcamento)

	//get Nota Entrega
	r.Post("/getNotasEntrega/", handlres.ListNotaEntrega)
	//get linhas Nota Entrega
	r.Post("/getLinhasNotasEntrega/", handlres.ListLinhasNotaEntrega)

	//get produtos 
	r.Post("/getProdutos/",handlres.ListProdutos)



 



	//facturadirecta
	//send Products
	r.Post("/insertProducts/",handlres.InserProduct)
	//send client
	r.Post("/insertClient/", handlres.InserClient)	
	//Send fatura
	r.Post("/insertFatura/", handlres.InserFatura)
	//Send fatura by Email
	r.Post("/sendByEmail/", handlres.SendByEmail)
	//Send Orcamento
	r.Post("/insertOrcamento/", handlres.InserOrcamento)

	//Send Fatura Rectificativa
	r.Post("/rectificarFatura/",handlres.InserRectificativa)
	//Send NotaEntrega
	r.Post("/insertNotaEntrega/", handlres.InserNotaEntrega)



	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html") // Serve o arquivo HTML
	})
	port := configs.GetServerPort()
	log.Printf("Server running on port %s", port)
	http.ListenAndServe(fmt.Sprintf(":%s", port), r)
}
