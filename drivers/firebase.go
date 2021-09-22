package drivers

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

type (
	AZFirebaseApp struct {
		app  *firebase.App
		auth *auth.Client
	}

	FirebaseUpdateOpts struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		PhoneNumber   string `json:"phone_number"`
		Password      string `json:"password"`
		DisplayName   string `json:"display_name"`
		PhotoURL      string `json:"photo_url"`
		Disabled      bool   `json:"disabled"`
	}
)

const (
// ProjectID = ""
)

var ()

func NewFirebaseApp(credentailsPath string) (*AZFirebaseApp, error) {
	azfp := new(AZFirebaseApp)
	err := azfp.init(credentailsPath)
	if err != nil {
		return nil, err
	}
	log.Println("firebase driver init")
	return azfp, nil
}

func (azfb *AZFirebaseApp) init(credentailsPath string) error {
	opt := option.WithCredentialsFile(credentailsPath)
	ctx := context.Background()
	var err error
	azfb.app, err = firebase.NewApp(ctx, nil, opt)
	if err != nil {
		// log.Fatal(err)
		log.Println("AZFirebaseApp init error:", err)
		return err
	}
	return azfb.authorize()
}

func (azfb *AZFirebaseApp) authorize() error {
	var err error
	ctx := context.Background()
	azfb.auth, err = azfb.app.Auth(ctx)
	if err != nil {
		log.Println("AZFirebaseApp authorize error:", err)
		return err
	}
	return nil
}

func (azfb *AZFirebaseApp) GetUser(email string) (*auth.UserRecord, error) {
	ctx := context.Background()
	return azfb.auth.GetUserByEmail(ctx, email)
}

func (azfb *AZFirebaseApp) InitCustomToken(uid string, customData map[string]interface{}) (string, error) {
	ctx := context.Background()
	return azfb.auth.CustomTokenWithClaims(ctx, uid, customData)
}

func (azfb *AZFirebaseApp) InitNewUser(email, pwd, displayName string) (*auth.UserRecord, error) {
	ctx := context.Background()
	userParams := (&auth.UserToCreate{}).Email(email).Password(pwd).DisplayName(displayName).Disabled(false)
	return azfb.auth.CreateUser(ctx, userParams)
}

func (azfb *AZFirebaseApp) UpdateUser(uid string, updateData FirebaseUpdateOpts) (*auth.UserRecord, error) {
	ctx := context.Background()
	userUpdateRecord := &auth.UserToUpdate{}

	if updateData.Email != "" {
		userUpdateRecord = userUpdateRecord.Email(updateData.Email)
	}
	if updateData.Password != "" {
		userUpdateRecord = userUpdateRecord.Password(updateData.Password)
	}
	if updateData.DisplayName != "" {
		userUpdateRecord = userUpdateRecord.DisplayName(updateData.DisplayName)
	}
	if updateData.PhoneNumber != "" {
		userUpdateRecord = userUpdateRecord.PhoneNumber(updateData.PhoneNumber)
	}
	if updateData.PhotoURL != "" {
		userUpdateRecord = userUpdateRecord.PhotoURL(updateData.PhotoURL)
	}
	return azfb.auth.UpdateUser(ctx, uid, userUpdateRecord)
}

func (azfb *AZFirebaseApp) DeleteUser(uid string) error {
	ctx := context.Background()
	return azfb.auth.DeleteUser(ctx, uid)
}

func (azfb *AZFirebaseApp) VerifyClientToken(tokenString string, checkRevocation bool) (*auth.Token, error) {
	ctx := context.Background()

	if checkRevocation {
		return azfb.auth.VerifyIDTokenAndCheckRevoked(ctx, tokenString)
	}
	return azfb.auth.VerifyIDToken(ctx, tokenString)
}

func (azfb *AZFirebaseApp) RevokeRefreshTokens(uid string) error {
	ctx := context.Background()
	return azfb.auth.RevokeRefreshTokens(ctx, uid)
}
