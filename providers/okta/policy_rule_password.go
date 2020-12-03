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

	"github.com/okta/terraform-provider-okta/sdk"
)

type PolicyRulePasswordGenerator struct {
	OktaService
}

func (g PolicyRulePasswordGenerator) createResources(policyRuleList []sdk.PolicyRule, policyID string) []terraformutils.Resource {
	var resources []terraformutils.Resource

	//TODO: Implement filter and consolidate all policy types and rules. Look at aws as example.

	//TODO: Complete the call to the list listAllPolices and use the policyID and the policyRuleID.

	// for _, policy := range policyList {
	// 	for _, policyRule := range policyRuleList {
	// 		resourceOktaPolicyRuleType := ""

	// 		switch policyRule.Type {
	// 		case "MFA_ENROLL":
	// 			resourceOktaPolicyRuleType = "okta_policy_rule_mfa"
	// 		case "PASSWORD":
	// 			resourceOktaPolicyRuleType = "okta_policy_rule_password"
	// 		case "SIGN_ON":
	// 			resourceOktaPolicyRuleType = "okta_policy_rule_signon"
	// 		}

	// 		//TODO: GH https://github.com/okta/okta-sdk-golang/issues/197. PolicyRule does not have 'name' field to use.
	// 		resources = append(resources, terraformutils.NewResource(
	// 			policyRule.Id,
	// 			policyRule.Name,
	// 			resourceOktaPolicyRuleType,
	// 			"okta",
	// 			map[string]string{},
	// 			[]string{},
	// 			map[string]interface{}{},
	// 		))

	// 	}
	// }

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
		//NOTE: Decided to use the terraform-provider-okta/sdk implementation version because the 'name' field and cast coversion are supported.

		resources = append(resources, terraformutils.NewResource(
			policyRule.Id,
			policyRule.Name,
			resourceOktaPolicyRuleType,
			"okta",
			map[string]string{
				"policyid": policyID,
			},
			[]string{},
			map[string]interface{}{},
		))
	}

	return resources
}

func (g *PolicyRulePasswordGenerator) InitResources() error {
	// policyTypeList := [...]string{"PASSWORD"}
	// var policies = []*okta.Policy{}
	// var policyRules = []sdk.PolicyRule{}
	var resources []terraformutils.Resource

	policySet, err := getPoliciesPassword(g, "PASSWORD")
	if err != nil {
		return err
	}

	for _, policy := range policySet {
		policyRuleSet, err := getPolicyRulesPassword(g, policy.Id)
		if err != nil {
			return err
		}

		// policyRules = append(policyRules, policyRuleSet...)

		for _, policyRule := range policyRuleSet {
			resources = append(resources, terraformutils.NewResource(
				policyRule.Id,
				policyRule.Name,
				"okta_policy_rule_password",
				"okta",
				map[string]string{
					"policyid": policy.Id,
				},
				[]string{},
				map[string]interface{}{},
			))
		}

	}

	// for _, policyType := range policyTypeList {
	// 	policySet, err := getPoliciesPassword(g, policyType)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	for _, policy := range policySet {
	// 		policyRuleSet, err := getPolicyRulesPassword(g, policy.Id)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		policyRules = append(policyRules, policyRuleSet...)

	// 		resources = g.createResources(policyRules, policy.Id)
	// 		resources = append(resources, resources...)
	// 	}
	// }

	g.Resources = resources
	return nil
}

func getPolicyRulesPassword(g *PolicyRulePasswordGenerator, policyID string) ([]sdk.PolicyRule, error) {
	ctx, client, e := g.generateAPISupplementClient()
	if e != nil {
		return nil, e
	}

	output, resp, err := client.ListPolicyRules(ctx, policyID)
	if err != nil {
		return nil, e
	}

	for resp.HasNextPage() {
		var nextPolicySet []sdk.PolicyRule
		resp, err = resp.Next(ctx, &nextPolicySet)
		output = append(output, nextPolicySet...)
	}

	return output, nil
}

//NOTE: Code smell, this impl is also in policy.go
func getPoliciesPassword(g *PolicyRulePasswordGenerator, policyType string) ([]*okta.Policy, error) {
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
