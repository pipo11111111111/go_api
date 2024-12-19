package dbfunctions

import (
	"facturadirecta_api/db"
	"fmt"
	"log"
)

 
 

 
// Funcoes ORcamento
func GetOrcamento(filtro ProdutoFiltro) (orcamentos []Fatura, err error) {
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
		FROM public.facturadirecta_orc 
		
	),0)+1 AS ultimo_idfacturadirecta
FROM 
	public.orc f 
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
		var orcamento Fatura
		if err = rows.Scan(&orcamento.Id, &orcamento.CodcliB2b, &orcamento.Datcom, &orcamento.Ultimoid); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		orcamentos = append(orcamentos, orcamento)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orcamentos, nil
}

func GetLinhasOrcamento(id int) (linhasorcamentos []LinhasFatura, err error) {

	conn, err := db.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	query := fmt.Sprintf(`SELECT o.idpai,o.codart,o.desart,o.preunitliq,o.coduni,o.qtdart,p.codartfacturadirecta FROM public.orclin o
	left join public.facturadirecta_produtos p on o.codart = p.codart WHERE o.idpai = %d`, id)

	rows, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var linhaorcamento LinhasFatura
		if err = rows.Scan(&linhaorcamento.Idpai, &linhaorcamento.Codart, &linhaorcamento.Desart, &linhaorcamento.Preunit, &linhaorcamento.Coduni, &linhaorcamento.Qtdart,&linhaorcamento.Codartfacturadirecta); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		linhasorcamentos = append(linhasorcamentos, linhaorcamento)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return linhasorcamentos, nil
}