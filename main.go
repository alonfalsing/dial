package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/alecthomas/kong"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	var cli struct {
		DSN      string        `help:"Database DSN" arg:"" required:"" env:"MYSQL_DSN" default:"root:root@tcp(localhost:3306)/mysql"`
		Interval time.Duration `help:"Interval" default:"10s"`
		Timeout  time.Duration `help:"Dial timeout" default:"2m"`
	}

	kong.Parse(&cli)

	db, err := sql.Open("mysql", cli.DSN)
	if err != nil {
		log.Fatal(err)
	}

	timeout := time.After(cli.Timeout)
	i := 0
	for {
		select {
		case <-timeout:
			log.Fatalf("\rFailed to connect to database within %v", cli.Timeout)

		default:
			if err := db.Ping(); err != nil {
				if i++; i >= len(spinnerFrames) {
					i = 0
				}

				fmt.Printf("\r%s", spinnerFrames[i])
				time.Sleep(cli.Interval)
				continue
			}

			return
		}
	}
}
