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

type PolicyRuleGenerator struct {
	OktaService
}

func (g PolicyRuleGenerator) createResources(policyRuleList []*okta.PolicyRule, policyList []*okta.Policy) []terraformutils.Resource {
	var resources []terraformutils.Resource

	for _, policy := range policyList {
		for _, policyRule := range policyRuleList {
			resourceOktaPolicyRuleType := ""

			switch policyRule.Type {
			case "MFA_ENROLL":
				resourceOktaPolicyRuleType = "okta_policy_rule_mfa"
			case "PASSWORD":
				resourceOktaPolicyRuleType = "okta_policy_rule_password"
			case "SIGN_ON":
				resourceOktaPolicyRuleType = "okta_policy_rule_signon"
			}

			//TODO: GH https://github.com/okta/okta-sdk-golang/issues/197. PolicyRule does not have 'name' field to use.
			resources = append(resources, terraformutils.NewResource(
				policyRule.Id,
				policyRule.Id,
				resourceOktaPolicyRuleType,
				"okta",
				map[string]string{
					"policyId": policy.Id,
				},
				[]string{},
				map[string]interface{}{},
			))

		}
	}

	// //TODO: We need the policyId
	// for _, policyRule := range policyRuleList {
	// 	resourceOktaPolicyRuleType := ""

	// 	switch policyRule.Type {
	// 	case "MFA_ENROLL":
	// 		resourceOktaPolicyRuleType = "okta_policy_rule_mfa"
	// 	case "PASSWORD":
	// 		resourceOktaPolicyRuleType = "okta_policy_rule_password"
	// 	case "SIGN_ON":
	// 		resourceOktaPolicyRuleType = "okta_policy_rule_signon"
	// 	}

	// 	resources = append(resources, terraformutils.NewSimpleResource(
	// 		policyRule.Id,
	// 		policyRule.Id,
	// 		resourceOktaPolicyRuleType,
	// 		"okta",
	// 		[]string{}))

	// }
	return resources
}

func (g *PolicyRuleGenerator) InitResources() error {
	// policyTypeList := [...]string{"PASSWORD", "MFA_ENROLL", "OKTA_SIGN_ON"}

	// var policies = []*okta.Policy{}
	// var policyRules = []*okta.PolicyRule{}

	// for _, policyType := range policyTypeList {
	// 	policySet, err := getAllPolicies(g, policyType)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	for _, policy := range policySet {
	// 		policyRuleSet, err := getPolicyRules(g, policy.Id)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		policyRules = append(policyRules, policyRuleSet...)
	// 	}

	// 	policies = append(policies, policySet...)
	// }

	// g.Resources = g.createResources(policyRules, policies)
	return nil
}

func getPolicyRules(g *PolicyRuleGenerator, policyID string) ([]*okta.PolicyRule, error) {
	ctx, client, e := g.generateClient()
	if e != nil {
		return nil, e
	}

	output, resp, err := client.Policy.ListPolicyRules(ctx, policyID)
	if err != nil {
		return nil, e
	}

	for resp.HasNextPage() {
		var nextPolicySet []*okta.PolicyRule
		resp, err = resp.Next(ctx, &nextPolicySet)
		output = append(output, nextPolicySet...)
	}

	return output, nil
}

//NOTE: Code smell, this impl is also in policy.go
func getAllPolicies(g *PolicyRuleGenerator, policyType string) ([]*okta.Policy, error) {
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
