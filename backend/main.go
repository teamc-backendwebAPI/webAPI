package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type Ingredient struct {
	Text   string  `json:"text"`
	Weight float64 `json:"weight"`
}

type Recipe struct {
	Label       string       `json:"label"`
	Image       string       `json:"image"`
	Source      string       `json:"source"`
	URL         string       `json:"url"`
	Yield       float64      `json:"yield"`
	Ingredients []Ingredient `json:"ingredients"`
	Calories    float64      `json:"calories"`
	TotalTime   float64      `json:"totalTime"`
}

type RecipeHit struct {
	Recipe Recipe `json:"recipe"`
}

type EdamamResponse struct {
	Hits []RecipeHit `json:"hits"`
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	// index.htmlからレシピ名を取得
	name := r.FormValue("name")
	url := "https://api.edamam.com/api/recipes/v2?type=public&q=" + name + "&app_id=1f53f4d6&app_key=8cfa79ecfe3f0a623174bfa1bd2e2d4d"
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Error getting data from Edamam: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var edamamResponse EdamamResponse
	err = json.NewDecoder(resp.Body).Decode(&edamamResponse)
	if err != nil {
		http.Error(w, "Error decoding JSON from Edamam: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// レシピを取得できなかった場合
	if len(edamamResponse.Hits) == 0 {
		http.Error(w, "Error not found", http.StatusNotFound)
		return
	}

	// レシピを1つだけ取得する
	foundRecipe := edamamResponse.Hits[0].Recipe

	// HTMLのテンプレートに値を渡す
	tmpl := template.Must(template.ParseFiles("../frontend/index.html"))
	err = tmpl.Execute(w, foundRecipe)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
func topHandler(w http.ResponseWriter, r *http.Request) {
	//テンプレートを表示するだけ。submitHundlerの軽量版
	tmpl := template.Must(template.ParseFiles("../frontend/index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	// htmlから受け取った内容をweb APIに送信
	http.HandleFunc("/", topHandler)
	http.HandleFunc("/submit", submitHandler)
	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
