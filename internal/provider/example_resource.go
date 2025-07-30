// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ExampleResource{}
var _ resource.ResourceWithImportState = &ExampleResource{}

func NewExampleResource() resource.Resource {
	return &ExampleResource{}
}

// ExampleResource defines the resource implementation.
type ExampleResource struct {
	client *http.Client
}

// ExampleResourceModel describes the resource data model.
type ExampleResourceModel struct {
	ID                  types.String `tfsdk:"id"`
	SetIfNotNullOrEmpty types.String `tfsdk:"set_if_not_null_or_empty"` // Setter (input, write-only)
	Value               types.String `tfsdk:"value"`                    // Getter (computed, read-only)
}

func (r *ExampleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_example"
}

func (r *ExampleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Example resource mit write-only Setter und read-only Getter.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Example identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			// SETTER: write-only, optionaler String
			"set_if_not_null_or_empty": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Setter-Property: Wenn nicht null/leer, wird `value` auf diesen String gesetzt. Nach Apply wird dieses Feld im State auf null gesetzt.",
			},

			// GETTER: read-only Ergebnis
			"value": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					valueFromSetterOrStateModifier{},
				},
				MarkdownDescription: "Getter-Property. Wird nur aus dem Setter gesetzt; sonst bleibt der State-Wert erhalten.",
			},
		},
	}
}

func (r *ExampleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*http.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ExampleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ExampleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Dummy-ID setzen
	data.ID = types.StringValue(time.Now().UTC().Format(time.RFC3339Nano))

	// write-only: Setter nach Apply entfernen
	data.SetIfNotNullOrEmpty = types.StringNull()

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExampleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ExampleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExampleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ExampleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// write-only: Setter wieder nullen
	data.SetIfNotNullOrEmpty = types.StringNull()

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExampleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ExampleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *ExampleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// -----------------------------------------------------------------------------
// Plan-Modifier: `valueFromSetterOrStateModifier`
// -----------------------------------------------------------------------------
// Semantik:
//   - Wenn `set_if_not_null_or_empty` (Config) NICHT null/unknown UND trim != "",
//     dann wird `value` im Plan auf diesen String gesetzt.
//   - Andernfalls wird `value` aus dem bisherigen State übernommen (falls vorhanden).
//   - Bei Create ohne State und ohne Setter bleibt `value` unknown/null im Plan.
//
// -----------------------------------------------------------------------------
type valueFromSetterOrStateModifier struct{}

func (m valueFromSetterOrStateModifier) Description(_ context.Context) string {
	return "Setzt `value` aus dem Setter, wenn dieser nicht leer ist; sonst übernimmt `value` den vorhandenen State-Wert."
}

func (m valueFromSetterOrStateModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m valueFromSetterOrStateModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Früh beenden, falls bereits Fehler vorliegen
	if resp.Diagnostics.HasError() {
		return
	}

	// Setter aus der Config lesen
	var setter types.String
	if diags := req.Config.GetAttribute(ctx, path.Root("set_if_not_null_or_empty"), &setter); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Wenn der Setter bekannt und nicht leer ist, setze value auf den Setter-String
	if !setter.IsNull() && !setter.IsUnknown() {
		if s := strings.TrimSpace(setter.ValueString()); s != "" {
			resp.PlanValue = types.StringValue(s)
			return
		}
	}

	// Ansonsten: State beibehalten, wenn vorhanden
	if !req.StateValue.IsNull() && !req.StateValue.IsUnknown() {
		resp.PlanValue = req.StateValue
		return
	}

	// Kein State vorhanden (z. B. Create) und kein gültiger Setter:
	// -> PlanValue bleibt unknown/null (Terraform behandelt das korrekt).
}
