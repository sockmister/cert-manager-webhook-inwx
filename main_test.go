package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	logf "github.com/cert-manager/cert-manager/pkg/logs"
	"github.com/cert-manager/cert-manager/test/acme/dns"
	"github.com/cert-manager/cert-manager/test/acme/dns/server"
	"github.com/sockmister/cert-manager-webhook-inwx/test"
	extapi "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

var (
	zone      = "smueller18.de."
	zoneTwoFA = "smueller18mfa.de."
	fqdn      string
)

func TestRunSuite(t *testing.T) {

	if os.Getenv("TEST_ZONE_NAME") != "" {
		zone = os.Getenv("TEST_ZONE_NAME")
	}
	fqdn = "cert-manager-dns01-tests." + zone

	ctx := logf.NewContext(context.TODO(), logf.Log, t.Name())

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

	ctx := logf.NewContext(context.TODO(), logf.Log, t.Name())

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
		dns.SetPropagationLimit(time.Duration(60)*time.Second),
		dns.SetUseAuthoritative(false),
		dns.SetConfig(&extapi.JSON{
			Raw: d,
		}),
	)

	fixture.RunConformance(t)
}

func TestRunSuiteWithTwoFA(t *testing.T) {

	if os.Getenv("TEST_ZONE_NAME_WITH_TWO_FA") != "" {
		zoneTwoFA = os.Getenv("TEST_ZONE_NAME_WITH_TWO_FA")
	}

	fqdn = "cert-manager-dns01-tests." + zoneTwoFA

	ctx := logf.NewContext(nil, logf.Log, t.Name())

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
			Zones: []string{zoneTwoFA},
		},
	}

	if err := srv.Run(ctx); err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer srv.Shutdown()

	d, err := ioutil.ReadFile("testdata/config-otp.json")
	if err != nil {
		log.Fatal(err)
	}

	fixture := dns.NewFixture(&solver{},
		dns.SetResolvedZone(zoneTwoFA),
		dns.SetResolvedFQDN(fqdn),
		dns.SetAllowAmbientCredentials(false),
		dns.SetDNSServer(srv.ListenAddr()),
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

func TestRunSuiteWithSecretAndTwoFA(t *testing.T) {

	if os.Getenv("TEST_ZONE_NAME_WITH_TWO_FA") != "" {
		zoneTwoFA = os.Getenv("TEST_ZONE_NAME_WITH_TWO_FA")
	}
	fqdn = "cert-manager-dns01-tests-with-secret." + zoneTwoFA

	ctx := logf.NewContext(nil, logf.Log, t.Name())

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
			Zones: []string{zoneTwoFA},
		},
	}

	if err := srv.Run(ctx); err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	defer srv.Shutdown()

	d, err := ioutil.ReadFile("testdata/config-otp.secret.json")
	if err != nil {
		log.Fatal(err)
	}

	fixture := dns.NewFixture(&solver{},
		dns.SetResolvedZone(zoneTwoFA),
		dns.SetResolvedFQDN(fqdn),
		dns.SetAllowAmbientCredentials(false),
		dns.SetDNSServer(srv.ListenAddr()),
		dns.SetManifestPath("testdata/secret-inwx-credentials-otp.yaml"),
		dns.SetPropagationLimit(time.Duration(60)*time.Second),
		dns.SetUseAuthoritative(false),
		dns.SetConfig(&extapi.JSON{
			Raw: d,
		}),
	)

	fixture.RunConformance(t)
}
