package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	apiclient "github.com/oNaiPs/fyde-cli/client"
	apiauth "github.com/oNaiPs/fyde-cli/client/auth"
	apiusers "github.com/oNaiPs/fyde-cli/client/users"
	"github.com/oNaiPs/fyde-cli/models"
)

func main() {
	fmt.Println("fyde-cli stub. Hello world!")

	host := os.Getenv("FYDECLI_HOST")
	if host == "" {
		log.Fatalln("Missing FYDECLI_HOST environment variable")
	}

	username := os.Getenv("FYDECLI_USER")
	if username == "" {
		log.Fatalln("Missing FYDECLI_USER environment variable")
	}

	password := os.Getenv("FYDECLI_PASS")
	if password == "" {
		log.Fatalln("Missing FYDECLI_PASS environment variable")
	}

	transport := httptransport.New(os.Getenv("FYDECLI_HOST"), "/api/v1", nil)

	client := apiclient.New(transport, strfmt.Default)

	params := apiauth.NewSignInParams()
	params.WithBody(&models.SignInRequest{
		Email:    username,
		Password: password,
	})
	signInResponse, err := client.Auth.SignIn(params)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Logged in as", signInResponse.Payload.Data.Name)

	uparams := apiusers.NewListUsersParams()

	authWriter := FydeAPIKeyAuth(signInResponse.AccessToken,
		signInResponse.Client,
		signInResponse.UID)

	// make the request to get all users
	resp, err := client.Users.ListUsers(uparams, authWriter)
	if err != nil {
		log.Fatal(err)
	}

	for _, item := range resp.Payload {
		log.Println("User", item.Name, "has email", item.Email, "and enrollment status", item.EnrollmentStatus)
	}
}

// FydeAPIKeyAuth provides an API key auth info writer
func FydeAPIKeyAuth(accessToken, client, uid string) runtime.ClientAuthInfoWriter {
	return runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		err := r.SetHeaderParam("access-token", accessToken)
		if err != nil {
			return err
		}

		err = r.SetHeaderParam("client", client)
		if err != nil {
			return err
		}

		return r.SetHeaderParam("uid", uid)
	})
}
