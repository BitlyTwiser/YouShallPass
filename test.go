package main

import (
	"crypto/rand"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"net/http"
	"log"
	mRand "math/rand"
	"os"
	//"strconv"
	"strings"
	"time"
)

var (
	count int
	webServer bool
	wl    string
	list        []string
	bU          []string
	baseList    = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	baseSpecial = []string{"!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "_", "=", "+", "-", ">", "<", "?", ",", ".", "/", "|", "\\", " "}
	data        []string
	baseInt = []string{"1","2","3","4","5","6","7","8","9","10"}
)

type App struct {
	router *mux.Router
}

type passGen struct {
	list        []string
	specialList []string
	baseUpper   []string
	baseInts []string
}

func init() {
	flag.IntVar(&count, "length", 8, "Length of the password to generate.")
	flag.StringVar(&wl, "wordlist", "", "Location of wordlist to utilize for password generation.")
	flag.BoolVar(&webServer, "server", false, "Starts an Rest API to be queried for passwords.")
}

func main() {
	flag.Parse()
	if count <= 0 {
		fmt.Println("Count cannot be blank")
		os.Exit(0)
	}
	if count < 8 {
		fmt.Println("Be aware! Password length is less than the minumum recommended length of 8 characters.")
	}
	if webServer {
		err := godotenv.Load()
		if err != nil {
			log.Printf("Error loading environment variables. Error: %v", err)
		}
		run()
	}
	if len(wl) > 0 || wl != "" {
		pass := wordlist()
		fmt.Println(pass.genPass(count))
	} else {
		pass := getChars()
		fmt.Println(pass.genPass(count))
	}

}

func run() {
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
	a.router.HandleFunc("/password", genPassWeb).Methods(http.MethodGet)
	return a
}

func genPassWeb(w http.ResponseWriter, req *http.Request){

}

func getChars() passGen {
	var pass passGen

	for _, i := range baseList {
		pass.list = append(pass.list, i)
		pass.list = append(pass.list, strings.ToUpper(i))
		pass.baseUpper = append(pass.baseUpper, strings.ToUpper(i))
	}

	for _, num := range baseInt {
		pass.list = append(pass.list, num)
		pass.baseInts = append(pass.baseInts, num)
	}

	for _, v := range baseSpecial {
		pass.list = append(pass.list, v)
		pass.specialList = append(pass.specialList, v)
	}
	return pass

}

//Obtains words from wordlist and inserts into slice to obtain random values.
func wordlist() passGen {

	content, err := os.Open("/home/smallz/Documents/gitclones/YouShallPass/parsed.csv")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer content.Close()

	lines, err := csv.NewReader(content).ReadAll()
	if err != nil {
		log.Fatalf("Error: %v")
		os.Exit(0)
	}

	for _, line := range lines {
		data = append(data, line[0])
	}

	for _, char := range getChars().list {
		data = append(data, char)
	}

	return passGen{list: data, specialList: getChars().specialList, baseUpper: getChars().baseUpper, baseInts: getChars().baseInts}

}

//Generates random slice of bytes up to the numerical digit presented.
// The length will be a flag for length of password.. after this is passed ito the slice generator the integers therein will be utilized to obtain the values randomly.
func randGen() []byte {
	by := make([]byte, 20)
	_, err := rand.Read(by)
	if err != nil {
		log.Printf("error: %v", err)
	}
	return by
}

//Validates that a given password holds a letter, number, and special char.
func (l *passGen) passValidation(pass string, count int) string {
	mRand.Seed(time.Now().UnixNano())

	if len(pass) > count {
		pass = pass[0:count]
	}
	for _, val := range l.specialList {
		if strings.Contains(pass, val) {
			break
		} else {
			pass = strings.Replace(pass, string(pass[mRand.Intn(len(pass)-7)]), l.mathRand(l.specialList), -1)
			break
		}
	}
	for _, val := range l.baseUpper {
		if strings.Contains(pass, val) {
			break
		} else {
			pass = strings.Replace(pass, string(pass[mRand.Intn(len(pass)-7)]), l.mathRand(l.baseUpper), -1)
			break
		}
	}
	for _, num := range l.baseInts{
		if strings.Contains(pass, num) {
			break
		} else {
			pass = strings.Replace(pass, string(pass[mRand.Intn(len(pass)-7)]), l.mathRand(l.baseInts), -1)
			break
		}
	}
	return pass
}

func (l *passGen) mathRand(list []string) string {
	mRand.Seed(int64(randGen()[0]))
	return list[mRand.Intn(len(list))]
}

func (l *passGen) genPass(count int) string {
	var tempPass []string
	for i := 0; i < count; i++ {
		tempPass = append(tempPass, l.mathRand(l.list))
	}

	return l.passValidation(strings.Join(tempPass, ""), count)
}
