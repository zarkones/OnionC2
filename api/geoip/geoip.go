package geoip

import (
	"encoding/binary"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sort"
	"strconv"

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

// networks holds the sorted list of network entries.
var networks []NetworkEntry

// readGeonames reads the geoname CSV and returns a map from geoname_id to Country.
func readGeonames(filename string) (map[int]Country, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open geonames file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read() // Skip header
	if err != nil {
		return nil, fmt.Errorf("failed to read geonames header: %v", err)
	}

	geonames := make(map[int]Country)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read geonames row: %v", err)
		}
		geonameID, err := strconv.Atoi(row[0])
		if err != nil {
			log.Printf("Skipping invalid geoname_id: %s", row[0])
			continue
		}
		isoCode := row[4]
		name := row[5]
		geonames[geonameID] = Country{ISOCode: isoCode, Name: name}
	}
	return geonames, nil
}

// readNetworks reads the network CSV and returns a slice of NetworkEntry.
func readNetworks(filename string, geonames map[int]Country) ([]NetworkEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open networks file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read() // Skip header
	if err != nil {
		return nil, fmt.Errorf("failed to read networks header: %v", err)
	}

	var networks []NetworkEntry
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read networks row: %v", err)
		}
		networkStr := row[0]
		_, ipNet, err := net.ParseCIDR(networkStr)
		if err != nil {
			log.Printf("Skipping invalid CIDR: %s", networkStr)
			continue
		}
		regCountryIDStr := row[2] // registered_country_geoname_id
		if regCountryIDStr == "" {
			log.Printf("Skipping network with empty registered_country_geoname_id: %s", networkStr)
			continue
		}
		regCountryID, err := strconv.Atoi(regCountryIDStr)
		if err != nil {
			log.Printf("Skipping invalid registered_country_geoname_id: %s", regCountryIDStr)
			continue
		}
		country, ok := geonames[regCountryID]
		if !ok {
			log.Printf("No country found for geoname_id %d in network %s", regCountryID, networkStr)
			continue
		}
		networks = append(networks, NetworkEntry{Net: ipNet, Country: country})
	}
	return networks, nil
}

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
	i := sort.Search(len(networks), func(i int) bool {
		return ipToUint32(networks[i].Net.IP) > ipUint
	})
	// Check the previous network (if it exists)
	if i > 0 && networks[i-1].Net.Contains(ip) {
		country := networks[i-1].Country
		return country.Name, country.ISOCode, nil
	}
	return "", "", fmt.Errorf("IP not found in any network")
}

//go:embed networks.json
var geoData []byte

func Init() (err error) {
	if err := json.Unmarshal(geoData, &networks); err != nil {
		return err
	}
	sort.Slice(networks, func(i, j int) bool {
		return ipToUint32(networks[i].Net.IP) < ipToUint32(networks[j].Net.IP)
	})
	return nil
}
