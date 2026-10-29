package ipgeolocation

import (
	"context"
	"reflect"
	"sort"
	"testing"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

func TestPluginTables(t *testing.T) {
	p := Plugin(context.Background())
	var names []string
	for name, table := range p.TableMap {
		names = append(names, name)
		if table.List == nil || table.List.Hydrate == nil {
			t.Fatalf("%s has no query handler", name)
		}
		seen := map[string]bool{}
		for _, column := range table.Columns {
			if seen[column.Name] {
				t.Fatalf("duplicate %s column in %s", column.Name, name)
			}
			seen[column.Name] = true
		}
		if name != "ipgeolocation_asn" && (len(table.List.KeyColumns) != 1 || table.List.KeyColumns[0].Name != "ip" || table.List.KeyColumns[0].Require != plugin.Required) {
			t.Fatalf("%s must require an IP", name)
		}
	}
	sort.Strings(names)
	want := []string{"ipgeolocation_abuse", "ipgeolocation_asn", "ipgeolocation_ip", "ipgeolocation_security"}
	if p.Name != "ipgeolocation" || !reflect.DeepEqual(names, want) || p.ConnectionConfigSchema.NewInstance() == nil {
		t.Fatalf("unexpected plugin definition: %v", names)
	}
}

func TestFractionalTimezoneOffsets(t *testing.T) {
	for _, offset := range []float64{5.5, 5.75, -3.5, 0} {
		row := map[string]interface{}{"time_zone": map[string]interface{}{"offset": offset, "offset_with_dst": offset + 1}}
		found := 0
		for _, c := range ipGeolocationColumns() {
			if c.Name != "timezone_offset" && c.Name != "timezone_offset_with_dst" {
				continue
			}
			found++
			if c.Type != proto.ColumnType_DOUBLE {
				t.Fatalf("%s must preserve fractional hours", c.Name)
			}
			got, err := c.Transform.Execute(context.Background(), &transform.TransformData{HydrateItem: row})
			want := offset
			if c.Name == "timezone_offset_with_dst" {
				want++
			}
			if err != nil || got != want {
				t.Fatalf("%s=%v, want %v, error=%v", c.Name, got, want, err)
			}
		}
		if found != 2 {
			t.Fatal("both timezone offset columns must exist")
		}
	}
}

func TestNormaliseASN(t *testing.T) {
	for input, want := range map[string]string{"15169": "15169", "AS15169": "15169", " as15169 ": "15169", "": ""} {
		if got := normaliseASN(input); got != want {
			t.Fatalf("%q => %q, want %q", input, got, want)
		}
	}
	_, err := listASN(testContext(), &plugin.QueryData{}, nil)
	if err == nil {
		t.Fatal("an ASN query without an IP or ASN must fail")
	}
}
