package deploygate

import "testing"

func Test_Client(t *testing.T) {
	config := &Config{
		apiKey: "api key",
	}

	err := config.initClient()

	if err != nil {
		t.Fatal(err)
	}

	if config.client == nil {
		t.Fatalf("Nil Client.")
	}
}
