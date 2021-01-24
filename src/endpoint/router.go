package endpoint

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

const urlPrefix = "/doe"

type Router interface {
	ListenAndServe() error
}

func NewRouter(params RouterParams) (Router, error) {
	middlewareManager := negroni.New()
	middlewareManager.Use(negroni.NewRecovery())
	r, err := newRouter(params)
	if err != nil {
		return nil, err
	}
	middlewareManager.UseHandler(r)
	server := &http.Server{
		Addr:    os.Getenv(HostURLEnvKey),
		Handler: middlewareManager,
	}
	params.Logger.Infof("==== Client Microservice started on port %s", server.Addr)
	return server, nil
}

// newRouter initializes clientAPI microservice and returns router
func newRouter(params RouterParams) (*mux.Router, error) {
	srv := params.RPC
	router := mux.NewRouter().PathPrefix(urlPrefix).Subrouter()
	router.Use(commonMiddleware)
	{
		router.HandleFunc("/upload", srv.UploadData).Methods(http.MethodPost)
		router.HandleFunc("/ports/{port-id}", srv.GetOne).Methods(http.MethodGet)
		router.HandleFunc("/ports", srv.GetAll).Methods(http.MethodGet)
	}
	var corsRouter = mux.NewRouter()
	{
		corsRouter.PathPrefix(urlPrefix).Handler(negroni.New(
			cors.New(cors.Options{
				AllowedMethods: []string{http.MethodGet, http.MethodPost},
			}),
			negroni.Wrap(router),
		))
	}
	return corsRouter, nil
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
