package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {

	var port, err = os.LookupEnv("API_PORT")
	var congressUrl = os.Getenv("CONGRESS_API_URL")
	var congressApiKey = os.Getenv("CONGRESS_API_KEY")
	var targetURL, _ = url.Parse(congressUrl)

	if err != true {
		log.Fatal("Set environment variable API_PORT")
	}

	proxy := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		outreq := *r
		log.Default().Println(outreq.URL)
		//log.Default().Println(outreq.URL.)

		//outreq.URL = targetURL.JoinPath(outreq.URL.String())
		outreq.Host = targetURL.Host

		query := outreq.URL.Query()

		// Add a new query parameter
		query.Add("api_key", congressApiKey)

		// Update the URL with the modified query parameters
		outreq.URL.RawQuery = query.Encode()

		requrl := fmt.Sprintf("%s%s", targetURL.String(), outreq.URL.String())
		log.Default().Println(requrl)

		outresp, err := callCongressApi(requrl, outreq.Method)
		if err != nil {
			log.Default().Println(err.Error())
			http.Error(w, "Server Error", http.StatusInternalServerError)
			return
		}
		defer outresp.Body.Close()

		for k, v := range outresp.Header {
			w.Header()[k] = v
		}

		// Copy response status
		w.WriteHeader(outresp.StatusCode)

		// Copy response body
		_, err = io.Copy(w, outresp.Body)
		if err != nil {
			log.Println(err)
		}
	})

	log.Default().Printf("listening on port %v", port)

	http.ListenAndServe(":"+port, proxy)

}

func callCongressApi(url string, method string) (*http.Response, error) {

	var req, _ = http.NewRequest(method, url, nil)
	var res, err = http.DefaultClient.Do(req)

	return res, err
}
