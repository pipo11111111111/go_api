package dbfunctions

import (
	"facturadirecta_api/db"
	"fmt"
	"log"
)

 
// Funcoes Faturas
func GetFatura(filtro ProdutoFiltro) (faturas []Fatura, err error) {
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
		FROM public.facturadirecta_fvc 
		
	),0)+1 AS ultimo_idfacturadirecta
FROM 
	public.fvc f 
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
		var fatura Fatura
		if err = rows.Scan(&fatura.Id, &fatura.CodcliB2b, &fatura.Datcom, &fatura.Ultimoid); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		faturas = append(faturas, fatura)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return faturas, nil
}
func GetLinhasFatura(id int) (linhasfaturas []LinhasFatura, err error) {

	conn, err := db.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	query := fmt.Sprintf(`SELECT f.idpai,f.codart,f.desart,f.preunitliq,f.coduni,f.qtdart,p.codartfacturadirecta FROM public.fvclin f 
left join public.facturadirecta_produtos p on f.codart = p.codart WHERE idpai = %d`, id)

	rows, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var linhafatura LinhasFatura
		if err = rows.Scan(&linhafatura.Idpai, &linhafatura.Codart, &linhafatura.Desart, &linhafatura.Preunit, &linhafatura.Coduni, &linhafatura.Qtdart, &linhafatura.Codartfacturadirecta); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		linhasfaturas = append(linhasfaturas, linhafatura)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return linhasfaturas, nil
}

func GetIdFaturaDireta(filtro ProdutoFiltro) (idfacturadirectas []IdFacturaDirecta, err error) {
	conn, err := db.OpenConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	query := `  SELECT 
	 fc.idfacturadirecta,
     c.email
FROM 
	public.facturadirecta_fvc fc
    inner join public.fvc f on fc.idpai = f.id 
    inner join public.clientes c on f.codcli = c.codcli
 
 `
	// Adiciona a cláusula WHERE se o filtro `Mywhere` estiver presente
	if filtro.Mywhere != "" {
		query += " where fc.idpai=" + filtro.Mywhere
	}

	// Adiciona a cláusula ORDER BY se o filtro `MyOrder` estiver presente
	if filtro.MyOrder != "" {
		query += " " + filtro.MyOrder
	} else {
		// Ordenação padrão se não for especificada
		query += ""
	}

	rows, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var idfacturadirecta IdFacturaDirecta
		if err = rows.Scan(&idfacturadirecta.IdFacturaDirecta, &idfacturadirecta.Email); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		idfacturadirectas = append(idfacturadirectas, idfacturadirecta)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return idfacturadirectas, nil
}
func GetFaturaRectificativa(filtro ProdutoFiltro) (faturas []FaturaRectificativa, err error) {
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
		FROM public.facturadirecta_fvc_rectificativa 
		
	),0)+1 AS ultimo_idfacturadirecta,
    fv.idfacturadirecta,
    fv.nomefactura
FROM 
	public.fvc f 
INNER JOIN 
	public.facturadirecta_clientes c 
ON 
	f.codcli = c.codcli
inner join 
		public.facturadirecta_fvc fv 
        on f.id = fv.idpai
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
		var fatura FaturaRectificativa
		if err = rows.Scan(&fatura.Id, &fatura.CodcliB2b, &fatura.Datcom, &fatura.Ultimoid,&fatura.Idfacturadirecta,&fatura.Nomefactura); err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		faturas = append(faturas, fatura)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return faturas, nil
}