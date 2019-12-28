module gitlab.com/smueller18/cert-manager-webhook-inwx

go 1.13

require (
	github.com/go-logr/logr v0.1.0
	github.com/jetstack/cert-manager v0.12.0
	github.com/miekg/dns v0.0.0-20170721150254-0f3adef2e220
	github.com/nrdcg/goinwx v0.6.1
	google.golang.org/appengine v1.5.0
	k8s.io/apiextensions-apiserver v0.0.0-20191114105449-027877536833
	k8s.io/apimachinery v0.0.0-20191028221656-72ed19daf4bb
	k8s.io/client-go v0.0.0-20191114101535-6c5935290e33
	k8s.io/component-base v0.0.0-20191114102325-35a9586014f7
	k8s.io/klog v0.4.0
)

replace github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.4
