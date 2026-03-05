package github

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccGithubOrganizationNetworkConfiguration(t *testing.T) {
	networkSettingsID := testAccConf.testOrgNetworkSettingsID
	if networkSettingsID == "" {
		t.Skip("GH_TEST_ORG_NETWORK_SETTINGS_ID not set")
	}
	t.Run("creates organization network configuration without error", func(t *testing.T) {
		config := fmt.Sprintf(`
		resource "github_organization_network_configuration" "test" {
			name                  = "test-network-configuration"
			compute_service       = "actions"
			network_settings_ids  = ["%s"]
		}`, networkSettingsID)

		resource.Test(t, resource.TestCase{
			PreCheck:          func() { skipUnlessHasPaidOrgs(t) },
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config: config,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							"github_organization_network_configuration.test",
							tfjsonpath.New("name"),
							knownvalue.NotNull(),
						),
						statecheck.ExpectKnownValue(
							"github_organization_network_configuration.test",
							tfjsonpath.New("compute_service"),
							knownvalue.StringExact("actions"),
						),
						statecheck.ExpectKnownValue(
							"github_organization_network_configuration.test",
							tfjsonpath.New("network_settings_ids"),
							knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact(networkSettingsID),
							}),
						),
						statecheck.ExpectKnownValue(
							"github_organization_network_configuration.test",
							tfjsonpath.New("id"),
							knownvalue.NotNull(),
						),
						statecheck.ExpectKnownValue(
							"github_organization_network_configuration.test",
							tfjsonpath.New("created_on"),
							knownvalue.NotNull(),
						),
					},
				},
			},
		})
	})

	t.Run("updates organization network configuration without error", func(t *testing.T) {
		name := "test-network-config-one"
		computeService := "actions"

		updatedName := "test-network-config-two"
		updatedComputeService := "actions"

		config := `
resource "github_organization_network_configuration" "test" {
	name                  = "%s"
	compute_service       = "%s"
	network_settings_ids  = ["%s"]
}
`

		resource.Test(t, resource.TestCase{
			PreCheck:          func() { skipUnlessHasPaidOrgs(t) },
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config: fmt.Sprintf(config, name, computeService, networkSettingsID),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							"github_organization_network_configuration.test",
							tfjsonpath.New("name"),
							knownvalue.StringExact(name),
						),
						statecheck.ExpectKnownValue(
							"github_organization_network_configuration.test",
							tfjsonpath.New("compute_service"),
							knownvalue.StringExact(computeService),
						),
						statecheck.ExpectKnownValue(
							"github_organization_network_configuration.test",
							tfjsonpath.New("network_settings_ids"),
							knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact(networkSettingsID),
							}),
						),
					},
				},
				{
					Config: fmt.Sprintf(config, updatedName, updatedComputeService, networkSettingsID),
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							"github_organization_network_configuration.test",
							tfjsonpath.New("name"),
							knownvalue.StringExact(updatedName),
						),
						statecheck.ExpectKnownValue(
							"github_organization_network_configuration.test",
							tfjsonpath.New("compute_service"),
							knownvalue.StringExact(updatedComputeService),
						),
						statecheck.ExpectKnownValue(
							"github_organization_network_configuration.test",
							tfjsonpath.New("network_settings_ids"),
							knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact(networkSettingsID),
							}),
						),
					},
				},
			},
		})
	})

	t.Run("imports organization network configuration without error", func(t *testing.T) {
		name := "test-network-config-import"
		computeService := "actions"

		config := fmt.Sprintf(`
		resource "github_organization_network_configuration" "test" {
			name                  = "%s"
			compute_service       = "%s"
			network_settings_ids  = ["%s"]
		}`, name, computeService, networkSettingsID)

		resource.Test(t, resource.TestCase{
			PreCheck:          func() { skipUnlessHasPaidOrgs(t) },
			ProviderFactories: providerFactories,
			Steps: []resource.TestStep{
				{
					Config: config,
					ConfigStateChecks: []statecheck.StateCheck{
						statecheck.ExpectKnownValue(
							"github_organization_network_configuration.test",
							tfjsonpath.New("name"),
							knownvalue.StringExact(name),
						),
						statecheck.ExpectKnownValue(
							"github_organization_network_configuration.test",
							tfjsonpath.New("compute_service"),
							knownvalue.StringExact(computeService),
						),
					},
				},
				{
					ResourceName:      "github_organization_network_configuration.test",
					ImportState:       true,
					ImportStateVerify: true,
				},
			},
		})
	})
}
