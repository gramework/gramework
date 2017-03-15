// Copyright 2015-2017 Matt Holt and caddy contributors
// Copyright 2017 Kirill Danshin and gramework contributors

package gramework

import (
	"crypto/tls"
	"reflect"
	"testing"
)

func TestConvertTLSConfigProtocolVersions(t *testing.T) {
	// same min and max protocol versions
	config := &Config{
		Enabled:            true,
		ProtocolMinVersion: tls.VersionTLS12,
		ProtocolMaxVersion: tls.VersionTLS12,
	}
	err := config.buildStandardTLSConfig()
	if err != nil {
		t.Fatalf("Did not expect an error, but got %v", err)
	}
	if got, want := config.tlsConfig.MinVersion, uint16(tls.VersionTLS12); got != want {
		t.Errorf("Expected min version to be %x, got %x", want, got)
	}
	if got, want := config.tlsConfig.MaxVersion, uint16(tls.VersionTLS12); got != want {
		t.Errorf("Expected max version to be %x, got %x", want, got)
	}
}

func TestConvertTLSConfigPreferServerCipherSuites(t *testing.T) {
	// prefer server cipher suites
	config := Config{Enabled: true, PreferServerCipherSuites: true}
	err := config.buildStandardTLSConfig()
	if err != nil {
		t.Fatalf("Did not expect an error, but got %v", err)
	}
	if got, want := config.tlsConfig.PreferServerCipherSuites, true; got != want {
		t.Errorf("Expected PreferServerCipherSuites==%v but got %v", want, got)
	}
}

func TestMakeTLSConfigTLSEnabledDisabledError(t *testing.T) {
	// verify handling when Enabled is true and false
	configs := []*Config{
		{Enabled: true},
		{Enabled: false},
	}
	_, err := MakeTLSConfig(configs)
	if err == nil {
		t.Fatalf("Expected an error, but got %v", err)
	}
}

func TestConvertTLSConfigCipherSuites(t *testing.T) {
	// ensure cipher suites are unioned and
	// that TLS_FALLBACK_SCSV is prepended
	configs := []*Config{
		{Enabled: true, Ciphers: []uint16{0xc02c, 0xc030}},
		{Enabled: true, Ciphers: []uint16{0xc012, 0xc030, 0xc00a}},
		{Enabled: true, Ciphers: nil},
	}

	expectedCiphers := [][]uint16{
		{tls.TLS_FALLBACK_SCSV, 0xc02c, 0xc030},
		{tls.TLS_FALLBACK_SCSV, 0xc012, 0xc030, 0xc00a},
		append([]uint16{tls.TLS_FALLBACK_SCSV}, defaultCiphers...),
	}

	for i, config := range configs {
		err := config.buildStandardTLSConfig()
		if err != nil {
			t.Errorf("Test %d: Expected no error, got: %v", i, err)
		}
		if !reflect.DeepEqual(config.tlsConfig.CipherSuites, expectedCiphers[i]) {
			t.Errorf("Test %d: Expected ciphers %v but got %v",
				i, expectedCiphers[i], config.tlsConfig.CipherSuites)
		}

	}
}

func TestStorageForNoURL(t *testing.T) {
	c := &Config{}
	if _, err := c.StorageFor(""); err == nil {
		t.Fatal("Expected error on empty URL")
	}
}

func TestStorageForBadURL(t *testing.T) {
	c := &Config{}
	if _, err := c.StorageFor("http://192.168.0.%31/"); err == nil {
		t.Fatal("Expected error for bad URL")
	}
}

func TestStorageForDefault(t *testing.T) {
	c := &Config{}
	s, err := c.StorageFor("example.com")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := s.(*FileStorage); !ok {
		t.Fatalf("Unexpected storage type: %#v", s)
	}
}
