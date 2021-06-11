package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tapvanvn/go-wspubsub/entity"
	"github.com/tapvanvn/go-wspubsub/runtime"
	"github.com/tapvanvn/go-wspubsub/server"
	"github.com/tapvanvn/go-wspubsub/utility"
)

func main() {
	var port = utility.MustGetEnv("PORT")

	rootPath, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {

		panic(err)
	}
	runtime.RootPath = rootPath

	//MARK: init system config
	jsonFile2, err := os.Open(rootPath + "/config/config.json")

	if err != nil {
		panic(err)
	}

	defer jsonFile2.Close()
	configData, err := ioutil.ReadAll(jsonFile2)
	systemConfig := entity.Config{}

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(configData, &systemConfig)
	if err != nil {
		panic(err)
	}
	runtime.Config = &systemConfig

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		server.ServeWs(w, r)
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("i am ok"))
	})

	fmt.Println("listen on port", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
