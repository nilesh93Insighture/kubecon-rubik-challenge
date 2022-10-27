package services

import (
	"context"
	"fmt"
	"os"

	"github.com/mailgun/mailgun-go/v4"
)

func SendEmail(ctx context.Context, email string, challengeWon bool) (string, error) {
	mg := mailgun.NewMailgun(os.Getenv("domain"), os.Getenv("apiKey"))

	m := mg.NewMessage(
		fmt.Sprintf("Choreo Rubik Challenge <postmaster@%s>", os.Getenv("domain")),
		"Choreo Rubik Challenge",
		"",
		email,
	)

	if challengeWon {
		m.SetTemplate("rubik-challenge-won")
	} else {
		m.SetTemplate("rubik-challenge-lost")
	}

	_, id, err := mg.Send(ctx, m)

	return id, err
}
