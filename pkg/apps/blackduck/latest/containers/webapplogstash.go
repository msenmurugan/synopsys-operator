/*
Copyright (C) 2019 Synopsys, Inc.

Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements. See the NOTICE file
distributed with this work for additional information
regarding copyright ownership. The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied. See the License for the
specific language governing permissions and limitations
under the License.
*/

package containers

import (
	"fmt"

	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	"github.com/blackducksoftware/horizon/pkg/components"
	apputil "github.com/blackducksoftware/synopsys-operator/pkg/apps/util"
	"github.com/blackducksoftware/synopsys-operator/pkg/util"
)

// GetWebappLogstashDeployment will return the webapp and logstash deployment
func (c *Creater) GetWebappLogstashDeployment(webAppImageName string, logstashImageName string) (*components.Deployment, error) {
	podName := "webapp"
	containerName := "logstash"
	label := fmt.Sprintf("%s-%s", podName, containerName)

	webappEnvs := []*horizonapi.EnvConfig{c.getHubConfigEnv(), c.getHubDBConfigEnv()}
	webappEnvs = append(webappEnvs, &horizonapi.EnvConfig{Type: horizonapi.EnvVal, NameOrPrefix: "HUB_MAX_MEMORY", KeyOrVal: c.hubContainerFlavor.WebappHubMaxMemory})

	webappVolumeMounts := c.getWebappVolumeMounts()

	webappContainerConfig := &util.Container{
		ContainerConfig: &horizonapi.ContainerConfig{Name: podName, Image: webAppImageName,
			PullPolicy: horizonapi.PullAlways, MinMem: c.hubContainerFlavor.WebappMemoryLimit, MaxMem: c.hubContainerFlavor.WebappMemoryLimit, MinCPU: c.hubContainerFlavor.WebappCPULimit},
		EnvConfigs:   webappEnvs,
		VolumeMounts: webappVolumeMounts,
		PortConfig:   []*horizonapi.PortConfig{{ContainerPort: webappPort, Protocol: horizonapi.ProtocolTCP}},
	}

	if c.blackDuck.Spec.LivenessProbes {
		webappContainerConfig.LivenessProbeConfigs = []*horizonapi.ProbeConfig{{
			ActionConfig: horizonapi.ActionConfig{
				Type: horizonapi.ActionTypeCommand,
				Command: []string{
					"/usr/local/bin/docker-healthcheck.sh",
					"https://127.0.0.1:8443/api/health-checks/liveness",
					"/opt/blackduck/hub/hub-webapp/security/root.crt",
					"/opt/blackduck/hub/hub-webapp/security/blackduck_system.crt",
					"/opt/blackduck/hub/hub-webapp/security/blackduck_system.key",
				},
			},
			Delay:           360,
			Interval:        30,
			Timeout:         10,
			MinCountFailure: 1000,
		}}
	}

	logstashVolumeMounts := c.getLogstashVolumeMounts()

	logstashContainerConfig := &util.Container{
		ContainerConfig: &horizonapi.ContainerConfig{Name: containerName, Image: logstashImageName,
			PullPolicy: horizonapi.PullAlways, MinMem: c.hubContainerFlavor.LogstashMemoryLimit, MaxMem: c.hubContainerFlavor.LogstashMemoryLimit, MinCPU: "", MaxCPU: ""},
		EnvConfigs:   []*horizonapi.EnvConfig{c.getHubConfigEnv()},
		VolumeMounts: logstashVolumeMounts,
		PortConfig:   []*horizonapi.PortConfig{{ContainerPort: logstashPort, Protocol: horizonapi.ProtocolTCP}},
	}

	if c.blackDuck.Spec.LivenessProbes {
		logstashContainerConfig.LivenessProbeConfigs = []*horizonapi.ProbeConfig{{
			ActionConfig: horizonapi.ActionConfig{
				Type:    horizonapi.ActionTypeCommand,
				Command: []string{"/usr/local/bin/docker-healthcheck.sh", "http://localhost:9600/"},
			},
			Delay:           240,
			Interval:        30,
			Timeout:         10,
			MinCountFailure: 1000,
		}}
	}
	podConfig := &util.PodConfig{
		Volumes:             c.getWebappLogtashVolumes(),
		Containers:          []*util.Container{webappContainerConfig, logstashContainerConfig},
		Labels:              c.GetVersionLabel(label),
		NodeAffinityConfigs: c.GetNodeAffinityConfigs(label),
		ServiceAccount:      util.GetResourceName(c.blackDuck.Name, util.BlackDuckName, "service-account"),
	}

	if c.blackDuck.Spec.RegistryConfiguration != nil && len(c.blackDuck.Spec.RegistryConfiguration.PullSecrets) > 0 {
		podConfig.ImagePullSecrets = c.blackDuck.Spec.RegistryConfiguration.PullSecrets
	}

	apputil.ConfigurePodConfigSecurityContext(podConfig, c.blackDuck.Spec.SecurityContexts, "blackduck-webapp", c.config.IsOpenshift)

	return util.CreateDeploymentFromContainer(
		&horizonapi.DeploymentConfig{Namespace: c.blackDuck.Spec.Namespace, Name: util.GetResourceName(c.blackDuck.Name, util.BlackDuckName, label), Replicas: util.IntToInt32(1)},
		podConfig, c.GetLabel(label))
}

