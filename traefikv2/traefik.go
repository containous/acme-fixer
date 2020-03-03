package traefikv2

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/containous/acme-fixer/check"
	"github.com/containous/traefik/v2/pkg/provider/acme"
)

// Process Removes revoked certificates.
func Process(filename string, dryRun bool) error {
	data, err := readFile(filename)
	if err != nil {
		return err
	}

	err = saveFile(filename+".bak.json", data)
	if err != nil {
		return err
	}

	revoked := getRevokedDomains(data)

	filter(data, revoked)

	if dryRun {
		for name := range revoked {
			fmt.Println("Affected domain:", name)
		}
		return saveFile(filename+".dryrun.json", data)
	}

	return saveFile(filename, data)
}

func getRevokedDomains(data map[string]*acme.StoredData) map[string]struct{} {
	results := map[string]struct{}{}

	for _, storedData := range data {
		for _, certificate := range storedData.Certificates {
			if certificate.Domain.Main != "" {
				safe, err := check.IsSafe(certificate.Domain.Main, true)
				if err != nil {
					log.Println(err)
					continue
				}

				if !safe {
					results[certificate.Domain.Main] = struct{}{}
				}
			}

			for _, san := range certificate.Domain.SANs {
				safe, err := check.IsSafe(san, true)
				if err != nil {
					log.Println(err)
					continue
				}

				if !safe {
					results[certificate.Domain.Main] = struct{}{}
				}
			}
		}
	}

	return results
}

func filter(data map[string]*acme.StoredData, revoked map[string]struct{}) {
	safeCerts := make(map[string][]*acme.CertAndStore)

	for name, storedData := range data {
		for _, certificate := range storedData.Certificates {
			_, ok := revoked[certificate.Domain.Main]
			if !ok {
				safeCerts[name] = append(safeCerts[name], certificate)
			}
		}
	}

	for name := range data {
		data[name].Certificates = safeCerts[name]
	}
}

func readFile(filename string) (map[string]*acme.StoredData, error) {
	file, err := os.Open(filepath.Clean(filename))
	if err != nil {
		return nil, err
	}

	defer file.Close()

	data := map[string]*acme.StoredData{}

	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func saveFile(filename string, data map[string]*acme.StoredData) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(data)
}
