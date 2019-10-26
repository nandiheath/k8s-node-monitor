package dns

import (
	"context"
	"fmt"
	"github.com/nandiheath/k8s-node-monitor/internal/config"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/option"
	"log"
	"regexp"
)

type DNS struct {
}

func New() *DNS {
	return &DNS{}
}

func (d *DNS) UpdateDNS(addresses []string) {

	ctx := context.Background()
	dnsService, err := dns.NewService(ctx, option.WithCredentialsJSON(config.GoogleCredential))
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	records, err := dnsService.ResourceRecordSets.List(config.GCloudProjectID, config.GCloudDNSManagedZone).Do()
	if err != nil {
		log.Fatalf("unable to list the current dns records. %v", err)
	}

	var recordsToRemove []*dns.ResourceRecordSet
	var recordsToAdd []*dns.ResourceRecordSet
	var aRecords []*dns.ResourceRecordSet
	var srvRecord *dns.ResourceRecordSet
	srvDomainName := fmt.Sprintf("_mongodb._tcp.%s.%s", config.RegionName, config.GCloudDomainName)
	r, _ := regexp.Compile(fmt.Sprintf("db\\d.%s.%s", config.RegionName, config.GCloudDomainName))
	for _, v := range records.Rrsets {
		fmt.Printf("Records: [%s]: %s\n", v.Type, v.Name)
		if v.Name == fmt.Sprintf("%s.", srvDomainName) {
			srvRecord = v
		}
		if v.Type == "A" {
			if r.MatchString(v.Name) {
				aRecords = append(aRecords, v)
			}
		}
	}

	needUpdate := false

	if len(aRecords) == 0 {
		// cname records does not exists. going to add them
		needUpdate = true
	} else {
		//check if any of the cname a;ready
		for _, r := range aRecords {
			found := false
			for _, v := range addresses {
				if v == r.Rrdatas[0] {
					found = true
				}
			}
			if !found {
				fmt.Printf("%s no longer exists\n", r.Rrdatas[0])
				needUpdate = true
			}
		}
	}

	if needUpdate {
		var rrdata []string
		for _, r := range aRecords {
			recordsToRemove = append(recordsToRemove, r)
		}
		for i, v := range addresses {
			// maximum 5 records
			if i == 5 {
				break
			}
			host := fmt.Sprintf("db%d.%s.%s", i, config.RegionName, config.GCloudDomainName)
			rrdata = append(rrdata, fmt.Sprintf("0 5 30017 %s.", host))

			recordsToAdd = append(recordsToAdd, &dns.ResourceRecordSet{
				Name:    fmt.Sprintf("%s.", host),
				Type:    "A",
				Ttl:     600,
				Rrdatas: []string{fmt.Sprintf("%s", v)},
			})
		}

		recordsToAdd = append(recordsToAdd, &dns.ResourceRecordSet{
			Name:    fmt.Sprintf("%s.", srvDomainName),
			Type:    "SRV",
			Ttl:     3600,
			Rrdatas: rrdata,
		})

		recordsToRemove = append(recordsToRemove, srvRecord)

		c := dns.Change{
			Deletions: recordsToRemove,
			Additions: recordsToAdd,
		}

		changeService := dns.NewChangesService(dnsService)

		call := changeService.Create(config.GCloudProjectID, config.GCloudDNSManagedZone, &c)
		newC, err := call.Do()
		if err != nil {
			log.Fatalf("Unable to parse client secret file to config: %v", err)
		}
		fmt.Printf("change applied. id: %s", newC.Id)
	}

}
