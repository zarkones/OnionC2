package geoip

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sort"

	_ "embed"
)

// Country holds the country ISO code and name.
type Country struct {
	ISOCode string `json:"i"`
	Name    string `json:"n"`
}

// NetworkEntry holds an IP network and its associated country.
type NetworkEntry struct {
	Net     *net.IPNet
	Country Country
}

// MarshalJSON custom marshaler for NetworkEntry to serialize N as a string.
func (ne NetworkEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Net     string  `json:"n"`
		Country Country `json:"c"`
	}{
		Net:     ne.Net.String(),
		Country: ne.Country,
	})
}

// UnmarshalJSON custom unmarshaler for NetworkEntry to deserialize string back to N.
func (ne *NetworkEntry) UnmarshalJSON(data []byte) error {
	var aux struct {
		Net     string  `json:"n"`
		Country Country `json:"c"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	_, ipNet, err := net.ParseCIDR(aux.Net)
	if err != nil {
		return fmt.Errorf("invalid CIDR: %s", aux.Net)
	}
	ne.Net = ipNet
	ne.Country = aux.Country
	return nil
}

// Networks holds the sorted list of network entries.
var Networks []NetworkEntry

// ipToUint32 converts an IPv4 address to a uint32 for comparison.
// Returns zero on error.
func ipToUint32(ip net.IP) uint32 {
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	if ip == nil {
		log.Println("not an IPv4 address:", ip)
		return 0
	}
	return binary.BigEndian.Uint32(ip)
}

// IpToCountry returns the country name and ISO code for a given IP address.
func IpToCountry(ipStr string) (string, string, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "", "", fmt.Errorf("invalid IP address")
	}
	ip = ip.To4()
	if ip == nil {
		return "", "", fmt.Errorf("only IPv4 addresses are supported")
	}
	ipUint := ipToUint32(ip)

	// Binary search to find the first network where Net.IP > ip
	i := sort.Search(len(Networks), func(i int) bool {
		return ipToUint32(Networks[i].Net.IP) > ipUint
	})
	// Check the previous network (if it exists)
	if i > 0 && Networks[i-1].Net.Contains(ip) {
		country := Networks[i-1].Country
		return country.Name, country.ISOCode, nil
	}
	return "", "", fmt.Errorf("IP not found in any network")
}

//go:embed networks.json
var geoData []byte

func Init() (err error) {
	if err := json.Unmarshal(geoData, &Networks); err != nil {
		return err
	}
	sort.Slice(Networks, func(i, j int) bool {
		return ipToUint32(Networks[i].Net.IP) < ipToUint32(Networks[j].Net.IP)
	})
	return nil
}
