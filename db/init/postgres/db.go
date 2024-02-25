package postgres

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "pinterest"
)

func Connect() (*sql.DB, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	return sql.Open("postgres", psqlconn)
}

func Disconnect(db *sql.DB) {
	db.Close()
}

func FatalIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func RegisterUser(email, nickname, password string) {
	sqlInsertRequest := `INSERT INTO public.Users ("email", "nickname", "password") VALUES($1, $2, $3)`
	db, errConn := Connect()
	FatalIfError(errConn)
	defer Disconnect(db)
	_, errWrite := db.Exec(sqlInsertRequest, email, nickname, password)
	FatalIfError(errWrite)
}

func GetUserByEmail(email string) *User {
	sqlSelectQuery := `SELECT user_id, email, nickname, "password" FROM public.Users WHERE email=$1`
	db, errConn := Connect()
	FatalIfError(errConn)
	defer Disconnect(db)
	var user = &User{}
	db.QueryRow(sqlSelectQuery, email).Scan(&user.User_id, &user.Email, &user.Nickname, &user.Password)
	return user
}
