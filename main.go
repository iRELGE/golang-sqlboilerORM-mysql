package main

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/thanos-io/thanos/pkg/runutil"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"rabie.com/sqlb/models"
)

func main() {
	//open db (regular sql open call)
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/golangportfolio")
	dieIf(err)
	//close deferred
	defer runutil.CloseWithErrCapture(&err, db, "close db")
	//check if db is connected
	err = db.Ping()
	fmt.Println("connected")
	//create a normal test structure you can choose any attribute you want to insert and negligee others if it's not required
	u := &models.Test{Name: "john", LastName: "go"}
	//insert test to db caling Insert ,infer the column that should be insert if is a primary key or has default value
	err = u.Insert(context.Background(), db, boil.Infer())
	dieIf(err)
	//print the test that has been insert
	fmt.Println("user id:", u.UserID)
	//get one row back by calling one it return the first one
	got, err := models.Tests().One(context.Background(), db)
	dieIf(err)

	println("got user:", got.Name)
	//find a particular test
	found, err := models.FindTest(context.Background(), db, got.UserID)
	dieIf(err)

	fmt.Println("found user:", found.UserID)
	found.UserID = 31
	//update a testand boil.Whitelist() tell us whish collumn we gonna update
	rows, err := found.Update(context.Background(), db, boil.Whitelist(models.TestColumns.Name))
	dieIf(err)

	fmt.Println("updeted row:", rows)
	//find user by id where id equal an id supposed
	foundAgain, err := models.Tests(qm.Where("user_id = ?", found.UserID)).One(context.Background(), db)
	dieIf(err)

	fmt.Println("user found:", foundAgain.UserID, foundAgain.Name)
	//check if the user exist or not TestExists will return a boolean value
	exists, err := models.TestExists(context.Background(), db, foundAgain.UserID)
	dieIf(err)

	fmt.Println("user:", foundAgain.Name, "user exists:", exists)
	//it count how much row in table it can be personalized how mush row has specified things
	count, err := models.Tests(qm.Where("name = ?", "john")).Count(context.Background(), db)
	dieIf(err)

	println("number of rows:", count)
}
func dieIf(err error) {
	if err != nil {
		panic(err)
	}
}
