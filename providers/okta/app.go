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
	"context"
	"fmt"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

type AppGenerator struct {
	OktaService
}

//NOTE:
// SIGNON_DEFAULTS = {
//     "bookmark": "BOOKMARK",
//     "template_basic_auth": "BASIC_AUTH",
//     "template_swa": "BROWSER_PLUGIN",
//     "template_swa3field": "BROWSER_PLUGIN",
//     "template_sps": "SECURE_PASSWORD_STORE",
//     "oidc_client": "OPENID_CONNECT",
//     "template_wsfed": "WS_FEDERATION",
// }

func (g AppGenerator) createResources(appList []*okta.Application) []terraformutils.Resource {
	var resources []terraformutils.Resource
	for _, app := range appList {
		resourceOktaAppName := ""

		// app := app.(*okta.Application)

		switch app.SignOnMode {
		case "OPENID_CONNECT":
			resourceOktaAppName = "okta_app_oauth"
		case "SAML_1_1", "SAML_2_0":
			resourceOktaAppName = "okta_app_saml"
		case "WS_FEDERATION":
			fmt.Printf("WARN: Does not support WS_FEDERATION type")
		case "BOOKMARK":
			resourceOktaAppName = "okta_app_bookmark"
		case "AUTO_LOGIN":
			resourceOktaAppName = "okta_app_auto_login"
		case "BASIC_AUTH":
			resourceOktaAppName = "okta_app_basic_auth"
		case "SECURE_PASSWORD_STORE":
			resourceOktaAppName = "okta_app_secure_password_store"
		case "BROWSER_PLUGIN":
			if app.Name == "template_swa" {
				resourceOktaAppName = "okta_app_swa"
			} else if app.Name == "template_swa3field" {
				resourceOktaAppName = "okta_app_three_field"
			} else {
				fmt.Printf("ERROR: Not supported application type %s\n", app.Name)
			}
		default:
			{
				fmt.Printf("ERROR: Not supported Sign On Mode type %s\n", app.SignOnMode)
			}
		}

		resources = append(resources, terraformutils.NewSimpleResource(
			app.Id,
			app.Name,
			resourceOktaAppName,
			"okta",
			[]string{}))
	}
	return resources
}

// Generate Terraform Resources from Okta API,
func (g *AppGenerator) InitResources() error {
	ctx, client, e := g.generateClient()
	if e != nil {
		return e
	}

	apps, err := getAllApplications(ctx, client)
	if err != nil {
		return err
	}

	g.Resources = g.createResources(apps)
	return nil
}

func getAllApplications(ctx context.Context, client *okta.Client) ([]*okta.Application, error) {
	apps, resp, err := client.Application.ListApplications(ctx, nil)
	if err != nil {
		return nil, err
	}
	resultingApps := make([]*okta.Application, len(apps))
	for i := range apps {
		resultingApps[i] = apps[i].(*okta.Application)
	}
	for resp.HasNextPage() {
		var nextApps []*okta.Application
		resp, err = resp.Next(ctx, &nextApps)
		if err != nil {
			return nil, err
		}
		for i := range nextApps {
			resultingApps = append(resultingApps, nextApps[i])
		}
	}

	var supportedApps []*okta.Application
	for _, app := range resultingApps {
		//NOTE: Okta provider does not support the following app types:
		if app.Name == "template_wsfed" ||
			app.Name == "template_swa_two_page" ||
			app.Name == "okta_enduser" ||
			app.Name == "okta_browser_plugin" ||
			app.Name == "saasure" {
			continue
		}
		supportedApps = append(supportedApps, app)
	}
	return supportedApps, nil
}
