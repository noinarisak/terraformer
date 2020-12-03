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

type GroupGenerator struct {
	OktaService
}

func (g GroupGenerator) createResources(groupList []*okta.Group) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for _, group := range groupList {
		resources = append(resources, terraformutils.NewSimpleResource(
			group.Id,
			group.Id,
			"okta_group",
			"okta",
			[]string{}))
	}
	return resources
}

func (g *GroupGenerator) InitResources() error {
	ctx, client, e := g.generateClient()
	if e != nil {
		return e
	}

	output, resp, err := client.Group.ListGroups(ctx, nil)
	if err != nil {
		return e
	}

	for resp.HasNextPage() {
		var nextGroupSet []*okta.Group
		resp, err = resp.Next(ctx, &nextGroupSet)
		output = append(output, nextGroupSet...)
	}

	g.Resources = g.createResources(output)
	return nil
}
