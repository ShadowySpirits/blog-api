package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func ForwardHandler(writer http.ResponseWriter, request *http.Request) {
	client := &http.Client{}
	path := strings.Replace(request.URL.Path, "/disqus-proxy", "", 1)
	disqusReq, err := http.NewRequest("GET", "https://disqus.com"+path+"?"+request.URL.RawQuery, nil)
	if err != nil {
		log.Printf("%v", err)
		writer.Write([]byte(fmt.Sprintf("%v", err)))
	}
	header := map[string]string{
		"Host":            "disqus.com",
		"User-Agent":      request.Header.Get("User-Agent"),
		"Accept":          request.Header.Get("Accept"),
		"Accept-Language": request.Header.Get("Accept-Language"),
		"Accept-Encoding": request.Header.Get("Accept-Encoding"),
		"Connection":      "keep-alive",
		"Cache-Control":   "max-age=0",
	}
	if _, ok := request.Header["Referer"]; ok {
		header["Referer"] = request.Header.Get("Referer")
	}
	if _, ok := request.Header["Origin"]; ok {
		header["Origin"] = request.Header.Get("Origin")
	}

	disqusReq.Header = make(http.Header)
	for k, v := range header {
		disqusReq.Header.Add(k, v)
	}

	log.Println("proxy request to: ", disqusReq.URL.String())
	log.Println("proxy request header: ", disqusReq.Header)

	resp, err := client.Do(disqusReq)
	defer resp.Body.Close()
	if err != nil {
		log.Printf("%v", err)
		writer.Write([]byte(fmt.Sprintf("%v", err)))
	}

	log.Println("proxy response header: ", resp.Header)
	for k, _ := range resp.Header {
		writer.Header().Add(k, resp.Header.Get(k))
	}
	body, _ := ioutil.ReadAll(resp.Body)
	writer.Write(body)
}
