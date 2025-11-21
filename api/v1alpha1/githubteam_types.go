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
	"strings"

	vaultutils "github.com/redhat-cop/vault-config-operator/api/v1alpha1/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// GitHubTeamSpec defines the desired state of GitHubTeam
type GitHubTeamSpec struct {
	// Connection represents the information needed to connect to Vault. This operator uses the standard Vault environment variables to connect to Vault. If you need to override those settings and for example connect to a different Vault instance, you can do with this section of the CR.
	// +kubebuilder:validation:Optional
	Connection *vaultutils.VaultConnection `json:"connection,omitempty"`

	// Authentication is the kube auth configuration to be used to execute this request
	// +kubebuilder:validation:Required
	Authentication vaultutils.KubeAuthConfiguration `json:"authentication,omitempty"`

	// Path at which to make the configuration.
	// The final path in Vault will be {[spec.authentication.namespace]}/auth/{spec.path}/map/teams/{spec.teamName}.
	// The authentication role must have the following capabilities = [ "create", "read", "update", "delete"] on that path.
	// +kubebuilder:validation:Required
	Path vaultutils.Path `json:"path,omitempty"`

	// TeamName is the GitHub team name in "slugified" format
	// +kubebuilder:validation:Required
	TeamName string `json:"teamName,omitempty"`

	// Policies is the comma separated list of policies to assign to this team
	// +kubebuilder:validation:Required
	// +listType=set
	Policies []string `json:"policies,omitempty"`
}

// GitHubTeamStatus defines the observed state of GitHubTeam
type GitHubTeamStatus struct {
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GitHubTeam is the Schema for the githubteams API
type GitHubTeam struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GitHubTeamSpec   `json:"spec,omitempty"`
	Status GitHubTeamStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GitHubTeamList contains a list of GitHubTeam
type GitHubTeamList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GitHubTeam `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GitHubTeam{}, &GitHubTeamList{})
}

var _ vaultutils.VaultObject = &GitHubTeam{}
var _ vaultutils.ConditionsAware = &GitHubTeam{}

func (d *GitHubTeam) GetVaultConnection() *vaultutils.VaultConnection {
	return d.Spec.Connection
}

func (d *GitHubTeam) GetKubeAuthConfiguration() *vaultutils.KubeAuthConfiguration {
	return &d.Spec.Authentication
}

func (d *GitHubTeam) GetPath() string {
	return vaultutils.CleansePath("auth/" + string(d.Spec.Path) + "/map/teams/" + d.Spec.TeamName)
}

func (d *GitHubTeam) GetPayload() map[string]interface{} {
	return map[string]interface{}{
		"value": strings.Join(d.Spec.Policies, ","),
	}
}

func (d *GitHubTeam) IsEquivalentToDesiredState(payload map[string]interface{}) bool {
	desiredPayload := d.GetPayload()
	return reflect.DeepEqual(desiredPayload, payload)
}

func (d *GitHubTeam) IsInitialized() bool {
	return true
}

func (d *GitHubTeam) IsValid() (bool, error) {
	return true, nil
}

func (d *GitHubTeam) IsDeletable() bool {
	return true
}

func (d *GitHubTeam) PrepareInternalValues(context context.Context, object client.Object) error {
	return nil
}

func (d *GitHubTeam) PrepareTLSConfig(context context.Context, object client.Object) error {
	return nil
}

func (m *GitHubTeam) GetConditions() []metav1.Condition {
	return m.Status.Conditions
}

func (m *GitHubTeam) SetConditions(conditions []metav1.Condition) {
	m.Status.Conditions = conditions
}
