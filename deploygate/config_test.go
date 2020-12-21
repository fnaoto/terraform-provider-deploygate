package deploygate

import "testing"

func Test_Client(t *testing.T) {
	config := &Config{
		apiKey: "api key",
	}

	client, err := config.Client()

	if err != nil {
		t.Fatal(err)
	}

	if client == nil {
		t.Fatalf("Nil Client.")
	}
}
