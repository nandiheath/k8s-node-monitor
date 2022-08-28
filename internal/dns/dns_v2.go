package dns

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/nandiheath/k8s-node-monitor/internal/config"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/option"
)

// DNSConfig https://docs.mongodb.com/manual/reference/connection-string/#dns-seed-list-connection-format
type DNSConfig struct {
	Priority int
	Weight   int
	Port     int
}

// UpdateDNSV2 upsert the DNS record with the node ips. this assumes mongos-0/-1/-2 service are setup
func (d *DNS) UpdateDNSV2(addresses []string, dnsConfigs []DNSConfig) {

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
	srvDomainName := fmt.Sprintf("_mongodb._tcp.%s.v2.%s", config.RegionName, config.GCloudDomainName)
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
		log.Printf("there is no existing A records. will create them and update SRV record as well")
		// a records does not exists. going to add them
		needUpdate = true
	} else if srvRecord == nil {
		log.Printf("there is no existing SRV record")
		needUpdate = true
	} else {
		//check if any of the cname a ready
		for _, r := range aRecords {
			found := false
			for _, v := range addresses {
				if v == r.Rrdatas[0] {
					found = true
				}
			}
			if !found {
				log.Printf("node %s no longer exists. will update the A records\n", r.Rrdatas[0])
				needUpdate = true
			}
		}
	}

	if needUpdate {
		log.Printf("going to update the DNS record ")
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
			recordsToAdd = append(recordsToAdd, &dns.ResourceRecordSet{
				Name:    fmt.Sprintf("%s.", host),
				Type:    "A",
				Ttl:     600,
				Rrdatas: []string{fmt.Sprintf("%s", v)},
			})

			for _, dnsConfig := range dnsConfigs {
				rrdata = append(rrdata, fmt.Sprintf("%d %d %d %s.", dnsConfig.Priority, dnsConfig.Weight, dnsConfig.Port, host))
			}
		}

		recordsToAdd = append(recordsToAdd, &dns.ResourceRecordSet{
			Name:    fmt.Sprintf("%s.", srvDomainName),
			Type:    "SRV",
			Ttl:     3600,
			Rrdatas: rrdata,
		})

		if srvRecord != nil {
			recordsToRemove = append(recordsToRemove, srvRecord)
		}

		c := dns.Change{
			Deletions: recordsToRemove,
			Additions: recordsToAdd,
		}

		log.Printf("records to add:")
		for _, set := range recordsToAdd {
			log.Printf("%+v", set)
		}

		log.Printf("records to remove:")
		for _, set := range recordsToRemove {
			log.Printf("%+v", set)
		}


		changeService := dns.NewChangesService(dnsService)

		return
		call := changeService.Create(config.GCloudProjectID, config.GCloudDNSManagedZone, &c)
		newC, err := call.Do()
		if err != nil {
			log.Fatalf("Unable to parse client secret file to config: %v", err)
		}
		log.Printf("DNS change applied. id: %s", newC.Id)
	} else {
		log.Printf("nothing to update")
	}

}
