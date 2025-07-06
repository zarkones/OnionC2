package apictrl

import (
	"api/geoip"
	"api/models"
	"encoding/json"
	"log"
	"net/http"
	"sort"
)

func Deduplicate(slice []geoip.Country) []geoip.Country {
	ids := make(map[string]bool)
	l := []geoip.Country{}
	for _, item := range slice {
		if _, v := ids[item.ISOCode]; !v {
			ids[item.ISOCode] = true
			l = append(l, item)
		}
	}
	return l
}

func GetOrigins() func(w http.ResponseWriter, r *http.Request) {
	origins := make([]geoip.Country, len(geoip.Networks))
	for i, n := range geoip.Networks {
		origins[i] = n.Country
	}
	origins = Deduplicate(origins)

	sort.Slice(origins, func(i, j int) bool {
		if origins[i].Name == "unknown" {
			return true
		}
		return origins[i].Name < origins[j].Name
	})

	origins = append([]geoip.Country{
		{
			ISOCode: "unknown",
			Name:    "unknown",
		},
	}, origins...)

	serializedOrigins, err := json.Marshal(origins)
	if err != nil {
		log.Println("fatal: geoip origins data failed to be serialized", err)
		return func(w http.ResponseWriter, r *http.Request) {
			_, _, reject := authenticateAndAuthorize(w, r, models.PERMISSION_AGENTS_LIST, nil)
			if reject {
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		_, _, reject := authenticateAndAuthorize(w, r, models.PERMISSION_AGENTS_LIST, nil)
		if reject {
			return
		}
		w.Write(serializedOrigins)
	}
}
