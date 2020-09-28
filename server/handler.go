package server

import (
	"log"
	"net/http"

	"github.com/Eldius/learning-go/poc-reverse-proxy/config"
)

var cfg *config.RoutesConfig

func HandleRequest(w http.ResponseWriter, req *http.Request) {
	r := cfg.GetRoute(req)
	if r != nil {
		r.Redirect(w, req)
	}
}

func refreshRoutes() {
	_cfg, err := config.LoadRoutes()
	if err != nil {
		log.Println("Failed to load routes")
		log.Fatal(err.Error())
	}
	cfg = &_cfg
}

func Start(host string, port int) {

	refreshRoutes()
	mux := http.NewServeMux()

	mux.HandleFunc("/", HandleRequest)
}
