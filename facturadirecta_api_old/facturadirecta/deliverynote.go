package facturadirecta_api

import (
	"facturadirecta_api/configs"
	"facturadirecta_api/db"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func CreateDeliveryNote(contentI DeliveyNoteContent, ultimoid int, id int) error {

	
	conf := configs.GetB2B()

	account_id := conf.Account_id
	url := fmt.Sprintf("https://app.facturadirecta.com/api/%s/deliveryNotes", account_id)
	clientRequest := DeliveyNoteRequest{
		Content: contentI,
	}
	// Marshal do Payload
	payload, err := json.Marshal(clientRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Criar requisição POST
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	api := conf.Api_key
	value := conf.Value
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set(""+api, ""+value)

	// Enviar requisição
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	// Verificar status da resposta
	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create Delivery Note, status code: %d, response: %s", res.StatusCode, string(body))
	}
	var responseData struct {
		Content struct {
			Main struct {
				Title string `json:"title"`
			}
			Uuid string `json:"uuid"`
		} `json:"content"`
	}
	if err := json.Unmarshal(body, &responseData); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}
	// Now insert the name and companyId into your PostgreSQL database
	conn, err := db.OpenConnection()
	if err != nil {
		return fmt.Errorf("failed to open database connection: %v", err)
	}
	defer conn.Close()

	query := `INSERT INTO facturadirecta_ncc (id,idpai,idfacturadirecta,nomenotaentraga) VALUES ($1, $2,$3,$4)`
	_, err = conn.Exec(query, ultimoid, id, responseData.Content.Uuid, responseData.Content.Main.Title)
	if err != nil {
		return fmt.Errorf("failed to insert data into database: %v", err)
	}
	fmt.Println("Delivery Note created successfully:", string(body))
	return nil
}