package main

import (
	"crypto/rand"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"html/template"
	"log"
	mRand "math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	count       int
	webServer   bool
	wl          string
	special     bool
	upper       bool
	list        []string
	baseList    = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	baseSpecial = []string{"!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "_", "=", "+", "-", ">", "<", "?", ",", ".", "/", "|", "\\", " "}
	baseInt     = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10"}
	tpl         *template.Template
)

type App struct {
	router *mux.Router
}

type passGen struct {
	list        []string
	specialList []string
	baseUpper   []string
	baseInts    []string
}

func init() {
	flag.IntVar(&count, "length", 8, "Length of the password to generate.")
	flag.StringVar(&wl, "wordlist", "", "Location of wordlist to utilize for password generation.")
	flag.BoolVar(&webServer, "server", false, "Starts an Rest API to be queried for passwords. Also will generate a web UI running on port 8080")
	flag.BoolVar(&upper, "upper", true, "Determines if one desires to have upper case characters within the generated password.")
	flag.BoolVar(&special, "special", true, "Determines if one desires to DISABLE special characters within the generated password.")
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	flag.Parse()
	if count <= 0 {
		fmt.Println("Count cannot be blank")
		os.Exit(0)
	}
	if count < 8 {
		s, err := unquoteCodePoint("\\U0001F631")
		if err != nil {
			log.Printf("Error converting unicode... Error: %v", err)
		}
		fmt.Printf("%s: Be aware! Password length is less than the minumum recommended length of 8 characters, we are quitting because of this.\n", s)
		os.Exit(0)
	}
	if webServer {
		err := godotenv.Load()
		if err != nil {
			log.Printf("Error loading environment variables. Error: %v", err)
		}
		go run()
		web()
	}
	if len(wl) > 0 || wl != "" {
		pass := wordlist()
		fmt.Println(pass.genPass(count))
	} else {
		pass := getChars()
		fmt.Println(pass.genPass(count))
	}

}

func web() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		err := tpl.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			log.Printf("Error executing template: %v", err)
		}
	})
	fmt.Println("Listening on port", os.Getenv("WEBPORT"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("WEBPORT")), nil))
}

func unquoteCodePoint(s string) (string, error) {
	r, err := strconv.ParseInt(strings.TrimPrefix(s, "\\U"), 16, 32)
	return string(r), err
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
	a.router.HandleFunc("/password", genPassWeb).Methods(http.MethodPost)
	return a
}

func genPassWeb(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		log.Printf("Error parsing form. Error: %v", err)
	}

	for k, v := range req.Form {
		switch strings.ToLower(k) {
		case "length":
			c, err := strconv.Atoi(v[0])
			if err != nil {
				log.Printf("error converting string to integer value. Error: %v", err)
				http.Error(w, "Error utilizing provided length value - Forcing usage of 8 characters.", http.StatusBadRequest)
			}
			if c < 8 {
				log.Printf("Value was to small, forcing value upgrade to 8 chars. Initial Value: %v", c)
				count = 8
			} else {
				count = c
			}
		case "specchar":
			s, err := strconv.ParseBool(v[0])
			if err != nil {
				log.Printf("Error parsing bool value. Value given: %v", v)
				http.Error(w, fmt.Sprintf("Error with provided value for special characters. Value Given: %v", err), http.StatusBadRequest)
			} else {
				special = s
			}
		case "uppercase":
			s, err := strconv.ParseBool(v[0])
			if err != nil {
				log.Printf("Error parsing bool value. Value given: %v", v)
				http.Error(w, fmt.Sprintf("Error with provided value for special characters. Value Given: %v", err), http.StatusBadRequest)
			} else {
				upper = s
			}
		default:
			log.Printf("It appears a non usable value was provided. Value: %v", v)
			http.Error(w, "No valid values provided. Please provide at minumum a length value.", http.StatusBadRequest)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	pass := wordlist()
	p, _ := json.Marshal(map[string]string{"Password": pass.genPass(count)})
	w.Write(p)
}

func getChars() passGen {
	var pass passGen

	for _, i := range baseList {
		pass.list = append(pass.list, i)
		if upper {
			pass.list = append(pass.list, strings.ToUpper(i))
			pass.baseUpper = append(pass.baseUpper, strings.ToUpper(i))
		} else {
			continue
		}
	}

	for _, num := range baseInt {
		pass.list = append(pass.list, num)
		pass.baseInts = append(pass.baseInts, num)
	}

	for _, v := range baseSpecial {
		if special {
			pass.list = append(pass.list, v)
			pass.specialList = append(pass.specialList, v)
		} else {
			break
		}
	}

	return pass

}

//Obtains words from wordlist and inserts into slice to obtain random values.
func wordlist() passGen {

	var data []string

	content, err := os.Open("/home/smallz/Documents/gitclones/YouShallPass/wordlist.csv")
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
	fmt.Println(l.baseUpper)
	fmt.Println(l.specialList)
	if len(pass) > count {
		pass = pass[0:count]
	}

	for _, val := range l.specialList {
		if special {
			if strings.Contains(pass, val) {
				break
			} else {
				pass = strings.Replace(pass, string(pass[mRand.Intn(len(pass)-7)]), l.mathRand(l.specialList), -1)
				break
			}
		} else {
			if strings.Contains(pass, val) {
				if strings.Contains(pass, val) {
					pass = strings.Replace(pass, pass[strings.IndexRune(pass, val)], l.mathRand(baseList), -1)
					continue
			}
		}
	}

	for _, val := range l.baseUpper {
		if upper {
			if strings.Contains(pass, val) {
				break
			} else {
				pass = strings.Replace(pass, string(pass[mRand.Intn(len(pass)-7)]), l.mathRand(l.baseUpper), -1)
				break
			}
		} else {
			//WE do not want a random value to be replaced here, we need the index of the value to exchange
			//Use strings.indexrune to determine index of the value and then replace.
			if strings.Contains(pass, val) {
				pass = strings.Replace(pass, pass[strings.IndexRune(pass, val)], l.mathRand(baseList), -1)
				continue
			}
		}
	}
	for _, num := range l.baseInts {
		if strings.Contains(pass, num) {
			break
		} else {
			pass = strings.Replace(pass, string(pass[mRand.Intn(len(pass)-7)]), l.mathRand(l.baseInts), -1)
			break
		}
	}
	fmt.Println(pass)
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
