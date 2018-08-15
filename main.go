package main

import (
  "fmt"
  "github.com/gocolly/colly"
  "github.com/go-chi/chi"
  "net/http"
  "html/template"
  "os"
  "path/filepath"
  "strings"
)

type Person struct{
  Username string
}
type ScrappedData struct {
  Data [][]string
}

func main() {
  port := os.Getenv("PORT")

  if port == "" {
    log.Fatal("$PORT must be set")
  }

  fmt.Println("Starting...")

  var tpl = template.Must(template.ParseGlob("templates/*.html"))  
  r := chi.NewRouter()
  workDir, _ := os.Getwd()
  filesDir := filepath.Join(workDir, "static")
  FileServer(r, "/static", http.Dir(filesDir))

  r.Get("/", func(w http.ResponseWriter, r *http.Request){
    data,size := fetchData()
    var finalData[][] string
    finalData = resizeData(data, size)
    d := ScrappedData{Data: finalData}
    tpl.ExecuteTemplate(w, "index.html", d)
  })
  http.ListenAndServe(":" + port, r)
}

func fetchData()([100][2]string, int) {
  var data[100][2] string
  index := 0

  c := colly.NewCollector(
    colly.AllowedDomains("www.fixkick.com"),
  )

  c.OnHTML("a[href]", func(e *colly.HTMLElement){
    link := e.Attr("href")
    //fmt.Printf("Link found: %q -> %s\n", e.Text, link)
    if (len(e.Text) > 5) && (index < 100) {
      data[index][0] = e.Text
      data[index][1] = e.Request.AbsoluteURL(link)
      index = index + 1
    }
  })

  c.OnRequest(func(r *colly.Request){
    fmt.Println("Visiting", r.URL.String())
  })

  c.Visit("http://www.fixkick.com/")
  return data, index
}

func resizeData(ar [100][2] string, size int)[][]string{
  var d [][]string
  d = make([][]string, size)
  for i := 0; i < size; i++ {
    var r []string
    r = make([]string, 2)
    r[0] = ar[i][0]
    r[1] = ar[i][1]
    d[i] = r
  }
  return d
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

