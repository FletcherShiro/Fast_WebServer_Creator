package main

import "fmt"
import "log"
import "os"
import "html/template"
import "io/ioutil"
import "net/http"
import "strconv"
import "encoding/json"

// Config structure
var RenderingFolder = "templates"

type Config struct {
	WS struct {
		Port int
		RenderingFolder string
	}
}


// The main runtime of the server

func main() {
	fmt.Println(`
	__      __        _      ___                               
	\ \    / /  ___  | |__  / __|  ___   _ _  __ __  ___   _ _ 
	 \ \/\/ /  / -_) | '_ \ \__ \ / -_) | '_| \ V / / -_) | '_|
	  \_/\_/   \___| |_.__/ |___/ \___| |_|    \_/  \___| |_|  `)
	fmt.Println()
	fmt.Println()
	log.Println("/ SERVER: CONFIG]: Starting loading all configurations...")

	file, err := os.Open("config.json")
	if err != nil {
		log.Println("/ SERVER: CONFIG]: Can't open config file, creating it and restarting...")
		os.Create("config.json")
		ioutil.WriteFile("config.json", []byte(`{
    "WS": {
        "Port": 3000,
        "RenderingFolder": "templates"
    }
}`), 0644)
		main()
		return
	}

	defer file.Close()
	decoder := json.NewDecoder(file)
	Config := Config{}
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatal("/ SERVER: CONFIG]: Can't decode config JSON: ", err)
	}

	RenderingFolder = Config.WS.RenderingFolder

	log.Println("/ SERVER: CONFIG]: All configurations loaded success")
	log.Println("/ SERVER: WS: RENDERING]: Starting loading all documents...")

	file, err = os.Open(RenderingFolder + "/index.html")
	if err != nil {
		os.Mkdir(RenderingFolder, 0755)
		log.Println("/ SERVER: WS: RENDERING]: Can't open \"" + RenderingFolder + "/index.html\" file, creating it and restarting...")
		os.Create(RenderingFolder + "/index.html")
		ioutil.WriteFile(RenderingFolder + "/index.html", []byte(`<h1>Hello World!</h1>`), 0644)
		main()
		return
	}

	file, err = os.Open(RenderingFolder + "/404.html")
	if err != nil {
		os.Mkdir(RenderingFolder, 0755)
		log.Println("/ SERVER: WS: RENDERING]: Can't open \"" + RenderingFolder + "/404.html\" file, creating it and restarting...")
		os.Create(RenderingFolder + "/404.html")
		ioutil.WriteFile(RenderingFolder + "/404.html", []byte(`<h1>404 - Not Found</h1>`), 0644)
		main()
		return
	}

	log.Println("/ SERVER: WS: RENDERING]: All main files rendered success")
	log.Println("/ SERVER: WS]: Starting WebServer on port", strconv.Itoa(Config.WS.Port) + "...")

	log.Println("/ SERVER: WS]: Server started on port", Config.WS.Port)

	http.HandleFunc("/", Request(index))
	http.ListenAndServe(":" + strconv.Itoa(Config.WS.Port), nil)
}


// Handlers for HTTP request

func Request(req http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("/ SERVER: WS]: %s %s %s\n", r.RemoteAddr, r.Method, r.URL)
        req(w, r)
    }
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
    w.WriteHeader(status)
    if status == http.StatusNotFound {
        tmpl, err := template.ParseFiles(RenderingFolder + "/404.html")

		if err != nil {
			log.Fatal("/ SERVER: WS]: ", err)
		}

		tmpl.Execute(w, nil)
    }
}


// Templates

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
      errorHandler(w, r, http.StatusNotFound)
      return
    }

	tmpl, err := template.ParseFiles(RenderingFolder + "/index.html")

	if err != nil {
		log.Fatal("/ SERVER: WS]: ", err)
	}

	tmpl.Execute(w, nil)
}