// getWebappLogtashVolumes will return the webapp and logstash volumes
func (c *Creater) getWebappLogtashVolumes() []*components.Volume {
	webappSecurityEmptyDir, _ := util.CreateEmptyDirVolumeWithoutSizeLimit("dir-webapp-security")
	var webappVolume *components.Volume
	if c.blackDuck.Spec.PersistentStorage {
		webappVolume, _ = util.CreatePersistentVolumeClaimVolume("dir-webapp", c.getPVCName("webapp"))
	} else {
		webappVolume, _ = util.CreateEmptyDirVolumeWithoutSizeLimit("dir-webapp")
	}

	var logstashVolume *components.Volume
	if c.blackDuck.Spec.PersistentStorage {
		logstashVolume, _ = util.CreatePersistentVolumeClaimVolume("dir-logstash", c.getPVCName("logstash"))
	} else {
		logstashVolume, _ = util.CreateEmptyDirVolumeWithoutSizeLimit("dir-logstash")
	}

	volumes := []*components.Volume{webappSecurityEmptyDir, webappVolume, logstashVolume, c.getDBSecretVolume()}
	// Mount the HTTPS proxy certificate if provided
	if len(c.blackDuck.Spec.ProxyCertificate) > 0 {
		volumes = append(volumes, c.getProxyVolume())
	}

	return volumes
}

// getLogstashVolumeMounts will return the Logstash volume mounts
func (c *Creater) getLogstashVolumeMounts() []*horizonapi.VolumeMountConfig {
	volumesMounts := []*horizonapi.VolumeMountConfig{
		{Name: "dir-logstash", MountPath: "/var/lib/logstash/data"},
	}
	return volumesMounts
}

// getWebappVolumeMounts will return the Webapp volume mounts
func (c *Creater) getWebappVolumeMounts() []*horizonapi.VolumeMountConfig {
	volumesMounts := []*horizonapi.VolumeMountConfig{
		{Name: "db-passwords", MountPath: "/tmp/secrets/HUB_POSTGRES_ADMIN_PASSWORD_FILE", SubPath: "HUB_POSTGRES_ADMIN_PASSWORD_FILE"},
		{Name: "db-passwords", MountPath: "/tmp/secrets/HUB_POSTGRES_USER_PASSWORD_FILE", SubPath: "HUB_POSTGRES_USER_PASSWORD_FILE"},
		{Name: "dir-webapp", MountPath: "/opt/blackduck/hub/hub-webapp/ldap"},
		{Name: "dir-webapp-security", MountPath: "/opt/blackduck/hub/hub-webapp/security"},
		{Name: "dir-logstash", MountPath: "/opt/blackduck/hub/logs"},
	}

	// Mount the HTTPS proxy certificate if provided
	if len(c.blackDuck.Spec.ProxyCertificate) > 0 {
		volumesMounts = append(volumesMounts, &horizonapi.VolumeMountConfig{
			Name:      "proxy-certificate",
			MountPath: "/tmp/secrets/HUB_PROXY_CERT_FILE",
			SubPath:   "HUB_PROXY_CERT_FILE",
		})
	}

	return volumesMounts
}

// GetWebAppService will return the webapp service
func (c *Creater) GetWebAppService() *components.Service {
	return util.CreateService(util.GetResourceName(c.blackDuck.Name, util.BlackDuckName, "webapp"), c.GetLabel("webapp-logstash"), c.blackDuck.Spec.Namespace, webappPort, webappPort, horizonapi.ServiceTypeServiceIP, c.GetVersionLabel("webapp-logstash"))
}

// GetLogStashService will return the logstash service
func (c *Creater) GetLogStashService() *components.Service {
	return util.CreateService(util.GetResourceName(c.blackDuck.Name, util.BlackDuckName, "logstash"), c.GetLabel("webapp-logstash"), c.blackDuck.Spec.Namespace, logstashPort, logstashPort, horizonapi.ServiceTypeServiceIP, c.GetVersionLabel("webapp-logstash"))
}
