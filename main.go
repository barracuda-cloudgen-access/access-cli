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
	apidevices "github.com/oNaiPs/fyde-cli/client/devices"
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

	authWriter := FydeAPIKeyAuth(signInResponse.AccessToken,
		signInResponse.Client,
		signInResponse.UID)

	// make the request to get all users
	uparams := apiusers.NewListUsersParams()
	uresp, err := client.Users.ListUsers(uparams, authWriter)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("User list:")
	for _, item := range uresp.Payload {
		log.Println("User", item.Name, "has email", item.Email, "and enrollment status", item.EnrollmentStatus)
	}

	// make the request to get all devices
	dparams := apidevices.NewListDevicesParams()
	dresp, err := client.Devices.ListDevices(dparams, authWriter)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Device list:")
	for _, item := range dresp.Payload {
		log.Println("Device", item.ID, "has status", item.Status, "and belongs to user", item.User.Name)
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
