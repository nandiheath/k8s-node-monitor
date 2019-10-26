package config

import (
	"encoding/base64"
	_ "github.com/joho/godotenv/autoload"
	"os"
)

var (
	GoogleCredential     []byte
	GCloudProjectID      string
	GCloudDomainName     string
	GCloudDNSManagedZone string
	RegionName           string
)

func init() {
	b64String := os.Getenv("GOOGLE_CREDENTIAL_JSON_B64")
	if b64String == "" {
		panic("GOOGLE_CREDENTIAL_JSON_B64 is not set!")
	}

	GoogleCredential, _ = base64.StdEncoding.DecodeString(b64String)
	GCloudProjectID = os.Getenv("GCLOUD_DNS_PROJECT_ID")
	GCloudDomainName = os.Getenv("GCLOUD_DNS_DOMAIN_NAME")
	GCloudDNSManagedZone = os.Getenv("GCLOUD_DNS_MANAGED_ZONE")
	RegionName = os.Getenv("REGION_NAME")
}
