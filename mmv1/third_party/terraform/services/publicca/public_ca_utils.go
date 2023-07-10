package publicca

import (
	"crypto"
	"encoding/base64"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
	"log"
	"net/url"
)

const (
	GCPDirectoryProduction = "https://dv.acme-v02.api.pki.goog/directory"
	GCPDirectoryStaging    = "https://dv.acme-v02.test-api.pki.goog/directory"
)

type AcmeUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *AcmeUser) GetEmail() string {
	return u.Email
}
func (u *AcmeUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *AcmeUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

func createNewAccountUsingEab(email string, isStagingEnv bool, privateKeyPem string, keyId string, hmacEncoded string) (string, error) {
	baseUrl := getBaseUrl(isStagingEnv)
	user := AcmeUser{
		Email: email,
		key:   getPrivateKey(privateKeyPem),
	}
	config := lego.NewConfig(&user)
	config.CADirURL = baseUrl
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		log.Fatalf("[ERROR] couldn't create acme client: %s", err)
	}
	eabOptions := registration.RegisterEABOptions{
		TermsOfServiceAgreed: true,
		Kid:                  keyId,
		HmacEncoded:          decodeHmacKey(hmacEncoded),
	}
	account, err := client.Registration.RegisterWithExternalAccountBinding(eabOptions)
	if err != nil {
		log.Fatalf("[ERROR] couldn't register account: %s", err)
	}
	log.Printf("[DEBUG] Account created, URL: %s", account.URI)
	return account.URI, nil
}

func deactivateAccount(accountUri string, email string, privateKeyPem string) error {
	user := AcmeUser{
		Email: email,
		key:   getPrivateKey(privateKeyPem),
	}
	config := lego.NewConfig(&user)
	config.CADirURL = getCADirUrlFromAccountUri(accountUri)
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		log.Fatalf("[ERROR] couldn't create acme client: %s", err)
	}
	user.Registration, err = client.Registration.ResolveAccountByKey()
	if err != nil {
		log.Fatalf("[ERROR] couldn't find existing acccount: %s", err)
	}
	err = client.Registration.DeleteRegistration()
	if err != nil {
		log.Fatalf("[ERROR] couldn't deactivate account: %s", err)
	}
	log.Printf("[DEBUG] Account deactivated.")
	return nil
}

func getBaseUrl(isStagingEnv bool) string {
	if isStagingEnv {
		return GCPDirectoryStaging
	}
	return GCPDirectoryProduction
}

func getPrivateKey(privateKeyPem string) crypto.PrivateKey {
	_key, err := certcrypto.ParsePEMPrivateKey([]byte(privateKeyPem))
	if err != nil {
		log.Fatalf("[ERROR] couldn't parse private key: %v", err)
	}
	return _key
}

func decodeHmacKey(hmacB64 string) string {
	result, err := base64.StdEncoding.DecodeString(hmacB64)
	if err != nil {
		log.Fatalf("[ERROR] hmac is not in base64 format: %s", err)
	}
	return string(result)
}

func getCADirUrlFromAccountUri(accountUri string) string {
	parse, err := url.Parse(accountUri)
	if err != nil {
		log.Fatalf("couldn't parse the account URI")
	}
	parse.Path = "/directory"
	log.Printf("[DEBUG] Found CA Directory from Account URI: %s", parse.String())
	return parse.String()
}
