// Package cmd implements access-cli commands
package cmd

/*
Copyright © 2020 Barracuda Networks, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	goruntime "runtime"

	"github.com/gbl08ma/httpcache"
	"github.com/gbl08ma/httpcache/diskcache"
	"github.com/motemen/go-loghttp"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	apiclient "github.com/barracuda-cloudgen-access/access-cli/client"
)

var cfgFile string
var authFile string

var cfgViper *viper.Viper
var authViper *viper.Viper

type globalInfo struct {
	Transport        *httptransport.Runtime
	Client           *apiclient.CloudGenAccessConsole
	AuthWriter       runtime.ClientAuthInfoWriter
	VerboseLevel     int
	WriteFiles       bool
	FetchPerPage     int
	DefaultRangeSize int
	CurrentTenant    string
	FilterData       map[*cobra.Command]*filterData
	InputData        map[*cobra.Command]*inputData
	MultiOpData      map[*cobra.Command]*multiOpData
}

var global globalInfo

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   ApplicationName,
	Short: "Command-line client for the CloudGen Access Console",
	Long:  ApplicationName + ` allows access to all Console APIs from the command line`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(versionInfo *VersionInformation) {
	version = *versionInfo
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().SortFlags = false
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initAuthConfig)
	cobra.OnInitialize(initClient)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	d := filepath.Join(getUserConfigPath(), ConfigFileName)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is "+d+")")
	d = filepath.Join(getUserConfigPath(), AuthFileName)
	rootCmd.PersistentFlags().StringVar(&authFile, "auth", "", "credentials file (default is "+d+")")
	rootCmd.PersistentFlags().IntVarP(&global.VerboseLevel, "verbose", "v", 0, "verbose output level, higher levels are more verbose")

	rootCmd.PersistentFlags().SetNormalizeFunc(aliasNormalizeFunc)

	rootCmd.SetOut(os.Stdout)
}

func aliasNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	switch name {
	case "record-start":
		name = "range-start"
	case "record-end":
		name = "range-end"
	case "format":
		name = "output"
	}
	return pflag.NormalizedName(name)
}

func initClient() {
	endpoint := authViper.GetString(ckeyAuthEndpoint)
	if endpoint == "" {
		return
	}

	transport := http.DefaultTransport

	schemes := []string{"https"}
	insecureUseHTTP := authViper.GetBool(ckeyAuthUseInsecureHTTP)
	if insecureUseHTTP {
		fmt.Fprintln(os.Stderr, "WARNING: HTTP, instead of HTTPS, is being used for API communication. THIS IS INSECURE.")
		schemes = []string{"http"}
	}

	insecureSkipVerify := authViper.GetBool(ckeyAuthSkipTLSVerify)
	if insecureSkipVerify && !insecureUseHTTP {
		fmt.Fprintln(os.Stderr, "WARNING: TLS certificate verification is being skipped for the endpoint. THIS IS INSECURE.")
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	if authViper.GetBool(ckeyAuthUseCache) {
		path := cfgViper.GetString(ckeyCachePath)
		path = filepath.Join(path, "httpcache")
		cache := diskcache.New(path)
		// wrap transport in httpcache
		cachedTransport := httpcache.NewTransport(cache)
		cachedTransport.Transport = transport
		transport = cachedTransport
	}

	// the order of these two wrappings is important for output to make more sense when VerboseLevel > 2
	if global.VerboseLevel > 1 {
		// wrap transport in loghttp
		transport = &loghttp.Transport{
			Transport: transport,
		}
	}
	if global.VerboseLevel > 2 {
		transport = &dumpRequestResponseTransport{
			T: transport,
		}
	}

	// setUserAgentTransport must wrap at the end,
	// otherwise the updated user agent does not show in the dumpRequestResponseTransport dumps
	transport = &setUserAgentTransport{
		T: transport,
	}

	global.Transport = httptransport.New(endpoint, "/api/v1", schemes)
	global.Transport.Transport = transport

	global.Client = apiclient.New(global.Transport, strfmt.Default)
	global.FetchPerPage = cfgViper.GetInt(ckeyRecordsPerGetRequest)
	if global.FetchPerPage > 100 {
		fmt.Fprintf(os.Stderr, "WARNING: %s setting exceeds limit of 100. Limiting to 100.\n", ckeyRecordsPerGetRequest)
		global.FetchPerPage = 100
	} else if global.FetchPerPage < 1 {
		fmt.Fprintf(os.Stderr, "WARNING: %s setting is invalid. Setting to 50.\n", ckeyRecordsPerGetRequest)
		global.FetchPerPage = 50
	}

	global.DefaultRangeSize = cfgViper.GetInt(ckeyDefaultRangeSize)
	if global.DefaultRangeSize < 1 {
		fmt.Fprintf(os.Stderr, "WARNING: %s setting is invalid. Setting to 20.\n", ckeyDefaultRangeSize)
		global.DefaultRangeSize = 20
	}

	switch authViper.GetString(ckeyAuthMethod) {
	case authMethodBearerToken:
		accessToken := authViper.GetString(ckeyAuthAccessToken)
		client := authViper.GetString(ckeyAuthClient)
		uid := authViper.GetString(ckeyAuthUID)
		global.AuthWriter = AccessAPIKeyAuth(accessToken, client, uid)
	default:
	}

	global.CurrentTenant = authViper.GetString(ckeyAuthCurrentTenant)
}

// AccessAPIKeyAuth provides an API key auth info writer
func AccessAPIKeyAuth(accessToken, client, uid string) runtime.ClientAuthInfoWriter {
	return runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		err := r.SetHeaderParam("access-token", accessToken)
		if err != nil {
			return err
		}

		err = r.SetHeaderParam("client", client)
		if err != nil {
			return err
		}

		return r.SetHeaderParam("uid", uid)
	})
}

type setUserAgentTransport struct {
	T http.RoundTripper
}

func (t *setUserAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", fmt.Sprintf("access-cli/%s (%s; %s)", version.Version, goruntime.GOOS, goruntime.GOARCH))
	return t.T.RoundTrip(req)
}

type dumpRequestResponseTransport struct {
	T http.RoundTripper
}

func (t *dumpRequestResponseTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	b, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s\n", string(b))

	res, err := t.T.RoundTrip(req)
	if err != nil {
		return res, err
	}
	b, err = httputil.DumpResponse(res, true)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s\n", string(b))
	return res, err
}
