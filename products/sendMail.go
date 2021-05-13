package products

import (
	"fmt"
	"net/smtp"
)

// smtpServer data to smtp server.
type smtpServer struct {
	host string
	port string
}

// Address URI to smtp server.
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

func sendMail(message []byte) {
	// Sender data.
	from := "gogreen.golive@gmail.com"
	password := "GOliverun2"

	// Receiver email address.
	to := []string{"yenyenjoe@yahoo.com.sg", "gogreen.golive@gmail.com"}

	// smtp server configuration.
	smtpServer := smtpServer{host: "smtp.gmail.com", port: "587"}

	// // Message.
	// message = []byte("This is a really unimaginative message, I know.")

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpServer.host)

	// Sending email.
	err := smtp.SendMail(smtpServer.Address(), auth, from, to, message)
	if err != nil {
		fmt.Println("Failed sending email", err)
		return
	}

	fmt.Println("Email Sent!")
}
