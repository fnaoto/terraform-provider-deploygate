package deploygate

import (
	"context"
	"os"
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
		if os.Getenv("DG_API_KEY") == "" {
			t.Fatal("DG_API_KEY must be set for acceptance tests")
		}

		if os.Getenv("DG_USER_NAME") == "" && os.Getenv("DG_ORGANIZATION_NAME") == "" {
			t.Fatal("DG_USER_NAME or DG_ORGANIZATION_NAME must be set for acceptance tests")
		}

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
