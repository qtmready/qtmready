// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 
//
// This software is made available by Breu, Inc., under the terms of the 
// Breu Community License Agreement ("BCL Agreement"), version 1.0, found at  
// https://www.breu.io/license/community. By installating, downloading, 
// accessing, using or distrubting any of the software, you agree to the  
// terms of the license agreement. 
//
// The above copyright notice and the subsequent license agreement shall be 
// included in all copies or substantial portions of the software. 
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, 
// IMPLIED, STATUTORY, OR OTHERWISE, AND SPECIFICALLY DISCLAIMS ANY WARRANTY OF 
// MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE 
// SOFTWARE. 
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT 
// LIMITED TO, LOST PROFITS OR ANY CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, 
// OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, ARISING 
// OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY  
// APPLICABLE LAW. 

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
	"go.breu.io/ctrlplane/internal/shared"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")

	loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Login to ctrlplane.ai",
		Long:  "Login to ctrlplane.ai",
		Run:   loginRun,
	}
)

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
	request, _ := http.NewRequest("POST", shared.Service.CLI.BaseURL+"/auth/login", bytes.NewBuffer(marshalled))
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
