// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 

package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"go.breu.io/ctrlplane/internal/api/auth"
)

const BaseUrl = "http://localhost:8000" // TODO: use a better way to do this
var ErrInvalidCredentials = errors.New("invalid credentials")

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to ctrlplane.ai",
	Long:  "Login to ctrlplane.ai",
	Run:   loginRun,
}

func init() {
	rootCmd.AddCommand(loginCmd)
}

func loginRun(cmd *cobra.Command, args []string) {
PROMPT:
	email, err := promptEmail()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	password, err := promptPassword()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	token, err := doLogin(email, password)
	if err != nil {
		fmt.Println(err)
		goto PROMPT
	}

	fmt.Println("Access Token: ", token.AccessToken)
	fmt.Println("Refresh Token: ", token.RefreshToken)
}

func promptEmail() (string, error) {
	prompt := promptui.Prompt{
		Label: "Please enter your email address registered with ctrlplane.ai",
		Validate: func(input string) error {
			if len(input) < 3 {
				return errors.New("invalid email")
			}
			return nil
		},
	}
	return prompt.Run()
}

func promptPassword() (string, error) {
	prompt := promptui.Prompt{
		Label: "Please enter your password",
		Mask:  '*',
	}
	return prompt.Run()
}

func doLogin(email string, password string) (*auth.TokenResponse, error) {
	data := auth.LoginRequest{Email: email, Password: password}
	token := &auth.TokenResponse{}

	marshalled, _ := json.Marshal(data)
	request, _ := http.NewRequest("POST", BaseUrl+"/auth/login", bytes.NewBuffer(marshalled))
	request.Header.Set("User-Agent", "ctrlplane-cli/0.0.1")
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return token, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return token, ErrInvalidCredentials
	}

	fmt.Println("Login successful")
	if err := json.NewDecoder(response.Body).Decode(token); err != nil {
		return token, err
	}

	return token, nil
}
