package traefikv1

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/containous/acme-fixer/check"
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

	data = filter(data, revoked)

	if dryRun {
		for name := range revoked {
			fmt.Println("Affected domain:", name)
		}
		return saveFile(filename+".dryrun.json", data)
	}

	return saveFile(filename, data)
}

func getRevokedDomains(data StoredData) map[string]struct{} {
	results := map[string]struct{}{}

	for _, certificate := range data.Certificates {
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
	return results
}

func filter(data StoredData, revoked map[string]struct{}) StoredData {
	var safeCerts []*Certificate
	for _, certificate := range data.Certificates {
		_, ok := revoked[certificate.Domain.Main]
		if !ok {
			safeCerts = append(safeCerts, certificate)
		}
	}

	data.Certificates = safeCerts

	return data
}

func readFile(filename string) (StoredData, error) {
	file, err := os.Open(filepath.Clean(filename))
	if err != nil {
		return StoredData{}, err
	}

	defer file.Close()

	data := StoredData{}

	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return StoredData{}, err
	}

	return data, nil
}

func saveFile(filename string, data StoredData) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(data)
}
