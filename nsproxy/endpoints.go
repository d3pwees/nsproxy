package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unixvoid/nsproxy/pkg/nslog"
)

func asyncClusterListener() {
	// async listener gets its own redis connection
	redisClient, _ := initRedisConnection()

	// first boot, remove live file and state file
	redisClient.Del(fmt.Sprintf("index:live"))

	// format the string to be :port
	port := fmt.Sprint(":", config.Clustermanager.Port)
	nslog.Info.Println("started async cluster listener on port", config.Clustermanager.Port)

	router := mux.NewRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clusterHandler(w, r, redisClient)
	}).Methods("POST")
	router.HandleFunc("/dns", func(w http.ResponseWriter, r *http.Request) {
		dnsHandler(w, r, redisClient)
	}).Methods("POST")
	router.HandleFunc("/dns", func(w http.ResponseWriter, r *http.Request) {
		dnsHostsHandler(w, r, redisClient)
	}).Methods("GET")
	router.HandleFunc("/dns/rm", func(w http.ResponseWriter, r *http.Request) {
		dnsRmHandler(w, r, redisClient)
	}).Methods("POST")
	router.HandleFunc("/clusterspec", func(w http.ResponseWriter, r *http.Request) {
		apiClusterSpecHandler(w, r, redisClient)
	}).Methods("POST")
	router.HandleFunc("/hostspec", func(w http.ResponseWriter, r *http.Request) {
		apiHostSpecHandler(w, r, redisClient)
	}).Methods("POST")
	router.HandleFunc("/dnsspec", func(w http.ResponseWriter, r *http.Request) {
		apiDnsSpecHandler(w, r, redisClient)
	}).Methods("POST")
	router.HandleFunc("/hosts", func(w http.ResponseWriter, r *http.Request) {
		apiHostsHandler(w, r, redisClient)
	}).Methods("GET")
	router.HandleFunc("/clusters", func(w http.ResponseWriter, r *http.Request) {
		apiClustersHandler(w, r, redisClient)
	}).Methods("GET")

	// check index:master against index:live
	clusterDiff(redisClient)
	log.Fatal(http.ListenAndServe(port, router))
}
