// Copyright 2015-2017 Matt Holt and caddy contributors
// Copyright 2017 Kirill Danshin and gramework contributors

package gramework

import (
	"bytes"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
)

// setupTLS sets up the TLS configuration and installs certificates that
// are specified by the user in the config file. All the automatic HTTPS
// stuff comes later outside of this function.
func (app *App) setupTLS(config *Config) error {
	if config == nil {
		return fmt.Errorf("no Config provided")
	}

	config.Enabled = true
	config.ACMEEmail = app.TLSEmails[rand.Intn(len(app.TLSEmails))]

	err := loadCertsInDir("./grametls-cache.d")
	if err != nil {
		return err
	}

	SetDefaultTLSParams(config)

	return nil
}

// loadCertsInDir loads all the certificates/keys in dir, as long as
// the file ends with .pem. This method of loading certificates is
// modeled after haproxy, which expects the certificate and key to
// be bundled into the same file:
// https://cbonte.github.io/haproxy-dconv/configuration-1.5.html#5.1-crt
//
// This function may write to the log as it walks the directory tree.
func loadCertsInDir(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("[WARNING] Unable to traverse into %s; skipping", path)
			return nil
		}
		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(strings.ToLower(info.Name()), ".pem") {
			certBuilder, keyBuilder := new(bytes.Buffer), new(bytes.Buffer)
			var foundKey bool // use only the first key in the file

			bundle, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			for {
				// Decode next block so we can see what type it is
				var derBlock *pem.Block
				derBlock, bundle = pem.Decode(bundle)
				if derBlock == nil {
					break
				}

				if derBlock.Type == "CERTIFICATE" {
					// Re-encode certificate as PEM, appending to certificate chain
					pem.Encode(certBuilder, derBlock)
				} else if derBlock.Type == "EC PARAMETERS" {
					// EC keys generated from openssl can be composed of two blocks:
					// parameters and key (parameter block should come first)
					if !foundKey {
						// Encode parameters
						pem.Encode(keyBuilder, derBlock)

						// Key must immediately follow
						derBlock, bundle = pem.Decode(bundle)
						if derBlock == nil || derBlock.Type != "EC PRIVATE KEY" {
							return fmt.Errorf("%s: expected elliptic private key to immediately follow EC parameters", path)
						}
						pem.Encode(keyBuilder, derBlock)
						foundKey = true
					}
				} else if derBlock.Type == "PRIVATE KEY" || strings.HasSuffix(derBlock.Type, " PRIVATE KEY") {
					// RSA key
					if !foundKey {
						pem.Encode(keyBuilder, derBlock)
						foundKey = true
					}
				} else {
					return fmt.Errorf("%s: unrecognized PEM block type: %s", path, derBlock.Type)
				}
			}

			certPEMBytes, keyPEMBytes := certBuilder.Bytes(), keyBuilder.Bytes()
			if len(certPEMBytes) == 0 {
				return fmt.Errorf("%s: failed to parse PEM data", path)
			}
			if len(keyPEMBytes) == 0 {
				return fmt.Errorf("%s: no private key block found", path)
			}

			err = cacheUnmanagedCertificatePEMBytes(certPEMBytes, keyPEMBytes)
			if err != nil {
				return fmt.Errorf("%s: failed to load cert and key: %v", path, err)
			}
			log.Printf("[INFO] Successfully loaded TLS assets from %s", path)
		}
		return nil
	})
}
