package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Jalenarms1/caters-go/internal/db"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}

	if err := db.SetDb(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("DB Connected")
}

func main() {

	fmt.Println("Hello there")

	addr := os.Getenv("LISTEN_ADDR")

	server := NewServer(addr)

	server.run()
}
