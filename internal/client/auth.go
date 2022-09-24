// Copyright © 2022, Breu Inc. <info@breu.io>. All rights reserved. 

package client

import (
	"go.breu.io/ctrlplane/internal/api/auth"
)

func (client *Client) Login(email, password string) (*auth.TokenResponse, error) {
	// data := auth.LoginRequest{Email: email, Password: password}
	// token := &auth.TokenResponse{}

	// marshalled, _ := json.Marshal(data)
	// request, _ := http.NewRequest("POST", BaseUrl+"/auth/login", bytes.NewBuffer(marshalled))
	// request.Header.Set("User-Agent", "ctrlplane-cli/0.0.1")
	// request.Header.Set("Content-Type", "application/json")

	// client := &http.Client{}
	// response, err := client.Do(request)

	// if err != nil {
	// 	return token, err
	// }

	// defer response.Body.Close()

	// if response.StatusCode != 200 {
	// 	return token, errors.New("invalid credentials")
	// }

	// body, _ := ioutil.ReadAll(response.Body)

	// err = json.Unmarshal(body, &token)

	// if err != nil {
	// 	return token, err
	// }

	// return token, nil
	return nil, nil
}
