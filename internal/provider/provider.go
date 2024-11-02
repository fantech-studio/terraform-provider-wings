package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	tffunc "github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"fantech.dev/terraform-provider-wings/internal/model"
)

const (
	headerKeyID       = "X-API-KEY-ID"
	headerKey         = "X-API-KEY"
	headerUA          = "User-Agent"
	headerContentType = "Content-Type"
)

const (
	applicationJSON = "application/json"
)

var _ provider.Provider = &WingsProvider{}

type WingsProvider struct {
	version string
	config  *config
}

type wingsProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	APIKeyID types.String `tfsdk:"api_key_id"`
	APIKey   types.String `tfsdk:"api_key"`
}

func (p *WingsProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "wings"
}

func (p *WingsProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Wings.",
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				Required: true,
			},
			"api_key_id": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
			"api_key": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *WingsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var cfg wingsProviderModel
	diags := req.Config.Get(ctx, &cfg)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if cfg.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown Wings API Endpoint",
			"The provider cannot create the Wings API client as there is an unknown configuration value for the Wings API endpoint. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the WINGS_ENDPOINT environment variable.",
		)
	}

	if cfg.APIKeyID.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key_id"),
			"Unknown Wings API Key ID",
			"The provider cannot create the Wings API client as there is an unknown configuration value for the Wings API Key ID. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the WINGS_API_KEY_ID environment variable.",
		)
	}

	if cfg.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown Wings API Key",
			"The provider cannot create the Wings API client as there is an unknown configuration value for the Wings API Key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the WINGS_API_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	endpoint := os.Getenv("WINGS_ENDPOINT")
	apiKeyID := os.Getenv("WINGS_API_KEY_ID")
	apiKey := os.Getenv("WINGS_API_KEY")

	if !cfg.Endpoint.IsNull() {
		endpoint = cfg.Endpoint.ValueString()
	}

	if !cfg.APIKeyID.IsNull() {
		apiKeyID = cfg.APIKeyID.ValueString()
	}

	if !cfg.APIKey.IsNull() {
		apiKey = cfg.APIKey.ValueString()
	}

	if endpoint == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Missing Wings API Endpoint",
			"The provider cannot create the Wings API client as there is a missing or empty value for the Wings API endpoint. "+
				"Set the endpoint value in the configuration or use the WINGS_ENDPOINT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiKeyID == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key_id"),
			"Missing Wings API Key",
			"The provider cannot create the Wings API client as there is a missing or empty value for the Wings API Key. "+
				"Set the api_key value in the configuration or use the WINGS_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing Wings API Key",
			"The provider cannot create the Wings API client as there is a missing or empty value for the Wings API Key. "+
				"Set the api_key value in the configuration or use the WINGS_API_KEY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "wings_endpoint", endpoint)
	ctx = tflog.SetField(ctx, "wings_api_key_id", apiKeyID)
	ctx = tflog.SetField(ctx, "wings_api_key", apiKey)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "wings_api_key")

	tflog.Debug(ctx, "Creating Wings client")

	if p.config == nil {
		retryClient := retryablehttp.NewClient()
		retryClient.RetryMax = 5
		rc := retryClient.StandardClient()
		p.config = &config{
			ua:       "terraform-provider-wings",
			keyID:    apiKeyID,
			key:      apiKey,
			endpoint: endpoint,
			client:   rc,
		}
	}

	resp.DataSourceData = p.config
	resp.ResourceData = p.config

	tflog.Info(ctx, "Configured Wings client", map[string]any{"success": true})
}

func (p *WingsProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *WingsProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewValueResource,
	}
}

func (p *WingsProvider) Functions(_ context.Context) []func() tffunc.Function {
	return []func() tffunc.Function{
		NewUnixTimeConverterFunc,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &WingsProvider{
			version: version,
		}
	}
}

type config struct {
	ua       string
	keyID    string
	key      string
	endpoint string
	client   *http.Client
}

func (c *config) GetValue(ctx context.Context, id string) (*model.Value, error) {
	u, err := url.JoinPath(c.endpoint, "values", id)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, %s", resp.StatusCode, string(b))
	}

	value := new(model.Value)
	err = json.NewDecoder(resp.Body).Decode(value)
	return value, err
}

func (c *config) CreateValue(ctx context.Context, value *model.Value) (*model.Value, error) {
	u, err := url.JoinPath(c.endpoint, "values")
	if err != nil {
		return nil, err
	}

	j, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(j))
	if err != nil {
		return nil, err
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, %s", resp.StatusCode, string(b))
	}

	v := new(model.Value)
	err = json.NewDecoder(resp.Body).Decode(v)
	return v, err
}

func (c *config) UpdateValue(ctx context.Context, value *model.Value) (*model.Value, error) {
	u, err := url.JoinPath(c.endpoint, "values", value.ID)
	if err != nil {
		return nil, err
	}

	j, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPut, u, bytes.NewReader(j))
	if err != nil {
		return nil, err
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, %s", resp.StatusCode, string(b))
	}

	v := new(model.Value)
	err = json.NewDecoder(resp.Body).Decode(v)
	return v, err
}

func (c *config) DeleteValue(ctx context.Context, id string) error {
	u, err := url.JoinPath(c.endpoint, "values", id)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, u, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode != http.StatusNotFound {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, %s", resp.StatusCode, string(b))
	}

	return nil
}

func (c *config) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req.Header.Set(headerKeyID, c.keyID)
	req.Header.Set(headerKey, c.key)
	req.Header.Set(headerUA, c.ua)
	req.Header.Set(headerContentType, applicationJSON)
	req.WithContext(ctx)
	return c.client.Do(req)
}
