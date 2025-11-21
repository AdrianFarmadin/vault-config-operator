/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"context"
	"reflect"

	vaultutils "github.com/redhat-cop/vault-config-operator/api/v1alpha1/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GitHubAuthEngineConfigSpec defines the desired state of GitHubAuthEngineConfig
type GitHubAuthEngineConfigSpec struct {
	// Connection represents the information needed to connect to Vault. This operator uses the standard Vault environment variables to connect to Vault. If you need to override those settings and for example connect to a different Vault instance, you can do with this section of the CR.
	// +kubebuilder:validation:Optional
	Connection *vaultutils.VaultConnection `json:"connection,omitempty"`

	// Authentication is the kube auth configuration to be used to execute this request
	// +kubebuilder:validation:Required
	Authentication vaultutils.KubeAuthConfiguration `json:"authentication,omitempty"`

	// Path at which to make the configuration.
	// The final path in Vault will be {[spec.authentication.namespace]}/auth/{spec.path}/config.
	// The authentication role must have the following capabilities = [ "create", "read", "update", "delete"] on that path.
	// +kubebuilder:validation:Required
	Path vaultutils.Path `json:"path,omitempty"`

	// Organization is the organization users must be part of to authenticate.
	// +kubebuilder:validation:Required
	Organization string `json:"organization,omitempty"`

	// OrganizationID is the ID of the organization users must be part of. Vault will attempt to fetch and set this value if it is not provided.
	// +kubebuilder:validation:Optional
	OrganizationID int `json:"organizationID,omitempty"`

	// BaseURL is the API endpoint to use. Useful if you are running GitHub Enterprise or an API-compatible authentication server.
	// +kubebuilder:validation:Optional
	BaseURL string `json:"baseURL,omitempty"`

	// TokenTTL is the incremental lifetime for generated tokens. This current value of this will be referenced at renewal time.
	// +kubebuilder:validation:Optional
	TokenTTL metav1.Duration `json:"tokenTTL,omitempty"`

	// TokenMaxTTL is the maximum lifetime for generated tokens. This current value of this will be referenced at renewal time.
	// +kubebuilder:validation:Optional
	TokenMaxTTL metav1.Duration `json:"tokenMaxTTL,omitempty"`

	// TokenPolicies is the list of token policies to encode onto generated tokens. Depending on the auth method, this list may be supplemented by user/group/other values.
	// +kubebuilder:validation:Optional
	// +listType=set
	TokenPolicies []string `json:"tokenPolicies,omitempty"`

	// TokenBoundCIDRs is the list of CIDR blocks; if set, specifies blocks of IP addresses which can authenticate successfully, and ties the resulting token to these blocks as well.
	// +kubebuilder:validation:Optional
	// +listType=set
	TokenBoundCIDRs []string `json:"tokenBoundCIDRs,omitempty"`

	// TokenExplicitMaxTTL if set, will encode an explicit max TTL onto the token. This is a hard cap even if token_ttl and token_max_ttl would otherwise allow a renewal.
	// +kubebuilder:validation:Optional
	TokenExplicitMaxTTL metav1.Duration `json:"tokenExplicitMaxTTL,omitempty"`

	// TokenNoDefaultPolicy if set, the default policy will not be set on generated tokens; otherwise it will be added to the policies set in token_policies.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=false
	TokenNoDefaultPolicy bool `json:"tokenNoDefaultPolicy,omitempty"`

	// TokenNumUses is the maximum number of times a generated token may be used (within its lifetime); 0 means unlimited.
	// +kubebuilder:validation:Optional
	// +kubebuilder:default=0
	TokenNumUses int `json:"tokenNumUses,omitempty"`

	// TokenPeriod is the maximum allowed period value when a periodic token is requested from this role.
	// +kubebuilder:validation:Optional
	TokenPeriod metav1.Duration `json:"tokenPeriod,omitempty"`

	// TokenType is the type of token that should be generated. Can be service, batch, or default to use the mount's tuned default.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Enum={"","service","batch","default","default-service","default-batch"}
	TokenType string `json:"tokenType,omitempty"`
}

