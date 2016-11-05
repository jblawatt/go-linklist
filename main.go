package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// FILENAME is the Constant for JSON File
const FILENAME = "links.json"

// Link ... blablabla
type Link struct {
	ID    string   `json:"id"`
	Group string   `json:"group"`
	Title string   `json:"title"`
	Link  string   `json:"link"`
	// Tags  []string `json:"tags"`
}

func initialData() {
	err := ioutil.WriteFile(FILENAME, []byte("[]"), os.ModePerm)
	if err != nil {
		log.Println("Error creating file:", err)
		os.Exit(1)
	}
	log.Println("file does not exists. creating", FILENAME)
}

func add(link *Link) {
	uniqueID := uuid.NewV4()
	link.ID = uniqueID.String()
	list := append(load(), link)
	save(list)
}

func save(links []*Link) {
	data, err := json.Marshal(links)
	if err != nil {
		log.Println("Error writing file:", err)
		os.Exit(1)
	}
	ioutil.WriteFile(FILENAME, data, os.ModePerm)
	log.Println("writing to file", FILENAME)
}

func load() []*Link {
	if _, err := os.Stat(FILENAME); err != nil {
		initialData()
	}
	file, err := ioutil.ReadFile(FILENAME)
	if err != nil {
		log.Println("Error reading file:", err)
		os.Exit(1)
	}
	var links []*Link
	json.Unmarshal(file, &links)
	log.Println("loading from file", FILENAME)
	return links
}

func removeDuplicates(xs *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *xs {
		if !found[x] {
			found[x] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}

func mainHandler(context *gin.Context) {
	tmpl, _ := template.ParseFiles("index.html")
	list := load()

	var groups []string

	for _, link := range list {
		groups = append(groups, link.Group)
	}

	removeDuplicates(&groups)

	data := map[string]interface{}{
		"Links":  list,
		"Groups": groups,
	}
	tmpl.ExecuteTemplate(context.Writer, "index.html", data)
}

func removeHandler(c *gin.Context) {
	uniqueID := c.Param("id")
	links := load()
	var removeLink *Link
	var keepLinks []*Link

	for _, link := range links {
		if link.ID == uniqueID {
			removeLink = link
		} else {
			keepLinks = append(keepLinks, link)
		}
	}

	save(keepLinks)
	c.JSON(http.StatusOK, removeLink)

}

func addHandler(context *gin.Context) {
	var link *Link
	context.BindJSON(&link)
	add(link)
	context.JSON(http.StatusCreated, link)
}

func main() {

	router := gin.Default()

	router.GET("/", mainHandler)
	router.POST("/remove/:id", removeHandler)
	router.POST("/add", addHandler)

	router.Run("127.0.0.1:3000")
}
