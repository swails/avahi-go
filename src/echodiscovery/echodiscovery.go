// This is ressponsible for discovering the Amazon Echo devices on the LAN. This
// is currently done by doing DNS discovery on the local network looking for the
// _workgroup._tcp service and taking the MAC address from the instance name.
// The MAC address is matched against the MAC addresses assigned to Amazon
package echodiscovery

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/oleksandr/bonjour"
)

// This web service keeps an up-to-date database of MAC address-to-vendor
// mappings, allowing us to look at the _workstation._tcp service which the Echo
// seems to broadcast and see if the IP address it publishes originates from an
// Amazon device. This can be fooled by MAC spoofing, but such spoofing is
// harder to do for an embedded device like the Echo, and anyone that would
// spoof a MAC address isn't likely going to need the assistance we intend to
// provide here.
const (
	macLookupHost = "http://api.macvendors.com/"
	echoService   = "_workstation._tcp"
	echoDomain    = "local."
)

var macFinder *regexp.Regexp

func init() {
	macFinder = regexp.MustCompile("([A-Fa-f0-9]{2}[:-][A-Fa-f0-9]{2}[:-][A-Fa-f0-9]{2}[:-][A-Fa-f0-9]{2}[:-][A-Fa-f0-9]{2}[:-][A-Fa-f0-9]{2})")
}

// Returns true if a device with an Amazon MAC address is detected on the
// network broadasting the service _workstation._tcp and false otherwise. This
// function is blocking
func NetworkHasAmazonDevice(timeout time.Duration) (bool, error) {
	resolver, err := bonjour.NewResolver(nil)
	if err != nil {
		return false, fmt.Errorf("could not create resolver: %v", err)
	}

	results := make(chan *bonjour.ServiceEntry)

	// Send the "stop browsing" signal after the desired timeout
	go func() {
		time.Sleep(timeout)
		resolver.Exit <- true
	}()

	err = resolver.Browse(echoService, echoDomain, results)

	if err != nil {
		return false, fmt.Errorf("could not browse DNS services: %v", err)
	}

	for e := range results {
		macAddress := macFinder.FindString(e.Instance)
		if macAddress == "" {
			continue
		}
		resp, err := http.Get(macLookupHost + macAddress)
		if err != nil || resp.StatusCode != 200 {
			continue
		}
		body, err := ioutil.ReadAll(resp.Body)

		manufacturer := strings.Trim(strings.ToUpper(string(body)), "\t\n\r ")
		if strings.HasPrefix(manufacturer, "AMAZON") {
			return true, nil
		}
	}

	return false, nil
}
