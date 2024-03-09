package main

import (
	"WEBAPI/auth"
	"encoding/json"
	"math"
	"net/http"
	"sort"
	"strconv"

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
	RoundedCalories int
}

type RecipeHit struct {
	Recipe Recipe `json:"recipe"`
}

type EdamamResponse struct {
	Hits []RecipeHit `json:"hits"`
}

var recipes []Recipe

func submitHandler(c *gin.Context) {
	name := c.PostForm("name")
	sortCalories := c.PostForm("sortCalories") // ソートオプションの値を取得

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

	recipes := make([]Recipe, len(edamamResponse.Hits))
	for i := 0; i < len(edamamResponse.Hits); i++ {
			recipes[i] = edamamResponse.Hits[i].Recipe
	}

	for i := 0; i < len(edamamResponse.Hits); i++ {
		recipe := edamamResponse.Hits[i].Recipe
		recipe.RoundedCalories = int(math.Round(recipe.Calories))
		recipes[i] = recipe
	}

	recipes = nil
	for i := 0; i < len(edamamResponse.Hits); i++ {
		recipes = append(recipes, edamamResponse.Hits[i].Recipe)
	}

	// カロリーでソートするかどうかチェック
	if sortCalories == "on" {
			// カロリーでレシピを昇順にソート
			sort.Slice(recipes, func(i, j int) bool {
					return recipes[i].Calories < recipes[j].Calories
			})
	}

	c.HTML(http.StatusOK, "index.html", gin.H{"recipes": recipes})
}


func topHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func recipeHandler(c *gin.Context) {
	indexStr := c.Param("index")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid index")
		return
	}

	if index >= len(recipes) || index < 0 {
		c.String(http.StatusNotFound, "Recipe not found")
		return
	}

	recipe := recipes[index]
	c.HTML(http.StatusOK, "recipe.html", gin.H{"recipe": recipe})
}

func main() {
	r := gin.Default()
	auth.DbInit()
	// frontendディレクトリの中身を読み込む
	r.LoadHTMLGlob("../frontend/*")
	r.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "signup.html", nil)
	})
	r.POST("/signup", auth.SignUpUser)

	r.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{})
	})
	r.POST("/login", auth.LoginUser)


	r.GET("/", topHandler)

	r.POST("/submit", submitHandler)
	r.GET("/recipe/:index", recipeHandler)
	r.Run(":8080")
}
