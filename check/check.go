package check

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

const successMsg = "The certificate currently available on %s is OK."

// IsSafe check if a domain is safe.
func IsSafe(domain string, debug bool) (bool, error) {
	trimmed := strings.TrimPrefix(domain, "*.")

	values := url.Values{}
	values.Set("fqdn", trimmed)

	resp, err := http.DefaultClient.PostForm("https://unboundtest.com/caaproblem/checkhost", values)
	if err != nil {
		return false, fmt.Errorf("client error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return false, fmt.Errorf("domain not found: %s", domain)
	}

	all, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("fail to read body: %w", err)
	}

	if debug {
		log.Println("Response:", string(all))
	}

	return strings.Contains(string(all), fmt.Sprintf(successMsg, trimmed)), nil
}
