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
	router.HandleFunc("/data/addstring", AddDataString)
	router.HandleFunc("/data/addfile", AddDataFile)
	router.HandleFunc("/data/cat", CatData)
	router.HandleFunc("/data/get", GetData)
	router.HandleFunc("/data/persistant", DataPersistant)
	router.HandleFunc("/data/traverse", TraverseTable)
	router.HandleFunc("/data/close", CloseData)

	router.HandleFunc("/image/create", CreateImage)
	router.HandleFunc("/image/addstring", AddImageString)
	router.HandleFunc("/image/addfile", AddImageFile)
	router.HandleFunc("/image/catbyidx", CatImageByIdx)
	router.HandleFunc("/image/catbycid", CatImageByCid)
	router.HandleFunc("/image/getbyidx", GetImageByIdx)
	router.HandleFunc("/image/getbycid", GetImageByCid)
	router.HandleFunc("/image/getroothash", GetMTreeRootHash)
	router.HandleFunc("/image/getheight", GetHeightByIdx)
	router.HandleFunc("/image/getimagelist", GetImageList)
	router.HandleFunc("/image/gc", GarbageCollection)
	router.HandleFunc("/image/close", CloseImage)

	// Start the HTTP server on port 8080
	http.ListenAndServe(":23333", router)
}
