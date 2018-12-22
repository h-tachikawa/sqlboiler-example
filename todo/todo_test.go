package todo

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"main/models"
	"os"
	"testing"
	"time"
)

var app App

func TestMain(m *testing.M) {
	mySqlInfo := fmt.Sprintf("%s:%s@tcp/%s?charset=utf8&parseTime=True&loc=Local",
		"testuser", "testpass", "testdb",
	)
	db, err := sql.Open("mysql", mySqlInfo)
	boil.DebugMode = true
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	app = App{
		db: db,
	}

	os.Exit(m.Run())
}

func TestFetchTaskAndUser(t *testing.T) {
	user, err := app.FetchTaskAndUser()
	if err != nil {
		t.Fatal(err)
	}
	jsonBytes, err := json.Marshal(user)
	if err != nil {
		t.Fatal("JSON Marshall error:", err)
	}
	out := new(bytes.Buffer)
	_ = json.Indent(out, jsonBytes, "", "    ")
	t.Log(out.String())
}

func TestFetchTasksAndUser(t *testing.T) {
	tau, err := app.FetchTasksAndUser()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range tau {
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			t.Fatal("JSON Marshall error:", err)
		}
		out := new(bytes.Buffer)
		_ = json.Indent(out, jsonBytes, "", "    ")
		t.Log(out.String())
	}
}

func TestFetchUnfinished(t *testing.T) {
	todos, err := app.FetchUnfinished()
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range todos {
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			t.Fatal("JSON Marshall error:", err)
		}
		out := new(bytes.Buffer)
		_ = json.Indent(out, jsonBytes, "", "    ")
		t.Log(out.String())
	}
}

func TestStore(t *testing.T) {
	now := time.Now().UTC()
	todos := []*models.Todo{
		&models.Todo{
			ID:       1,
			Title:    "Sample ToDo 1",
			DueDate:  now,
			Note:     null.StringFrom("note1..."),
			Finished: false,
		},
		&models.Todo{
			ID:       2,
			Title:    "Sample ToDo 2",
			DueDate:  now,
			Note:     null.StringFrom("note2..."),
			Finished: false,
		},
		&models.Todo{
			ID:       3,
			Title:    "Sample ToDo 3",
			Note:     null.StringFrom("note3..."),
			Finished: false,
		},
	}

	if err := app.Store(todos); err != nil {
		t.Fatal(err)
	}

	for _, v := range todos {
		t.Log(fmt.Sprintf("%+v", v))
	}
}

func TestFinish(t *testing.T) {
	if err := app.Finish([]int64{1, 2}); err != nil {
		t.Fatal(err)
	}
}
