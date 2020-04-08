package main

import (
	"fmt"
	"log"
	"flags"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

/*
Todo:
[]: Add wordlist flag to ingest a wordlist and consume words from the wordlist to utilize.
[]: Combine the arrays into a single array to withdraw data from.
[]: Flag for length
[]: Flag for any chars not desired.
[]: flag for web server? We will eventually have a UI, so we can have a web serer and a standard CLI app?
*/

var count int
var webServer bool

type App struct {
	router *mux.Router
}

func init(){
	flag.IntVar(&count, "length", 8, "Length of the password to generate.")
	flag.BoolVar(*webServer, "server", false, "Starts a web server to be queried for passwords.")
}


func run() {
	flag.Parse()

	//Load new router
	mux := newRouter()
	port := os.Getenv("PORT")
	fmt.Println("Listening on ", port)
	serv := &http.Server{
		Addr: fmt.Sprintf(":%v", port),
		Handler: handlers.CORS(
			handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
			handlers.AllowedOrigins([]string{"*"}),
			handlers.AllowedHeaders([]string{
				"Content-Type", "X-Requested-With", "Access-Control-Allow-Origin", "Origin",
				"Accept", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization",
			}),
		)(mux.router),
	}

	log.Fatal(serv.ListenAndServe())
}

func (a App) ServeHTTP(w http.ResponseWriter, r *http.Request) {}


func newRouter() App {
	a := App{router: mux.NewRouter()}
	a.router.Use(mux.CORSMethodMiddleware(a.router))
	a.router.Handle("/favicon.ico", http.NotFoundHandler())
	a.router.Handle("/password", genPasswordWeb).Methods(http.MethodGet)
	return a
}

func main(){

	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading environment variables. Error: %v", err)
	}
	run()
}

func genGuid(){

}