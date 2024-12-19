package facturadirecta_api

//Cliente
type ContactRequest struct {
	Content Content `json:"content"`
}
type Content struct {
	Type string `json:"type"`
	Main Main   `json:"main"`
}

type Main struct {
	FiscalID string  `json:"fiscalId"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Phone    string  `json:"phone"`
	Address  string  `json:"address"`
	ZipCode  string  `json:"zipcode"`
	City     string  `json:"city"`
	Country  string  `json:"country"`
	Currency string  `json:"currency"`
	Accounts Account `json:"accounts"`
}

type Account struct {
	Client string `json:"client"`
}

//Produtos
type ProductRequest struct {
	Content ProductContent `json:"content"`
}
type ProductContent struct {
	Type string      `json:"type"`
	Main ProductMain `json:"main"`
}

type ProductMain struct {
	Sku      string `json:"sku"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
	Sales    Sales  `json:"sales"`
}
type Sales struct {
	Price       float64  `json:"price"`
	Description string   `json:"description"`
	Tax         []string `json:"tax"`
	Account     string   `json:"account"`
}

//Invoice
type InvoiceRequest struct {
	Content InvoiceContent `json:"content"`
}

type InvoiceContent struct {
	Type string      `json:"type"`
	Main InvoiceMain `json:"main"`
}
type InvoiceMain struct {
	DocNumber        DocNumber `json:"docNumber"`
	CorrectedInvoice string    `json:"correctedInvoice,omitempty"`
	Contact          string    `json:"contact"`
	Currency         string    `json:"currency"`
	Notes            string    `json:"notes,omitempty"`
	Lines            []Lines   `json:"lines"`
}
type Lines struct {
	Document string `json:"document,omitempty"`
	Account  string `json:"account,omitempty"`

	Text      string   `json:"text"`
	Quantity  float64  `json:"quantity"`
	UnitPrice float64  `json:"unitPrice"`
	Tax       []string `json:"tax"`
}
type DocNumber struct {
	Series string `json:"series"`
}

//Email
type EmailRequest struct {
	To []string `json:"to"`
}

//DeliveyNote
type DeliveyNoteRequest struct {
	Content DeliveyNoteContent `json:"content"`
}

type DeliveyNoteContent struct {
	Type string          `json:"type"`
	Main DeliveyNoteMain `json:"main"`
}
type DeliveyNoteMain struct {
	DocNumber DeliveyNoteDocNumber `json:"docNumber"`
	BaseState string               `json:"baseState"`
	Contact   string               `json:"contact"`
	Currency  string               `json:"currency"`

	Lines []DeliveyNoteLines `json:"lines"`
}
type DeliveyNoteLines struct {
	Document  string   `json:"document,omitempty"`
	Account   string   `json:"account,omitempty"`
	Text      string   `json:"text"`
	Quantity  float64  `json:"quantity"`
	UnitPrice float64  `json:"unitPrice"`
	Tax       []string `json:"tax"`
}
type DeliveyNoteDocNumber struct {
	Series string `json:"series"`
}

//Estimate
type EstimatesRequest struct {
	Content EstimatesContent `json:"content"`
}

type EstimatesContent struct {
	Type string        `json:"type"`
	Main EstimatesMain `json:"main"`
}
type EstimatesMain struct {
	DocNumber EstimatesDocNumber `json:"docNumber"`
	BaseState string             `json:"baseState"`
	Contact   string             `json:"contact"`
	Currency  string             `json:"currency"`

	Lines []EstimatesLines `json:"lines"`
}
type EstimatesLines struct {
	Document string `json:"document,omitempty"`
	Account  string `json:"account,omitempty"`

	Text      string   `json:"text"`
	Quantity  float64  `json:"quantity"`
	UnitPrice float64  `json:"unitPrice"`
	Tax       []string `json:"tax"`
}
type EstimatesDocNumber struct {
	Series string `json:"series"`
}
