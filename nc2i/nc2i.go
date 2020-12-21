package nc2i

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/shibukawa/configdir"
	"github.com/spf13/cobra"

	"github.com/omecodes/common/utils/log"
	"github.com/omecodes/common/utils/prompt"
)

var (
	Cmd                 *cobra.Command
	addr                string
	dbURI               string
	mailerURI           string
	managerEmail        string
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
	flags.StringVar(&dbURI, "db", "", "Database URI (MySQL is highly encouraged for production environment)")
	flags.StringVar(&mailerURI, "mailer", "", "Mailer source name")
	flags.StringVar(&managerEmail, "email", "", "Notification email. Or set it in NC2I_MAILER env variable")
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

	if mailerURI == "" {
		mailerURI = os.Getenv("NC2I_MAILER")
		if mailerURI == "" {
			fmt.Println("mailer uri is required")
			_ = cmd.Help()
			os.Exit(-1)
		}
	}

	if managerEmail == "" {
		managerEmail = os.Getenv("NC2I_MAN_EMAIL")
		if managerEmail == "" {
			fmt.Println("manager email is required")
			_ = cmd.Help()
			os.Exit(-1)
		}
	}

	if os.Getenv("NC2I_DB") != "" {
		dbURI = os.Getenv("NC2I_DB")
		if dbURI == "" {
			dbURI = "bome:bome@(127.0.0.1:3306)/bome?charset=utf8"
		}
	}

	srv := Server{
		MailerSourceName: mailerURI,
		Email:            managerEmail,
		DataDir:          dataDir,
		ResDir:           resDir,
		DBUri:            dbURI,
		BindAddr:         addr,
		TLSCertFilename:  tlsCertFilename,
		TLSKeyFilename:   tlsKeyFilename,
		TLSSelfSigned:    tlsSelfSigned,
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
