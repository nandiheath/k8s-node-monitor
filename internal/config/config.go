package config

import (
	"encoding/base64"
	_ "github.com/joho/godotenv/autoload"
	"os"
)

type config struct {
	GoogleCredential     []byte
	GCloudProjectID      string
	GCloudSrvRecordName  string
	GCloudDNSManagedZone string
}

var Config config

func init() {
	Config = config{}
	b64String := os.Getenv("GOOGLE_CREDENTIAL_JSON_B64")
	if b64String == "" {
		panic("GOOGLE_CREDENTIAL_JSON_B64 is not set!")
	}

	Config.GoogleCredential, _ = base64.StdEncoding.DecodeString(b64String)
	Config.GCloudProjectID = os.Getenv("GCLOUD_DNS_PROJECT_ID")
	Config.GCloudSrvRecordName = os.Getenv("GCLOUD_SRV_RECORD_NAME")
    Config.GCloudDNSManagedZone = os.Getenv("GCLOUD_DNS_MANAGED_ZONE")
}