// GitHubAuthEngineConfigStatus defines the observed state of GitHubAuthEngineConfig
type GitHubAuthEngineConfigStatus struct {
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GitHubAuthEngineConfig is the Schema for the githubauthengineconfigs API
type GitHubAuthEngineConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GitHubAuthEngineConfigSpec   `json:"spec,omitempty"`
	Status GitHubAuthEngineConfigStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GitHubAuthEngineConfigList contains a list of GitHubAuthEngineConfig
type GitHubAuthEngineConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GitHubAuthEngineConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GitHubAuthEngineConfig{}, &GitHubAuthEngineConfigList{})
}

var _ vaultutils.VaultObject = &GitHubAuthEngineConfig{}
var _ vaultutils.ConditionsAware = &GitHubAuthEngineConfig{}

func (d *GitHubAuthEngineConfig) GetVaultConnection() *vaultutils.VaultConnection {
	return d.Spec.Connection
}

func (d *GitHubAuthEngineConfig) GetKubeAuthConfiguration() *vaultutils.KubeAuthConfiguration {
	return &d.Spec.Authentication
}

func (d *GitHubAuthEngineConfig) GetPath() string {
	return vaultutils.CleansePath("auth/" + string(d.Spec.Path) + "/config")
}

func (d *GitHubAuthEngineConfig) GetPayload() map[string]interface{} {
	payload := map[string]interface{}{
		"organization": d.Spec.Organization,
	}

	if d.Spec.OrganizationID != 0 {
		payload["organization_id"] = d.Spec.OrganizationID
	}

	if d.Spec.BaseURL != "" {
		payload["base_url"] = d.Spec.BaseURL
	}

	if d.Spec.TokenTTL.Duration > 0 {
		payload["token_ttl"] = d.Spec.TokenTTL.Duration.String()
	}

	if d.Spec.TokenMaxTTL.Duration > 0 {
		payload["token_max_ttl"] = d.Spec.TokenMaxTTL.Duration.String()
	}

	if len(d.Spec.TokenPolicies) > 0 {
		payload["token_policies"] = d.Spec.TokenPolicies
	}

	if len(d.Spec.TokenBoundCIDRs) > 0 {
		payload["token_bound_cidrs"] = d.Spec.TokenBoundCIDRs
	}

	if d.Spec.TokenExplicitMaxTTL.Duration > 0 {
		payload["token_explicit_max_ttl"] = d.Spec.TokenExplicitMaxTTL.Duration.String()
	}

	if d.Spec.TokenNoDefaultPolicy {
		payload["token_no_default_policy"] = d.Spec.TokenNoDefaultPolicy
	}

	if d.Spec.TokenNumUses != 0 {
		payload["token_num_uses"] = d.Spec.TokenNumUses
	}

	if d.Spec.TokenPeriod.Duration > 0 {
		payload["token_period"] = d.Spec.TokenPeriod.Duration.String()
	}

	if d.Spec.TokenType != "" {
		payload["token_type"] = d.Spec.TokenType
	}

	return payload
}

func (d *GitHubAuthEngineConfig) IsEquivalentToDesiredState(payload map[string]interface{}) bool {
	desiredPayload := d.GetPayload()
	return reflect.DeepEqual(desiredPayload, payload)
}

func (d *GitHubAuthEngineConfig) IsInitialized() bool {
	return true
}

func (d *GitHubAuthEngineConfig) IsValid() (bool, error) {
	return true, nil
}

func (d *GitHubAuthEngineConfig) IsDeletable() bool {
	return false
}

func (d *GitHubAuthEngineConfig) PrepareInternalValues(context context.Context, object client.Object) error {
	return nil
}

func (d *GitHubAuthEngineConfig) PrepareTLSConfig(context context.Context, object client.Object) error {
	return nil
}

func (m *GitHubAuthEngineConfig) GetConditions() []metav1.Condition {
	return m.Status.Conditions
}

func (m *GitHubAuthEngineConfig) SetConditions(conditions []metav1.Condition) {
	m.Status.Conditions = conditions
}
