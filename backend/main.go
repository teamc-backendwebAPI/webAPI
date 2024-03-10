package main

import (
	"WEBAPI/auth"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

var store = NewMemoryStore()

type Ingredient struct {
	Text   string  `json:"text"`
	Weight float64 `json:"weight"`
}

type Recipe struct {
	Label           string       `json:"label"`
	Image           string       `json:"image"`
	Source          string       `json:"source"`
	URL             string       `json:"url"`
	Yield           float64      `json:"yield"`
	Ingredients     []Ingredient `json:"ingredients"`
	Calories        float64      `json:"calories"`
	TotalTime       float64      `json:"totalTime"`
	RoundedCalories int
}

type MemoryStore struct {
	Recipes map[int]Recipe
	mu      sync.RWMutex
}

type RecipeHit struct {
	Recipe Recipe `json:"recipe"`
}

type EdamamResponse struct {
	Hits []RecipeHit `json:"hits"`
}


func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		Recipes: make(map[int]Recipe),
	}
}

// SaveRecipe はレシピをMemoryStoreに保存します。
func (store *MemoryStore) SaveRecipe(recipe Recipe, i int) {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.Recipes[i] = recipe
}

// GetRecipe は指定されたIDのレシピをMemoryStoreから取得します。
func (store *MemoryStore) GetRecipe(id int) (Recipe, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	recipe, exists := store.Recipes[id]
	return recipe, exists
}

var separateRecipes []Recipe

func submitHandler(c *gin.Context) {
	log.Println("submitHandler")
	name := c.PostForm("name")
	sortCalories := c.PostForm("sortCalories") // ソートオプションの値を取得

	// ページネーションのパラメータを取得
	page, err := strconv.Atoi(c.DefaultPostForm("page", "1"))
	if err != nil {
		log.Println("Error getting page or pageSize: ", err)
		page = 1
	}
	pageSize, err := strconv.Atoi(c.DefaultPostForm("pageSize", "6"))
	if err != nil {
		log.Println("Error getting page or pageSize: ", err)
		pageSize = 6
	}


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
		store.SaveRecipe(recipe, i)
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

	// ページネーションの情報を計算
	prevPage := page - 1
	if prevPage < 1 {
		prevPage = 1
	}
	nextPage := page + 1
	if nextPage > len(recipes) {
		nextPage = len(recipes)
	}

	// ページネーションの範囲を計算
	start := (page - 1) * pageSize
	end := start + pageSize
	if end > len(recipes) {
		end = len(recipes)
	}
	// ページネーションの範囲でレシピをフィルタリング
	separateRecipes = recipes[start:end]
	
	fmt.Println("page:", page)
	fmt.Println("prevPage:", prevPage)
	fmt.Println("nextPage:", nextPage)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"recipes": separateRecipes,
		"page": page,
		"prevPage": prevPage,
		"nextPage": nextPage,
	})
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

	if recipe, exists := store.GetRecipe(index); exists {
		c.HTML(http.StatusOK, "recipe.html", gin.H{"recipe": recipe})
	}
}
func apiDocumentationHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "apiDocument.html", nil)
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
	r.GET("/api-documentation", apiDocumentationHandler)
	r.Run(":8080")
}
