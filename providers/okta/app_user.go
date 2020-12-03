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

type AppUserGenerator struct {
	OktaService
}

type OktaAppUser struct {
	appId   string
	appUser *okta.AppUser
}

func (g AppUserGenerator) createResources(appUserList []*OktaAppUser) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for _, appUser := range appUserList {
		resources = append(resources, terraformutils.NewResource(
			appUser.appUser.Id,
			appUser.appUser.Id,
			"okta_app_user",
			"okta",
			map[string]string{
				"app_id":   appUser.appId,
				"user_id":  appUser.appUser.Id,
				"username": appUser.appUser.Credentials.UserName,
			},
			[]string{},
			map[string]interface{}{},
		))
	}
	return resources
}

func (g *AppUserGenerator) InitResources() error {
	ctx, client, e := g.generateClient()
	if e != nil {
		return e
	}

	apps, err := getAllApplications(ctx, client)
	if err != nil {
		return err
	}

	// NOTE: Odd reason this cause all arrays for record.
	// for resp.HasNextPage() {
	// 	var nextAppUserSet []*okta.AppUser
	// 	resp, err = resp.Next(ctx, &nextAppUserSet)
	// 	output = append(output, nextAppUserSet...)
	// }

	//BUG: When jan.doe is assign to more then apps it blows up with duplicated error.

	var oktaAppUsers []*OktaAppUser
	for _, app := range apps {

		//NOTE: Left this example when I only wanted to test a single app to user association.

		// if app.Name == "template_swa" {
		// 	appUsers, _, err := client.Application.ListApplicationUsers(ctx, app.Id, nil)
		// 	if err != nil {
		// 		return err
		// 	}

		// 	for _, appUser := range appUsers {
		// 		oktaAppUser := &OktaAppUser{
		// 			appId:   app.Id,
		// 			appUser: appUser,
		// 		}
		// 		oktaAppUsers = append(oktaAppUsers, oktaAppUser)
		// 	}
		// }

		appUsers, _, err := client.Application.ListApplicationUsers(ctx, app.Id, nil)
		if err != nil {
			return err
		}

		for _, appUser := range appUsers {
			oktaAppUser := &OktaAppUser{
				appId:   app.Id,
				appUser: appUser,
			}
			oktaAppUsers = append(oktaAppUsers, oktaAppUser)
		}
	}

	g.Resources = g.createResources(oktaAppUsers)
	return nil
}
