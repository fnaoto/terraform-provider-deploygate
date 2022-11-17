package deploygate

import (
	"testing"

	go_deploygate "github.com/fnaoto/go_deploygate"
)

func Test_Client(t *testing.T) {
	config := &Config{
		clientConfig: go_deploygate.ClientConfig{
			ApiKey: "api_key",
		},
	}

	err := config.initClient()

	if err != nil {
		t.Fatal(err)
	}

	if config.client == nil {
		t.Fatalf("Nil Client.")
	}
}
