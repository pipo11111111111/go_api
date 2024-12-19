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

func CreateProduct(contact ProductContent, codcli string) error {
	conf := configs.GetB2B()

	account_id := conf.Account_id
	url := fmt.Sprintf("https://app.facturadirecta.com/api/%s/products", account_id)

	clientRequest := ProductRequest{
		Content: contact,
	}

	payload, err := json.Marshal(clientRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal contact: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	api := conf.Api_key
	value := conf.Value
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set(""+api, ""+value)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create producto, status code: %d, response: %s", res.StatusCode, string(body))
	}
	// Parse the JSON response
	var responseData struct {
		Content struct {
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

	query := `INSERT INTO facturadirecta_produtos (codart,codartfacturadirecta) VALUES ($1, $2)`
	_, err = conn.Exec(query, codcli, responseData.Content.Uuid)
	if err != nil {
		return fmt.Errorf("failed to insert data into database: %v", err)
	}
	fmt.Println("Producto created successfully:", string(body))
	return nil
}
