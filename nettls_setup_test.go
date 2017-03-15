// Copyright 2015-2017 Matt Holt and caddy contributors
// Copyright 2017 Kirill Danshin and gramework contributors

package gramework

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Write test certificates to disk before tests, and clean up
	// when we're done.
	err := ioutil.WriteFile(certFile, testCert, 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(keyFile, testKey, 0644)
	if err != nil {
		os.Remove(certFile)
		log.Fatal(err)
	}

	result := m.Run()

	os.Remove(certFile)
	os.Remove(keyFile)
	os.Exit(result)
}

const (
	certFile = "test_cert.pem"
	keyFile  = "test_key.pem"
)

var testCert = []byte(`-----BEGIN CERTIFICATE-----
MIIBkjCCATmgAwIBAgIJANfFCBcABL6LMAkGByqGSM49BAEwFDESMBAGA1UEAxMJ
bG9jYWxob3N0MB4XDTE2MDIxMDIyMjAyNFoXDTE4MDIwOTIyMjAyNFowFDESMBAG
A1UEAxMJbG9jYWxob3N0MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEs22MtnG7
9K1mvIyjEO9GLx7BFD0tBbGnwQ0VPsuCxC6IeVuXbQDLSiVQvFZ6lUszTlczNxVk
pEfqrM6xAupB7qN1MHMwHQYDVR0OBBYEFHxYDvAxUwL4XrjPev6qZ/BiLDs5MEQG
A1UdIwQ9MDuAFHxYDvAxUwL4XrjPev6qZ/BiLDs5oRikFjAUMRIwEAYDVQQDEwls
b2NhbGhvc3SCCQDXxQgXAAS+izAMBgNVHRMEBTADAQH/MAkGByqGSM49BAEDSAAw
RQIgRvBqbyJM2JCJqhA1FmcoZjeMocmhxQHTt1c+1N2wFUgCIQDtvrivbBPA688N
Qh3sMeAKNKPsx5NxYdoWuu9KWcKz9A==
-----END CERTIFICATE-----
`)

var testKey = []byte(`-----BEGIN EC PARAMETERS-----
BggqhkjOPQMBBw==
-----END EC PARAMETERS-----
-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIGLtRmwzYVcrH3J0BnzYbGPdWVF10i9p6mxkA4+b2fURoAoGCCqGSM49
AwEHoUQDQgAEs22MtnG79K1mvIyjEO9GLx7BFD0tBbGnwQ0VPsuCxC6IeVuXbQDL
SiVQvFZ6lUszTlczNxVkpEfqrM6xAupB7g==
-----END EC PRIVATE KEY-----
`)
