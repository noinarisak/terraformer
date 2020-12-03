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
)

type IDPGenerator struct {
	OktaService
}

func (g IDPGenerator) createResources(idpList []*okta.IdentityProvider) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for _, idp := range idpList {

		resourceOktaIdpType := ""

		switch idp.Type {
		case "OIDC":
			{
				resourceOktaIdpType = "oidc"
			}
		case "SAML2":
			{
				resourceOktaIdpType = "saml"
			}
		case "APPLE", "FACEBOOK", "GOOGLE", "LINKEDIN", "MICROSOFT":
			{
				resourceOktaIdpType = "social"
			}
		}

		resources = append(resources, terraformutils.NewSimpleResource(
			idp.Id,
			idp.Name,
			"okta_idp_"+resourceOktaIdpType,
			"okta",
			[]string{}))
	}
	return resources
}

// Generate Terraform Resources from Okta API,
func (g *IDPGenerator) InitResources() error {
	ctx, client, e := g.generateClient()
	if e != nil {
		return e
	}

	output, resp, err := client.IdentityProvider.ListIdentityProviders(ctx, nil)
	if err != nil {
		return e
	}

	for resp.HasNextPage() {
		var nextIDPSet []*okta.IdentityProvider
		resp, err = resp.Next(ctx, &nextIDPSet)
		output = append(output, nextIDPSet...)
	}

	g.Resources = g.createResources(output)
	return nil
}
