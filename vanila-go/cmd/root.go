package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	h "github.com/muhfaris/lib-go/vanila-go/handler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	dbPool  *sql.DB
)

var rootCmd = &cobra.Command{
	Use:   "vanila-go",
	Short: "",
	Long:  "long description of application",
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize for init package

		r := mux.NewRouter()
		var url = r.PathPrefix("/api").Subrouter()

		var apiv1 = url.PathPrefix("/v1").Subrouter()
		apiv1.HandleFunc("/home", h.HandlerApiHome)

		// SSL HTTPS config
		/*
			tlsconfig := &tls.config{
				preferserverciphersuites: true,
				curvepreferences: []tls.curveid{
					tls.curvep256,
					tls.x25519,
				},
				minversion: tls.versiontls12,
				ciphersuites: []uint16{
					tls.tls_ecdhe_ecdsa_with_aes_256_gcm_sha384,
					tls.tls_ecdhe_rsa_with_aes_256_gcm_sha384,
					tls.tls_ecdhe_ecdsa_with_chacha20_poly1305, // go 1.8 only
					tls.tls_ecdhe_rsa_with_chacha20_poly1305,   // go 1.8 only
					tls.tls_ecdhe_ecdsa_with_aes_128_gcm_sha256,
					tls.tls_ecdhe_rsa_with_aes_128_gcm_sha256,
				},
			}
		*/
		port := fmt.Sprintf(":%s", viper.GetString("app.port"))
		s := &http.Server{
			Handler:      r,
			WriteTimeout: 5 * time.Second,
			ReadTimeout:  5 * time.Second,
			Addr:         port,
			//tlsconfig:    tlsconfig,
		}

		log.Println("application running in:", port)
		//s.listenandservetls("", "")
		//log.fatal(s.listenandservetls("localhost.crt", "localhost.key"))
		log.Fatal(s.ListenAndServe())

	},
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settijjkjjj:gs.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/vanila-go.toml)")
}

func initConfig() {
	viper.SetConfigType("toml")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		//search config with name ".name"
		viper.SetConfigName("config")
		viper.AddConfigPath("./configs")

		err := viper.ReadInConfig()
		if err != nil {
			fmt.Println("Error can not read the config file:", err)
		}
	}

	//read env
	viper.AutomaticEnv()

	// if a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("using config file:", viper.ConfigFileUsed())
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
