package api

import (
	"database/sql"
	"elnewsAPI/loads"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx"
	"log"
	"net/http"
	"strconv"
)

type PostsUnion struct {

	Title 		string
	Link		string
	CreateDate	string			`db:"create_date" json:"create_date"`
	ModifyDate	sql.NullString	`db:"modify_date" json:"modify_date"`
	Author		sql.NullString
	ShortParagraph	sql.NullString	`db:"short_paragraph" json:"short_paragraph"`
	Category	sql.NullString
	FeedImage	sql.NullString	`db:"feed_image" json:"feed_image"`
	OriginalLink	sql.NullString	`db:"original_link" json:"original_link"`
	FeedImageSaved	sql.NullString	`db:"feed_image_saved" json:"feed_image_saved"`
	Factor		int64			`json:"-"`
	TableType		sql.NullString	`db:"table_type" json:"table_type"`

}

type PostsUnionSet struct {

	Posts	[]PostsUnion	`json:"posts"`
	LastID	int64			`json:"lastID"`

}


func PostsUnionJSON (w http.ResponseWriter, r *http.Request) {

	//loads.EnableCors(&w)

	w.Header().Set("Content-Type", "application/json")

	lastIDFrom := r.FormValue("last_id")

	var postsUnionArray []PostsUnion
	var lastID int64

	if lastIDFrom == "" {
		postsUnionArray, lastID = GetFirstPostsUnionFromDB(15)

	} else {
		l1, err := strconv.Atoi(lastIDFrom)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Can't find enough data from request"))
			return
		} else {
			postsUnionArray, lastID, err = GetNextPostsUnionFromDB(15, int64(l1))
			if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("400 - Can't find enough data from request"))
			return
		}
	}

	postsUnionSet := &PostsUnionSet{
		Posts: postsUnionArray,
		LastID: lastID,
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(postsUnionSet)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Something bad happened!"))
		return
	}
	return

}


func GetFirstPostsUnionFromDB(limitNum int) ([]PostsUnion, int64, error) {

	dbPostgres, err := pgx.Connect(loads.PgConfigLoaded)
	if err != nil {
		return nil, nil, err
	}
	defer dbPostgres.Close()

	var postsUnionArray []PostsUnion

	sqlQuery := "SELECT title, link, create_date, modify_date, short_paragraph, category, " +
		"feed_image, original_link, feed_image_saved, factor, table_type from cheditor " +
		"UNION ALL SELECT title, link, create_date, modify_date, short_paragraph, category, " +
		"feed_image, original_link, feed_image_saved, factor, table_type from rssnews " +
		"order by factor desc limit $1"

	rows, err := dbPostgres.Query(sqlQuery, limitNum)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var lastID int64

	for rows.Next() {
		var p PostsUnion
		err = rows.Scan(&p.Title, &p.Link, &p.CreateDate, &p.ModifyDate, &p.ShortParagraph, &p.Category, &p.FeedImage,
			&p.OriginalLink, &p.FeedImageSaved, &p.Factor, &p.TableType)
		if err != nil {
			return nil, nil, err
		}

		postsUnionArray = append(postsUnionArray, p)
		lastID = p.Factor

	}

	return postsUnionArray, lastID, err

}

func GetNextPostsUnionFromDB(limitNum int, lastID int64) ([]PostsUnion, int64, error){

	dbPostgres, err := pgx.Connect(loads.PgConfigLoaded)
	if err != nil {
		return nil, nil, err
	}
	defer dbPostgres.Close()

	var postsUnionArray []PostsUnion

	sqlQuery := "SELECT title, link, create_date, modify_date, short_paragraph, category, " +
		"feed_image, original_link, feed_image_saved, factor, table_type from cheditor " +
		"WHERE factor < $1 " +
		"UNION ALL SELECT title, link, create_date, modify_date, short_paragraph, category, " +
		"feed_image, original_link, feed_image_saved, factor, table_type from rssnews " +
		"WHERE factor < $2 ORDER BY factor DESC LIMIT $3"

	rows, err := dbPostgres.Query(sqlQuery, lastID, lastID, limitNum)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p PostsUnion
		err = rows.Scan(&p.Title, &p.Link, &p.CreateDate, &p.ModifyDate, &p.ShortParagraph, &p.Category, &p.FeedImage,
			&p.OriginalLink, &p.FeedImageSaved, &p.Factor, &p.TableType)
		if err != nil {
			return nil, nil, err
		}

		postsUnionArray = append(postsUnionArray, p)
		lastID = p.Factor

	}

	return postsUnionArray, lastID, err
}
