package uiplugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// TestIncidentRulesDefaultsAreValid checks that the embedded incident rules
// satisfy the cluster-health-analyzer rules.yaml schema constraints.
func TestIncidentRulesDefaultsAreValid(t *testing.T) {
	var cfg struct {
		Rules []struct {
			Name       string `yaml:"name"`
			Strategies []struct {
				Match         map[string]any      `yaml:"match"`
				Within        string              `yaml:"within"`
				From          string              `yaml:"from"`
				ExcludeAlerts []map[string]string `yaml:"excludeAlerts"`
			} `yaml:"strategies"`
		} `yaml:"rules"`
	}
	require.NoError(t, yaml.Unmarshal([]byte(incidentRulesConfig), &cfg))
	require.NotEmpty(t, cfg.Rules)

	cm := newIncidentRulesConfigMap("ns")
	assert.Equal(t, "cluster-health-analyzer-default-rules", cm.Name)
	assert.Equal(t, incidentRulesConfig, cm.Data["rules.yaml"])

	role := newIncidentRulesRole("ns")
	require.Len(t, role.Rules, 1)
	assert.ElementsMatch(t, []string{"cluster-health-analyzer-default-rules", "cluster-health-analyzer-custom-rules"}, role.Rules[0].ResourceNames)
	assert.ElementsMatch(t, []string{"get", "list", "watch"}, role.Rules[0].Verbs)

	rb := newIncidentRulesRoleBinding("ns", "monitoring-plugin-sa")
	assert.Equal(t, "monitoring-plugin-sa", rb.Subjects[0].Name)
	assert.Equal(t, "cluster-health-analyzer-rules", rb.RoleRef.Name)
}
