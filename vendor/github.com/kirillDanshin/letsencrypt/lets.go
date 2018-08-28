// Copyright 2016 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package letsencrypt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"

	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/time/rate"

	"github.com/xenolf/lego/acme"
)

const debug = false

// A Manager m takes care of obtaining and refreshing a collection of TLS certificates
// obtained by LetsEncrypt.org.
//  The zero Manager is not yet registered with LetsEncrypt.org and has no TLS certificates
// but is nonetheless ready for use.
// See the package comment for an overview of how to use a Manager.
type Manager struct {
	mu           sync.Mutex
	state        state
	rateLimit    *rate.Limiter
	newHostLimit *rate.Limiter
	certCache    map[string]*cacheEntry
	certTokens   map[string]*tls.Certificate
	watchChan    chan struct{}
}

// state is the serializable state for the Manager.
// It also implements acme.User.
type state struct {
	Email string
	Reg   *acme.RegistrationResource
	Key   string
	key   *ecdsa.PrivateKey
	Hosts []string
	Certs map[string]stateCert
}

func (s *state) GetEmail() string                            { return s.Email }
func (s *state) GetRegistration() *acme.RegistrationResource { return s.Reg }
func (s *state) GetPrivateKey() crypto.PrivateKey            { return s.key }

type stateCert struct {
	Cert string
	Key  string
}

func (cert stateCert) toTLS() (*tls.Certificate, error) {
	c, err := tls.X509KeyPair([]byte(cert.Cert), []byte(cert.Key))
	if err != nil {
		return nil, err
	}
	return &c, err
}

type cacheEntry struct {
	host string
	m    *Manager

	mu         sync.Mutex
	cert       *tls.Certificate
	timeout    time.Time
	refreshing bool
	err        error
}

func (m *Manager) init() {
	m.mu.Lock()
	if m.certCache == nil {
		m.rateLimit = rate.NewLimiter(rate.Every(1*time.Minute), 20)
		m.newHostLimit = rate.NewLimiter(rate.Every(3*time.Hour), 20)
		m.certCache = map[string]*cacheEntry{}
		m.certTokens = map[string]*tls.Certificate{}
		m.watchChan = make(chan struct{}, 1)
		m.watchChan <- struct{}{}
	}
	m.mu.Unlock()
}

// Watch returns the manager's watch channel,
// which delivers a notification after every time the
// manager's state (as exposed by Marshal and Unmarshal) changes.
// All calls to Watch return the same watch channel.
//
// The watch channel includes notifications about changes
// before the first call to Watch, so that in the pattern below,
// the range loop executes once immediately, saving
// the result of setup (along with any background updates that
// may have raced in quickly).
//
//	m := new(letsencrypt.Manager)
//	setup(m)
//	go backgroundUpdates(m)
//	for range m.Watch() {
//		save(m.Marshal())
//	}
//
func (m *Manager) Watch() <-chan struct{} {
	m.init()
	m.updated()
	return m.watchChan
}

func (m *Manager) updated() {
	select {
	case m.watchChan <- struct{}{}:
	default:
	}
}

func (m *Manager) CacheFile(name string) error {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	f.Close()
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return err
	}
	if len(data) > 0 {
		if err := m.Unmarshal(string(data)); err != nil {
			return err
		}
	}
	go func() {
		for range m.Watch() {
			err := ioutil.WriteFile(name, []byte(m.Marshal()), 0600)
			if err != nil {
				log.Printf("writing letsencrypt cache: %v", err)
			}
		}
	}()
	return nil
}

// Registered reports whether the manager has registered with letsencrypt.org yet.
func (m *Manager) Registered() bool {
	m.init()
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.registered()
}

func (m *Manager) registered() bool {
	return m.state.Reg != nil && m.state.Reg.Body.Agreement != ""
}

// Register registers the manager with letsencrypt.org, using the given email address.
// Registration may require agreeing to the letsencrypt.org terms of service.
// If so, Register calls prompt(url) where url is the URL of the terms of service.
// Prompt should report whether the caller agrees to the terms.
// A nil prompt func is taken to mean that the user always agrees.
// The email address is sent to LetsEncrypt.org but otherwise unchecked;
// it can be omitted by passing the empty string.
//
// Calling Register is only required to make sure registration uses a
// particular email address or to insert an explicit prompt into the
// registration sequence. If the manager is not registered, it will
// automatically register with no email address and automatic
// agreement to the terms of service at the first call to Cert or GetCertificate.
func (m *Manager) Register(email string, prompt func(string) bool) error {
	m.init()
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.register(email, prompt)
}

