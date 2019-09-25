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

	sayingTodayArray := GetSayingTodayFromDB()

	err := json.NewEncoder(w).Encode(sayingTodayArray)
	if err != nil {
		log.Println(err)
	}

}

func GetSayingTodayFromDB() []SayingToday {

	dbPostgres, err := pgx.Connect(loads.PgConfigLoaded)
	if err != nil {
		log.Println(err.Error())
	}

	var sayingTodayArray []SayingToday

	// To make available to pass parameter to sql query, do not recommend to save sql query in string variable.
	// Better to hard code query in db.Query("")
	rows, err := dbPostgres.Query("select id, title, author, description, reset_datetime, origin_id from sayingtoday")
	if err != nil {
		log.Println(err.Error())
	}
	defer rows.Close()

	//var lastID int

	for rows.Next() {
		var s SayingToday
		err = rows.Scan(&s.ID, &s.Title, &s.Author, &s.Description, &s.ResetDatetime, &s.OriginID)
		if err != nil {
			log.Println(err.Error())
		}

		sayingTodayArray = append(sayingTodayArray, s)
		//lastID = s.ID

	}

	//log.Println(reflect.TypeOf(n[1].FeedImageSaved.String))
	//log.Println(rssNewsArray)

	return sayingTodayArray
}
