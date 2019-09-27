/*
 * Copyright (C) 2019 Synopsys, Inc.
 *
 *  Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 *  with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 *  under the License.
 */

package polaris

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"strings"
)

func GetPolarisProvisionComponents(baseUrl string, polarisConf Polaris) (map[string]runtime.Object, error) {
	content, err := GetBaseYaml(baseUrl, "polaris", polarisConf.Version, "organization-provision-job.yaml")
	if err != nil {
		return nil, err
	}

	// regex patching
	content = strings.ReplaceAll(content, "${NAMESPACE}", polarisConf.Namespace)
	content = strings.ReplaceAll(content, "${ENVIRONMENT_NAME}", polarisConf.Namespace)
	content = strings.ReplaceAll(content, "${POLARIS_ROOT_DOMAIN}", polarisConf.EnvironmentDNS)
	content = strings.ReplaceAll(content, "${IMAGE_PULL_SECRETS}", polarisConf.ImagePullSecrets)
	content = strings.ReplaceAll(content, "${ORG_DESCRIPTION}", polarisConf.OrganizationDetails.OrganizationProvisionOrganizationDescription)
	content = strings.ReplaceAll(content, "${ORG_NAME}", polarisConf.OrganizationDetails.OrganizationProvisionOrganizationName)
	content = strings.ReplaceAll(content, "${ADMIN_NAME}", polarisConf.OrganizationDetails.OrganizationProvisionAdminName)
	content = strings.ReplaceAll(content, "${ADMIN_USERNAME}", polarisConf.OrganizationDetails.OrganizationProvisionAdminUsername)
	content = strings.ReplaceAll(content, "${ADMIN_EMAIL}", polarisConf.OrganizationDetails.OrganizationProvisionAdminEmail)
	content = strings.ReplaceAll(content, "${SEAT_COUNT}", polarisConf.OrganizationDetails.OrganizationProvisionLicenseSeatCount)
	content = strings.ReplaceAll(content, "${TYPE}", polarisConf.OrganizationDetails.OrganizationProvisionLicenseType)
	content = strings.ReplaceAll(content, "${RESULTS_START_DATE}", polarisConf.OrganizationDetails.OrganizationProvisionResultsStartDate)
	content = strings.ReplaceAll(content, "${RESULTS_END_DATE}", polarisConf.OrganizationDetails.OrganizationProvisionResultsEndDate)
	content = strings.ReplaceAll(content, "${RETENTION_START_DATE}", polarisConf.OrganizationDetails.OrganizationProvisionRetentionStartDate)
	content = strings.ReplaceAll(content, "${RETENTION_END_DATE}", polarisConf.OrganizationDetails.OrganizationProvisionRetentionEndDate)

	if len(polarisConf.Repository) != 0 {
		content = strings.ReplaceAll(content, "gcr.io/snps-swip-staging", polarisConf.Repository)
	}

	mapOfUniqueIdToBaseRuntimeObject := ConvertYamlFileToRuntimeObjects(content)

	mapOfUniqueIdToBaseRuntimeObject["Secret.coverity-license"] = &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "coverity-license",
			Namespace: polarisConf.Namespace,
			Labels: map[string]string{
				"environment": polarisConf.Namespace,
			},
		},
		Data: map[string][]byte{
			"license": []byte(polarisConf.Licenses.Coverity),
		},
		Type: corev1.SecretTypeOpaque,
	}

	return mapOfUniqueIdToBaseRuntimeObject, nil
}
