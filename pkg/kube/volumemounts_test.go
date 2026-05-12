package kube

import (
	"testing"

	"github.com/openshift/splunk-forwarder-operator/config"
	"github.com/openshift/splunk-forwarder-operator/internal/testutil"
	corev1 "k8s.io/api/core/v1"
)

func TestGetVolumeMounts(t *testing.T) {
	var mountPropagationMode = corev1.MountPropagationHostToContainer
	testInstance := testutil.NewSplunkForwarderCR().WithImageTag(testutil.ImageTag).Build()

	type args struct {
		useHECToken bool
	}
	tests := []struct {
		name string
		args args
		want []corev1.VolumeMount
	}{
		{
			name: "Returns volume mounts for universal forwarder with auth secret",
			args: args{
				useHECToken: false,
			},
			want: []corev1.VolumeMount{
				{
					Name:      config.SplunkAuthSecretName,
					MountPath: "/opt/splunkforwarder/etc/apps/splunkauth/default",
				},
				{
					Name:      config.SplunkAuthSecretName,
					MountPath: "/opt/splunkforwarder/etc/apps/splunkauth/local",
				},
				{
					Name:      config.SplunkAuthSecretName,
					MountPath: "/opt/splunkforwarder/etc/apps/splunkauth/metadata",
				},
				{
					Name:      "osd-monitored-logs-local",
					MountPath: "/opt/splunkforwarder/etc/apps/osd_monitored_logs/local",
				},
				{
					Name:      "osd-monitored-logs-metadata",
					MountPath: "/opt/splunkforwarder/etc/apps/osd_monitored_logs/metadata",
				},
				{
					Name:      "splunk-state",
					MountPath: "/opt/splunkforwarder/var/lib",
				},
				{
					Name:             "host",
					MountPath:        "/host",
					MountPropagation: &mountPropagationMode,
					ReadOnly:         true,
				},
			},
		},
		{
			name: "Returns volume mounts for HEC token configuration without mTLS auth secret",
			args: args{
				useHECToken: true,
			},
			want: []corev1.VolumeMount{
				{
					Name:      "splunk-config",
					MountPath: "/opt/splunkforwarder/etc/system/local",
				},
				{
					Name:      "osd-monitored-logs-local",
					MountPath: "/opt/splunkforwarder/etc/apps/osd_monitored_logs/local",
				},
				{
					Name:      "osd-monitored-logs-metadata",
					MountPath: "/opt/splunkforwarder/etc/apps/osd_monitored_logs/metadata",
				},
				{
					Name:      "splunk-state",
					MountPath: "/opt/splunkforwarder/var/lib",
				},
				{
					Name:             "host",
					MountPath:        "/host",
					MountPropagation: &mountPropagationMode,
					ReadOnly:         true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetVolumeMounts(testInstance, tt.args.useHECToken)
			testutil.DeepEqualWithDiff(t, tt.want, got)
		})
	}
}
