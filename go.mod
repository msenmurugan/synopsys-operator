module github.com/blackducksoftware/synopsys-operator

go 1.13

require (
	github.com/Azure/go-autorest/autorest/adal v0.8.2 // indirect
	github.com/DATA-DOG/go-sqlmock v1.4.1 // indirect
	github.com/blackducksoftware/horizon v0.0.0-20190625151958-16cafa9109a3
	github.com/docker/spdystream v0.0.0-20181023171402-6480d4af844c // indirect
	github.com/elazarl/goproxy v0.0.0-20190911111923-ecfe977594f1 // indirect
	github.com/gin-gonic/gin v1.4.0
	github.com/gobuffalo/packr v1.30.1
	github.com/google/go-cmp v0.4.0
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/googleapis/gnostic v0.3.1 // indirect
	github.com/gophercloud/gophercloud v0.3.0 // indirect
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.3
	github.com/imdario/mergo v0.3.7
	github.com/jmoiron/sqlx v1.2.0 // indirect
	github.com/juju/errors v0.0.0-20190806202954-0232dcc7464d
	github.com/juju/loggo v0.0.0-20190526231331-6e530bcce5d8 // indirect
	github.com/juju/testing v0.0.0-20190723135506-ce30eb24acd2 // indirect
	github.com/lib/pq v1.2.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/onsi/ginkgo v1.10.1
	github.com/onsi/gomega v1.7.0
	github.com/openshift/api v0.0.0-20200217161739-c99157bc6492
	github.com/openshift/client-go v0.0.0-20200116152001-92a2713fa240
	github.com/pkg/errors v0.9.1
	github.com/rubenv/sql-migrate v0.0.0-20200212082348-64f95ea68aa3 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2 // indirect
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22 // indirect
	helm.sh/helm v2.16.3+incompatible
	helm.sh/helm/v3 v3.1.0
	k8s.io/api v0.17.3
	k8s.io/apiextensions-apiserver v0.17.2
	k8s.io/apimachinery v0.17.3
	k8s.io/client-go v0.17.2
	k8s.io/helm v2.16.3+incompatible // indirect
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/yaml v1.1.0
)

replace (
	github.com/blackducksoftware/horizon => github.com/blackducksoftware/horizon v0.0.0-20190625151958-16cafa9109a3
	github.com/docker/spdystream => github.com/docker/spdystream v0.0.0-20181023171402-6480d4af844c // indirect
	github.com/evanphx/json-patch => github.com/evanphx/json-patch v4.5.0+incompatible // indirect
	github.com/gin-gonic/gin => github.com/gin-gonic/gin v1.4.0
	github.com/golang/mock => github.com/golang/mock v1.2.0 // indirect
	github.com/google/go-cmp => github.com/google/go-cmp v0.3.0
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.3.0 // indirect
	github.com/gophercloud/gophercloud => github.com/gophercloud/gophercloud v0.3.0 // indirect
	github.com/imdario/mergo => github.com/imdario/mergo v0.3.7
	github.com/juju/errors => github.com/juju/errors v0.0.0-20190806202954-0232dcc7464d
	github.com/lib/pq => github.com/lib/pq v1.2.0
	github.com/mitchellh/go-homedir => github.com/mitchellh/go-homedir v1.1.0
	github.com/onsi/ginkgo => github.com/onsi/ginkgo v1.7.0
	github.com/onsi/gomega => github.com/onsi/gomega v1.4.3
	github.com/sirupsen/logrus => github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra => github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag => github.com/spf13/pflag v1.0.3
	github.com/spf13/viper => github.com/spf13/viper v1.4.0
	github.com/stretchr/testify => github.com/stretchr/testify v1.3.0
	gopkg.in/inf.v0 => gopkg.in/inf.v0 v0.9.1 // indirect
	sigs.k8s.io/yaml => sigs.k8s.io/yaml v1.1.0
)
