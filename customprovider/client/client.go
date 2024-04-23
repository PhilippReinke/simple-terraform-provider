package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"terraform-provider-customprovider/model"
)

const (
	itemsURL = "http://localhost:8080/items"
)

func CreateItem(name string) (*model.Item, error) {
	reqUrl := fmt.Sprintf("%v/%v", itemsURL, name)
	resp, err := http.Post(reqUrl, "application/json", strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var item *model.Item
	err = json.Unmarshal(body, &item)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return nil, err
	}

	return item, nil
}

func ReadItems() (*model.Items, error) {
	resp, err := http.Get(itemsURL)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var items *model.Items
	err = json.Unmarshal(body, &items)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return nil, err
	}

	return items, nil
}

func ReadItem(id string) (*model.Item, error) {
	reqUrl := fmt.Sprintf("%v/%v", itemsURL, id)
	resp, err := http.Get(reqUrl)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var item *model.Item
	err = json.Unmarshal(body, &item)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return nil, err
	}

	return item, nil
}

func UpdateItem(id, updatedName string) (*model.Item, error) {
	reqUrl := fmt.Sprintf("%v/%v/%v", itemsURL, id, updatedName)
	// data := []byte(`{"data": "some-data"}`)
	req, err := http.NewRequest(http.MethodPut, reqUrl, bytes.NewBuffer(nil))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var item *model.Item
	err = json.Unmarshal(body, &item)
	if err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return nil, err
	}

	return item, nil
}

func DeleteItem(id string) error {
	reqUrl := fmt.Sprintf("%v/%v", itemsURL, id)
	// data := []byte(`{"data": "some-data"}`)
	req, err := http.NewRequest(http.MethodDelete, reqUrl, bytes.NewBuffer(nil))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Could not delete item. Status code %v", resp.StatusCode)
	}

	return nil
}
