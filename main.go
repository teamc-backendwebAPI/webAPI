package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Recipe struct {
    ID          int      `json:"id"`
    Name        string   `json:"name"`
    Description string   `json:"description"`
    ImageURL    string   `json:"image"`
    Categories  []string `json:"categories"`
    Ingredients []struct {
        Name   string `json:"name"`
        Amount string `json:"amount"`
    } `json:"ingredients"`
    Steps        []string `json:"steps"`
    Nutrition    struct {
        Calories     string `json:"calories"`
        Protein      string `json:"protein"`
        Fat          string `json:"fat"`
        Carbohydrates string `json:"carbohydrates"`
    } `json:"nutrition"`
    Difficulty string `json:"difficulty"`
    Time       struct {
        Prep string `json:"prep"`
        Cook string `json:"cook"`
    } `json:"time"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}

type RecipesContainer struct {
    Recipes []Recipe `json:"recipes"`
}

var container RecipesContainer

func init() {
    file, err := ioutil.ReadFile("./sample.json")
    if err != nil {
        fmt.Println("Error reading file:", err)
        os.Exit(1)
    }

    err = json.Unmarshal(file, &container)
    if err != nil {
        fmt.Println("Error decoding JSON:", err)
        os.Exit(1)
    }
}

func recipeHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    recipeName := r.URL.Query().Get("name")
    if recipeName == "" {
        http.Error(w, "Missing recipe name", http.StatusBadRequest)
        return
    }

    for _, recipe := range container.Recipes {
        if recipe.Name == recipeName {
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(recipe)
            return
        }
    }

    http.NotFound(w, r)
}


func main() {
    http.HandleFunc("/recipe", recipeHandler)
    fmt.Println("Server is running on port 8080")
    http.ListenAndServe(":8080", nil)
}
