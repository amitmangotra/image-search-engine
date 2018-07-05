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
	"sort"
	"strings"
)

func main() {

	// Clarifai API Key
	apiKey := "c65bf7947a1d42c29d97a976dd2cd342"

	// Model ID for General Predict Model
	generalModelID := "aaa03c23b3724a16a56b629203edc62c"

	// API Domain and Endpoint
	apiDomain := "https://api.clarifai.com/v2/models/"
	endpoint := "/outputs"

	/***************************************************************/

	/*
	 * Structs for API Response
	 */

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

	type Image struct {
		URL string `json:"url"`
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

	type APIResponse struct {
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

	/***************************************************************/

	/*
	 * Structs for constructing data to make an API call
	 * to Clarifai's Predict Endpoint
	 */

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

	/***************************************************************/

	/*
	 * Slice to store all the image URLs
	 */
	var urls []string

	/***************************************************************/

	/*
	 * Struct defined to store Image URL
	 * and concept's Probability for that image
	 */
	type Tuple struct {
		URL         string
		Probability float64
	}

	/*
	 * Map to store the inverted index consisting of
	 * the concept name as keys and a slice of tuple/pair of
	 * image and concept's probability as values
	 */
	invertedIndex := make(map[string][]Tuple)

	/***************************************************************/

	/*
	 * Opening the given images file (images.txt)
	 * and logging error
	 */
	file, err := os.Open("images.txt")
	if err != nil {
		log.Fatal(err)
	}

	/*
	 * Closing the opened file
	 */
	defer file.Close()

	/*
	 * Initializing Scanner to read file
	 */
	scanner := bufio.NewScanner(file)

	/*
	 * Scanning through the lines in file
	 * and appending each line in file (i.e. URL) to a slice
	 */
	for scanner.Scan() {
		url := scanner.Text()
		urls = append(urls, url)
	}

	/***************************************************************/

	/*
	 * Looping through each image URL, constructing
	 * the API data using Payload struct and making a Post request to
	 * Clarifai's Predict Service Endpoint
	 */

	for index, url := range urls {
		// Constructing the data using Payload struct and each url
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
		// Marshalling the data to feed as request body
		payloadBytes, err := json.Marshal(pl)
		// Handling Error
		if err != nil {
			panic(err)
		}

		// Reading from Marshalled data to generate body
		body := bytes.NewReader(payloadBytes)

		// Making an HTTP POST request to Clarifai's Predict Service
		req, err := http.NewRequest("POST", apiDomain+generalModelID+endpoint, body)
		// Handling Error
		if err != nil {
			panic(err)
		}

		// Setting HTTP request Headers
		req.Header.Set("Authorization", "Key "+apiKey)
		req.Header.Set("Content-Type", "application/json")

		// Making an HTTP Request using the set Headers and storing Response
		resp, err := http.DefaultClient.Do(req)
		// Handling Error
		if err != nil {
			panic(err)
		}
		// Closing the response body
		defer resp.Body.Close()

		// Reading the response body and storing it
		result, _ := ioutil.ReadAll(resp.Body)

		/*
		 * Parsing JSON Response by Unmarshalling
		 * it to APIResponse
		 */
		var apiResponse APIResponse
		e := json.Unmarshal([]byte(result), &apiResponse)
		// Handling Error
		if e != nil {
			panic(e)
		}

		/*
		 * Iterating through Parsed JSON Response to filter out all concepts
		 * for a given image and storing each concept name into a map as key with it's
		 * input image url and conceptas value to construct an INVERTED INDEX
		 */
		for _, item := range apiResponse.Outputs[0].Data.Concepts {
			// Appending each concept to a slice to image URLs
			tuple := Tuple{apiResponse.Outputs[0].Input.Data.Image.URL, item.Value}
			invertedIndex[item.Name] = append(invertedIndex[item.Name], tuple)
		}

		// Log information after each image url gets predicted and each concept related to it gets indexed
		fmt.Println("Concepts in image", index+1, "are predicted and each concept is indexed")
	}

	/***************************************************************/

	/*
	 * Interface to interact with Search Engine
	 */

	fmt.Println()
	fmt.Println("Search Engine is ready !!")
	fmt.Println()

	// Halting Measure for User's Interaction
	flag := true

	/*
	 * Reading User's Input as long as he/she wants
	 * HALTING CONDITION: flag set to 'false'
	 */
	for flag {

		// Initialize a reader to read from stdin
		reader := bufio.NewReader(os.Stdin)

		// User Interaction with system
		fmt.Println("What would you like to search today?")
		text, _ := reader.ReadString('\n')

		/*
		 * Checking if the user's entered term exists in INVERTED INDEX.
		 * If the term exists, the corresponding image URLs are printed out, otherwise,
		 * "No Image Found" is printed out.
		 */
		if val, ok := invertedIndex[strings.ToLower(strings.TrimRight(text, "\n"))]; ok {
			// Sorting the result by most probable images
			sort.Slice(val, func(i, j int) bool {
				return val[i].Probability > val[j].Probability
			})

			fmt.Println("\nFollowing are image URLs for the given keyword:")
			count := 0
			for _, item := range val {
				if count < 10 {
					fmt.Println(item.URL)
					count++
				} else {
					break
				}
			}
		} else {
			fmt.Println("\nNo Image found")
		}

		// Checking if user wants to continue searching
		fmt.Println("\nDo you want to continue?\nPress Y for Yes or any other alphabet for No:")
		yesOrNo, _ := reader.ReadString('\n')
		if strings.TrimRight(yesOrNo, "\n") != "Y" {
			flag = false
		}
		fmt.Println()
	}

}
