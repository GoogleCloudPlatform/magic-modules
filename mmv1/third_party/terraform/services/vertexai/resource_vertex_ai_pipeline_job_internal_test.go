package vertexai

import (
	"testing"
)

func TestUnitVertexAIPipelineJob_PipelineSpecDiffSuppress(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		old      string
		new      string
		expected bool
	}{
		{
			name:     "both empty",
			old:      "",
			new:      "",
			expected: true,
		},
		{
			name:     "old empty new set",
			old:      "",
			new:      `{"pipelineInfo":{"name":"test"}}`,
			expected: false,
		},
		{
			name:     "identical JSON with different key order",
			old:      `{"schemaVersion":"2.1.0","pipelineInfo":{"name":"test"}}`,
			new:      `{"pipelineInfo":{"name":"test"},"schemaVersion":"2.1.0"}`,
			expected: true,
		},
		{
			name:     "API returns both deploymentSpec and deploymentConfig when config only has deploymentSpec",
			old:      `{"deploymentConfig":{"@type":"type.googleapis.com/ml_pipelines.PipelineDeploymentConfig","executors":{"exec":{"container":{"image":"python:3.11"}}}},"deploymentSpec":{"executors":{"exec":{"container":{"image":"python:3.11"}}}},"pipelineInfo":{"name":"test"},"schemaVersion":"2.1.0"}`,
			new:      `{"deploymentSpec":{"executors":{"exec":{"container":{"image":"python:3.11"}}}},"pipelineInfo":{"name":"test"},"schemaVersion":"2.1.0"}`,
			expected: true,
		},
		{
			name:     "different image in deploymentSpec is not suppressed",
			old:      `{"deploymentConfig":{"@type":"type.googleapis.com/ml_pipelines.PipelineDeploymentConfig","executors":{"exec":{"container":{"image":"python:3.10"}}}},"deploymentSpec":{"executors":{"exec":{"container":{"image":"python:3.10"}}}},"pipelineInfo":{"name":"test"},"schemaVersion":"2.1.0"}`,
			new:      `{"deploymentSpec":{"executors":{"exec":{"container":{"image":"python:3.11"}}}},"pipelineInfo":{"name":"test"},"schemaVersion":"2.1.0"}`,
			expected: false,
		},
		{
			name:     "different pipelineInfo is not suppressed",
			old:      `{"pipelineInfo":{"name":"old-name"},"schemaVersion":"2.1.0"}`,
			new:      `{"pipelineInfo":{"name":"new-name"},"schemaVersion":"2.1.0"}`,
			expected: false,
		},
		{
			name:     "invalid JSON returns false",
			old:      `{invalid}`,
			new:      `{"pipelineInfo":{"name":"test"}}`,
			expected: false,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := vertexAIPipelineJobPipelineSpecDiffSuppress("pipeline_spec", tc.old, tc.new, nil)
			if got != tc.expected {
				t.Errorf("vertexAIPipelineJobPipelineSpecDiffSuppress(%q, %q) = %v, want %v", tc.old, tc.new, got, tc.expected)
			}
		})
	}
}
