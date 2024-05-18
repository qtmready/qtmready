package main

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log/slog"
	"net/http"
	"os/exec"

	ghi "github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/uuid"

	// gh "github.com/google/go-github/v62/github"

	"go.breu.io/quantm/internal/providers/github"
	"go.breu.io/quantm/internal/shared"
)

func main() {
	block, _ := pem.Decode([]byte(github.Instance().PrivateKey))
	key, _ := x509.ParsePKCS1PrivateKey(block.Bytes)

	shared.Logger().Info("decoded ....", slog.Any("block", block), slog.Any("key", key))

	client, _ := ghi.New(http.DefaultTransport, github.Instance().AppID, 50886707, []byte(github.Instance().PrivateKey))
	token, _ := client.Token(context.Background())

	shared.Logger().Info("token ....", slog.Any("token", token))

	url := fmt.Sprintf("https://git:%s@github.com/breuHQ/governance.git", token)

	id, _ := uuid.NewV7()

	cmd := exec.Command("git", "clone", url, fmt.Sprintf("/tmp/%s", id.String()))
	err := cmd.Run()

	if err != nil {
		shared.Logger().Error("error cloning repo", slog.Any("error", err.Error()))
	}
}

// git clone https://git:<token>@github.com/owner/repo.git
