package dbfunctions

import (
	"facturadirecta_api/db"

	"log"
)

// Funcoes Produtos
func GetProdutos(filtro ProdutoFiltro) (produtos []Produto, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	query := `  SELECT a.codart, a.codartext, a.precotabela
FROM public.artigos a
INNER JOIN public.fvclin f ON f.codart = a.codart AND f.idpai = 2023101114990

UNION

SELECT a.codart, a.codartext, a.precotabela
FROM public.artigos a
INNER JOIN public.orclin o ON o.codart = a.codart AND o.idpai = 2018100125685

UNION

SELECT a.codart, a.codartext, a.precotabela
FROM public.artigos a
INNER JOIN public.ncclin n ON n.codart = a.codart AND n.idpai = 2016100017135;
    
    
    `
	// Adiciona a cláusula WHERE se o filtro `Mywhere` estiver presente
	if filtro.Mywhere != "" {
		query += "  " + filtro.Mywhere
	}

	// Adiciona a cláusula ORDER BY se o filtro `MyOrder` estiver presente
	if filtro.MyOrder != "" {
		query += " " + filtro.MyOrder
	} else {
		// Ordenação padrão se não for especificada
		query += " "
	}

	// Exibe a query final (útil para debugging)
	log.Printf("Executing query: %s", query)

	// Executa a query
	rows, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var produto Produto
		if err = rows.Scan(&produto.Codart, &produto.Codartext, &produto.Precotabela); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		produtos = append(produtos, produto)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return produtos, nil
}
