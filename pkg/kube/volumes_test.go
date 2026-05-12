package kube

import (
	"testing"

	"github.com/openshift/splunk-forwarder-operator/config"
	"github.com/openshift/splunk-forwarder-operator/internal/testutil"
	corev1 "k8s.io/api/core/v1"
)

func TestGetVolumes(t *testing.T) {
	type args struct {
		mountHost     bool
		mountSecret   bool
		mountHECToken bool
		instanceName  string
	}
	tests := []struct {
		name string
		args args
		want []corev1.Volume
	}{
		{
			name: "Returns volumes with host mount and without secret",
			args: args{
				mountHost:    true,
				mountSecret:  false,
				instanceName: testutil.InstanceName,
			},
			want: []corev1.Volume{
				testutil.NewConfigMapVolume("osd-monitored-logs-local"),
				testutil.NewConfigMapVolume("osd-monitored-logs-metadata"),
				testutil.NewHostPathVolume("splunk-state", testutil.SplunkStatePath),
				testutil.NewHostPathVolume("host", testutil.HostRootPath),
				testutil.NewConfigMapVolume(testutil.InstanceName + "-internalsplunk"),
			},
		},
		{
			name: "Returns volumes with both host mount and auth secret",
			args: args{
				mountHost:    true,
				mountSecret:  true,
				instanceName: testutil.InstanceName,
			},
			want: []corev1.Volume{
				testutil.NewConfigMapVolume("osd-monitored-logs-local"),
				testutil.NewConfigMapVolume("osd-monitored-logs-metadata"),
				testutil.NewHostPathVolume("splunk-state", testutil.SplunkStatePath),
				testutil.NewHostPathVolume("host", testutil.HostRootPath),
				testutil.NewSecretVolume(config.SplunkAuthSecretName, config.SplunkAuthSecretName),
			},
		},
		{
			name: "Returns volumes without host mount or secret",
			args: args{
				mountHost:    false,
				mountSecret:  false,
				instanceName: testutil.InstanceName,
			},
			want: []corev1.Volume{
				testutil.NewConfigMapVolume("osd-monitored-logs-local"),
				testutil.NewConfigMapVolume("osd-monitored-logs-metadata"),
				testutil.NewHostPathVolume("splunk-state", testutil.SplunkStatePath),
				testutil.NewConfigMapVolume(testutil.InstanceName + "-hfconfig"),
				testutil.NewConfigMapVolume(testutil.InstanceName + "-internalsplunk"),
			},
		},
		{
			name: "Returns volumes with secret but without host mount",
			args: args{
				mountHost:    false,
				mountSecret:  true,
				instanceName: testutil.InstanceName,
			},
			want: []corev1.Volume{
				testutil.NewConfigMapVolume("osd-monitored-logs-local"),
				testutil.NewConfigMapVolume("osd-monitored-logs-metadata"),
				testutil.NewHostPathVolume("splunk-state", testutil.SplunkStatePath),
				testutil.NewConfigMapVolume(testutil.InstanceName + "-hfconfig"),
				testutil.NewSecretVolume(config.SplunkAuthSecretName, config.SplunkAuthSecretName),
			},
		},
		{
			name: "Returns HEC token volume configuration when token secret is enabled",
			args: args{
				mountHost:     true,
				mountSecret:   true,
				mountHECToken: true,
				instanceName:  testutil.InstanceName,
			},
			want: []corev1.Volume{
				testutil.NewConfigMapVolume("osd-monitored-logs-local"),
				testutil.NewConfigMapVolume("osd-monitored-logs-metadata"),
				testutil.NewHostPathVolume("splunk-state", testutil.SplunkStatePath),
				testutil.NewHostPathVolume("host", testutil.HostRootPath),
				testutil.NewSecretVolume(config.SplunkHECTokenSecretName, config.SplunkHECTokenSecretName),
				{
					Name: "splunk-config",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetVolumes(tt.args.mountHost, tt.args.mountSecret, tt.args.mountHECToken, tt.args.instanceName)
			testutil.DeepEqualWithDiff(t, tt.want, got)
		})
	}
}
