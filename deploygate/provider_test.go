package deploygate

import (
	"context"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
	"testing"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	ProviderNameDG  = "deploygate"
	FixtureBasePath = "fixtures"
)

var testDGProvider *schema.Provider
var testDGProviderConfigure sync.Once
var testDGProviders map[string]*schema.Provider
var testDGConfigs map[string]*Config

func Test_DGPreCheck(t *testing.T) {
	initProvider(t)
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

func initProvider(t *testing.T) map[string]*schema.Provider {
	testDGProvider = Provider()
	testDGProvider.ConfigureFunc = providerConfigureVCR(testDGProvider, t)

	testDGProviders = map[string]*schema.Provider{
		ProviderNameDG: testDGProvider,
	}

	return testDGProviders
}

func providerConfigureVCR(p *schema.Provider, t *testing.T) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		config := &Config{
			apiKey: d.Get("api_key").(string),
		}

		err := config.initClient()
		if err != nil {
			return nil, err
		}

		fixture := filepath.Join(FixtureBasePath, t.Name(), strconv.Itoa(len(testDGConfigs)))

		rec, err := recorder.New(fixture)
		if err != nil {
			return nil, err
		}

		rec.AddSaveFilter(func(i *cassette.Interaction) error {
			delete(i.Request.Headers, "Authorization")
			i.Response.Headers = make(http.Header)
			return nil
		})

		config.client.HttpClient.Transport = rec

		testDGConfigs[fixture] = config

		return config, nil
	}
}

func closeVCR(t *testing.T) {
	for _, cfg := range testDGConfigs {
		err := cfg.client.HttpClient.Transport.(*recorder.Recorder).Stop()
		if err != nil {
			t.Error(err)
		}
	}
}

func testWithVCR(t *testing.T, c resource.TestCase) {
	testDGConfigs = make(map[string]*Config)

	providers := initProvider(t)
	c.Providers = providers
	defer closeVCR(t)
	resource.Test(t, c)
}
