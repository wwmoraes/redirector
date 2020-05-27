package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
)

var targetURL = ""

func printError(err error) {
	if err != nil {
		log.Print(err)
	}
}

func fatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func serve_https(wg *sync.WaitGroup) {
	host := os.Getenv("HTTPS_HOST")
	port := getEnv("HTTPS_PORT", "8081")

	keyFilePath := os.Getenv("KEY_FILE")
	certFilePath := os.Getenv("CERT_FILE")

	if key, ok := os.LookupEnv("KEY"); ok {
		log.Println("creating temporary key file...")
		keyFile, err := ioutil.TempFile("", "")
		printError(err)
		keyData, err := base64.StdEncoding.DecodeString(key)
		printError(err)
		_, err = keyFile.Write(keyData)
		printError(err)
		keyFilePath = keyFile.Name()
	}

	if cert, ok := os.LookupEnv("CERT"); ok {
		log.Println("creating temporary cert file...")
		certData, err := base64.StdEncoding.DecodeString(cert)
		printError(err)
		certFile, err := ioutil.TempFile("", "")
		printError(err)
		_, err = certFile.Write(certData)
		printError(err)
		certFilePath = certFile.Name()
	}

	log.Printf("** Redirecting HTTPS from %s:%s to %s **", host, port, targetURL)
	if err := http.ListenAndServeTLS(fmt.Sprintf("%s:%s", host, port), certFilePath, keyFilePath, nil); err != nil {
		if os.IsNotExist(err) {
			log.Println("error: unable to start HTTPS - did you set KEY_FILE and CERT_FILE?")
			log.Println(err)
		} else if err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}

	wg.Done()
}

func serve_http(wg *sync.WaitGroup) {
	host := os.Getenv("HTTP_HOST")
	port := getEnv("HTTP_PORT", "8080")

	log.Printf("** Redirecting HTTP from %s:%s to %s **", host, port, targetURL)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), nil); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	wg.Done()
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.TLS == nil {
		r.URL.Scheme = "http"
	} else {
		r.URL.Scheme = "https"
	}

	log.Printf("received request from %s to %s://%s%s", r.RemoteAddr, r.URL.Scheme, r.Host, r.RequestURI)
	w.Header().Add("Location", fmt.Sprintf("https://%s%s", targetURL, r.URL.Path))
}

func main() {
	envUrl, ok := os.LookupEnv("URL")
	if !ok || len(envUrl) == 0 || !IsUrl(envUrl) {
		log.Fatal("set URL environment variable with the desired domain (e.g. https://artero.dev)")
	}

	targetURL = envUrl

	http.HandleFunc("/", RedirectHandler)

	var wg sync.WaitGroup
	wg.Add(2)

	go serve_http(&wg)
	go serve_https(&wg)

	wg.Wait()
}
