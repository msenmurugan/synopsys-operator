package v1

import (
	"fmt"

	"github.com/blackducksoftware/horizon/pkg/components"
	blackduckapi "github.com/blackducksoftware/synopsys-operator/pkg/api/blackduck/v1"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/blackduck/components/rc/utils"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/database/postgres"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/store"
	"github.com/blackducksoftware/synopsys-operator/pkg/apps/types"
	apputils "github.com/blackducksoftware/synopsys-operator/pkg/apps/utils"
	"github.com/blackducksoftware/synopsys-operator/pkg/protoform"
	"github.com/blackducksoftware/synopsys-operator/pkg/util"
	"k8s.io/client-go/kubernetes"
)

// BdReplicationController holds the Black Duck RC configuration
type BdReplicationController struct {
	*types.ReplicationController
	config     *protoform.Config
	kubeClient *kubernetes.Clientset
	blackDuck  *blackduckapi.Blackduck
}

func init() {
	store.Register(blackduck.BlackDuckPostgresRCV1, NewBdReplicationController)
}

// GetRc returns the RC
func (c *BdReplicationController) GetRc() (*components.ReplicationController, error) {
	containerConfig, ok := c.Containers[blackduck.PostgresContainerName]
	if !ok {
		return nil, fmt.Errorf("couldn't find container %s", blackduck.PostgresContainerName)
	}

	name := apputils.GetResourceName(c.blackDuck.Name, util.BlackDuckName, "postgres")

	var pvcName string
	if c.blackDuck.Spec.PersistentStorage {
		pvcName = utils.GetPVCName("postgres", c.blackDuck)
	}

	p := &postgres.Postgres{
		Name:                   name,
		Namespace:              c.blackDuck.Spec.Namespace,
		PVCName:                pvcName,
		Port:                   int32(5432),
		Image:                  containerConfig.Image,
		MinCPU:                 util.Int32ToInt(containerConfig.MinCPU),
		MaxCPU:                 util.Int32ToInt(containerConfig.MaxCPU),
		MinMemory:              util.Int32ToInt(containerConfig.MinMem),
		MaxMemory:              util.Int32ToInt(containerConfig.MaxMem),
		Database:               "blackduck",
		User:                   "blackduck",
		PasswordSecretName:     apputils.GetResourceName(c.blackDuck.Name, util.BlackDuckName, "db-creds"),
		UserPasswordSecretKey:  "HUB_POSTGRES_ADMIN_PASSWORD_FILE",
		AdminPasswordSecretKey: "HUB_POSTGRES_POSTGRES_PASSWORD_FILE",
		MaxConnections:         300,
		SharedBufferInMB:       1024,
		EnvConfigMapRefs:       []string{apputils.GetResourceName(c.blackDuck.Name, util.BlackDuckName, "db-config")},
		Labels:                 apputils.GetLabel("postgres", c.blackDuck.Name),
		IsOpenshift:            c.config.IsOpenshift,
	}

	return p.GetPostgresReplicationController()
}

// NewBdReplicationController returns the Black Duck RC configuration
func NewBdReplicationController(replicationController *types.ReplicationController, config *protoform.Config, kubeClient *kubernetes.Clientset, cr interface{}) (types.ReplicationControllerInterface, error) {
	blackDuck, ok := cr.(*blackduckapi.Blackduck)
	if !ok {
		return nil, fmt.Errorf("unable to cast the interface to Black Duck object")
	}
	return &BdReplicationController{ReplicationController: replicationController, config: config, kubeClient: kubeClient, blackDuck: blackDuck}, nil
}
