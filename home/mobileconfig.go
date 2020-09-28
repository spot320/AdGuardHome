package home

import (
	"net/http"

	"howett.net/plist"
)

type DNSSettings struct {
	DNSProtocol     string
	ServerURL       string `plist:",omitempty"`
	ServerName      string `plist:",omitempty"`
	ServerAddresses []string
}

type PayloadContent = struct {
	Name               string
	PayloadDescription string
	PayloadDisplayName string
	PayloadIdentifier  string
	PayloadType        string
	PayloadUUID        string
	PayloadVersion     int
	DNSSettings        DNSSettings
}

type MobileConfig = struct {
	PayloadContent           []PayloadContent
	PayloadDescription       string
	PayloadDisplayName       string
	PayloadIdentifier        string
	PayloadRemovalDisallowed bool
	PayloadType              string
	PayloadUUID              string
	PayloadVersion           int
}

func getMobileConfig(w http.ResponseWriter, d DNSSettings) []byte {
	data := MobileConfig{
		PayloadContent: []PayloadContent{{
			Name:               "AdGuard DNS over HTTPS",
			PayloadDescription: "Configures device to use AdGuard DNS",
			PayloadDisplayName: "AdGuard DNS",
			PayloadIdentifier:  "com.apple.dnsSettings.managed.767A11FC-31D2-4950-815E-B37B15448CA2",
			PayloadType:        "com.apple.dnsSettings.managed",
			PayloadUUID:        "767A11FC-31D2-4950-815E-B37B15448CA2",
			PayloadVersion:     1,
			DNSSettings:        d,
		}},
		PayloadDescription:       "Adds AdGuard DNS toBig Sur and iOS 14 or newer systems",
		PayloadDisplayName:       "AdGuard DNS",
		PayloadIdentifier:        "E3E3CB8B-C59E-486B-A713-D765328DB2A2",
		PayloadRemovalDisallowed: false,
		PayloadType:              "Configuration",
		PayloadUUID:              "F2609BEA-93D6-4966-8487-33713DBCB644",
		PayloadVersion:           1,
	}

	mobileconfig, err := plist.MarshalIndent(data, plist.XMLFormat, "\t")

	if err != nil {
		httpError(w, http.StatusInternalServerError, "plist.MarshalIndent: %s", err)
	}

	return mobileconfig
}

func handleMobileConfig(w http.ResponseWriter, d DNSSettings) {
	mobileconfig := getMobileConfig(w, d)

	w.Header().Set("Content-Type", "application/xml")
	w.Write(mobileconfig)
}

func handleMobileConfigDoh(w http.ResponseWriter, r *http.Request) {
	handleMobileConfig(w, DNSSettings{
		DNSProtocol:     "HTTPS",
		ServerURL:       "https://dns.adguard.com/dns-query",
		ServerAddresses: []string{getDNSEncryption().https},
	})
}

func handleMobileConfigDot(w http.ResponseWriter, r *http.Request) {
	handleMobileConfig(w, DNSSettings{
		DNSProtocol:     "TLS",
		ServerName:      "dns.adguard.com",
		ServerAddresses: []string{getDNSEncryption().tls},
	})
}