func (m *Manager) register(email string, prompt func(string) bool) error {
	if m.registered() {
		return fmt.Errorf("already registered")
	}
	m.state.Email = email
	if m.state.key == nil {
		key, err := newKey()
		if err != nil {
			return fmt.Errorf("generating key: %v", err)
		}
		Key, err := marshalKey(key)
		if err != nil {
			return fmt.Errorf("generating key: %v", err)
		}
		m.state.key = key
		m.state.Key = string(Key)
	}

	c, err := acme.NewClient(letsEncryptURL, &m.state, acme.EC256)
	//
	log.SetOutput(ioutil.Discard) // disable the logger
	acme.Logger = nil
	//
	if err != nil {
		return fmt.Errorf("create client: %v", err)
	}
	reg, err := c.Register()
	if err != nil {
		return fmt.Errorf("register: %v", err)
	}

	m.state.Reg = reg
	if reg.Body.Agreement == "" {
		if prompt != nil && !prompt(reg.TosURL) {
			return fmt.Errorf("did not agree to TOS")
		}
		if err := c.AgreeToTOS(); err != nil {
			return fmt.Errorf("agreeing to TOS: %v", err)
		}
	}

	m.updated()

	return nil
}

// Marshal returns an encoding of the manager's state,
// suitable for writing to disk and reloading by calling Unmarshal.
// The state includes registration status, the configured host list
// from SetHosts, and all known certificates, including their private
// cryptographic keys.
// Consequently, the state should be kept private.
func (m *Manager) Marshal() string {
	m.init()
	m.mu.Lock()
	js, err := json.MarshalIndent(&m.state, "", "\t")
	m.mu.Unlock()
	if err != nil {
		panic("unexpected json.Marshal failure")
	}
	return string(js)
}

// Unmarshal restores the state encoded by a previous call to Marshal
// (perhaps on a different Manager in a different program).
func (m *Manager) Unmarshal(enc string) error {
	m.init()
	var st state
	if err := json.Unmarshal([]byte(enc), &st); err != nil {
		return err
	}
	if st.Key != "" {
		key, err := unmarshalKey(st.Key)
		if err != nil {
			return err
		}
		st.key = key
	}
	m.mu.Lock()
	m.state = st
	m.mu.Unlock()
	for host, cert := range m.state.Certs {
		c, err := cert.toTLS()
		if err != nil {
			log.Printf("letsencrypt: ignoring entry for %s: %v", host, err)
			continue
		}
		m.certCache[host] = &cacheEntry{host: host, m: m, cert: c}
	}
	m.updated()
	return nil
}

// SetHosts sets the manager's list of known host names.
// If the list is non-nil, the manager will only ever attempt to acquire
// certificates for host names on the list.
// If the list is nil, the manager does not restrict the hosts it will
// ask for certificates for.
func (m *Manager) SetHosts(hosts []string) {
	m.init()
	m.mu.Lock()
	m.state.Hosts = append(m.state.Hosts[:0], hosts...)
	m.mu.Unlock()
	m.updated()
}

// GetCertificate can be placed a tls.Config's GetCertificate field to make
// the TLS server use Let's Encrypt certificates.
// Each time a client connects to the TLS server expecting a new host name,
// the TLS server's call to GetCertificate will trigger an exchange with the
// Let's Encrypt servers to obtain that certificate, subject to the manager rate limits.
//
// As noted in the Manager's documentation comment,
// to obtain a certificate for a given host name, that name
// must resolve to a computer running a TLS server on port 443
// that obtains TLS SNI certificates by calling m.GetCertificate.
// In the standard usage, then, installing m.GetCertificate in the tls.Config
// both automatically provisions the TLS certificates needed for
// ordinary HTTPS service and answers the challenges from LetsEncrypt.org.
func (m *Manager) GetCertificate(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	m.init()

	host := clientHello.ServerName

	if debug {
		log.Printf("GetCertificate %s", host)
	}

	if strings.HasSuffix(host, ".acme.invalid") {
		m.mu.Lock()
		cert := m.certTokens[host]
		m.mu.Unlock()
		if cert == nil {
			return nil, fmt.Errorf("unknown host")
		}
		return cert, nil
	}

	if host == "" {
		host = localhost
	}

	return m.Cert(host)
}

