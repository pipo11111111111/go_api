package dbfunctions

import (
	"facturadirecta_api/db"
	"log"
)

// Funcoes CLiente
func GetClient(filtro ProdutoFiltro) (clientes []Cliente, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	query := ` SELECT 
	COALESCE(codcli,0)AS codcli,
		COALESCE(nif, ' ') AS nif,
		COALESCE(nomcli, ' ') AS nomcli,
        COALESCE(email, ' ') AS email,
        COALESCE(telefone, ' ') AS telefone,
        COALESCE(morcli, ' ') AS morcli,
        COALESCE(codpostal, ' ') AS codpostal,
        COALESCE(localidade, ' ') AS localidade, 
        COALESCE(codpai, ' ') AS codpai
        
    FROM 
        public.clientes 
    
    
    `
	// Adiciona a cláusula WHERE se o filtro `Mywhere` estiver presente
	if filtro.Mywhere != "" {
		query += " where codcli=" + filtro.Mywhere
	}

	// Adiciona a cláusula ORDER BY se o filtro `MyOrder` estiver presente
	if filtro.MyOrder != "" {
		query += " " + filtro.MyOrder
	} else {
		// Ordenação padrão se não for especificada
		query += " ORDER BY codcli"
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
		var cliente Cliente
		if err = rows.Scan(&cliente.Codcli, &cliente.Nif, &cliente.Nomcli, &cliente.Email, &cliente.Telefone, &cliente.Morcli, &cliente.Codpostal, &cliente.Localidade, &cliente.Codpai); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		clientes = append(clientes, cliente)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return clientes, nil
}
