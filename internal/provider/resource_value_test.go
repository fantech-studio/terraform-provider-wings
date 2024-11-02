package provider

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

//go:embed testdata/bool.json
var boolTestdata string

func TestAccResourceWingsValue_BoolValue(t *testing.T) {
	mock := httpmock.NewMockTransport()
	mock.RegisterResponder(
		http.MethodGet,
		"http://localhost:8018/values/test-bool-value",
		httpmock.NewStringResponder(200, boolTestdata),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/values",
		httpmock.NewStringResponder(200, boolTestdata),
	)
	mock.RegisterResponder(
		http.MethodPut,
		"http://localhost:8018/values/test-bool-value",
		httpmock.NewStringResponder(200, boolTestdata),
	)
	mock.RegisterResponder(
		http.MethodDelete,
		"http://localhost:8018/values/test-bool-value",
		httpmock.NewStringResponder(204, boolTestdata),
	)

	client := &http.Client{
		Transport: mock,
	}
	cfg := &config{
		endpoint: "http://localhost:8018",
		client:   client,
	}
	client.Transport = mock

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(cfg),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccResourceBool(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("wings_value.test-bool-value", "value_id", "test-bool-value"),
					resource.TestCheckResourceAttr("wings_value.test-bool-value", "enabled", "true"),
					resource.TestCheckResourceAttr("wings_value.test-bool-value", "description", "test bool value"),
					resource.TestCheckResourceAttr("wings_value.test-bool-value", "default_variant", "off"),
					resource.TestCheckResourceAttr("wings_value.test-bool-value", "bool.#", "2"),
					resource.TestCheckResourceAttr("wings_value.test-bool-value", "bool.0.variant", "on"),
					resource.TestCheckResourceAttr("wings_value.test-bool-value", "bool.0.value", "true"),
					resource.TestCheckResourceAttr("wings_value.test-bool-value", "bool.1.variant", "off"),
					resource.TestCheckResourceAttr("wings_value.test-bool-value", "bool.1.value", "false"),
					resource.TestCheckResourceAttr("wings_value.test-bool-value", "targeting.#", "2"),
					resource.TestCheckResourceAttr("wings_value.test-bool-value", "targeting.0.variant", "on"),
					resource.TestCheckResourceAttr("wings_value.test-bool-value", "targeting.0.expr", "env == 'dev'"),
					resource.TestCheckResourceAttr("wings_value.test-bool-value", "targeting.1.variant", "on"),
					resource.TestCheckResourceAttr("wings_value.test-bool-value", "targeting.1.expr", "userId == 'XXX'"),
					resource.TestCheckResourceAttr("wings_value.test-bool-value", "test.#", "1"),
					resource.TestCheckResourceAttr("wings_value.test-bool-value", "test.0.variables", "{\"count\":1,\"env\":\"test\"}"),
					resource.TestCheckResourceAttr("wings_value.test-bool-value", "test.0.expected", "on"),
				),
			},
		},
	})
}

//go:embed testdata/int.json
var intTestdata string

func TestAccResourceWingsValue_IntValue(t *testing.T) {
	mock := httpmock.NewMockTransport()
	mock.RegisterResponder(
		http.MethodGet,
		"http://localhost:8018/values/test-integer-value",
		httpmock.NewStringResponder(200, intTestdata),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/values",
		httpmock.NewStringResponder(200, intTestdata),
	)
	mock.RegisterResponder(
		http.MethodPut,
		"http://localhost:8018/values/test-integer-value",
		httpmock.NewStringResponder(200, intTestdata),
	)
	mock.RegisterResponder(
		http.MethodDelete,
		"http://localhost:8018/values/test-integer-value",
		httpmock.NewStringResponder(204, intTestdata),
	)

	client := &http.Client{
		Transport: mock,
	}
	cfg := &config{
		endpoint: "http://localhost:8018",
		client:   client,
	}
	client.Transport = mock

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(cfg),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccResourceInt(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("wings_value.test-integer-value", "value_id", "test-integer-value"),
					resource.TestCheckResourceAttr("wings_value.test-integer-value", "enabled", "true"),
					resource.TestCheckResourceAttr("wings_value.test-integer-value", "description", "test integer value"),
					resource.TestCheckResourceAttr("wings_value.test-integer-value", "default_variant", "one"),
					resource.TestCheckResourceAttr("wings_value.test-integer-value", "int.#", "1"),
					resource.TestCheckResourceAttr("wings_value.test-integer-value", "int.0.variant", "one"),
					resource.TestCheckResourceAttr("wings_value.test-integer-value", "int.0.value", "1"),
					resource.TestCheckResourceAttr("wings_value.test-integer-value", "targeting.#", "0"),
				),
			},
		},
	})
}

//go:embed testdata/string.json
var stringTestdata string

