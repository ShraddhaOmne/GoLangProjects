package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Nerzal/gocloak/v13"
)

func main() {
	// Replace with your Keycloak details
	keycloakHost := "http://localhost:18080" // Keycloak base URL
	realm := "MFRealm"
	clientID := "mfclient"
	clientSecret := "myclient-secret" // only if confidential client
	username := "sam"
	password := "aaa111"

	client := gocloak.NewClient(keycloakHost)

	// Use Go's context
	ctx := context.Background()

	// Login with username/password
	token, err := client.Login(ctx, clientID, clientSecret, realm, username, password)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	fmt.Println("Access Token:", token.AccessToken)
	fmt.Println("Refresh Token:", token.RefreshToken)

	// Verify the token
	rptResult, err := client.RetrospectToken(ctx, token.AccessToken, clientID, clientSecret, realm)
	if err != nil {
		log.Fatalf("Token introspection failed: %v", err)
	}
	if !*rptResult.Active {
		log.Fatal("Token is not active ❌")
	}
	fmt.Println("Token is valid ✅")
}
