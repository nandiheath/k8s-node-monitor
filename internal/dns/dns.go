package dns

import (
	"context"
	"fmt"
	"github.com/nandiheath/k8s-node-monitor/internal/config"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/option"
	"log"
)

type DNS struct {
}

func New() *DNS {
	return &DNS{}
}

func (d *DNS) UpdateDNS(addresses []string) {

	ctx := context.Background()
	dnsService, err := dns.NewService(ctx, option.WithCredentialsJSON(config.Config.GoogleCredential))
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	records, err := dnsService.ResourceRecordSets.List(config.Config.GCloudProjectID, config.Config.GCloudDNSManagedZone).Do()
	if err != nil {
		log.Fatalf("unable to list the current dns records. %v", err)
	}

	var recordsToRemove []*dns.ResourceRecordSet
	var recordsToAdd []*dns.ResourceRecordSet

	for _, v := range records.Rrsets {
		fmt.Printf("Records: %s: %s\n", v.Type, v.Name)
		if v.Name == config.Config.GCloudSrvRecordName {
			fmt.Printf("record to change found")
			recordsToRemove = append(recordsToRemove, v)
		}
	}

	var rrdata []string

	for _, v := range addresses {
		rrdata = append(rrdata, fmt.Sprintf("0 5 27017 %s.", v))
	}

	recordsToAdd = append(recordsToAdd, &dns.ResourceRecordSet{
		Name:    config.Config.GCloudSrvRecordName,
		Type:    "SRV",
		Ttl:     3600,
		Rrdatas: rrdata,
	})
	c := dns.Change{
		Deletions: recordsToRemove,
		Additions: recordsToAdd,
	}

	changeService := dns.NewChangesService(dnsService)

	call := changeService.Create(config.Config.GCloudProjectID, config.Config.GCloudDNSManagedZone, &c)
	newC, err := call.Do()
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	fmt.Printf("change applied. id: %s", newC.Id)

}
