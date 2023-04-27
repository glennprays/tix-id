package main

import (
	"log"
	"time"
	"tix-id/config"
	"tix-id/docs"
	"tix-id/models"
	"tix-id/routes"

	"github.com/claudiu/gocron"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	docs.SwaggerInfo.Title = "TIX-ID API Documentation"
	docs.SwaggerInfo.Description = "Ticketing Application"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	go cronTicketExpiry()

	r := routes.SetupRouter()
	r.Run(":8080")

}

func cronTicketExpiry() {
	s := gocron.NewScheduler()

	// Connect to database
	db := config.ConnectDB()

	// Ensure the database connection is closed when the function returns
	defer db.Close()
	s.Every(10).Seconds().Do(func() {
		//get all payments
		rows, err := db.Query("SELECT * FROM payment")
		if err != nil {
			log.Println(err)
			return
		}

		//loop through all payments
		for rows.Next() {
			var payment models.Payment
			var created_at time.Time
			//get payment data
			err := rows.Scan(&payment.ID, &payment.Amount, &payment.Status, &created_at)
			if err != nil {
				log.Println(err)
				return
			}
			//if status pending after 5 mins, change to expired
			t1 := created_at
			t2 := time.Now()
			if (payment.Status == "pending") && (t2.Sub(t1).Minutes() >= 1) {
				_, err := db.Exec("UPDATE payment SET payment_status = 'failed' WHERE id = ?", payment.ID)
				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	})
	<-s.Start()
}
