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
	testDGProviderConfigure.Do(func() {
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
