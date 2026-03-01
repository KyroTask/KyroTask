package services

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/bif12/kyrotask/internal/config"
	"google.golang.org/api/option"
)

var FirebaseAuthClient *auth.Client

func InitFirebaseAuth() {
	ctx := context.Background()

	var app *firebase.App
	var err error

	if config.AppConfig.FirebaseCredentialsPath != "" {
		opt := option.WithCredentialsFile(config.AppConfig.FirebaseCredentialsPath)
		app, err = firebase.NewApp(ctx, nil, opt)
		if err != nil {
			log.Fatalf("Error initializing Firebase App with credentials: %v", err)
		}
		log.Println("✅ Firebase Admin SDK initialized successfully with credentials.")
	} else {
		// Try to initialize without explicit credentials (might use default application credentials)
		// Or skip initialization if this is just dev mode
		log.Println("⚠️ FIREBASE_CREDENTIALS_PATH not set. Firebase Auth will be skipped.")
		return
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("Error getting Firebase Auth client: %v", err)
	}

	FirebaseAuthClient = client
}

// VerifyIDToken verifies the Firebase ID token and returns the UID and Email.
func VerifyIDToken(idToken string) (*auth.Token, error) {
	if FirebaseAuthClient == nil {
		// Mock token behavior for development without credentials
		log.Println("WARNING: Mocking Firebase Auth token verification because Admin SDK is not initialized")
		return &auth.Token{
			UID: "mock_firebase_user_" + idToken[:min(10, len(idToken))],
			Claims: map[string]interface{}{
				"email": "mockuser@example.com",
				"name":  "Mock User",
			},
		}, nil
	}

	ctx := context.Background()
	token, err := FirebaseAuthClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}
