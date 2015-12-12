package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var listen *string = flag.String("listen", ":8000", "\"[address]:<port>\" to bind to. [default: \":8000\"]")
var username *string = flag.String("username", "", "GitHub username to fetch keys for. [required]")
var ttl *int64 = flag.Int64("ttl", 86400, "Time in seconds to cache GitHub keys for. [default: 86400 (one day)]")

var cache []string = make([]string, 0)
var expire int64 = 0

func fetchKeys() error {
	fmt.Printf("Fetching keys for GitHub user \"%s\"\n", *username)
	var resp *http.Response
	var err error
	var uri string = fmt.Sprintf("https://api.github.com/users/%s/keys", *username)
	resp, err = http.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var keys []map[string]interface{} = make([]map[string]interface{}, 0)
	err = json.Unmarshal(body, &keys)
	if err != nil {
		return err
	}

	var newCache []string = make([]string, 0)
	for _, key := range keys {
		newCache = append(newCache, key["key"].(string))
	}
	var newExpire int64 = time.Now().UTC().Unix() + *ttl
	cache = newCache
	expire = newExpire

	return nil
}

func logRequest(r *http.Request) {
	fmt.Printf("\"%s\"\t\"%s\"\t\"%s\"\t\"%s\"\t\"%s\"\t\"%s\"", r.Method, r.URL, r.Proto, r.Host, r.RemoteAddr, r.RequestURI)
}

func handle(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var err error
	var now int64 = time.Now().UTC().Unix()
	if now >= expire {
		err = fetchKeys()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if len(cache) > 0 {
		w.WriteHeader(http.StatusOK)
		for _, key := range cache {
			fmt.Fprintf(w, "%s\n", key)
		}
	} else {
		w.WriteHeader(http.StatusNoContent)
	}
}

func main() {
	flag.Parse()
	if *username == "" {
		log.Fatal("Must provide `-username` parameter. e.g. `github-keys -username \"github-username\"`")
	}

	var err error
	err = fetchKeys()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Starting server on \"%s\" for GitHub user \"%s\"\n", *listen, *username)
	http.HandleFunc("/", handle)
	err = http.ListenAndServe(*listen, nil)
	if err != nil {
		log.Fatal(err)
	}
}
