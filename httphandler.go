package gsub

import (
	"html"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func welcomePage(w http.ResponseWriter, r *http.Request) {

	htmltxt := `<h1> Welcome to my page </h1> You will find something interesting apps. Browse local files  <a href={{.Prefix}}> {{.Prefix}} </a>. <tiny>{{.Author}}</tiny> </br>
	  The <a href="/api/apps">apps</a>  and
		<a href="/api/apps/2"/> /api/apps/2 </a>are sample API handlers`
	tmpl := template.New("welcome") //  ParseFiles("template.txt")
	tmpl.Parse(htmltxt)

	// tmpl, err := template.New("template.txt") //  ParseFiles("template.txt")
	// if err != nil {
	// 	log.Print(err)
	// }
	type Dummy struct {
		Author string
		Prefix string
	}
	info := Dummy{"ABCD", "GOVIND"}
	// tmpl.Execute(w, nil)
	tmpl.ExecuteTemplate(w, tmpl.Name(), info)

}
func serverootfiles(w http.ResponseWriter, r *http.Request) {

	log.Print("Root is here ", r.URL.Path, " escpaed = ", html.EscapeString(r.URL.Path))
	// fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	fname := "./web" + strings.TrimPrefix(r.URL.Path, "/www")
	// if fname != "" {
	http.ServeFile(w, r, fname)
	log.Print("Will server this file ", fname)
	// } else {
	// 	http.ServeFile(w, r, "./")
	// }

}
