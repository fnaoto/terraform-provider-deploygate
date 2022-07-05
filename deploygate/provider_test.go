package deploygate

import (
	"context"
	"sync"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	ProviderNameDG = "deploygate"
)

var testDGProvider *schema.Provider

var testDGProviderConfigure sync.Once

var testDGProviders map[string]*schema.Provider

func Test_DGPreCheck(t *testing.T) {
	initProvider("")
	// Test for config
	configRaw := map[string]interface{}{
		"api_key": "dummy",
	}
	testDGProviderConfigure.Do(func() {
		err := testDGProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(configRaw))
		if err != nil {
			t.Fatal(err)
		}
	})

	// Test for environment variables
	testDGProviderConfigure.Do(func() {
		t.Setenv("DG_API_KEY", "dummy")
		err := testDGProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil))
		if err != nil {
			t.Fatal(err)
		}
	})
}

func initProvider(fixture string) map[string]*schema.Provider {
	testDGProvider = Provider()
	testDGProvider.ConfigureFunc = providerConfigureVCR(testDGProvider, fixture)

	testDGProviders = map[string]*schema.Provider{
		ProviderNameDG: testDGProvider,
	}

	return testDGProviders
}

func providerConfigureVCR(p *schema.Provider, fixture string) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		config := Config{
			apiKey: d.Get("api_key").(string),
		}

		meta, err := config.Client()
		if err != nil {
			return nil, err
		}

		rec, err := recorder.New(fixture)
		if err != nil {
			return nil, err
		}

		meta.client.HttpClient.Transport = rec

		defer rec.Stop()

		return meta, nil
	}
}
