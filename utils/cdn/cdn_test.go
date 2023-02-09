package cdn

import "testing"

func TestFetchCloudflare(t *testing.T) {
	fetcher := CdnFetcher{}
	fetcher.fetchCloudFlare()
}

func TestFetchGcore(t *testing.T) {
	fetcher := CdnFetcher{}
	fetcher.fetchGcore()
}

func TestFetchCfoudfront(t *testing.T) {
	fetcher := CdnFetcher{}
	fetcher.fetchCloudfront()
}
