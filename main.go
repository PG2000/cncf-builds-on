package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"time"
	"log"
	"encoding/json"
	"sort"
	"os"
)

type Result struct {
	Repo string
	Url  string
}

type Language struct {
	Key   string
	Value int
}

type Response map[string]interface{}

type LanguageList []Language

func (p LanguageList) Len() int           { return len(p) }
func (p LanguageList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p LanguageList) Less(i, j int) bool { return p[i].Value < p[j].Value }

var some = make(map[string]int)

var ACCESS_TOKEN = ""

func main() {

	ACCESS_TOKEN = os.Getenv("GITHUB_ACCESS_TOKEN")

	if ACCESS_TOKEN == "" {
		panic("No Access Token")
	}

	cncfProjects := []string{
		"Azure/acs-engine",
		"azure/brigade",
		"azure/draft",
		"Microsoft/service-fabric",
		"kubernetes/kube-deploy",
		"GoogleCloudPlatform/skaffold",
		"aws/chalice",
		"awslabs/serverless-application-model",
		"prestodb/presto",
		"alibaba/pouch",
		"clearcontainers/runtime",
		"intelsdi-x/snap",
		"contiv/netplugin",
		"mysql/mysql-server",
		"oracle/terraform-kubernetes-installer",
		"wercker/wercker",
		"Netflix/eureka",
		"Netflix/Hystrix",
		"Netflix/zuul",
		"Netflix/ribbon",
		"spinnaker/spinnaker",
		"vmware/dispatch",
		"vmware/harbor",
		"argoproj/argo",
		"3scale/apicast",
		"ansible/ansible",
		"ceph/ceph",
		"coreos/clair",
		"kubernetes-incubator/bootkube",
		"coreos/etcd",
		"coreos/flannel",
		"gluster/glusterfs",
		"infinispan/infinispan",
		"OpenSCAP/openscap",
		"openshift/origin",
		"projectatomic/atomic",
		"coreos/tectonic-installer",
		"square/keywhiz",
		"airbnb/synapse",
		"lyft/confidant",
		"Yelp/paasta",
		"orientechnologies/orientdb",
		"mongodb/mongo",
		"cyberark/conjur",
		"concourse/concourse",
		"projectriff/riff",
		"spring-cloud/spring-cloud-function",
		"spring-cloud/spring-cloud-sleuth",
		"pinterest/knox",
		"docker/compose",
		"docker/distribution",
		"docker/swarm",
		"docker/infrakit",
		"linuxkit/linuxkit",
		"couchbase/manifest",
		"joyent/containerpilot",
		"joyent/smartos-live",
		"joyent/triton-kubernetes",
		"joyent/smartos-live",
		"joyent/manta",
		"puppetlabs/puppet",
		"chef/chef",
		"habitat-sh/habitat",
		"elastic/elasticsearch",
		"MariaDB/server",
		"antirez/redis",
		"hashicorp/consul",
		"hashicorp/nomad",
		"hashicorp/packer",
		"hashicorp/terraform",
		"hashicorp/vault",
		"influxdata/influxdb",
		"druid-io/druid",
		"apprenda/kismatic",
		"cockroachdb/cockroach",
		"draios/sysdig",
		"bustle/shep",
		"nuclio/nuclio",
		"gitlabhq/gitlabhq",
		"gitlabhq/gitlab-runner",
		"Kong/kong",
		"nginx/nginx",
		"fission/fission",
		"heptio/contour",
		"heptio/aws-quickstart",
		"streamsets/datacollector",
		"rancher/rancher",
		"saltstack/salt",
		"minio/minio",
		"projectcalico/calico",
		"scylladb/scylla",
		"pingcap/tidb",
		"weaveworks/flux",
		"weaveworks/weave",
		"weaveworks/scope",
		"kubernetes/kubeadm",
		"runconduit/conduit",
		"cfengine/core",
		"sensu/sensu",
		"getsentry/sentry",
		"nficano/python-lambda",
		"attic-labs/noms",
		"anchore/anchore-engine",
		"arangodb/arangodb",
		"crate/crate",
		"ManageIQ/manageiq",
		"bigchaindb/bigchaindb",
		"open-io/oio-sds",
		"Graylog2/graylog2-server",
		"serverless/serverless",
		"solo-io/gloo",
		"cloud66/habitus",
		"pachyderm/pachyderm",
		"stdlib/lib",
		"flynn/flynn",
		"containous/traefik",
		"digitalrebar/provision",
		"magneticio/vamp",
		"datawire/ambassador",
		"datawire/telepresence",
		"drone/drone",
		"apache/apex-core",
		"apache/brooklyn-server",
		"apache/incubator-heron",
		"apache/nifi",
		"apache/incubator-openwhisk",
		"apache/rocketmq",
		"apache/spark",
		"apache/storm",
		"apache/thrift",
		"apache/zookeeper",
		"apex/apex",
		"AppScale/appscale",
		"arc-repos/arc-functions",
		"apache/avro",
		"cloudfoundry/bosh",
		"buildkite/agent",
		"juju-solutions/bundle-canonical-kubernetes",
		"apache/carbondata",
		"apache/cassandra",
		"centreon/centreon",
		"cilium/cilium",
		"claudiajs/claudia",
		"cloudfoundry/cli",
		"cloudfoundry-incubator/kubo-deployment",
		"cloudevents/spec",
		"cloudify-cosmo/cloudify-manager",
		"containernetworking/cni",
		"container-storage-interface/spec",
		"containerd/containerd",
		"Huawei/containerops",
		"coredns/coredns",
		"kubernetes-incubator/cri-o",
		"envoyproxy/envoy",
		"apache/flink",
		"TIBCOSoftware/flogo",
		"fluent/fluentd",
		"fnproject/fn",
		"theforeman/foreman",
		"metrue/fx",
		"gocd/gocd",
		"grafana/grafana",
		"graphite-project/graphite-web",
		"grpc/grpc",
		"Miserlou/Zappa",
		"apache/hadoop",
		"haproxy/haproxy",
		"kubernetes/helm",
		"Icinga/icinga2",
		"inwinstack/kube-ansible",
		"istio/istio",
		"jaegertracing/jaeger",
		"jenkinsci/jenkins",
		"jhipster/generator-jhipster",
		"juju/juju",
		"apache/kafka",
		"kata-containers/runtime",
		"kontena/kontena",
		"kontena/pharos-cluster",
		"kinvolk/kube-spawn",
		"kubeless/kubeless",
		"kubernetes/kubernetes",
		"kubevirt/kubevirt",
		"kubicorn/kubicorn",
		"leo-project/leofs",
		"linkerd/linkerd",
		"lxc/lxd",
		"mantl/mantl",
		"apache/mesos",
		"midonet/midonet",
		"kubernetes/minikube",
		"NagiosEnterprises/nagioscore",
		"nats-io/gnatsd",
		"motdotla/node-lambda",
		"theupdateframework/notary",
		"open-policy-agent/opa",
		"openservicebrokerapi/servicebroker",
		"openvswitch/ovs",
		"OAI/OpenAPI-Specification",
		"Juniper/contrail-controller",
		"openebs/openebs",
		"openfaas/faas",
		"open-lambda/open-lambda",
		"openmessaging/openmessaging-java",
		"opensds/opensds",
		"openstack/openstack",
		"openstack-infra/zuul",
		"opentracing/opentracing-go",
		"OpenTSDB/opentsdb",
		"pilosa/pilosa",
		"portainer/portainer",
		"SUSE/Portus",
		"postgres/postgres",
		"prometheus/prometheus",
		"apache/incubator-pulsar",
		"rabbitmq/rabbitmq-server",
		"rethinkdb/rethinkdb",
		"codedellemc/rexray",
		"rkt/rkt",
		"romana/romana",
		"rook/rook",
		"opencontainers/runc",
		"rundeck/rundeck",
		"hyperhq/runv",
		"samsung-cnct/kraken",
		"sheepdog/sheepdog",
		"singularityware/singularity",
		"apache/incubator-skywalking",
		"mweagle/Sparta",
		"spiffe/spiffe",
		"spiffe/spire",
		"supergiant/supergiant",
		"openstack/swift",
		"theupdateframework/tuf",
		"travis-ci/travis-web",
		"tsuru/tsuru",
		"poseidon/typhoon",
		"vitessio/vitess",
		"openzipkin/zipkin",
	}

	projectLength := len(cncfProjects)

	fmt.Printf("Scanning %d projects\n", projectLength)

	ch := make(chan Response)

	list := createList(cncfProjects)

	for range cncfProjects {
		result := <-list
		go makeRequest(result.Url, ch)
	}

	for range cncfProjects {
		language := <-ch

		for i, _ := range language {
			if _, ok := some[i]; ok {
				some[i] += 1
			} else {
				some[i] = 1
			}
		}

	}

	p := make(LanguageList, len(some))

	sortPairListByValue(p)

	for _, v := range p {
		percent := (float64(v.Value) / float64(projectLength)) * 100
		fmt.Printf("Lang: %s, Value: %.d Percentage: %0.2f\n", v.Key, v.Value, percent)
	}

}
func sortPairListByValue(p LanguageList) {
	i := 0
	for k, v := range some {
		p[i] = Language{k, v}
		i++
	}
	sort.Sort(p)
}

func createList(cncfProjects [] string) <-chan Result {
	resultChannel := make(chan Result, len(cncfProjects))

	go func() {
		for _, url := range cncfProjects {
			resultChannel <- Result{url, getUrlForRepo(url)}
		}
	}()

	return resultChannel
}

func makeRequest(url string, ch chan<- Response) {
	fmt.Printf("Make request for %10s\n", url)
	spaceClient := http.Client{
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Authorization", "token "+ACCESS_TOKEN)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "mozilla")

	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	var response Response
	json.Unmarshal([]byte(body), &response)

	ch <- response
}

func getUrlForRepo(repo string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/languages", repo)
}
