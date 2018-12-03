package main

import (
	"context"
	"log"

	"github.com/aneshas/tx/example"
	"github.com/aneshas/tx/example/postgres"
)

func main() {
	db, err := postgres.NewDB("testdb", "root", "root", "localhost:5432")
	check(err)

	defer db.Close()

	svc := example.NewAccountService(
		postgres.NewAccount(db.DB), // Use AccountRepository postgres implementation
	)

	{
		// Start your web server or cli application here and use
		// svc (Account application service), eg. from http handler:

		req := example.TransferReq{
			SrcID:  456,
			DestID: 123,
			Amount: 3000,
		}

		err = svc.TransferMoney(context.Background(), &req)
		if err != nil {
			log.Fatalf("could not transfer money: %v", err)
		}

		log.Println("Money transfered successfully")
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
