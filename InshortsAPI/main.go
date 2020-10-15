package main

// Importimg the required pakages
import (   
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "time"
    "io/ioutil"
    "strings"
    "strconv"
    "sync"
)
// Defining the data structure for article
type Article struct {          
	Id            int       `json:"Id"`
    Title         string    `json:"Title"`
    Subtitle      string    `json:"SubTitle"`
    Content       string    `json:"content"`
    Created_at    time.Time `json:Timestamp`
}

// Creating slices of articles 
var allArticles []Article

var mutex sync.Mutex
var wg sync.WaitGroup


func homePage(w http.ResponseWriter, r *http.Request) {

	// fmt.Println(r.URL.Path)
	if r.URL.Path == "/" {
 		fmt.Fprintf(w, "Welcome to Inshorts API")
        return
    }
	// splitting the url into a list
	path := strings.Split(r.URL.Path, "/")
	// fmt.Println(path)

	if len(path) == 3 {
		if path[1] == "articles" {
			id, err := strconv.Atoi(path[2])
			if err == nil {
				getArticlebyId(w, r, id)
				return
			}
		}
	}
	http.Error(w, "404 not found.", http.StatusNotFound)
	
}

func createArticle(w http.ResponseWriter, r *http.Request, wg *sync.WaitGroup) {

	mutex.Lock()

	var newArticle Article
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter New article with title, subtitle and content only")
	}
	n := len(allArticles)
	prevId := allArticles[n-1].Id

	json.Unmarshal(reqBody, &newArticle)
	newArticle.Created_at = time.Now()
	newArticle.Id = prevId + 1

	// Adding the newly created article to the array of articles
	allArticles = append(allArticles, newArticle)

	// Return the 201 created status code
	w.WriteHeader(http.StatusCreated)
	// Return the newly created event
	fmt.Fprintf(w, "New article is created.")
	json.NewEncoder(w).Encode(newArticle)

	mutex.Unlock()
	wg.Done()
}

func listallArticles(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(allArticles)
	return

}
func createandlistArticles(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/articles" {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }
    // /articles url is used  for creating new artcles and listing all articles in memory so here switch method is used in case of post and get request
    switch r.Method {
    case "POST":
    	wg.Add(1)
    	createArticle(w, r, &wg)
    	wg.Wait()
    case "GET":
    	listallArticles(w, r)
    default:
        fmt.Fprintf(w, "Sorry, only GET and POST methods are supported at this endpoint.")
    }

}

func getArticlebyId(w http.ResponseWriter, r *http.Request, id int) {
	// response the article which match with the article's id in get request

	key := id
    for _, article := range allArticles {
        if article.Id == key {
            json.NewEncoder(w).Encode(article)
            return
        }
    }
    fmt.Fprintf(w, "Id doesn't exists")
}
// This function will return a list of strings which contains title, subtitle, content of an anticle in it
func combinedArticle(a *Article) ([]string) {
	// creating a res string slice which will be returned
	res := make([]string, 0)
	temp := make([]string, 0)

	titleList := strings.Split(a.Title, " ")
	subtitleList := strings.Split(a.Subtitle, " ")
	contentList := strings.Split(a.Content, " ")

	temp = append(temp, titleList...)
	temp = append(temp, subtitleList...)
	temp = append(temp, contentList...)

	for _, e := range temp {
			res = append(res, strings.ToLower(e))
		}
	return res
}
func searchQuery(w http.ResponseWriter, r *http.Request) {

	keys, ok := r.URL.Query()["q"]
    
    if !ok || len(keys[0]) < 1 {
        log.Println("Url Param 'key' is missing")
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }
    // Extracting the first the key in the search query and converting it to lowercase
    key := strings.ToLower(keys[0])
    fmt.Println(key)

    var searchedArticles []Article // 
    for _, article := range allArticles {
    	combinedList := combinedArticle(&article)
    	for _, e := range combinedList {
    		if e == key {
    			searchedArticles = append(searchedArticles, article)
    			break
    		}
    	}
    }
    n := len(searchedArticles)
    if n!=0 {
    	fmt.Fprintf(w, "All searched articles")
    	json.NewEncoder(w).Encode(searchedArticles)
    	return
    }
  
    fmt.Fprintf(w, "Article not found")
   
}

func handleRequests() {
	// Get an article using Id is handled inside  the homePage Function
	http.HandleFunc("/", homePage)
	http.HandleFunc("/articles", createandlistArticles)
	http.HandleFunc("/articles/search", searchQuery)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {

	allArticles = []Article{
	
	Article{
		Id:          1,
		Title:       "Golang",
		Subtitle:    "Introduction to Golang",
		Content:     "Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.",
		Created_at:   time.Now(),
	},
}
	handleRequests()
	
}

