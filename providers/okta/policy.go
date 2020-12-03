// Copyright 2021 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package okta

import (
	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/okta/okta-sdk-golang/v2/okta"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

type PolicyGenerator struct {
	OktaService
}

func (g PolicyGenerator) createResources(policyList []*okta.Policy) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for _, policy := range policyList {
		resourceOktaPolicyType := ""

		switch policy.Type {
		case "MFA_ENROLL":
			resourceOktaPolicyType = "okta_policy_mfa"
		case "PASSWORD":
			resourceOktaPolicyType = "okta_policy_password"
		case "OKTA_SIGN_ON":
			resourceOktaPolicyType = "okta_policy_signon"
		}

		resources = append(resources, terraformutils.NewSimpleResource(
			policy.Id,
			policy.Name,
			resourceOktaPolicyType,
			"okta",
			[]string{}))
	}
	return resources
}

// Generate Terraform Resources from Okta API,
func (g *PolicyGenerator) InitResources() error {
	policyTypeList := [...]string{"PASSWORD", "MFA_ENROLL", "OKTA_SIGN_ON"}
	var output = []*okta.Policy{}

	for _, policyType := range policyTypeList {
		policySet, err := getPolicies(g, policyType)
		if err != nil {
			return err
		}
		output = append(output, policySet...)
	}

	g.Resources = g.createResources(output)
	return nil
}

func getPolicies(g *PolicyGenerator, policyType string) ([]*okta.Policy, error) {
	ctx, client, e := g.generateClient()
	if e != nil {
		return nil, e
	}

	qp := query.NewQueryParams(query.WithType(policyType))
	output, resp, err := client.Policy.ListPolicies(ctx, qp)
	if err != nil {
		return nil, e
	}

	for resp.HasNextPage() {
		var nextPolicySet []*okta.Policy
		resp, err = resp.Next(ctx, &nextPolicySet)
		output = append(output, nextPolicySet...)
	}

	return output, nil
}
