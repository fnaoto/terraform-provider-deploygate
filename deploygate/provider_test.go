package deploygate

import (
	"context"
	"sync"
	"testing"

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

func init() {
	testDGProvider = Provider()

	testDGProviders = map[string]*schema.Provider{
		ProviderNameDG: testDGProvider,
	}
}
