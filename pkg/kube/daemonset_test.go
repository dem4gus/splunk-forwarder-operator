package kube

import (
	"testing"

	sfv1alpha1 "github.com/openshift/splunk-forwarder-operator/api/v1alpha1"
	"github.com/openshift/splunk-forwarder-operator/internal/testutil"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// daemonSetInstance produces (a pointer to) an expected DaemonSet produced by GenerateDaemonSet.
// Parameters;
// - sfInstance: SplunkForwarder instance under test.
func expectedDaemonSet(instance *sfv1alpha1.SplunkForwarder) *appsv1.DaemonSet {
	var (
		expectedRunAsUID                      int64 = 0
		expectedTerminationGracePeriodSeconds int64 = 10
		expectedPriority                      int32 = 2000001000
	)
	expectedIsPrivContainer := true
	expectedPriorityClassName := "system-node-critical"

	var sfImage string
	if instance.Spec.ImageDigest == "" {
		sfImage = testutil.Image + ":" + testutil.ImageTag
	} else {
		sfImage = testutil.Image + "@" + testutil.ImageDigest
	}

	// Expected volumes with auth secret (heavy forwarder not implemented)
	expectedVolumes := []corev1.Volume{
		testutil.NewConfigMapVolume("osd-monitored-logs-local"),
		testutil.NewConfigMapVolume("osd-monitored-logs-metadata"),
		testutil.NewHostPathVolume("splunk-state", testutil.SplunkStatePath),
		testutil.NewHostPathVolume("host", testutil.HostRootPath),
		testutil.NewSecretVolume("splunk-auth", "splunk-auth"),
	}

	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testutil.InstanceName + "-ds",
			Namespace: testutil.InstanceNamespace,
			Labels: map[string]string{
				"app": testutil.InstanceName,
			},
			Annotations: map[string]string{
				"genVersion": "10",
			},
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": "splunk-forwarder",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "splunk-forwarder",
					Namespace: testutil.InstanceNamespace,
					Labels: map[string]string{
						"name": "splunk-forwarder",
					},
				},
				Spec: corev1.PodSpec{
					PriorityClassName: expectedPriorityClassName,
					Priority:          &expectedPriority,
					NodeSelector: map[string]string{
						"kubernetes.io/os": "linux",
					},

					ServiceAccountName: "default",
					Tolerations: []corev1.Toleration{
						{
							Operator: corev1.TolerationOpExists,
						},
					},
					TerminationGracePeriodSeconds: &expectedTerminationGracePeriodSeconds,

					Containers: []corev1.Container{
						{
							Name:            "splunk-uf",
							ImagePullPolicy: corev1.PullAlways,
							Image:           sfImage,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8089,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							Resources:              corev1.ResourceRequirements{},
							TerminationMessagePath: "/dev/termination-log",
							Env: []corev1.EnvVar{
								{
									Name:  "SPLUNK_ACCEPT_LICENSE",
									Value: "yes",
								},
								{
									Name: "HOSTNAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "spec.nodeName",
										},
									},
								},
							},

							VolumeMounts: GetVolumeMounts(instance, false),

							SecurityContext: &corev1.SecurityContext{
								Privileged: &expectedIsPrivContainer,
								RunAsUser:  &expectedRunAsUID,
							},
						},
					},
					Volumes: expectedVolumes,
				},
			},
		},
	}
}

func TestGenerateDaemonSet(t *testing.T) {
	tests := []struct {
		name        string
		instance    *sfv1alpha1.SplunkForwarder
		useHECToken bool
	}{
		// TODO: The following configurations should be invalid and produce a predictable error:
		// - NewSplunkForwarderCR() with neither tag nor digest
		//   (Can't make sf pull spec when neither tag nor digest is present.)
		{
			name: "Generates DaemonSet using image digest when ImageDigest is specified",
			instance: testutil.NewSplunkForwarderCR().
				WithImageDigest(testutil.ImageDigest).
				WithGeneration(10).
				Build(),
		},
		{
			name: "Generates DaemonSet using image tag when ImageTag is specified",
			instance: testutil.NewSplunkForwarderCR().
				WithImageTag(testutil.ImageTag).
				WithGeneration(10).
				Build(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := expectedDaemonSet(tt.instance)
			actual := GenerateDaemonSet(tt.instance, tt.useHECToken)
			testutil.DeepEqualWithDiff(t, expected, actual)
		})
	}
}
