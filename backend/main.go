package main

import (
	"encoding/json"
	"net/http"
	"github.com/gin-gonic/gin"
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

func submitHandler(c *gin.Context) {
	name := c.PostForm("name")
	url := "https://api.edamam.com/api/recipes/v2?type=public&q=" + name + "&app_id=1f53f4d6&app_key=8cfa79ecfe3f0a623174bfa1bd2e2d4d"
	resp, err := http.Get(url)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error getting data from Edamam: "+err.Error())
		return
	}
	defer resp.Body.Close()

	var edamamResponse EdamamResponse
	err = json.NewDecoder(resp.Body).Decode(&edamamResponse)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error decoding JSON from Edamam: "+err.Error())
		return
	}

	if len(edamamResponse.Hits) == 0 {
		c.String(http.StatusNotFound, "Error not found")
		return
	}

	//foundRecipe := edamamResponse.Hits[0].Recipe
	//c.HTML(http.StatusOK, "index.html", foundRecipe)

	recipes := make([]Recipe, len(edamamResponse.Hits))
	for i := 0; i < len(edamamResponse.Hits); i++ {
		recipes[i] = edamamResponse.Hits[i].Recipe
	}
	c.HTML(http.StatusOK, "index.html", gin.H{"recipes": recipes})
}

func topHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func main() {
	r := gin.Default()
	// frontendディレクトリの中身を読み込む
	r.LoadHTMLGlob("../frontend/*")
	r.GET("/", topHandler)

	r.POST("/submit", submitHandler)
	r.Run(":8080")
}