// Cert returns the certificate for the given host name, obtaining a new one if necessary.
//
// As noted in the documentation for Manager and for the GetCertificate method,
// obtaining a certificate requires that m.GetCertificate be associated with host.
// In most servers, simply starting a TLS server with a configuration referring
// to m.GetCertificate is sufficient, and Cert need not be called.
//
// The main use of Cert is to force the manager to obtain a certificate
// for a particular host name ahead of time.
func (m *Manager) Cert(host string) (*tls.Certificate, error) {
	host = strings.ToLower(host)
	if debug {
		log.Printf("Cert %s", host)
	}

	m.init()
	m.mu.Lock()
	if !m.registered() {
		m.register("", nil)
	}

	ok := false
	if m.state.Hosts == nil {
		ok = true
	} else {
		for _, h := range m.state.Hosts {
			if host == h {
				ok = true
				break
			}
		}
	}
	if !ok {
		m.mu.Unlock()
		return nil, fmt.Errorf("unknown host")
	}

	// Otherwise look in our cert cache.
	entry, ok := m.certCache[host]
	if !ok {
		if host == localhost {
			c, err := m.certSelfSigned()
			if err != nil {
				return nil, err
			}

			entry = &cacheEntry{host: host, m: m, cert: c}
		} else {
			r := m.rateLimit.Reserve()
			ok := r.OK()
			if ok {
				ok = m.newHostLimit.Allow()
				if !ok {
					r.Cancel()
				}
			}
			if !ok {
				m.mu.Unlock()
				return nil, fmt.Errorf("rate limited")
			}
			entry = &cacheEntry{host: host, m: m}
		}

		m.certCache[host] = entry
	}
	m.mu.Unlock()

	entry.mu.Lock()
	defer entry.mu.Unlock()
	entry.init()
	if entry.err != nil {
		return nil, entry.err
	}

	return entry.cert, nil
}

const localhost = "localhost"

func (m *Manager) certSelfSigned() (*tls.Certificate, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to generate serial number: %s", err)
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},

		NotBefore: time.Now(),
		NotAfter:  time.Now().AddDate(0, 1, 0),

		DNSNames: []string{"localhost"},

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	cert, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	if err != nil {
		return nil, err
	}

	if m.state.Certs == nil {
		m.state.Certs = make(map[string]stateCert)
	}

	m.state.Certs[localhost] = stateCert{
		Cert: encodePEM(cert),
		Key:  encodePEM(priv),
	}

	m.updated()

	return &tls.Certificate{
		Certificate: [][]byte{cert},
		PrivateKey:  priv,
	}, nil
}

func (e *cacheEntry) init() {
	if e.err != nil && time.Now().Before(e.timeout) {
		return
	}
	if e.cert != nil {
		if e.timeout.IsZero() {
			t, err := certRefreshTime(e.cert)
			if err != nil {
				e.err = err
				e.timeout = time.Now().Add(1 * time.Minute)
				e.cert = nil
				return
			}
			e.timeout = t
		}
		if time.Now().After(e.timeout) && !e.refreshing {
			e.refreshing = true
			go e.refresh()
		}
		return
	}

	cert, refreshTime, err := e.m.verify(e.host)
	e.m.mu.Lock()
	e.m.certCache[e.host] = e
	e.m.mu.Unlock()
	e.install(cert, refreshTime, err)
}

func (e *cacheEntry) install(cert *tls.Certificate, refreshTime time.Time, err error) {
	e.cert = nil
	e.timeout = time.Time{}
	e.err = nil

	if err != nil {
		e.err = err
		e.timeout = time.Now().Add(1 * time.Minute)
		return
	}

	e.cert = cert
	e.timeout = refreshTime
}

