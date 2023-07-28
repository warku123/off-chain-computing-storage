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
	router.HandleFunc("/image/addstring", AddImageString)
	router.HandleFunc("/image/addfile", AddImageFile)
	router.HandleFunc("/image/catbyidx", CatImageByIdx)
	router.HandleFunc("/image/cattbycid", CatImageByCid)
	router.HandleFunc("/image/getbyidx", GetImageByIdx)
	router.HandleFunc("/image/getbycid", GetImageByCid)
	router.HandleFunc("/image/getroothash", GetMTreeRootHash)
	router.HandleFunc("/image/gettimestamp", GetTimeStampByIdx)
	router.HandleFunc("/image/close", CloseImage)

	// Start the HTTP server on port 8080
	http.ListenAndServe(":8080", router)
}
