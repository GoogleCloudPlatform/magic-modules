package monitoring

import "testing"

func TestMonitoringDashboardDiffSuppress(t *testing.T) {
	tests := []struct {
		name string
		old  string
		new  string
		want bool
	}{
		{
			name: "identical dashboard",
			old:  `{"displayName":"Dashboard","gridLayout":{"widgets":[{"blank":{}}]}}`,
			new:  `{"displayName":"Dashboard","gridLayout":{"widgets":[{"blank":{}}]}}`,
			want: true,
		},
		{
			name: "API fields and defaults",
			old: `{
				"name": "projects/my-project/dashboards/abc",
				"etag": "1234",
				"displayName": "Dashboard",
				"mosaicLayout": {
					"tiles": [
						{"height": 4, "width": 4},
						{"widget": {"xyChart": {
							"dataSets": [{"targetAxis": "Y1"}],
							"yAxis": {"scale": "LINEAR"}
						}}},
						{"widget": {"timeSeriesTable": {"metricVisualization": "NUMBER"}}}
					]
				}
			}`,
			new: `{
				"displayName": "Dashboard",
				"mosaicLayout": {
					"tiles": [
						{"height": 4, "width": 4, "xPos": 0, "yPos": 0},
						{"widget": {"xyChart": {"dataSets": [{}], "yAxis": {}}}},
						{"widget": {"timeSeriesTable": {}}}
					]
				}
			}`,
			want: true,
		},
		{
			name: "exported computed fields",
			old:  `{"name":"projects/my-project/dashboards/abc","etag":"old","displayName":"Dashboard"}`,
			new:  `{"name":"projects/my-project/dashboards/abc","etag":"new","displayName":"Dashboard"}`,
			want: true,
		},
		{
			name: "scalar array preserved",
			old:  `{"mosaicLayout":{"tiles":[{"widget":{"xyChart":{"dataSets":[{"timeSeriesQuery":{"timeSeriesFilter":{"aggregation":{"groupByFields":["metric.label.status"]}}}}]}}}]}}`,
			new:  `{"mosaicLayout":{"tiles":[{"widget":{"xyChart":{"dataSets":[{"timeSeriesQuery":{"timeSeriesFilter":{"aggregation":{"groupByFields":["metric.label.status"]}}}}]}}}]}}`,
			want: true,
		},
		{
			name: "display name changed",
			old:  `{"name":"projects/my-project/dashboards/abc","etag":"1234","displayName":"Old"}`,
			new:  `{"displayName":"New"}`,
			want: false,
		},
		{
			name: "unsupported config field added",
			old:  `{"displayName":"Dashboard"}`,
			new:  `{"displayName":"Dashboard","description":"Not supported by the API"}`,
			want: false,
		},
		{
			name: "nonzero tile position added",
			old:  `{"mosaicLayout":{"tiles":[{"width":4,"height":4}]}}`,
			new:  `{"mosaicLayout":{"tiles":[{"width":4,"height":4,"xPos":1}]}}`,
			want: false,
		},
		{
			name: "nonzero tile position changed to zero",
			old:  `{"mosaicLayout":{"tiles":[{"width":4,"height":4,"xPos":4}]}}`,
			new:  `{"mosaicLayout":{"tiles":[{"width":4,"height":4,"xPos":0}]}}`,
			want: false,
		},
		{
			name: "widget removed",
			old:  `{"gridLayout":{"widgets":[{"blank":{}},{"text":{"content":"remove","format":"MARKDOWN"}}]}}`,
			new:  `{"gridLayout":{"widgets":[{"blank":{}}]}}`,
			want: false,
		},
		{
			name: "widget added",
			old:  `{"gridLayout":{"widgets":[{"blank":{}}]}}`,
			new:  `{"gridLayout":{"widgets":[{"blank":{}},{"text":{"content":"add","format":"MARKDOWN"}}]}}`,
			want: false,
		},
		{
			name: "nested type changed",
			old:  `{"gridLayout":{"widgets":[{"blank":{"nested":{"value":1}}}]}}`,
			new:  `{"gridLayout":{"widgets":[{"blank":{"nested":"value"}}]}}`,
			want: false,
		},
		{
			name: "invalid state JSON",
			old:  `{not-json`,
			new:  `{"displayName":"Dashboard"}`,
			want: false,
		},
		{
			name: "invalid config JSON",
			old:  `{"displayName":"Dashboard"}`,
			new:  `{not-json`,
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := monitoringDashboardDiffSuppress("dashboard_json", test.old, test.new, nil)
			if got != test.want {
				t.Errorf("monitoringDashboardDiffSuppress() = %t, want %t", got, test.want)
			}
		})
	}
}

func TestRemoveKeysAbsentFromConfigDoesNotMutateState(t *testing.T) {
	state := map[string]interface{}{
		"mosaicLayout": map[string]interface{}{
			"tiles": []interface{}{
				map[string]interface{}{
					"widget": map[string]interface{}{
						"xyChart": map[string]interface{}{
							"dataSets": []interface{}{
								map[string]interface{}{"targetAxis": "Y1"},
							},
						},
					},
				},
			},
		},
	}
	config := map[string]interface{}{
		"mosaicLayout": map[string]interface{}{
			"tiles": []interface{}{
				map[string]interface{}{
					"widget": map[string]interface{}{
						"xyChart": map[string]interface{}{
							"dataSets": []interface{}{map[string]interface{}{}},
						},
					},
				},
			},
		},
	}

	removeKeysAbsentFromConfig(state, config)

	tiles := state["mosaicLayout"].(map[string]interface{})["tiles"].([]interface{})
	widget := tiles[0].(map[string]interface{})["widget"].(map[string]interface{})
	dataSets := widget["xyChart"].(map[string]interface{})["dataSets"].([]interface{})
	if _, ok := dataSets[0].(map[string]interface{})["targetAxis"]; !ok {
		t.Fatal("removeKeysAbsentFromConfig() mutated state")
	}
}

func TestStripMonitoringDashboardZeroTilePositionsDoesNotMutateInput(t *testing.T) {
	dashboard := map[string]interface{}{
		"mosaicLayout": map[string]interface{}{
			"tiles": []interface{}{map[string]interface{}{"xPos": float64(0)}},
		},
	}

	stripMonitoringDashboardZeroTilePositions(dashboard)

	tiles := dashboard["mosaicLayout"].(map[string]interface{})["tiles"].([]interface{})
	if _, ok := tiles[0].(map[string]interface{})["xPos"]; !ok {
		t.Fatal("stripMonitoringDashboardZeroTilePositions() mutated input")
	}
}
