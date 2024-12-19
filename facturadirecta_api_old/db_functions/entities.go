package dbfunctions

import "database/sql"

func (ns NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte(`""`), nil // Retorna uma string vazia para valores nulos
	}
	return []byte(`"` + escapeJSONString(ns.String) + `"`), nil
}

// escapeJSONString escapa caracteres especiais em strings JSON
func escapeJSONString(s string) string {
	// Implementa o escape dos caracteres especiais manualmente
	var result []rune
	for _, r := range s {
		switch r {
		case '\b':
			result = append(result, '\\', 'b')
		case '\f':
			result = append(result, '\\', 'f')
		case '\n':
			result = append(result, '\\', 'n')
		case '\r':
			result = append(result, '\\', 'r')
		case '\t':
			result = append(result, '\\', 't')
		case '\\', '"':
			result = append(result, '\\', r)
		default:
			if r < 0x20 {
				// Escapa caracteres de controle
				result = append(result, '\\', 'u', '0', '0', '0', '0'+r/16, '0'+r%16)
			} else {
				result = append(result, r)
			}
		}
	}
	return string(result)
}

// GetOrDefault retorna o valor da NullString ou um valor padrão
func (ns NullString) GetOrDefault(defaultValue string) string {
	if !ns.Valid {
		return defaultValue
	}
	return ns.String
}

type NullString struct {
	sql.NullString
}

type Cliente struct {
	Codcli     int        `json:"codcli"`
	Nif        NullString `json:"nif"`
	Nomcli     NullString `json:"nomcli"`
	Email      NullString `json:"email"`
	Telefone   NullString `json:"telefone"`
	Morcli     NullString `json:"morcli"`
	Codpostal  NullString `json:"codpostal"`
	Localidade NullString `json:"localidade"`
	Codpai     NullString `json:"codpai"`
}

type Fatura struct {
	Id        int    `json:"id"`
	CodcliB2b string `json:"codclib2b"`
	Datcom    string `json:"datcom"`
	Ultimoid  int    `json:"ultimoid"`
}
type FaturaRectificativa struct {
	Id        int    `json:"id"`
	CodcliB2b string `json:"codclib2b"`
	Datcom    string `json:"datcom"`
	Ultimoid  int    `json:"ultimoid"`
	Idfacturadirecta string `json:"idfacturadirecta"`
	Nomefactura string `json:"nomefactura"`
}
type LinhasFatura struct {
	Idpai                int        `json:"idpai"`
	Codart               NullString `json:"codart"`
	Desart               NullString `json:"desart"`
	Preunit              float64    `json:"preunitliq"`
	Coduni               NullString `json:"coduni"`
	Qtdart               float64    `json:"qtdart"`
	Codartfacturadirecta NullString `json:"codartfacturadirecta"`
}

type IdFacturaDirecta struct {
	IdFacturaDirecta string `json:"idFacturaDirecta"`
	Email            string `json:"email"`
}

type Produto struct {
	Codart      string     `json:"codart"`
	Codartext   NullString `json:"codartext"`
	Precotabela float64    `json:"precotabela"`
}

// Estrutura para os filtros enviados no body da requisição
type ProdutoFiltro struct {
	Mywhere string `json:"mywhere"`
	MyOrder string `json:"myorder"`
}
/* 
type NotaCredito struct {
	Id        int    `json:"id"`
	CodcliB2b int    `json:"codclib2b"`
	Datcom    string `json:"datcom"`
	Ultimoid  int    `json:"ultimoid"`
}

type LinhasNotaCredito struct {
	Idpai   int        `json:"idpai"`
	Codart  NullString `json:"codart"`
	Desart  NullString `json:"desart"`
	Preunit float64    `json:"preunitliq"`
	Coduni  NullString `json:"coduni"`
	Qtdart  float64    `json:"qtdart"`
	Codartfacturadirecta NullString `json:"codartfacturadirecta"`
} */
