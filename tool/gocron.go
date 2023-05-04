package tool

import (
	"log"
	"time"
	"tix-id/config"
	"tix-id/models"

	"github.com/claudiu/gocron"
)

func CronTicketExpiry() {
	s := gocron.NewScheduler()

	// Connect to database
	db := config.ConnectDB()

	// Ensure the database connection is closed when the function returns
	defer db.Close()
	s.Every(10).Seconds().Do(func() {
		//get all payments
		rows, err := db.Query("SELECT * FROM payment where payment_status = 'pending'")
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
			//if status pending after 1 mins, change to expired
			t1 := created_at.Truncate(time.Second)
			t2 := time.Now().In(time.FixedZone("WIB", 7*60*60)).Add(-7 * time.Hour)
			log.Println("CRON Checking")
			if t2.Sub(t1).Minutes() > 1 {
				log.Println("FAILED Transaction ")
				log.Println(t1)
				log.Println(t2)
				log.Println(t2.Sub(t1).Minutes())
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
