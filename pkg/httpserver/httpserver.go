package httpserver

import (
	"database/sql"
	"net/http"
)

func Start() error {
	// db...
	// server...
	serv := newServer()


	if err := http.ListenAndServe(":5252", serv); err != nil {
		return err
	}
	return nil
}

func newDb(dbUrl string) (*sql.DB, error) {
	// ...
	return nil, nil
}