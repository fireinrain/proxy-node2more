package cdn

import "testing"

func TestFetchCloudflare(t *testing.T) {
	fetcher := CdnFetcher{}
	fetcher.FetchCloudFlare()
}

func TestFetchGcore(t *testing.T) {
	fetcher := CdnFetcher{}
	fetcher.FetchGcore()
}

func TestFetchCfoudfront(t *testing.T) {
	fetcher := CdnFetcher{}
	fetcher.FetchCloudfront()
}
