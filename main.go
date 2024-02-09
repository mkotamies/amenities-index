package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
)

type Hotel struct {
	Id    string
	Name string
}

func getHotels(conn *pgx.Conn) map[string][]Hotel{
	var hotels = make(map[string][]Hotel)
	rows, _ := conn.Query(context.Background(), "select * from hotel")

	for rows.Next() {
	var id int32
	var name string
	
  	err := rows.Scan(&id, &name)
	if err != nil {
		fmt.Println(err)
	}
	hotels["Hotels"] = append(hotels["Hotels"], Hotel{Id: strconv.FormatInt(int64(id), 10), Name: name} )
	}
	return hotels
}

func main() {
	fmt.Println("Go app...")

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())



	var hotels map[string][]Hotel

	h1 := func(w http.ResponseWriter, r *http.Request) {
		hotels = getHotels(conn) 
		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, hotels)
	}


	// define handlers
	http.HandleFunc("/", h1)

	log.Fatal(http.ListenAndServe("localhost:9000", nil))

}