func (e *cacheEntry) refresh() {
	e.m.rateLimit.Wait(context.Background())
	cert, refreshTime, err := e.m.verify(e.host)

	e.mu.Lock()
	defer e.mu.Unlock()
	e.refreshing = false
	if err == nil {
		e.install(cert, refreshTime, nil)
	}
}

func (m *Manager) verify(host string) (cert *tls.Certificate, refreshTime time.Time, err error) {
	c, err := acme.NewClient(letsEncryptURL, &m.state, acme.EC256)
	if err != nil {
		return
	}
	if err = c.SetChallengeProvider(acme.TLSSNI01, tlsProvider{m}); err != nil {
		return
	}
	c.SetChallengeProvider(acme.TLSSNI01, tlsProvider{m})
	c.ExcludeChallenges([]acme.Challenge{acme.HTTP01})
	acmeCert, errmap := c.ObtainCertificate([]string{host}, true, nil, true)
	if len(errmap) > 0 {
		if debug {
			log.Printf("ObtainCertificate %v => %v", host, errmap)
		}
		err = fmt.Errorf("%v", errmap)
		return
	}

	entryCert := stateCert{
		Cert: string(acmeCert.Certificate),
		Key:  string(acmeCert.PrivateKey),
	}
	cert, err = entryCert.toTLS()
	if err != nil {
		if debug {
			log.Printf("ObtainCertificate %v toTLS failure: %v", host, err)
		}
		return
	}
	if refreshTime, err = certRefreshTime(cert); err != nil {
		return
	}

	m.mu.Lock()
	if m.state.Certs == nil {
		m.state.Certs = make(map[string]stateCert)
	}
	m.state.Certs[host] = entryCert
	m.mu.Unlock()
	m.updated()

	return cert, refreshTime, nil
}

func certRefreshTime(cert *tls.Certificate) (time.Time, error) {
	xc, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		if debug {
			log.Printf("ObtainCertificate to X.509 failure: %v", err)
		}
		return time.Time{}, err
	}
	t := xc.NotBefore.Add(xc.NotAfter.Sub(xc.NotBefore) / 2)
	monthEarly := xc.NotAfter.Add(-30 * 24 * time.Hour)
	if t.Before(monthEarly) {
		t = monthEarly
	}
	return t, nil
}

// tlsProvider implements acme.ChallengeProvider for TLS handshake challenges.
type tlsProvider struct {
	m *Manager
}

func (p tlsProvider) Present(domain, token, keyAuth string) error {
	cert, dom, err := acme.TLSSNI01ChallengeCert(keyAuth)
	if err != nil {
		return err
	}

	p.m.mu.Lock()
	p.m.certTokens[dom] = &cert
	p.m.mu.Unlock()

	return nil
}

func (p tlsProvider) CleanUp(domain, token, keyAuth string) error {
	_, dom, err := acme.TLSSNI01ChallengeCert(keyAuth)
	if err != nil {
		return err
	}

	p.m.mu.Lock()
	delete(p.m.certTokens, dom)
	p.m.mu.Unlock()

	return nil
}

func marshalKey(key *ecdsa.PrivateKey) ([]byte, error) {
	data, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: data}), nil
}

func unmarshalKey(text string) (*ecdsa.PrivateKey, error) {
	b, _ := pem.Decode([]byte(text))
	if b == nil {
		return nil, fmt.Errorf("unmarshalKey: missing key")
	}
	if b.Type != "EC PRIVATE KEY" {
		return nil, fmt.Errorf("unmarshalKey: found %q, not %q", b.Type, "EC PRIVATE KEY")
	}
	k, err := x509.ParseECPrivateKey(b.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unmarshalKey: %v", err)
	}
	return k, nil
}

func newKey() (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
}

func encodePEM(data interface{}) string {
	b := &pem.Block{}
	switch key := data.(type) {
	case *ecdsa.PrivateKey:
		b.Type = "EC PRIVATE KEY"
		b.Bytes, _ = x509.MarshalECPrivateKey(key)
	case *rsa.PrivateKey:
		b.Type = "RSA PRIVATE KEY"
		b.Bytes = x509.MarshalPKCS1PrivateKey(key)
	case []byte:
		b.Type = "CERTIFICATE"
		b.Bytes = data.([]byte)
	}

	return string(pem.EncodeToMemory(b))
}
