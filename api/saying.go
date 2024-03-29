package api

import (
	"database/sql"
	"elnewsAPI/loads"
	"encoding/json"
	"github.com/jackc/pgx"
	"log"
	"net/http"
)

type SayingToday struct {

	ID int
	Title string
	Author sql.NullString
	Description sql.NullString
	ResetDatetime string `db:"reset_datetime" json:"reset_datetime"`
	OriginID int `db:"origin_id" json:"origin_id"`

}


func SayingTodayJSONArray(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	sayingTodayArray, err := GetSayingTodayFromDB()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(sayingTodayArray)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		return
	}
	return
}

func GetSayingTodayFromDB() ([]SayingToday, error) {

	dbPostgres, err := pgx.Connect(loads.PgConfigLoaded)
	if err != nil {
		return nil, err
	}

	var sayingTodayArray []SayingToday

	// To make available parameter to pass to sql query, not recommend to save sql query in string variable.
	// Better to hard code query in db.Query("")
	rows, err := dbPostgres.Query("select id, title, author, description, reset_datetime, origin_id from sayingtoday")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s SayingToday
		err = rows.Scan(&s.ID, &s.Title, &s.Author, &s.Description, &s.ResetDatetime, &s.OriginID)
		if err != nil {
			return nil, err
		}

		sayingTodayArray = append(sayingTodayArray, s)
		//lastID = s.ID

	}

	return sayingTodayArray, err
}
