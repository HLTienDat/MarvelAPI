package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"
)

const (
	apiBaseURL    = "https://gateway.marvel.com/v1/public"
	publicKey     = "d4071f34c2390064906650ae83f526ef"
	privateKey    = "0abb9a54651f151e7982c6f3d863e047e16431a0"
	charactersURL = "/characters"
)

type CharacterDataContainer struct {
	Data struct {
		Results []Character `json:"results"`
	} `json:"data"`
}

type Comic struct {
	Name string `json:"name"`
}

type Character struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Comics      struct {
		Items []Comic `json:"items"`
	} `json:"comics"`
}

func randomChar() rune {
	randomValue := rand.Intn(26) + 97
	return rune(randomValue)
}
func main() {
	ts := time.Now().Format("20240407130608")
	hash := generateHash(ts)
	limit := 10
	offset := 0
	randChar := string(randomChar())
	fmt.Println("Find 10 random characters that have their name start with letter:", randChar)
	url := fmt.Sprintf("%s%s?nameStartsWith=%s&apikey=%s&ts=%s&hash=%s&limit=%d&offset=%d", apiBaseURL, charactersURL, randChar, publicKey, ts, hash, limit, offset)

	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making API request: %s\n", err.Error())
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err.Error())
		return
	}

	var characterData CharacterDataContainer
	err = json.Unmarshal(body, &characterData)
	if err != nil {
		fmt.Printf("Error parsing JSON: %s\n", err.Error())
		return
	}

	for _, character := range characterData.Data.Results {
		fmt.Println("=================")
		fmt.Printf("Character: %s\n", character.Name)
		fmt.Printf("Discription: %s", character.Description)
		if character.Description == "" {
			fmt.Print("No information")
		}
		fmt.Println()
		fmt.Println("Comics:")
		if len(character.Comics.Items) == 0 {
			fmt.Println("No comic Available")
		}
		for i, comic := range character.Comics.Items {
			fmt.Printf("%v. %s\n", i, comic.Name)
		}
		fmt.Println()
	}
}

func generateHash(ts string) string {
	data := []byte(ts + privateKey + publicKey)
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}
