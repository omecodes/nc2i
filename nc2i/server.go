package nc2i

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/omecodes/bome"
	"github.com/omecodes/common/httpx"
	"github.com/omecodes/common/utils/log"
	"github.com/omecodes/libome/crypt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rakyll/statik/fs"
	"net"
	"net/http"
)

type Server struct {
	MailerSourceName string
	Email            string
	DataDir          string
	ResDir           string
	DBUri            string
	BindAddr         string
	TLSCertFilename  string
	TLSKeyFilename   string
	TLSSelfSigned    bool

	initialized bool
	listener    net.Listener

	messages        *bome.JSONList
	visitsInfoStore *bome.JSONList
	resFS           http.FileSystem

	Errors chan error
}

func (srv *Server) init() error {
	if srv.initialized {
		return nil
	}
	srv.initialized = true
	srv.Errors = make(chan error, 1)

	db, err := sql.Open(bome.MySQL, srv.DBUri)
	if err != nil {
		return err
	}

	srv.messages, err = bome.NewJSONList(db, bome.MySQL, "messages")
	if err != nil {
		return err
	}

	srv.visitsInfoStore, err = bome.NewJSONList(db, bome.MySQL, "visits")
	if err != nil {
		return err
	}

	srv.resFS, err = fs.New()
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) listen() error {
	var (
		err error
		tc  *tls.Config
	)

	if srv.TLSCertFilename != "" || srv.TLSKeyFilename != "" {
		if srv.TLSCertFilename == "" {
			log.Fatal("missing certificate file path")
		}

		if srv.TLSKeyFilename == "" {
			log.Fatal("missing key file path")
		}

		cert, err := crypt.LoadCertificate(srv.TLSCertFilename)
		if err != nil {
			log.Fatal("loading certificate", log.Err(err))
		}

		key, err := crypt.LoadPrivateKey(nil, srv.TLSKeyFilename)
		if err != nil {
			log.Fatal("loading key", log.Err(err))
		}

		tc = &tls.Config{
			Certificates: []tls.Certificate{{
				Certificate: [][]byte{cert.Raw},
				PrivateKey:  key,
			}},
		}

		// Add resolved certificate as clients CA root in case of self signed certificate
		if srv.TLSSelfSigned {
			pool := x509.NewCertPool()
			pool.AddCert(cert)
			tc.ClientCAs = pool
		}
	}

	if tc != nil {
		srv.listener, err = tls.Listen("tcp", srv.BindAddr, tc)
	} else {
		srv.listener, err = net.Listen("tcp", srv.BindAddr)
	}
	if err != nil {
		return err
	}

	address := srv.listener.Addr().String()
	log.Info("starting HTTP server", log.Field("address", address))

	return nil
}

// Start starts the web server
func (srv *Server) Start() error {

	if err := srv.init(); err != nil {
		return err
	}

	if err := srv.listen(); err != nil {
		return err
	}

	go func() {
		router := srv.getHTTPRouter()
		if err := http.Serve(srv.listener, router); err != nil {
			srv.Errors <- err
		}
	}()

	return nil
}

func (srv *Server) getHTTPRouter() http.Handler {
	router := mux.NewRouter()

	router.Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpx.Redirect(w, &httpx.RedirectURL{
			URL:         appRoute,
			Code:        http.StatusMovedPermanently,
			ContentType: "text/html",
		})
	})
	router.PathPrefix(appRoute).Subrouter().
		Name("web-app").
		Handler(http.StripPrefix(appRoute, srv.middleware(http.HandlerFunc(serveWebApp)))).Methods(http.MethodGet)

	router.Handle(messagesRoute, srv.middleware(http.HandlerFunc(saveMessage))).Methods(http.MethodPost)
	router.Handle(metricsRoute, promhttp.Handler())
	return router
}

func (srv *Server) updateHttpContext(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = contextWithDataDir(ctx, srv.DataDir)
		ctx = contextWithExternalResDir(ctx, srv.ResDir)
		ctx = contextWithVisitsInfoStore(ctx, srv.visitsInfoStore)
		ctx = contextWithMessages(ctx, srv.messages)
		ctx = contextWithResFS(ctx, srv.resFS)
		ctx = contextWithNotificationEmail(ctx, srv.Email)
		ctx = contextWithMailerSourceName(ctx, srv.MailerSourceName)

		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
	})
}

func (srv *Server) middleware(handler http.Handler) http.Handler {

	wrappers := []middleware{
		visits,
		srv.updateHttpContext,
		logHandler.Handle,
	}

	for _, m := range wrappers {
		handler = m(handler)
	}
	return handler
}

// Start stops the web server
func (srv *Server) Stop() error {
	return srv.listener.Close()
}
