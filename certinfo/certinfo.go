package certinfo

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"strings"
)

func GetCertsInfo(URLs string, printFullChain bool) string {
	UrlArr := strings.Split(URLs, " ")
	result := ""
	for _, url := range UrlArr {
		certStr, _, err := GetCertInfo(url, printFullChain)
		if err != nil {
			result += err.Error()
		}
		result += certStr
	}
	return result
}

func GetCertInfo(URL string, printFullChain bool) (string, []*x509.Certificate, error) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", URL+":443", conf)
	if err != nil {
		log.Println("Error in Dial", err)
		return "", nil, fmt.Errorf("check certificate error - cannot check cert from URL %s. Error: %e\n\n", URL, err)
	}
	defer conn.Close()
	certs := conn.ConnectionState().PeerCertificates
	result := ""
	var certsResult []*x509.Certificate
	for _, cert := range certs {
		if !printFullChain && cert.IsCA {
			continue
		}
		certsResult = append(certsResult, cert)
		result += fmt.Sprintf("Check certificate for domain: %s\n", URL)
		result += fmt.Sprintf("DNSNames: %s\n", cert.DNSNames)
		result += fmt.Sprintf("Issuer Name: %s\n", cert.Issuer)
		result += fmt.Sprintf("Expiry: %s\n", cert.NotAfter.Format("2006-01-02"))
		result += fmt.Sprintf("Common Name: %s\n", cert.Issuer.CommonName)

	}
	return result + "\n", certsResult, nil
}
