// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://mozilla.org/MPL/2.0/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: MPL-2.0

package provider

// Section below is generated&owned by "gen/generator.go". //template:begin imports
import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/netascode/go-fmc"
	"github.com/netascode/terraform-provider-fmc/internal/provider/helpers"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin model

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &PortGroupsResource{}
	_ resource.ResourceWithImportState = &PortGroupsResource{}
)

func NewPortGroupsResource() resource.Resource {
	return &PortGroupsResource{}
}

type PortGroupsResource struct {
	client *fmc.Client
}

func (r *PortGroupsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_port_groups"
}

func (r *PortGroupsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: helpers.NewAttributeDescription("This plural resource manages a bulk of Port Groups. The FMC API supports quick bulk creation for this resource, but the deletion/modification is done one-by-one. Updating and deleting `fmc_port_groups` can thus take much more time than creating it (even >500 times more time, i.e. >50000%, depending on the change size).").String,

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The id of the object",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The name of the FMC domain",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"items": schema.MapNestedAttribute{
				MarkdownDescription: helpers.NewAttributeDescription("Map of port groups. The key of the map is the name of the individual Port Group. ").String,
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("UUID of the managed Port Groups.").String,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"type": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("").String,
							Required:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Optional user-created description.").String,
							Optional:            true,
						},
						"overridable": schema.BoolAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("Indicates whether object values can be overridden.").String,
							Optional:            true,
						},
						"ports": schema.SetNestedAttribute{
							MarkdownDescription: helpers.NewAttributeDescription("").String,
							Required:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("UUID of the port (such as fmc_port.test.id, etc.).").String,
										Optional:            true,
									},
									"type": schema.StringAttribute{
										MarkdownDescription: helpers.NewAttributeDescription("").AddStringEnumDescription("ProtocolPortObject", "ICMPV6Object", "ICMPV4Object").String,
										Required:            true,
										Validators: []validator.String{
											stringvalidator.OneOf("ProtocolPortObject", "ICMPV6Object", "ICMPV4Object"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *PortGroupsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*FmcProviderData).Client
}

// End of section. //template:end model

// Section below is generated&owned by "gen/generator.go". //template:begin create

func (r *PortGroupsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PortGroups

	// Read plan
	diags := req.Plan.Get(ctx, &plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Set request domain if provided
	reqMods := [](func(*fmc.Req)){}
	if !plan.Domain.IsNull() && plan.Domain.ValueString() != "" {
		reqMods = append(reqMods, fmc.DomainName(plan.Domain.ValueString()))
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Create", plan.Id.ValueString()))

	// Prepare state to track creation process
	// Create request is split to multiple requests, where just subset of them may be successful
	state := PortGroups{}
	state.Items = make(map[string]PortGroupsItems, len(plan.Items))
	state.Id = types.StringValue(uuid.New().String())
	state.Domain = plan.Domain

	// Create object
	// Creation process is put in a separate function, as that same proces will be needed with `Update`
	plan, diags = r.createSubresources(ctx, state, plan, reqMods...)
	resp.Diagnostics.Append(diags...)

	tflog.Debug(ctx, fmt.Sprintf("%s: Create finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	helpers.SetFlagImporting(ctx, false, resp.Private, &resp.Diagnostics)
}

// End of section. //template:end create

// Section below is generated&owned by "gen/generator.go". //template:begin read

func (r *PortGroupsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state PortGroups

	// Read state
	diags := req.State.Get(ctx, &state)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Set request domain if provided
	reqMods := [](func(*fmc.Req)){}
	if !state.Domain.IsNull() && state.Domain.ValueString() != "" {
		reqMods = append(reqMods, fmc.DomainName(state.Domain.ValueString()))
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Read", state.Id.String()))

	// Get all objects from FMC
	urlPath := state.getPath() + "?expanded=true"
	res, err := r.client.Get(urlPath, reqMods...)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Failed to retrieve object (GET), got error: %s, %s", err, res.String()))
		return
	}

	imp, diags := helpers.IsFlagImporting(ctx, req)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// After `terraform import` we switch to a full read.
	if imp {
		state.fromBody(ctx, res)
	} else {
		state.fromBodyPartial(ctx, res)
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Read finished successfully", state.Id.ValueString()))

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	helpers.SetFlagImporting(ctx, false, resp.Private, &resp.Diagnostics)
}

// End of section. //template:end read

// Section below is generated&owned by "gen/generator.go". //template:begin update

func (r *PortGroupsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state PortGroups

	// Read plan
	diags := req.Plan.Get(ctx, &plan)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Read state
	diags = req.State.Get(ctx, &state)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Set request domain if provided
	reqMods := [](func(*fmc.Req)){}
	if !plan.Domain.IsNull() && plan.Domain.ValueString() != "" {
		reqMods = append(reqMods, fmc.DomainName(plan.Domain.ValueString()))
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Update", plan.Id.ValueString()))

	// DELETE
	// Delete objects (that are present in state, but missing in plan)
	var toDelete PortGroups
	toDelete.Items = make(map[string]PortGroupsItems)
	planOwnedIDs := make(map[string]string, len(plan.Items))

	// Prepare list of ID that are in plan
	for k, v := range plan.Items {
		planOwnedIDs[v.Id.ValueString()] = k
	}

	// Check if ID from state list is in plan as well. If not, mark it for delete
	for k, v := range state.Items {
		if _, ok := planOwnedIDs[v.Id.ValueString()]; !ok {
			toDelete.Items[k] = v
		}
	}

	// If there are objects marked to be deleted
	if len(toDelete.Items) > 0 {
		tflog.Debug(ctx, fmt.Sprintf("%s: Number of items to delete: %d", state.Id.ValueString(), len(toDelete.Items)))
		state, diags = r.deleteSubresources(ctx, state, toDelete, reqMods...)
		if diags != nil {
			resp.Diagnostics.Append(diags...)
			diags = resp.State.Set(ctx, &state)
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	// CREATE
	// Create new objects (objects that have missing IDs in plan)
	var toCreate PortGroups
	toCreate.Items = make(map[string]PortGroupsItems)
	// Scan plan for items with no ID
	for k, v := range plan.Items {
		if v.Id.IsUnknown() || v.Id.IsNull() {
			toCreate.Items[k] = v
		}
	}

	// If there are objects marked for create
	if len(toCreate.Items) > 0 {
		tflog.Debug(ctx, fmt.Sprintf("%s: Number of items to create: %d", state.Id.ValueString(), len(toCreate.Items)))
		state, diags = r.createSubresources(ctx, state, toCreate, reqMods...)
		if diags != nil {
			resp.Diagnostics.Append(diags...)
			diags = resp.State.Set(ctx, &state)
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	// UPDATE
	// Update objects (objects that have different definition in plan and state)
	var notEqual bool
	var toUpdate PortGroups
	toUpdate.Items = make(map[string]PortGroupsItems)

	for _, valueState := range state.Items {

		// Check if the ID from plan exists on list of ID owned by state
		if keyState, ok := planOwnedIDs[valueState.Id.ValueString()]; ok {

			// Check if items in state and plan are qual
			notEqual, diags = helpers.IsConfigUpdatingAt(ctx, req.Plan, req.State, path.Root("items").AtMapKey(keyState))
			if diags != nil {
				resp.Diagnostics.Append(diags...)
				diags = resp.State.Set(ctx, &state)
				resp.Diagnostics.Append(diags...)
				return
			}

			// If definitions differ, add object to update list
			if notEqual {
				toUpdate.Items[keyState] = plan.Items[keyState]
			}
		}
	}

	// If there are objects marked for update
	if len(toUpdate.Items) > 0 {
		tflog.Debug(ctx, fmt.Sprintf("%s: Number of items to update: %d", state.Id.ValueString(), len(toUpdate.Items)))
		state, diags = r.updateSubresources(ctx, state, toUpdate, reqMods...)
		if diags != nil {
			resp.Diagnostics.Append(diags...)
			diags = resp.State.Set(ctx, &state)
			resp.Diagnostics.Append(diags...)
			return
		}
	}
	plan = state

	tflog.Debug(ctx, fmt.Sprintf("%s: Update finished successfully", plan.Id.ValueString()))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// End of section. //template:end update

// Section below is generated&owned by "gen/generator.go". //template:begin delete

func (r *PortGroupsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state PortGroups

	// Read state
	diags := req.State.Get(ctx, &state)
	if resp.Diagnostics.Append(diags...); resp.Diagnostics.HasError() {
		return
	}

	// Set request domain if provided
	reqMods := [](func(*fmc.Req)){}
	if !state.Domain.IsNull() && state.Domain.ValueString() != "" {
		reqMods = append(reqMods, fmc.DomainName(state.Domain.ValueString()))
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Beginning Delete", state.Id.ValueString()))

	// Execute delete
	state, diags = r.deleteSubresources(ctx, state, state, reqMods...)
	resp.Diagnostics.Append(diags...)

	// Check if every element was removed
	if len(state.Items) > 0 {
		tflog.Debug(ctx, fmt.Sprintf("%s: Not all elements have been removed", state.Id.ValueString()))
		diags = resp.State.Set(ctx, &state)
		resp.Diagnostics.Append(diags...)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("%s: Delete finished successfully", state.Id.ValueString()))

	resp.State.RemoveResource(ctx)
}

// End of section. //template:end delete

// Section below is generated&owned by "gen/generator.go". //template:begin import
func (r *PortGroupsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import looks for string in the following format: <domain_name>,[<object1_name>,<object2_name>,...]
	// <domain_name> is optional
	// <object1_name>,<object2_name>,... is coma-separated list of object names
	var config PortGroups

	// Compile pattern for import command parsing
	var inputPattern = regexp.MustCompile(`^(?P<domain>[^\s,]*),*\[(?P<names>.*?),*\]`)

	// Parse parameter
	match := inputPattern.FindStringSubmatch(req.ID)

	// Check if regex matched
	if match == nil {
		resp.Diagnostics.AddError("Import error", "Failed to parse import parameters. Please provide import string in the following format: <domain_name>,[<object1_name>,<object2_name>,...]")
		return
	}

	// Extract values
	if tmpDomain := match[inputPattern.SubexpIndex("domain")]; tmpDomain != "" {
		config.Domain = types.StringValue(tmpDomain)
	}
	names := strings.Split(match[inputPattern.SubexpIndex("names")], ",")

	// Fill state with names of objects to import
	config.Items = make(map[string]PortGroupsItems, len(names))
	for _, v := range names {
		config.Items[v] = PortGroupsItems{}
	}

	// Generate new ID
	config.Id = types.StringValue(uuid.New().String())

	// Set filled in structure
	diags := resp.State.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set import flag
	helpers.SetFlagImporting(ctx, true, resp.Private, &resp.Diagnostics)
}

// End of section. //template:end import

// Section below is generated&owned by "gen/generator.go". //template:begin createSubresources
// createSubresources takes list of objects, splits them into bulks and creates them
// We want to save the state after each create event, to be able track already created resources
func (r *PortGroupsResource) createSubresources(ctx context.Context, state, plan PortGroups, reqMods ...func(*fmc.Req)) (PortGroups, diag.Diagnostics) {
	var idx = 0
	var bulk PortGroups
	bulk.Items = make(map[string]PortGroupsItems, bulkSizeCreate)

	tflog.Debug(ctx, fmt.Sprintf("%s: Creating bulk of objects", state.Id.ValueString()))

	// iterate over all items
	for k, v := range plan.Items {
		// count loops
		idx++

		// add object to current bulk
		bulk.Items[k] = v

		// If bulk size was reached or all entries have been processed
		if idx%bulkSizeCreate == 0 || idx == len(plan.Items) {

			// Parse body of the request to string
			body := bulk.toBody(ctx, PortGroups{})

			// Execute request
			urlPath := bulk.getPath() + "?bulk=true"
			res, err := r.client.Post(urlPath, body, reqMods...)
			if err != nil {
				return state, diag.Diagnostics{
					diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to create a bulk (POST) id: %s, got error: %s, %s", state.Id.ValueString(), err, res.String())),
				}
			}

			// Read result and save it to the state
			bulk.fromBodyUnknowns(ctx, res)
			for k, v := range bulk.Items {
				state.Items[k] = v
			}

			// Clear bulk item for next run
			bulk.Items = make(map[string]PortGroupsItems, bulkSizeCreate)
		}
	}

	return state, nil
}

// End of section. //template:end createSubresources

// Section below is generated&owned by "gen/generator.go". //template:begin deleteSubresources
// deleteSubresources takes list of objects and deletes them either in bulk, or one-by-one, depending on FMC version
func (r *PortGroupsResource) deleteSubresources(ctx context.Context, state, plan PortGroups, reqMods ...func(*fmc.Req)) (PortGroups, diag.Diagnostics) {
	objectsToRemove := plan.Clone()

	tflog.Debug(ctx, fmt.Sprintf("%s: Deleting bulk of objects", state.Id.ValueString()))
	tflog.Debug(ctx, fmt.Sprintf("%s: One-by-one deletion mode", state.Id.ValueString()))
	for k, v := range objectsToRemove.Items {
		// Check if the object was not already deleted
		if v.Id.IsNull() {
			delete(state.Items, k)
			continue
		}

		urlPath := state.getPath() + "/" + url.QueryEscape(v.Id.ValueString())
		res, err := r.client.Delete(urlPath, reqMods...)
		if err != nil {
			return state, diag.Diagnostics{
				diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("%s: Failed to delete object (DELETE) id %s, got error: %s, %s", state.Id.ValueString(), v.Id.ValueString(), err, res.String())),
			}
		}

		// Remove deleted item from state
		delete(state.Items, k)
	}

	return state, nil
}

// End of section. //template:end deleteSubresources

// Section below is generated&owned by "gen/generator.go". //template:begin updateSubresources

// updateSubresources take elements one-by-one and updates them, as bulks are not supported
func (r *PortGroupsResource) updateSubresources(ctx context.Context, state, plan PortGroups, reqMods ...func(*fmc.Req)) (PortGroups, diag.Diagnostics) {
	var tmpObject PortGroups
	tmpObject.Items = make(map[string]PortGroupsItems, 1)

	tflog.Debug(ctx, fmt.Sprintf("%s: Updating bulk of objects", state.Id.ValueString()))

	for k, v := range plan.Items {
		tmpObject.Items[k] = v

		body := tmpObject.toBodyNonBulk(ctx, state)
		urlPath := tmpObject.getPath() + "/" + url.QueryEscape(v.Id.ValueString())
		res, err := r.client.Put(urlPath, body, reqMods...)
		if err != nil {
			return state, diag.Diagnostics{
				diag.NewErrorDiagnostic("Client Error", fmt.Sprintf("Failed to update object (PUT) id %s, got error: %s, %s", state.Id.ValueString(), err, res.String())),
			}
		}

		// Update state
		state.Items[k] = v

		// Clear tmpObject.Items
		delete(tmpObject.Items, k)
	}

	return state, nil
}

// End of section. //template:end updateSubresources