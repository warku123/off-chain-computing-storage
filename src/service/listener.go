package service

import (
	"fmt"
	"net/http"
)

func Listener() {
	fmt.Println("Listener start")

	// Create a new HTTP router
	router := http.NewServeMux()

	router.HandleFunc("/data/create", CreateData)
	router.HandleFunc("/data/add", AddData)
	router.HandleFunc("/data/get", GetData)
	router.HandleFunc("/data/close", CloseData)

	router.HandleFunc("/image/create", CreateImage)
	router.HandleFunc("/image/add", AddImage)
	router.HandleFunc("/image/getbyidx", GetImageByIdx)
	router.HandleFunc("/image/getbycid", GetImageByCid)
	router.HandleFunc("/image/getroothash", GetMTreeRootHash)
	router.HandleFunc("/image/close", CloseImage)

	// Start the HTTP server on port 8080
	http.ListenAndServe(":8080", router)
}
