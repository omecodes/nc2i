package nc2i

import (
	"fmt"
	"github.com/omecodes/common/utils/prompt"
	"os"
	"path"
	"path/filepath"

	"github.com/shibukawa/configdir"
	"github.com/spf13/cobra"

	"github.com/omecodes/common/utils/log"
)

var (
	Cmd                 *cobra.Command
	addr                string
	dbURI               string
	externalResourceDir string

	tlsCertFilename string
	tlsKeyFilename  string
	tlsSelfSigned   bool

	dataDir string
	resDir  string
)

func init() {
	dirs := configdir.New("Ome", "NC2I")
	globalFolder := dirs.QueryFolders(configdir.Global)[0]
	dataDir = globalFolder.Path

	if err := os.MkdirAll(dataDir, os.ModePerm); err != nil {
		fmt.Println("failed to create NC2I dir:", err)
		os.Exit(-1)
	}

	// Set log file at the data dir root
	log.File = filepath.Join(dataDir, "run.log")

	// Set default resources dir at the data dir root too
	resDir = path.Join(dataDir, "res")
	if err := os.MkdirAll(resDir, os.ModePerm); err != nil {
		fmt.Println("failed to create NC2I resources dir:", err)
		os.Exit(-1)
	}

	Cmd = &cobra.Command{
		Use:   "nc2i",
		Short: "Runs nC2i web server",
		Run:   runNC2iServer,
	}

	flags := Cmd.PersistentFlags()
	flags.StringVar(&addr, "bind", ":80", "Bind address")
	flags.StringVar(&dbURI, "db", "nc2i:nc2i@(localhost:3306)/nc2i?charset=utf8", "Database URI (MySQL is highly encouraged for production environment)")
	flags.StringVar(&externalResourceDir, "res", resDir, "External resources folder")
	flags.StringVar(&tlsCertFilename, "tls-crt", "", "TLS certificate filename")
	flags.StringVar(&externalResourceDir, "tls-key", "", "TLS key filename")
	flags.BoolVar(&tlsSelfSigned, "tls-ca", false, "Set this flag for self signed certificate")
}

func runNC2iServer(cmd *cobra.Command, args []string) {
	if tlsCertFilename != "" || tlsKeyFilename != "" || tlsSelfSigned {
		if tlsCertFilename == "" || tlsKeyFilename == "" {
			fmt.Println("tls-crt and tls-key must both be provided")
			os.Exit(-1)
		}
	}

	srv := Server{
		DataDir:         dataDir,
		ResDir:          resDir,
		BindAddr:        addr,
		TLSCertFilename: tlsCertFilename,
		TLSKeyFilename:  tlsKeyFilename,
		TLSSelfSigned:   tlsSelfSigned,
	}

	err := srv.Start()
	if err != nil {
		fmt.Println("NC2I server failed to start:", err)
		os.Exit(-1)
	}

	defer func() {
		if err := srv.Stop(); err != nil {
			fmt.Println("NC2I server stopped with error: ", err)
		}
	}()

	// waiting for CTRL+C or server error
	select {
	case <-prompt.QuitSignal():

	case err = <-srv.Errors:
		fmt.Println("NC2I server error: ", err)
	}
}
