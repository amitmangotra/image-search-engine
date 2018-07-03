package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {

	// Clarifai API Key
	apiKey := "baa1c074beb54b21b388d3a8632a22b5"

	// Model ID for General Predict Model
	generalModelID := "aaa03c23b3724a16a56b629203edc62c"

	// API Domain and Endpoint
	apiDomain := "https://api.clarifai.com/v2/models/"
	endpoint := "/outputs"

	type InvertedIndex struct {
		Tag []string `json: "tag"`
	}

	type Image struct {
		URL string `json:"url"`
	}

	type Status struct {
		Code        int64  `json: "code"`
		Description string `json: "description"`
	}

	type ModelVersion struct {
		ID        string  `json: "id"`
		CreatedAt string  `json: "created_at"`
		Status    *Status `json: "status"`
	}

	type OutputInfo struct {
		Type    string `json: "type"`
		TypeExt string `json: "type_ext"`
	}

	type Model struct {
		ID           string        `json: "id"`
		Name         string        `json:"name"`
		CreatedAt    string        `json: "created_at"`
		AppID        string        `json: "app_id"`
		OutputInfo   *OutputInfo   `json: "output_info"`
		ModelVersion *ModelVersion `json: "model_version"`
		DisplayName  string        `json: "display_name"`
	}

	type Input struct {
		ID   string `json: "id"`
		Data struct {
			Image *Image `json: "image"`
		} `json: "data"`
	}

	type Concept struct {
		ID    string  `json: "id"`
		Value float64 `json: "value"`
		Name  string  `json: "name"`
		AppID string  `json: "app_id"`
	}

	type ConceptAPIResponse struct {
		Status  *Status `json: "status"`
		Outputs []struct {
			ID        string  `json: "id"`
			Status    *Status `json: "status"`
			CreatedAt string  `json: "created_at"`
			Model     *Model  `json: "model"`
			Input     *Input  `json: "input"`
			Data      struct {
				Concepts []*Concept `json: "concepts"`
			} `json: "data"`
		} `json: "outputs"`
	}

	type ImageBody struct {
		URL string `json:"url"`
	}

	type Data struct {
		Image ImageBody `json:"image"`
	}

	type Inputs struct {
		Data Data `json:"data"`
	}
	type Payload struct {
		Inputs []Inputs `json:"inputs"`
	}

	var urls []string

	file, err := os.Open("images.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		url := scanner.Text()
		urls = append(urls, url)
		// fmt.Println(strconv.CanBackquote(url))
	}

	for _, url := range urls {
		pl := Payload{
			Inputs: []Inputs{
				{
					Data: Data{
						Image: ImageBody{
							URL: url,
						},
					},
				},
			},
		}

		payloadBytes, err := json.MarshalIndent(pl, "", "  ")
		if err != nil {
			panic(err)
		}
		body := bytes.NewReader(payloadBytes)

		req, err := http.NewRequest("POST", apiDomain+generalModelID+endpoint, body)
		if err != nil {
			panic(err)
		}

		req.Header.Set("Authorization", "Key "+apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()

		result, _ := ioutil.ReadAll(resp.Body)

		var conceptResponse ConceptAPIResponse
		e := json.Unmarshal([]byte(result), &conceptResponse)
		if e != nil {
			panic(e)
		}

		// fmt.Println(conceptResponse.Outputs[0])
		for _, item := range conceptResponse.Outputs[0].Data.Concepts {
			fmt.Print(item.Name, " ")
		}
		fmt.Println("-----------------------------")
	}

}
