package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"nastenka_udalosti/internal/config"
	"nastenka_udalosti/internal/driver"
	"nastenka_udalosti/internal/handlers"
	"nastenka_udalosti/internal/helpers"
	"nastenka_udalosti/internal/models"
	"nastenka_udalosti/internal/render"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// Main je hlavní funkce aplikace
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Printf("Zapínám aplikaci na portu %s", portNumber)
	fmt.Println()

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// TODO: Co budu vkládat do relace
	gob.Register(models.Event{})
	gob.Register(models.User{})
	// gob.Register(map[string]int{})

	inProduction := flag.Bool("production", true, "Aplikace je v produkci")
	useCache := flag.Bool("cache", true, "Použit template cache")
	dbName := flag.String("dbname", "", "Název databáze")
	dbHost := flag.String("dbhost", "localhost", "Host databáze")
	dbUser := flag.String("dbuser", "", "Uživatel databáze")
	dbPass := flag.String("dbpass", "", "Heslo databáze")
	dbPort := flag.String("dbport", "", "Port databáze")
	dbSSL := flag.String("dbssl", "disable", "SSL nastavení databáze (disable, prefer, require)")

	flag.Parse()

	// fmt.Printf("DB Name: %s\n", *dbName)
	// fmt.Printf("DB User: %s\n", *dbUser)

	if *dbName == "" || *dbUser == "" {
		fmt.Println("Chybí požadovaný nastavení ")
		os.Exit(1)
	}

	// TODO: Změň na true až půjde do produkce
	app.InProduction = *inProduction
	app.UseCache = *useCache

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// Přípojení k databázi
	log.Println("Připojuji se k databázi...")

	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)
	db, err := driver.ConnectSQL(connectionString)
	if err != nil {
		log.Fatal("Nepodařilo se připojit k databázi! Končím...")
	}
	log.Println("Připojeno k databázi!")
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("nepodařilo se vytvořit template cache")
		return nil, err
	}

	app.TemplateCache = tc

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
