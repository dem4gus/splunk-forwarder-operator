package kube

import (
	"fmt"
	"testing"

	sfv1alpha1 "github.com/openshift/splunk-forwarder-operator/api/v1alpha1"
	"github.com/openshift/splunk-forwarder-operator/internal/testutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func TestGenerateConfigMaps(t *testing.T) {
	var testInstance = testutil.NewSplunkForwarderCR().
		WithGeneration(10).
		WithSplunkInputs([]sfv1alpha1.SplunkForwarderInputs{
			{
				Path:      "",
				Index:     "test-index",
				WhiteList: ".*log$",
				BlackList: ".*bak$",
			},
			{
				Path:      "/var/derp",
				Index:     "test-index",
				WhiteList: ".*log$",
				BlackList: ".*bak$",
			},
			{
				Path:       "/var/derp.text",
				SourceType: "text",
				WhiteList:  ".*log$",
				BlackList:  ".*bak$",
			},
		}).
		Build()
	type args struct {
		instance       *sfv1alpha1.SplunkForwarder
		namespacedName types.NamespacedName
		clusterid      string
	}
	tests := []struct {
		name string
		args args
		want []*corev1.ConfigMap
	}{
		{
			name: "Generates metadata and local ConfigMaps with correct Splunk inputs configuration",
			args: args{
				instance:       testInstance,
				namespacedName: types.NamespacedName{Namespace: testutil.InstanceNamespace, Name: testutil.InstanceName},
				clusterid:      "test",
			},
			want: []*corev1.ConfigMap{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "osd-monitored-logs-metadata",
						Namespace: testutil.InstanceNamespace,
						Labels: map[string]string{
							"app": testutil.InstanceName,
						},
						Annotations: map[string]string{
							"genVersion": "10",
						},
					},
					Data: map[string]string{
						"local.meta": `
[]
access = read : [ * ], write : [ admin ]
export = system
`,
					},
				},
				{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "osd-monitored-logs-local",
						Namespace: testInstance.Namespace,
						Labels: map[string]string{
							"app": testInstance.Name,
						},
						Annotations: map[string]string{
							"genVersion": "10",
						},
					},
					Data: map[string]string{
						"app.conf": `
[install]
state = enabled

[package]
check_for_updates = false

[ui]
is_visible = false
is_manageable = false
`,
						"inputs.conf": `[monitor:///var/derp]
sourcetype = _json
index = test-index
whitelist = .*log$
blacklist = .*bak$
_meta = clusterid::test
disabled = false

[monitor:///var/derp.text]
sourcetype = text
index = main
whitelist = .*log$
blacklist = .*bak$
_meta = clusterid::test
disabled = false

`,
						"props.conf": fmt.Sprintf(`
[_json]
TRUNCATE = %d
`, MaxEventSize),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateConfigMaps(tt.args.instance, tt.args.namespacedName, tt.args.clusterid)
			testutil.DeepEqualWithDiff(t, tt.want, got)
		})
	}
}
