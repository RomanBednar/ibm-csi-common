/**
 * Copyright 2021 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package e2e

import (
	"os"

	"github.com/IBM/ibm-csi-common/tests/e2e/testsuites"
	. "github.com/onsi/ginkgo"
	v1 "k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/kubernetes/test/e2e/framework"
)

const defaultSecret = ""

var testResultFile = os.Getenv("E2E_TEST_RESULT")
var err error
var fpointer *os.File

var _ = Describe("[ics-e2e] [sc] [with-deploy] Dynamic Provisioning for all SC with Deployment", func() {
	f := framework.NewDefaultFramework("ics-e2e-deploy")

	var (
		cs        clientset.Interface
		ns        *v1.Namespace
		secretKey string
	)

	secretKey = os.Getenv("E2E_SECRET_ENCRYPTION_KEY")
	if secretKey == "" {
		secretKey = defaultSecret
	}

	BeforeEach(func() {
		cs = f.ClientSet
		ns = f.Namespace
	})

	It("with 5iops sc: should create a pvc &pv, deployment resources, write and read to volume, delete the pod, write and read to volume again", func() {
		reclaimPolicy := v1.PersistentVolumeReclaimDelete
		fpointer, err = os.OpenFile(testResultFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer fpointer.Close()

		pod := testsuites.PodDetails{
			Cmd:      "echo 'hello world' >> /mnt/test-1/data && while true; do sleep 2; done",
			CmdExits: false,
			Volumes: []testsuites.VolumeDetails{
				{
					PVCName:       "ics-vol-5iops-",
					VolumeType:    "ibmc-vpc-block-5iops-tier",
					FSType:        "ext4",
					ClaimSize:     "15Gi",
					ReclaimPolicy: &reclaimPolicy,
					MountOptions:  []string{"rw"},
					VolumeMount: testsuites.VolumeMountDetails{
						NameGenerate:      "test-volume-",
						MountPathGenerate: "/mnt/test-",
					},
				},
			},
		}
		test := testsuites.DynamicallyProvisioneDeployWithVolWRTest{
			Pod: pod,
			PodCheck: &testsuites.PodExecCheck{
				Cmd:              []string{"cat", "/mnt/test-1/data"},
				ExpectedString01: "hello world\n",
				ExpectedString02: "hello world\nhello world\n", // pod will be restarted so expect to see 2 instances of string
			},
		}
		test.Run(cs, ns)
		if _, err = fpointer.WriteString("VPC-BLK-CSI-TEST: 5iops SC DEPLOYMENT Test: PASS\n"); err != nil {
			panic(err)
		}
	})

	It("with generalpurpose sc: should create a pvc & pv, deployment resources, write and read to volume, delete the pod, write and read to volume again", func() {
		reclaimPolicy := v1.PersistentVolumeReclaimDelete
		fpointer, err = os.OpenFile(testResultFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer fpointer.Close()
		pod := testsuites.PodDetails{
			Cmd:      "echo 'hello world' >> /mnt/test-1/data && while true; do sleep 2; done",
			CmdExits: false,
			Volumes: []testsuites.VolumeDetails{
				{
					PVCName:       "ics-vol-gp-",
					VolumeType:    "ibmc-vpc-block-general-purpose",
					FSType:        "ext4",
					ClaimSize:     "35Gi",
					ReclaimPolicy: &reclaimPolicy,
					MountOptions:  []string{"rw"},
					VolumeMount: testsuites.VolumeMountDetails{
						NameGenerate:      "test-volume-",
						MountPathGenerate: "/mnt/test-",
					},
				},
			},
		}
		test := testsuites.DynamicallyProvisioneDeployWithVolWRTest{
			Pod: pod,
			PodCheck: &testsuites.PodExecCheck{
				Cmd:              []string{"cat", "/mnt/test-1/data"},
				ExpectedString01: "hello world\n",
				ExpectedString02: "hello world\nhello world\n", // pod will be restarted so expect to see 2 instances of string
			},
		}
		test.Run(cs, ns)
		if _, err = fpointer.WriteString("VPC-BLK-CSI-TEST: GP SC DEPLOYMENT Test: PASS\n"); err != nil {
			panic(err)
		}
	})
})

var _ = Describe("[ics-e2e] [sc] [with-pods] Dynamic Provisioning for all SC with PODs", func() {
	f := framework.NewDefaultFramework("ics-e2e-pods")

	var (
		cs        clientset.Interface
		ns        *v1.Namespace
		secretKey string
	)

	secretKey = os.Getenv("E2E_SECRET_ENCRYPTION_KEY")
	if secretKey == "" {
		secretKey = defaultSecret
	}

	BeforeEach(func() {
		cs = f.ClientSet
		ns = f.Namespace
	})

	It("with 5iops sc: should create a pvc & pv, pod resources, write and read to volume", func() {
		reclaimPolicy := v1.PersistentVolumeReclaimDelete
		fpointer, err = os.OpenFile(testResultFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer fpointer.Close()
		pods := []testsuites.PodDetails{
			{
				// Cmd: "echo 'hello world' > /mnt/test-1/data && grep 'hello world' /mnt/test-1/data",
				Cmd:      "echo 'hello world' > /mnt/test-1/data && while true; do sleep 2; done",
				CmdExits: false,
				Volumes: []testsuites.VolumeDetails{
					{
						PVCName:       "ics-vol-5iops-",
						VolumeType:    "ibmc-vpc-block-5iops-tier",
						FSType:        "ext4",
						ClaimSize:     "15Gi",
						ReclaimPolicy: &reclaimPolicy,
						MountOptions:  []string{"rw"},
						VolumeMount: testsuites.VolumeMountDetails{
							NameGenerate:      "test-volume-",
							MountPathGenerate: "/mnt/test-",
						},
					},
				},
			},
		}
		test := testsuites.DynamicallyProvisionePodWithVolTest{
			Pods: pods,
			PodCheck: &testsuites.PodExecCheck{
				Cmd:              []string{"cat", "/mnt/test-1/data"},
				ExpectedString01: "hello world\n",
				ExpectedString02: "hello world\nhello world\n", // pod will be restarted so expect to see 2 instances of string
			},
		}
		test.Run(cs, ns)
		if _, err = fpointer.WriteString("VPC-BLK-CSI-TEST: 5iops SC POD Test: PASS\n"); err != nil {
			panic(err)
		}
	})

	It("with generalpurpose sc: should create a pvc & pv, pod resources, write and read to volume", func() {
		reclaimPolicy := v1.PersistentVolumeReclaimDelete
		fpointer, err = os.OpenFile(testResultFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer fpointer.Close()
		pods := []testsuites.PodDetails{
			{
				Cmd:      "echo 'hello world' > /mnt/test-1/data && while true; do sleep 2; done",
				CmdExits: false,
				Volumes: []testsuites.VolumeDetails{
					{
						PVCName:       "ics-vol-gp-",
						VolumeType:    "ibmc-vpc-block-general-purpose",
						FSType:        "ext4",
						ClaimSize:     "35Gi",
						ReclaimPolicy: &reclaimPolicy,
						MountOptions:  []string{"rw"},
						VolumeMount: testsuites.VolumeMountDetails{
							NameGenerate:      "test-volume-",
							MountPathGenerate: "/mnt/test-",
						},
					},
				},
			},
		}
		test := testsuites.DynamicallyProvisionePodWithVolTest{
			Pods: pods,
			PodCheck: &testsuites.PodExecCheck{
				Cmd:              []string{"cat", "/mnt/test-1/data"},
				ExpectedString01: "hello world\n",
				ExpectedString02: "hello world\nhello world\n", // pod will be restarted so expect to see 2 instances of string
			},
		}
		test.Run(cs, ns)
		if _, err = fpointer.WriteString("VPC-BLK-CSI-TEST: GP SC POD Test: PASS\n"); err != nil {
			panic(err)
		}
	})

	It("with custom sc: should create a pvc & pv, pod resources, write and read to volume", func() {

		// For Custom SC PVC name and Secret name should be same and in same NS
		secret := testsuites.NewSecret(cs, "ics-vol-block-custom", ns.Name, "800", "e2e test",
			"false", secretKey, "vpc.block.csi.ibm.io")
		secret.Create()
		defer secret.Cleanup()
		reclaimPolicy := v1.PersistentVolumeReclaimDelete
		fpointer, err = os.OpenFile(testResultFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer fpointer.Close()
		pods := []testsuites.PodDetails{
			{
				Cmd:      "echo 'hello world' > /mnt/test-1/data && while true; do sleep 2; done",
				CmdExits: false,
				Volumes: []testsuites.VolumeDetails{
					{
						PVCName:       "ics-vol-block-custom",
						VolumeType:    "ibmc-vpc-block-custom",
						FSType:        "ext4",
						ClaimSize:     "45Gi",
						ReclaimPolicy: &reclaimPolicy,
						MountOptions:  []string{"rw"},
						VolumeMount: testsuites.VolumeMountDetails{
							NameGenerate:      "test-volume-",
							MountPathGenerate: "/mnt/test-",
						},
					},
				},
			},
		}
		test := testsuites.DynamicallyProvisionePodWithVolTest{
			Pods: pods,
			PodCheck: &testsuites.PodExecCheck{
				Cmd:              []string{"cat", "/mnt/test-1/data"},
				ExpectedString01: "hello world\n",
				ExpectedString02: "hello world\nhello world\n", // pod will be restarted so expect to see 2 instances of string
			},
		}
		test.Run(cs, ns)
		if _, err = fpointer.WriteString("VPC-BLK-CSI-TEST: Custom SC Pod Test: PASS\n"); err != nil {
			panic(err)
		}
	})
})

var _ = Describe("[ics-e2e] [sc] [with-statefulset] Dynamic Provisioning using statefulsets", func() {
	f := framework.NewDefaultFramework("ics-e2e-statefulsets")

	var (
		cs clientset.Interface
		ns *v1.Namespace
	)
	BeforeEach(func() {
		cs = f.ClientSet
		ns = f.Namespace
	})
	It("With 5iops sc: should creat statefuleset resources, write and read to volume", func() {
		fpointer, err = os.OpenFile(testResultFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer fpointer.Close()

		var replicaCount = int32(2)
		headlessService := testsuites.NewHeadlessService(cs, "ics-e2e-service-", ns.Name, "test")
		service := headlessService.Create()
		defer headlessService.Cleanup()

		reclaimPolicy := v1.PersistentVolumeReclaimDelete
		pod := testsuites.PodDetails{
			Cmd: "echo 'hello world' > /mnt/test-1/data && while true; do sleep 2; done",
			Volumes: []testsuites.VolumeDetails{
				{
					PVCName:       "ics-vol-5iops",
					VolumeType:    "ibmc-vpc-block-5iops-tier",
					FSType:        "ext4",
					ClaimSize:     "20Gi",
					ReclaimPolicy: &reclaimPolicy,
					MountOptions:  []string{"rw"},
					VolumeMount: testsuites.VolumeMountDetails{
						NameGenerate:      "test-volume-",
						MountPathGenerate: "/mnt/test-",
					},
				},
			},
		}

		test := testsuites.StatefulsetWithVolWRTest{
			Pod: pod,
			PodCheck: &testsuites.PodExecCheck{
				Cmd:              []string{"cat", "/mnt/test-1/data"},
				ExpectedString01: "hello world\n",
			},
			Labels:       service.Labels,
			ReplicaCount: replicaCount,
			ServiceName:  service.Name,
		}
		test.Run(cs, ns, false)
		if _, err = fpointer.WriteString("VPC-BLK-CSI-TEST: With Statefulsets: PASS\n"); err != nil {
			panic(err)
		}
	})
})

var _ = Describe("[ics-e2e] [node-drain] [with-pods] Dynamic Provisioning using statefulsets", func() {
	f := framework.NewDefaultFramework("ics-e2e-statefulsets")

	var (
		cs clientset.Interface
		ns *v1.Namespace
	)
	BeforeEach(func() {
		cs = f.ClientSet
		ns = f.Namespace
	})
	It("With statefulset: should create one pod and attach volume dynamically. Write and read to volume. Next drain the node where pod is attached and wait for the pod to come up on second node. Read the volume again.", func() {

		fpointer, err = os.OpenFile(testResultFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer fpointer.Close()
		var replicaCount = int32(1)
		headlessService := testsuites.NewHeadlessService(cs, "ics-e2e-service-", ns.Name, "node-drain")
		service := headlessService.Create()
		defer headlessService.Cleanup()

		reclaimPolicy := v1.PersistentVolumeReclaimDelete
		pod := testsuites.PodDetails{
			Cmd: "echo 'hello world' > /mnt/test-1/data && while true; do sleep 2; done",
			Volumes: []testsuites.VolumeDetails{
				{
					PVCName:       "ics-vol-5iops",
					VolumeType:    "ibmc-vpc-block-5iops-tier",
					FSType:        "ext4",
					ClaimSize:     "20Gi",
					ReclaimPolicy: &reclaimPolicy,
					MountOptions:  []string{"rw"},
					VolumeMount: testsuites.VolumeMountDetails{
						NameGenerate:      "test-volume-",
						MountPathGenerate: "/mnt/test-",
					},
				},
			},
		}

		test := testsuites.StatefulsetWithVolWRTest{
			Pod: pod,
			PodCheck: &testsuites.PodExecCheck{
				Cmd:              []string{"cat", "/mnt/test-1/data"},
				ExpectedString01: "hello world\n",
			},
			Labels:       service.Labels,
			ReplicaCount: replicaCount,
			ServiceName:  service.Name,
		}
		test.Run(cs, ns, true)
		if _, err = fpointer.WriteString("VPC-BLK-CSI-TEST: Cordon and Drain Node: PASS\n"); err != nil {
			panic(err)
		}
	})
})

var _ = Describe("[ics-e2e] [resize] [pv] Dynamic Provisioning and resize pv", func() {
	f := framework.NewDefaultFramework("ics-e2e-pods")

	var (
		cs        clientset.Interface
		ns        *v1.Namespace
		secretKey string
	)

	secretKey = os.Getenv("E2E_SECRET_ENCRYPTION_KEY")
	if secretKey == "" {
		secretKey = defaultSecret
	}

	BeforeEach(func() {
		cs = f.ClientSet
		ns = f.Namespace
	})

	It("with 5iops sc: should create a pvc & pv, pod resources, and resize the volume", func() {
		reclaimPolicy := v1.PersistentVolumeReclaimDelete
		fpointer, err = os.OpenFile(testResultFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer fpointer.Close()
		pods := []testsuites.PodDetails{
			{
				Cmd:      "echo 'hello world' > /mnt/test-1/data && while true; do sleep 2; done",
				CmdExits: false,
				Volumes: []testsuites.VolumeDetails{
					{
						PVCName:       "ics-vol-5iops-",
						VolumeType:    "ibmc-vpc-block-5iops-tier",
						FSType:        "ext4",
						ClaimSize:     "260Gi",
						ReclaimPolicy: &reclaimPolicy,
						MountOptions:  []string{"rw"},
						VolumeMount: testsuites.VolumeMountDetails{
							NameGenerate:      "test-volume-",
							MountPathGenerate: "/mnt/test-",
						},
					},
				},
			},
		}
		test := testsuites.DynamicallyProvisionedResizeVolumeTest{
			Pods: pods,
			PodCheck: &testsuites.PodExecCheck{
				Cmd:              []string{"cat", "/mnt/test-1/data"},
				ExpectedString01: "hello world\n",
				ExpectedString02: "hello world\nhello world\n", // pod will be restarted so expect to see 2 instances of string
			},
			// ExpandVolSize is in Gi i.e, 2000Gi
			ExpandVolSizeG: 2000,
			ExpandedSize:   1900,
		}
		test.Run(cs, ns)
		if _, err = fpointer.WriteString("VPC-BLK-CSI-TEST: 5iops SC POD Test and resize volume: PASS\n"); err != nil {
			panic(err)
		}
	})
})
