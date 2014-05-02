package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"appengine"
	"appengine/memcache"
	"appengine/urlfetch"
)

func init() {
	http.HandleFunc("/kml", kmlHandler)
}

func kmlHandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	url := r.FormValue("url")

	if item, err := memcache.Get(c, url); err == memcache.ErrCacheMiss || err != nil {
		c.Infof("cache miss")

		client := urlfetch.Client(c)
		resp, err := client.Get(url)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%s", string(contents))
		item := &memcache.Item{
			Key:   url,
			Value: contents,
		}
		if err := memcache.Set(c, item); err != nil {
			c.Errorf("error setting item: %v", err)
		}
	} else {
		fmt.Fprintf(w, "%s", string(item.Value))
	}
}
