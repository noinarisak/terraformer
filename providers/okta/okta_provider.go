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

// Okta resources
// Completed:
//	okta_app_auto_login,
//  okta_app_basic_auth,
//  okta_app_bookmark,
//  okta_app_oauth,
//  okta_app_saml,
//	okta_app_secure_password_store,
//  okta_app_swa,
//  okta_app_three_field,
//  okta_group,
//  okta_user,
//	okta_auth_server,
//  okta_group_rule,
//  okta_event_hooks,
//  okta_trusted_origin,
//  okta_user_type
//  okta_template_sms,
//  okta_inline_hooks,
//  okta_idp_oidc,
//  okta_idp_saml,
//  okta_idp_social,
//  okta_policy_mfa,
//  okta_policy_password,
//  okta_policy_sign_on,
//  okta_factor,
//  okta_network_zone,
// Issues:
//	?okta_app_user (ListApplicationUsers(appID))
// Not Completed: *= Futher research, ?=Most likely doable, ^=Doable

//  ?okta_app_group_assignment (ListApplicationGroupAssignments(appID)) (PR ready)

//	*okta_app_user_base_schema (API Supplement)
//  *okta_app_user_schema (API Supplement)

//	*okta_auth_server_claim (API Supplement) (PR ready)
//  *okta_auth_server_policy (API Supplement) (PR ready)
//  *okta_auth_server_policy_rule (API Supplement) (PR ready)
//  *okta_auth_server_scope (API Supplement) (PR ready)

//	?okta_group_roles (ListGroupAssignedRoles(groupId))
//  ?okta_idp_saml_signing_key (ListIdentityProviderSigningKeys(idpID))

//  *okta_policy_rule_idp_discovery (API Supplement) (PR ready)
//  ^okta_policy_rule_mfa (ListPolicyRules( type="MFA_ENROLL")) (PR ready)
//  ^okta_policy_rule_password (ListPolicyRules( type="PASSWORD")) (PR ready)
//  ^okta_policy_rule_sign_on (ListPolicyRules( type="SIGN_ON")) (PR ready)

//	*okta_profile_mapping (API Supplement)

//	*okta_template_email (API Supplement) (PR ready)

//	*okta_user_base_schema (API Supplement)
//	*okta_user_schema (API Supplement)

package okta

import (
	"os"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/GoogleCloudPlatform/terraformer/terraformutils/providerwrapper"

	"github.com/pkg/errors"
	"github.com/zclconf/go-cty/cty"
)

type OktaProvider struct {
	terraformutils.Provider
	apiToken string
	baseURL  string
	orgName  string
}

func (p OktaProvider) GetResourceConnections() map[string]map[string][]string {
	return map[string]map[string][]string{
		"alerts": {"alert_notification_endpoints": []string{"alert_notification_endpoints", "id"}},
	}
}

func (p OktaProvider) GetProviderData(arg ...string) map[string]interface{} {
	return map[string]interface{}{
		"provider": map[string]interface{}{
			"okta": map[string]interface{}{
				"version": providerwrapper.GetProviderVersion(p.GetName()),
			},
		},
	}
}

func (p *OktaProvider) GetConfig() cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"api_token": cty.StringVal(p.apiToken),
		"base_url":  cty.StringVal(p.baseURL),
		"org_name":  cty.StringVal(p.orgName),
	})
}

// Init OktaProvider
func (p *OktaProvider) Init(args []string) error {
	apiToken := os.Getenv("OKTA_API_TOKEN")
	if apiToken == "" {
		return errors.New("set OKTA_API_TOKEN env var")
	}
	p.apiToken = apiToken

	baseURL := os.Getenv("OKTA_BASE_URL")
	if baseURL == "" {
		return errors.New("set OKTA_BASE_URL env var")
	}
	p.baseURL = baseURL

	orgName := os.Getenv("OKTA_ORG_NAME")
	if orgName == "" {
		return errors.New("set OKTA_ORG_NAME env var")
	}
	p.orgName = orgName

	return nil
}

func (p *OktaProvider) GetName() string {
	return "okta"
}

func (p *OktaProvider) InitService(serviceName string, verbose bool) error {
	var isSupported bool
	if _, isSupported = p.GetSupportedService()[serviceName]; !isSupported {
		return errors.New(p.GetName() + ": " + serviceName + " not supported service")
	}
	p.Service = p.GetSupportedService()[serviceName]
	p.Service.SetName(serviceName)
	p.Service.SetVerbose(verbose)
	p.Service.SetProviderName(p.GetName())
	p.Service.SetArgs(map[string]interface{}{
		"api_token": p.apiToken,
		"base_url":  p.baseURL,
		"org_name":  p.orgName,
	})
	return nil
}

// GetSupportedService return map of support service for Okta
func (p *OktaProvider) GetSupportedService() map[string]terraformutils.ServiceGenerator {
	return map[string]terraformutils.ServiceGenerator{
		"app": &AppGenerator{},
		//BUG:	Still has bug when assigning a User to multiple clients scenario.
		"app_user":     &AppUserGenerator{},
		"auth_server":  &AuthServerGenerator{},
		"event_hook":   &EventHookGenerator{},
		"factor":       &FactorGenerator{},
		"group":        &GroupGenerator{},
		"group_rule":   &GroupRuleGenerator{},
		"idp":          &IDPGenerator{},
		"inline_hook":  &InlineHookGenerator{},
		"network_zone": &NetworkZoneGenerator{},
		"policy":       &PolicyGenerator{},
		//BUG: Does not work. Getting 'The API returned an error: The request is missing a required parameter., Status: 400 Bad Request'
		//		It seems the request is going to "GET /api/v1/policies/ HTTP/1.1' and not "ET /api/v1/policies/{policyID}/rules/{ruleID} HTTP/1.1"
		//		So possible solution is related to this PR (https://bit.ly/3bdEofU) and to terraformutils.NewResource() method.
		// 	 	2/1/21 - Using NewResource() was the key and getting closer to getting the policy_rule_* working.
		// "policy_rule":    &PolicyRuleGenerator{},
		"policy_rule_password": &PolicyRulePasswordGenerator{},
		// "policy_rule_signon":   &PolicyRuleSignOnGenerator{},
		"template_sms":   &TempaleSMSGenerator{},
		"trusted_origin": &TrustedOriginGenerator{},
		"user":           &UserGenerator{},
		"user_type":      &UserTypeGenerator{},
	}
}
