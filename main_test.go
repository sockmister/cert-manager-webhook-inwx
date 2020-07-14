package main

import (
	logf "github.com/jetstack/cert-manager/pkg/logs"
	"github.com/jetstack/cert-manager/test/acme/dns"
	"github.com/jetstack/cert-manager/test/acme/dns/server"
	"gitlab.com/smueller18/cert-manager-webhook-inwx/test"
	"io/ioutil"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"log"
	"os"
	"testing"
	"time"
)

var (
	zone = "smueller18.de."
	fqdn string
)

func TestRunSuite(t *testing.T) {

	if os.Getenv("TEST_ZONE_NAME") != "" {
		zone = os.Getenv("TEST_ZONE_NAME")
	}
	fqdn = "cert-manager-dns01-tests." + zone

	ctx := logf.NewContext(nil, nil, t.Name())

	srv := &server.BasicServer{
		Handler: &test.Handler{
			Log: logf.FromContext(ctx, "dnsBasicServer"),
			TxtRecords: map[string][][]string{
				fqdn: {
					{},
					{},
					{"123d=="},
					{"123d=="},
				},
			},
			Zones: []string{zone},
		},
	}

	if err := srv.Run(ctx); err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer srv.Shutdown()

	d, err := ioutil.ReadFile("testdata/config.json")
	if err != nil {
		log.Fatal(err)
	}

	fixture := dns.NewFixture(&solver{},
		dns.SetResolvedZone(zone),
		dns.SetResolvedFQDN(fqdn),
		dns.SetAllowAmbientCredentials(false),
		dns.SetDNSServer(srv.ListenAddr()),
		dns.SetBinariesPath("kubebuilder/bin"),
		dns.SetPropagationLimit(time.Duration(60)*time.Second),
		dns.SetUseAuthoritative(false),
		// Set to false because INWX implementation deletes all records
		dns.SetStrict(false),
		dns.SetConfig(&extapi.JSON{
			Raw: d,
		}),
	)

	fixture.RunConformance(t)
}

func TestRunSuiteWithSecret(t *testing.T) {

	if os.Getenv("TEST_ZONE_NAME") != "" {
		zone = os.Getenv("TEST_ZONE_NAME")
	}
	fqdn = "cert-manager-dns01-tests-with-secret." + zone

	ctx := logf.NewContext(nil, nil, t.Name())

	srv := &server.BasicServer{
		Handler: &test.Handler{
			Log: logf.FromContext(ctx, "dnsBasicServerSecret"),
			TxtRecords: map[string][][]string{
				fqdn: {
					{},
					{},
					{"123d=="},
					{"123d=="},
				},
			},
			Zones: []string{zone},
		},
	}

	if err := srv.Run(ctx); err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer srv.Shutdown()

	d, err := ioutil.ReadFile("testdata/config.secret.json")
	if err != nil {
		log.Fatal(err)
	}

	fixture := dns.NewFixture(&solver{},
		dns.SetResolvedZone(zone),
		dns.SetResolvedFQDN(fqdn),
		dns.SetAllowAmbientCredentials(false),
		dns.SetDNSServer(srv.ListenAddr()),
		dns.SetManifestPath("testdata/secret-inwx-credentials.yaml"),
		dns.SetBinariesPath("kubebuilder/bin"),
		dns.SetPropagationLimit(time.Duration(60)*time.Second),
		dns.SetUseAuthoritative(false),
		dns.SetConfig(&extapi.JSON{
			Raw: d,
		}),
	)

	fixture.RunConformance(t)
}

