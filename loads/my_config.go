package loads

import (
	"encoding/json"
	"github.com/jackc/pgx"
	"log"
	"net/http"
	"os"
)

const ElDevelopment = true

var PgConfigLoaded pgx.ConnConfig

func LoadConfiguration() {

	file := "file_path_here"
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&PgConfigLoaded)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("config loaded")
	}
}

func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
