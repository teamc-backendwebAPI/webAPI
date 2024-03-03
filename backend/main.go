package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/template"
	"io"
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
	Steps     []string `json:"steps"`
	Nutrition struct {
		Calories      string `json:"calories"`
		Protein       string `json:"protein"`
		Fat           string `json:"fat"`
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


func submitHandler(w http.ResponseWriter, r *http.Request) {
	//index.htmlからレシピ名を取得
	name := r.FormValue("name")
	url := "https://api.edamam.com/api/recipes/v2?type=public&q=" + name + "&app_id=1f53f4d6&app_key=8cfa79ecfe3f0a623174bfa1bd2e2d4d"
	resp, err := http.Get(url)
	defer resp.Body.Close()
	
	if err != nil {
		fmt.Println("Error getting data from Edamam:", err)
		os.Exit(1)
	}


	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading data from Edamam:", err)
		os.Exit(1)
	}

	err = json.Unmarshal(body, &container)
	if err != nil {
		fmt.Println("Error decoding JSON from Edamam:", err)
		os.Exit(1)
	}

}

func recipeHandler(w http.ResponseWriter, r *http.Request) {
	//クエリパラメータからレシピ名を取得
	recipeName := r.URL.Query().Get("name")

	var foundRecipe Recipe
	for _, recipe := range container.Recipes {
		if recipe.Name == recipeName {
			foundRecipe = recipe
			break
		}
	}

	if foundRecipe.Name == "" {
		http.NotFound(w, r)
		return
	}
	//HTMLのテンプレートに値を渡す
	tmpl := template.Must(template.ParseFiles("../frontend/index.html"))
	err := tmpl.Execute(w, foundRecipe)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
	http.NotFound(w, r)

}

func main() {
	// htmlから受け取った内容をweb APIに送信
	http.HandleFunc("/submit", submitHandler)
	// web APIから受け取った内容をhtmlに送信
	http.HandleFunc("/recipe", recipeHandler)
	http.Handle("/", http.FileServer(http.Dir(".")))
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
