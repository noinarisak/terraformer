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

type AuthServerGenerator struct {
	OktaService
}

func (g AuthServerGenerator) createResources(authServerList []*okta.AuthorizationServer) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for _, authServer := range authServerList {

		resources = append(resources, terraformutils.NewSimpleResource(
			authServer.Id,
			authServer.Name,
			"okta_auth_server",
			"okta",
			[]string{}))
	}
	return resources
}

func (g *AuthServerGenerator) InitResources() error {
	ctx, client, e := g.generateClient()
	if e != nil {
		return e
	}

	output, resp, err := client.AuthorizationServer.ListAuthorizationServers(ctx, nil)
	if err != nil {
		return e
	}

	for resp.HasNextPage() {
		var nextAuthServerSet []*okta.AuthorizationServer
		resp, err = resp.Next(ctx, &nextAuthServerSet)
		output = append(output, nextAuthServerSet...)
	}

	g.Resources = g.createResources(output)
	return nil
}
