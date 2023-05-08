package connection

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
)

var Conn *pgx.Conn

func DatabaseConnect() {

	databaseUrl := "postgresql://postgres:P@sswo0rd@localhost:5432/db_blog"

	var err error
	Conn, err = pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Success connect to database 😎👌")

}
