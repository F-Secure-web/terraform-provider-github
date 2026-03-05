package github

import (
	"context"
	"errors"
	"net/http"
	"regexp"

	"github.com/google/go-github/v83/github"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceGithubOrganizationNetworkConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGithubOrganizationNetworkConfigurationCreate,
		ReadContext:   resourceGithubOrganizationNetworkConfigurationRead,
		UpdateContext: resourceGithubOrganizationNetworkConfigurationUpdate,
		DeleteContext: resourceGithubOrganizationNetworkConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.All(
					validation.StringLenBetween(1, 100),
					validation.StringMatch(
						regexp.MustCompile(`^[a-zA-Z0-9._-]+$`),
						"name may only contain upper and lowercase letters a-z, numbers 0-9, '.', '-', and '_'",
					),
				)),
				Description: "Name of the network configuration. Must be between 1 and 100 characters and may only contain upper and lowercase letters a-z, numbers 0-9, '.', '-', and '_'.",
			},
			"compute_service": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "none",
				ValidateDiagFunc: validateValueFunc([]string{"none", "actions"}),
				Description:      "The hosted compute service to use for the network configuration. Can be one of: 'none', 'actions'. Defaults to 'none'.",
			},
			"network_settings_ids": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "An array containing exactly one network settings ID. A network settings resource can only be associated with one network configuration at a time.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network configuration ID.",
			},
			"created_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the network configuration was created.",
			},
		},
	}
}

func resourceGithubOrganizationNetworkConfigurationCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	err := checkOrganization(m)
	if err != nil {
		return diag.FromErr(err)
	}

	meta := m.(*Owner)
	client := meta.v3client
	orgName := meta.name

	name := d.Get("name").(string)
	computeService := github.ComputeService(d.Get("compute_service").(string))
	networkSettingsIDs := expandStringList(d.Get("network_settings_ids").([]any))

	req := github.NetworkConfigurationRequest{
		Name:               github.Ptr(name),
		ComputeService:     &computeService,
		NetworkSettingsIDs: networkSettingsIDs,
	}

	tflog.Debug(ctx, "Creating network configuration", map[string]any{
		"name": name,
		"org":  orgName,
	})

	configuration, _, err := client.Organizations.CreateNetworkConfiguration(ctx, orgName, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(configuration.GetID())

	if err := d.Set("created_on", configuration.GetCreatedOn().Format("2006-01-02T15:04:05Z")); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGithubOrganizationNetworkConfigurationRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	err := checkOrganization(m)
	if err != nil {
		return diag.FromErr(err)
	}

	meta := m.(*Owner)
	client := meta.v3client
	orgName := meta.name
	ctx = tflog.SetField(ctx, "id", d.Id())
	networkID := d.Id()

	tflog.Debug(ctx, "Reading network configuration", map[string]any{
		"org": orgName,
	})

	configuration, resp, err := client.Organizations.GetNetworkConfiguration(ctx, orgName, networkID)
	if err != nil {
		var ghErr *github.ErrorResponse
		if errors.As(err, &ghErr) {
			if ghErr.Response.StatusCode == http.StatusNotFound {
				tflog.Info(ctx, "Removing network configuration from state because it no longer exists")
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	if resp.StatusCode == http.StatusNotModified {
		return nil
	}

	if err := d.Set("name", configuration.GetName()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("compute_service", configuration.GetComputeService()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("network_settings_ids", configuration.NetworkSettingsIDs); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_on", configuration.GetCreatedOn().Format("2006-01-02T15:04:05Z")); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGithubOrganizationNetworkConfigurationUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	err := checkOrganization(m)
	if err != nil {
		return diag.FromErr(err)
	}

	meta := m.(*Owner)
	client := meta.v3client
	orgName := meta.name
	ctx = tflog.SetField(ctx, "id", d.Id())
	networkID := d.Id()

	computeService := github.ComputeService(d.Get("compute_service").(string))
	networkSettingsIDs := expandStringList(d.Get("network_settings_ids").([]any))

	req := github.NetworkConfigurationRequest{
		Name:               github.Ptr(d.Get("name").(string)),
		ComputeService:     &computeService,
		NetworkSettingsIDs: networkSettingsIDs,
	}

	tflog.Debug(ctx, "Updating network configuration", map[string]any{
		"org": orgName,
	})

	_, _, err = client.Organizations.UpdateNetworkConfiguration(ctx, orgName, networkID, req)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGithubOrganizationNetworkConfigurationDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	err := checkOrganization(m)
	if err != nil {
		return diag.FromErr(err)
	}

	meta := m.(*Owner)
	client := meta.v3client
	orgName := meta.name
	ctx = tflog.SetField(ctx, "id", d.Id())
	networkID := d.Id()

	tflog.Debug(ctx, "Deleting network configuration", map[string]any{
		"org": orgName,
	})

	_, err = client.Organizations.DeleteNetworkConfigurations(ctx, orgName, networkID)
	if err != nil {
		var ghErr *github.ErrorResponse
		if errors.As(err, &ghErr) {
			if ghErr.Response.StatusCode == http.StatusNotFound {
				return nil
			}
		}
		return diag.FromErr(err)
	}

	return nil
}
