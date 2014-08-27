package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/garyburd/go-oauth/oauth"
	"github.com/garyburd/twitterstream"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
)

var (
	configPath  = flag.String("config", "config.json", "Path to configuration file containing the application's credentials.")
	accessToken oauth.Credentials
	oauthClient = oauth.Client{
		TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
		ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/auhorize",
		TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
	}
)

func readConfig() error {
	b, err := ioutil.ReadFile(*configPath)
	if err != nil {
		return err
	}
	var config = struct {
		Consumer, Access *oauth.Credentials
	}{
		&oauthClient.Credentials, &accessToken,
	}
	return json.Unmarshal(b, &config)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n %s keyword ...\n", os.Args[0], os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if err := readConfig(); err != nil {
		log.Fatalf("Error reading configration, %v", err)
	}

	ts, err := twitterstream.Open(
		&oauthClient,
		&accessToken,
		"https://stream.twitter.com/1.1/statuses/filter.json",
		url.Values{"track": {strings.Join(flag.Args(), ", ")}},
	)
	if err != nil {
		log.Fatal(err)
	}

	defer ts.Close()

	for ts.Err == nil {
		var t interface{}
		if err := ts.UnmarshalNext(&t); err != nil {
			log.Fatal(err)
		}
		log.Printf("unmarshal: %v", t)
	}
	log.Print(ts.Err)
}
