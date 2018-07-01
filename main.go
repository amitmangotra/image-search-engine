package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {

	// Clarifai API Key
	apiKey := "baa1c074beb54b21b388d3a8632a22b5"

	// Model ID for General Predict Model
	generalModelID := "aaa03c23b3724a16a56b629203edc62c"

	// API Domain and Endpoint
	apiDomain := "https://api.clarifai.com/v2/models/"
	endpoint := "/outputs"

	//invertedIndex map[string][]string
	var invertedIndex = make(map[string][]string)

	// type InvertedIndex struct {
	// 	Tag []string `json: "tag"`
	// }

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

	// type ConceptAPIResponse struct {
	// 	Status struct {
	// 		Code        int64  `json: "code"`
	// 		Description string `json: "description"`
	// 	} `json: "status"`
	// 	Outputs []struct {
	// 		Id     string `json: "id"`
	// 		Status struct {
	// 			Code        int64  `json: "code"`
	// 			Description string `json: "description"`
	// 		} `json: "status"`
	// 		CreatedAt string `json: "created_at"`
	// 		Model     struct {
	// 			ID         string `json: "id"`
	// 			Name       string `json:"name"`
	// 			CreatedAt  string `json: "created_at"`
	// 			AppID      string `json: "app_id"`
	// 			OutputInfo struct {
	// 				Type    string `json: "type"`
	// 				TypeExt string `json: "type_ext"`
	// 			} `json: "output_info"`
	// 			ModelVersion struct {
	// 				ID        string `json: "id"`
	// 				CreatedAt string `json: "created_at"`
	// 				Status    struct {
	// 					Code        int64  `json: "code"`
	// 					Description string `json: "description"`
	// 				} `json: "status"`
	// 			} `json: "model_version"`
	// 			DisplayName string `json: "display_name"`
	// 		} `json: "model"`
	// 		Input struct {
	// 			ID   string `json: "id"`
	// 			Data struct {
	// 				Image struct {
	// 					URL string `json:"url"`
	// 				} `json: "image"`
	// 			} `json: "data"`
	// 		} `json: "input"`
	// 		Data struct {
	// 			Concepts []struct {
	// 				Id    string  `json: "id"`
	// 				Value float64 `json: "value"`
	// 				Name  string  `json: "name"`
	// 				AppId string  `json: "app_id"`
	// 			} `json: "concepts"`
	// 		} `json: "data"`
	// 	} `json: "outputs"`
	// }

	body := strings.NewReader(`{"inputs": [{"data": {"image": {"url": "https://farm7.staticflickr.com/5769/21094803716_da3cea21b8_o.jpg"}}}]}`)

	req, err := http.NewRequest("POST", apiDomain+generalModelID+endpoint, body)
	if err != nil {
		panic(err)
		fmt.Println(err)
	}
	req.Header.Set("Authorization", "Key "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)

	var responseObject ConceptAPIResponse
	e := json.Unmarshal([]byte(result), &responseObject)
	if e != nil {
		panic(e)
	}

	// fmt.Println(result)
	for _, item := range responseObject.Outputs[0].Data.Concepts {
		invertedIndex[item.Name] = append(invertedIndex[item.Name], responseObject.Outputs[0].Input.Data.Image.URL)
	}

	for k, v := range invertedIndex {
		fmt.Println(k, "--", v)
	}

}