func TestAccResourceWingsValue_StringValue(t *testing.T) {
	mock := httpmock.NewMockTransport()
	mock.RegisterResponder(
		http.MethodGet,
		"http://localhost:8018/values/test-string-value",
		httpmock.NewStringResponder(200, stringTestdata),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/values",
		httpmock.NewStringResponder(200, stringTestdata),
	)
	mock.RegisterResponder(
		http.MethodPut,
		"http://localhost:8018/values/test-string-value",
		httpmock.NewStringResponder(200, stringTestdata),
	)
	mock.RegisterResponder(
		http.MethodDelete,
		"http://localhost:8018/values/test-string-value",
		httpmock.NewStringResponder(204, stringTestdata),
	)

	client := &http.Client{
		Transport: mock,
	}
	cfg := &config{
		endpoint: "http://localhost:8018",
		client:   client,
	}
	client.Transport = mock

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(cfg),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccResourceString(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("wings_value.test-string-value", "value_id", "test-string-value"),
					resource.TestCheckResourceAttr("wings_value.test-string-value", "enabled", "true"),
					resource.TestCheckResourceAttr("wings_value.test-string-value", "description", "test string value"),
					resource.TestCheckResourceAttr("wings_value.test-string-value", "default_variant", "key"),
					resource.TestCheckResourceAttr("wings_value.test-string-value", "string.#", "1"),
					resource.TestCheckResourceAttr("wings_value.test-string-value", "string.0.variant", "key"),
					resource.TestCheckResourceAttr("wings_value.test-string-value", "string.0.value", "test value"),
					resource.TestCheckResourceAttr("wings_value.test-string-value", "targeting.#", "0"),
				),
			},
		},
	})
}

//go:embed testdata/object.json
var objectTestdata string

func TestAccResourceWingsValue_ObjectValue(t *testing.T) {
	mock := httpmock.NewMockTransport()
	mock.RegisterResponder(
		http.MethodGet,
		"http://localhost:8018/values/test-json-value",
		httpmock.NewStringResponder(200, objectTestdata),
	)
	mock.RegisterResponder(
		http.MethodPost,
		"http://localhost:8018/values",
		httpmock.NewStringResponder(200, objectTestdata),
	)
	mock.RegisterResponder(
		http.MethodPut,
		"http://localhost:8018/values/test-json-value",
		httpmock.NewStringResponder(200, objectTestdata),
	)
	mock.RegisterResponder(
		http.MethodDelete,
		"http://localhost:8018/values/test-json-value",
		httpmock.NewStringResponder(204, objectTestdata),
	)

	client := &http.Client{
		Transport: mock,
	}
	cfg := &config{
		endpoint: "http://localhost:8018",
		client:   client,
	}
	client.Transport = mock

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: protoV6ProviderFactories(cfg),
		Steps: []resource.TestStep{
			{
				Config: providerConfig + testAccResourceObject(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("wings_value.test-json-value", "value_id", "test-json-value"),
					resource.TestCheckResourceAttr("wings_value.test-json-value", "enabled", "true"),
					resource.TestCheckResourceAttr("wings_value.test-json-value", "description", "test json value"),
					resource.TestCheckResourceAttr("wings_value.test-json-value", "default_variant", "json"),
					resource.TestCheckResourceAttr("wings_value.test-json-value", "object.#", "1"),
					resource.TestCheckResourceAttr("wings_value.test-json-value", "object.0.variant", "json"),
					resource.TestCheckResourceAttr("wings_value.test-json-value", "object.0.value", "{\"items\":[{\"content\":\"content1\",\"viewable\":true},{\"content\":\"content2\",\"viewable\":true},{\"content\":\"content3\",\"viewable\":false}]}"),
					resource.TestCheckResourceAttr("wings_value.test-json-value", "object.0.transform.#", "2"),
					resource.TestCheckResourceAttr("wings_value.test-json-value", "object.0.transform.0.expr", "{\"items\":items.map(item, item.viewable ? item : item.deleteKey([\"content\"]))}"),
					resource.TestCheckResourceAttr("wings_value.test-json-value", "object.0.transform.1.expr", "{\"items\":items.map(item, item.viewable ? item.selectKey([\"content\"]) : item)}"),
					resource.TestCheckResourceAttr("wings_value.test-json-value", "targeting.#", "0"),
				),
			},
		},
	})
}

func testAccResourceBool() string {
	return `
resource "wings_value" "test-bool-value" {
  value_id = "test-bool-value"
  enabled = true
  description = "test bool value"
  default_variant = "off"

  bool {
	variant = "on"
	value = true
  }

  bool {
	variant = "off"
	value = false
  }

  targeting {
    variant = "on"
    expr = "env == 'dev'"
  }

  targeting {
    variant = "on"
    expr = "userId == 'XXX'"
  }
	
  test {
	variables = jsonencode({
	  env = "test"
	  count = 1
	})
	expected = "on"
  }
}`
}

func testAccResourceString() string {
	return `
resource "wings_value" "test-string-value" {
  value_id = "test-string-value"
  enabled = true
  description = "test string value"
  default_variant = "key"

  string {
	variant = "key"
	value = "test value"
  }
}`
}

func testAccResourceObject() string {
	return `
resource "wings_value" "test-json-value" {
  value_id = "test-json-value"
  enabled = true
  description = "test json value"
  default_variant = "json"

  object {
	variant = "json"
	value = jsonencode({
	  "items": [
		{"viewable": true, "content": "content1"},
		{"viewable": true, "content": "content2"},
		{"viewable": false, "content": "content3"}
	  ]	
	})
	transform {
	  expr = "{\"items\":items.map(item, item.viewable ? item : item.deleteKey([\"content\"]))}"
	}
	transform {
	  expr = "{\"items\":items.map(item, item.viewable ? item.selectKey([\"content\"]) : item)}"
	}
  }
}`
}

func testAccResourceInt() string {
	return `
resource "wings_value" "test-integer-value" {
  value_id = "test-integer-value"
  enabled = true
  description = "test integer value"
  default_variant = "one"

  int {
	variant = "one"
	value = 1
  }
}`
}
