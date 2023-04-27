package controller

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"tix-id/models"

	"gopkg.in/gomail.v2"
)

func SendEmail(content string, receiverMail string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "No Reply <no-reply@example.com>")
	m.SetHeader("To", receiverMail)
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", content)
	email := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")
	fmt.Println(email, " ", password)
	d := gomail.NewDialer("smtp.gmail.com", 465, "if-21056@students.ithb.ac.id", "factcheck")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	log.Println("in here")
	log.Println(content)
}

func GenerateEmail(customer models.Customer, payment models.Payment, scheduleTicket models.ScheduleTicket) string {

	content := `<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<style>
			body {
				font-family: Arial, sans-serif;
				font-size: 14px;
				line-height: 1.5;
				color: #333;
				margin: 0;
				padding: 0;
			}
	
			.email-container {
				max-width: 600px;
				margin: 0 auto;
				padding: 20px;
				background-color: #f8f8f8;
				border: 1px solid #ddd;
				border-radius: 5px;
			}
	
			.email-header {
				text-align: center;
				margin-bottom: 30px;
			}
	
			.email-header img {
				max-width: 150px;
				height: auto;
			}
	
			.email-content {
				padding: 20px;
				background-color: #ffffff;
				border-radius: 5px;
			}
	
			.email-content h1 {
				font-size: 24px;
				margin-bottom: 20px;
			}
	
			.email-content p {
				margin-bottom: 15px;
			}
	
			.email-content ul {
				padding-left: 20px;
				margin-bottom: 15px;
			}
	
			.email-footer {
				text-align: center;
				margin-top: 30px;
			}
	
			.email-footer p {
				font-size: 12px;
				color: #777;
			}
		</style>
	</head>
	<body>
		<div class="email-container">
			<div class="email-header">
				<img src="your-logo.png" alt="Your App Logo">
			</div>
			<div class="email-content">`
	content += `<h1>TIX-ID</h1>
			<p>Hi, ` + customer.Name + `,</p>
			<p>Thank you for using TIX-ID, we hope you enjoyed our service. </p>
			`
	content += `<ul><li><strong>Amount Paid: </strong> ` + strconv.Itoa(int(payment.Amount)) + `</li><br>
				<strong>--------------------ORDER DETAILS--------------------</strong> <br>
				<li><strong>Ticket ID:</strong> ` + strconv.Itoa(scheduleTicket.ID) + `</li>
				<li><strong>` + scheduleTicket.Movie.Title + `</li>
				<li><strong>` + scheduleTicket.Branch.Name + `</li>
				<strong>-------------
				<li><strong>SHOWTIME   ` + scheduleTicket.Showtime.String() + `</li>
				<li><strong>SEAT       ` + scheduleTicket.Seat.Row + scheduleTicket.Seat.Number + `</li>
				<li><strong>COST       ` + strconv.Itoa(int(payment.Amount)) + `</li>
				<li><strong>Paid with X-Pay</li></ul>
						<p>Sold by PT BadBums Programmers (NPWP: 02.331.777-9.054.000 - Address: Gedung AIA Central Lt. 26, Jl. Jend. Sudirman Kav. 48A, Setiabudi, Jakarta Selatan, DKI Jakarta. Admin Fee and Discount, if any, are provided by ANONYMOUS.</p>
						<div class="email-footer">
							<p>Need help? Contact at: whatever@gpt.com</p>
						</div>
					</div>
				</body>
				</html>`

	return content
}
