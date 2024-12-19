package dbfunctions

import (
	"facturadirecta_api/db"
	"fmt"
	"log"
)

 
 

 
// Funcoes Notas Entrega
func GetNotaEntrega(filtro ProdutoFiltro) (notasEntrega []Fatura, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	query := ` SELECT 
	f.id,
	c.codclifacturadirecta,
	f.datcom,
	COALESCE((
		SELECT MAX(id) 
		FROM public.facturadirecta_ncc 
		
	),0)+1 AS ultimo_idfacturadirecta
FROM 
	public.ncc f 
INNER JOIN 
	public.facturadirecta_clientes c 
ON 
	f.codcli = c.codcli
 `
	// Adiciona a cláusula WHERE se o filtro `Mywhere` estiver presente
	if filtro.Mywhere != "" {
		query += " where f.id=" + filtro.Mywhere
	}

	// Adiciona a cláusula ORDER BY se o filtro `MyOrder` estiver presente
	if filtro.MyOrder != "" {
		query += " " + filtro.MyOrder
	} else {
		// Ordenação padrão se não for especificada
		query += " ORDER BY f.codcli limit 10 "
	}

	rows, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var notaEntrega Fatura
		if err = rows.Scan(&notaEntrega.Id, &notaEntrega.CodcliB2b, &notaEntrega.Datcom, &notaEntrega.Ultimoid); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		notasEntrega = append(notasEntrega, notaEntrega)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return notasEntrega, nil
}

func GetLinhasNotasEntrega(id int) (linhasnotasEntrega []LinhasFatura, err error) {

	conn, err := db.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	query := fmt.Sprintf(`SELECT idpai,codart,desart,preunitliq,coduni,qtdart FROM public.ncclin WHERE idpai = %d`, id)

	rows, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var linhasnotaEntrega LinhasFatura
		if err = rows.Scan(&linhasnotaEntrega.Idpai, &linhasnotaEntrega.Codart, &linhasnotaEntrega.Desart, &linhasnotaEntrega.Preunit, &linhasnotaEntrega.Coduni, &linhasnotaEntrega.Qtdart); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		linhasnotasEntrega = append(linhasnotasEntrega, linhasnotaEntrega)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return linhasnotasEntrega, nil
}