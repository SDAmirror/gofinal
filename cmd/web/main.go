package main
import (
	"snippetbox/pkg/models/postgres"
	"context"
	"flag"
	"github.com/jackc/pgx/v4/pgxpool"

	"log"
	"net/http"
	"os"
)

type application struct {
	errorLog *log.Logger
	infoLog *log.Logger
	snippets *postgres.SnippetModel
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	pool, err := pgxpool.Connect(context.Background(), "user=postgres password=123456 host=localhost port=5432 dbname=snippetbox sslmode=disable pool_max_conns=10")
	if err != nil {
		log.Fatalf("Unable to connection to database: %v\n", err)
	}

	defer pool.Close()


	app := &application{errorLog,infoLog,&postgres.SnippetModel{pool}}

	srv := &http.Server{
		Addr: *addr,
		ErrorLog: errorLog,
		Handler: app.routes(),
	}
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}