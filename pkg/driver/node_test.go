/*
Copyright 2023 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package driver

import (
	"context"
	"fmt"
	"github.com/kubernetes-sigs/aws-fsx-openzfs-csi-driver/pkg/driver/internal"
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/mock/gomock"

	cloudMock "github.com/kubernetes-sigs/aws-fsx-openzfs-csi-driver/pkg/cloud/mocks"
	driverMocks "github.com/kubernetes-sigs/aws-fsx-openzfs-csi-driver/pkg/driver/mocks"
)

func TestNodePublishVolume(t *testing.T) {

	var (
		fsVolumeId = "fs-0a2d0632b5ff567e9"
		volumeId   = "fsvol-0efb292807cc770ff"
		dnsName    = "fs-0a2d0632b5ff567e9.fsx.us-west-2.amazonaws.com"
		volumePath = "/fsx/testVol"
		targetPath = "/target/path"
		stdVolCap  = &csi.VolumeCapability{
			AccessType: &csi.VolumeCapability_Mount{
				Mount: &csi.VolumeCapability_MountVolume{},
			},
			AccessMode: &csi.VolumeCapability_AccessMode{
				Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
			},
		}
	)

	testCases := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "success: normal filesystem nfs mount",
			testFunc: func(t *testing.T) {
				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()

				mockMetadata := cloudMock.NewMockMetadataService(mockCtl)
				mockMounter := driverMocks.NewMockMounter(mockCtl)

				driver := &nodeService{
					metadata: mockMetadata,
					mounter:  mockMounter,
					inFlight: internal.NewInFlight(),
				}
				source := fmt.Sprintf("%s:%s", dnsName, rootVolumePath)

				ctx := context.Background()
				req := &csi.NodePublishVolumeRequest{
					VolumeId: fsVolumeId,
					VolumeContext: map[string]string{
						"DNSName":      dnsName,
						"VolumePath":   rootVolumePath,
						"ResourceType": fsType,
					},
					VolumeCapability: stdVolCap,
					TargetPath:       targetPath,
				}

				mockMounter.EXPECT().MakeDir(gomock.Eq(targetPath)).Return(nil)
				mockMounter.EXPECT().IsLikelyNotMountPoint(gomock.Eq(targetPath)).Return(true, nil)
				mockMounter.EXPECT().Mount(gomock.Eq(source), gomock.Eq(targetPath), gomock.Eq("nfs"), gomock.Any()).Return(nil)
				_, err := driver.NodePublishVolume(ctx, req)
				if err != nil {
					t.Fatalf("NodePublishVolume is failed: %v", err)
				}

				mockCtl.Finish()
			},
		},
		{
			name: "success: normal volume nfs mount",
			testFunc: func(t *testing.T) {
				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()

				mockMetadata := cloudMock.NewMockMetadataService(mockCtl)
				mockMounter := driverMocks.NewMockMounter(mockCtl)

				driver := &nodeService{
					metadata: mockMetadata,
					mounter:  mockMounter,
					inFlight: internal.NewInFlight(),
				}
				source := fmt.Sprintf("%s:%s", dnsName, volumePath)

				ctx := context.Background()
				req := &csi.NodePublishVolumeRequest{
					VolumeId: volumeId,
					VolumeContext: map[string]string{
						"DNSName":      dnsName,
						"VolumePath":   volumePath,
						"ResourceType": volType,
					},
					VolumeCapability: stdVolCap,
					TargetPath:       targetPath,
				}

				mockMounter.EXPECT().MakeDir(gomock.Eq(targetPath)).Return(nil)
				mockMounter.EXPECT().IsLikelyNotMountPoint(gomock.Eq(targetPath)).Return(true, nil)
				mockMounter.EXPECT().Mount(gomock.Eq(source), gomock.Eq(targetPath), gomock.Eq("nfs"), gomock.Any()).Return(nil)
				_, err := driver.NodePublishVolume(ctx, req)
				if err != nil {
					t.Fatalf("NodePublishVolume is failed: %v", err)
				}

				mockCtl.Finish()
			},
		},
		{
			name: "success: missing volumePath for file system mount in static provisioning, default 'fsx' used",
			testFunc: func(t *testing.T) {
				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()

				mockMetadata := cloudMock.NewMockMetadataService(mockCtl)
				mockMounter := driverMocks.NewMockMounter(mockCtl)

				driver := &nodeService{
					metadata: mockMetadata,
					mounter:  mockMounter,
					inFlight: internal.NewInFlight(),
				}
				source := fmt.Sprintf("%s:%s", dnsName, rootVolumePath)

				ctx := context.Background()
				req := &csi.NodePublishVolumeRequest{
					VolumeId: fsVolumeId,
					VolumeContext: map[string]string{
						"DNSName":      dnsName,
						"ResourceType": fsType,
					},
					VolumeCapability: stdVolCap,
					TargetPath:       targetPath,
				}

				mockMounter.EXPECT().MakeDir(gomock.Eq(targetPath)).Return(nil)
				mockMounter.EXPECT().IsLikelyNotMountPoint(gomock.Eq(targetPath)).Return(true, nil)
				mockMounter.EXPECT().Mount(gomock.Eq(source), gomock.Eq(targetPath), gomock.Eq("nfs"), gomock.Any()).Return(nil)
				_, err := driver.NodePublishVolume(ctx, req)
				if err != nil {
					t.Fatalf("NodePublishVolume is failed: %v", err)
				}

				mockCtl.Finish()
			},
		},
		{
			name: "success: normal nfs mount with read only mount",
			testFunc: func(t *testing.T) {
				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()

				mockMetadata := cloudMock.NewMockMetadataService(mockCtl)
				mockMounter := driverMocks.NewMockMounter(mockCtl)

				driver := &nodeService{
					metadata: mockMetadata,
					mounter:  mockMounter,
					inFlight: internal.NewInFlight(),
				}

				source := fmt.Sprintf("%s:%s", dnsName, rootVolumePath)

				ctx := context.Background()
				req := &csi.NodePublishVolumeRequest{
					VolumeId: fsVolumeId,
					VolumeContext: map[string]string{
						"DNSName":      dnsName,
						"VolumePath":   rootVolumePath,
						"ResourceType": fsType,
					},
					VolumeCapability: stdVolCap,
					TargetPath:       targetPath,
					Readonly:         true,
				}

				mockMounter.EXPECT().MakeDir(gomock.Eq(targetPath)).Return(nil)
				mockMounter.EXPECT().IsLikelyNotMountPoint(gomock.Eq(targetPath)).Return(true, nil)
				mockMounter.EXPECT().Mount(gomock.Eq(source), gomock.Eq(targetPath), gomock.Eq("nfs"), gomock.Eq([]string{"ro"})).Return(nil)
				_, err := driver.NodePublishVolume(ctx, req)
				if err != nil {
					t.Fatalf("NodePublishVolume is failed: %v", err)
				}

				mockCtl.Finish()
			},
		},
		{
			name: "success: normal nfs mount with mount options",
			testFunc: func(t *testing.T) {
				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()

				mockMetadata := cloudMock.NewMockMetadataService(mockCtl)
				mockMounter := driverMocks.NewMockMounter(mockCtl)

				driver := &nodeService{
					metadata: mockMetadata,
					mounter:  mockMounter,
					inFlight: internal.NewInFlight(),
				}

				source := fmt.Sprintf("%s:%s", dnsName, rootVolumePath)

				ctx := context.Background()
				req := &csi.NodePublishVolumeRequest{
					VolumeId: "volumeId",
					VolumeContext: map[string]string{
						"DNSName":      dnsName,
						"VolumePath":   rootVolumePath,
						"ResourceType": fsType,
					},
					VolumeCapability: &csi.VolumeCapability{
						AccessType: &csi.VolumeCapability_Mount{
							Mount: &csi.VolumeCapability_MountVolume{
								MountFlags: []string{"nfsvers=4.1", "rsize=1048576", "wsize=1048576", "timeo=600"},
							},
						},
						AccessMode: &csi.VolumeCapability_AccessMode{
							Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
						},
					},
					TargetPath: targetPath,
				}

				mockMounter.EXPECT().MakeDir(gomock.Eq(targetPath)).Return(nil)
				mockMounter.EXPECT().IsLikelyNotMountPoint(gomock.Eq(targetPath)).Return(true, nil)
				mockMounter.EXPECT().Mount(gomock.Eq(source), gomock.Eq(targetPath), gomock.Eq("nfs"), gomock.Eq([]string{"nfsvers=4.1", "rsize=1048576", "wsize=1048576", "timeo=600"})).Return(nil)
				_, err := driver.NodePublishVolume(ctx, req)
				if err != nil {
					t.Fatalf("NodePublishVolume is failed: %v", err)
				}

				mockCtl.Finish()
			},
		},
		{
			name: "fail: missing volume id",
			testFunc: func(t *testing.T) {
				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()

				mockMetadata := cloudMock.NewMockMetadataService(mockCtl)
				mockMounter := driverMocks.NewMockMounter(mockCtl)

				driver := &nodeService{
					metadata: mockMetadata,
					mounter:  mockMounter,
					inFlight: internal.NewInFlight(),
				}

				ctx := context.Background()
				req := &csi.NodePublishVolumeRequest{
					VolumeContext: map[string]string{
						"DNSName":      dnsName,
						"VolumePath":   rootVolumePath,
						"ResourceType": fsType,
					},
					VolumeCapability: stdVolCap,
					TargetPath:       targetPath,
				}

				_, err := driver.NodePublishVolume(ctx, req)
				if err == nil {
					t.Fatalf("NodePublishVolume is not failed: %v", err)
				}

				mockCtl.Finish()
			},
		},
		{
			name: "fail: missing volume capability",
			testFunc: func(t *testing.T) {
				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()

				mockMetadata := cloudMock.NewMockMetadataService(mockCtl)
				mockMounter := driverMocks.NewMockMounter(mockCtl)

				driver := &nodeService{
					metadata: mockMetadata,
					mounter:  mockMounter,
					inFlight: internal.NewInFlight(),
				}

				ctx := context.Background()
				req := &csi.NodePublishVolumeRequest{
					VolumeId: "volumeId",
					VolumeContext: map[string]string{
						"DNSName":      dnsName,
						"VolumePath":   rootVolumePath,
						"ResourceType": fsType,
					},
					TargetPath: targetPath,
				}

				_, err := driver.NodePublishVolume(ctx, req)
				if err == nil {
					t.Fatalf("NodePublishVolume is not failed: %v", err)
				}

				mockCtl.Finish()
			},
		},
		{
			name: "fail: unsupported volume capability",
			testFunc: func(t *testing.T) {
				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()

				mockMetadata := cloudMock.NewMockMetadataService(mockCtl)
				mockMounter := driverMocks.NewMockMounter(mockCtl)

				driver := &nodeService{
					metadata: mockMetadata,
					mounter:  mockMounter,
					inFlight: internal.NewInFlight(),
				}

				ctx := context.Background()
				req := &csi.NodePublishVolumeRequest{
					VolumeId: fsVolumeId,
					VolumeContext: map[string]string{
						"DNSName":      dnsName,
						"VolumePath":   rootVolumePath,
						"ResourceType": fsType,
					},
					VolumeCapability: &csi.VolumeCapability{
						AccessType: &csi.VolumeCapability_Mount{
							Mount: &csi.VolumeCapability_MountVolume{},
						},
						AccessMode: &csi.VolumeCapability_AccessMode{
							Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_READER_ONLY,
						},
					},
					TargetPath: targetPath,
				}

				_, err := driver.NodePublishVolume(ctx, req)
				if err == nil {
					t.Fatalf("NodePublishVolume is not failed: %v", err)
				}

				mockCtl.Finish()
			},
		},
		{
			name: "fail: missing dns name",
			testFunc: func(t *testing.T) {
				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()

				mockMetadata := cloudMock.NewMockMetadataService(mockCtl)
				mockMounter := driverMocks.NewMockMounter(mockCtl)

				driver := &nodeService{
					metadata: mockMetadata,
					mounter:  mockMounter,
					inFlight: internal.NewInFlight(),
				}

				ctx := context.Background()
				req := &csi.NodePublishVolumeRequest{
					VolumeId: fsVolumeId,
					VolumeContext: map[string]string{
						"VolumePath":   rootVolumePath,
						"ResourceType": fsType,
					},
					VolumeCapability: stdVolCap,
					TargetPath:       targetPath,
				}

				_, err := driver.NodePublishVolume(ctx, req)
				if err == nil {
					t.Fatalf("NodePublishVolume is not failed: %v", err)
				}

				mockCtl.Finish()
			},
		},
		{
			name: "fail: invalid volume type",
			testFunc: func(t *testing.T) {
				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()

				mockMetadata := cloudMock.NewMockMetadataService(mockCtl)
				mockMounter := driverMocks.NewMockMounter(mockCtl)

				driver := &nodeService{
					metadata: mockMetadata,
					mounter:  mockMounter,
					inFlight: internal.NewInFlight(),
				}

				ctx := context.Background()
				req := &csi.NodePublishVolumeRequest{
					VolumeId: fsVolumeId,
					VolumeContext: map[string]string{
						"DNSName":      dnsName,
						"VolumePath":   rootVolumePath,
						"ResourceType": "voolume",
					},
					VolumeCapability: stdVolCap,
					TargetPath:       targetPath,
				}

				_, err := driver.NodePublishVolume(ctx, req)
				if err == nil {
					t.Fatalf("NodePublishVolume is not failed: %v", err)
				}

				mockCtl.Finish()
			},
		},
		{
			name: "fail: missing volumePath when mounting volume",
			testFunc: func(t *testing.T) {
				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()

				mockMetadata := cloudMock.NewMockMetadataService(mockCtl)
				mockMounter := driverMocks.NewMockMounter(mockCtl)

				driver := &nodeService{
					metadata: mockMetadata,
					mounter:  mockMounter,
					inFlight: internal.NewInFlight(),
				}

				ctx := context.Background()
				req := &csi.NodePublishVolumeRequest{
					VolumeId: volumeId,
					VolumeContext: map[string]string{
						"DNSName":      dnsName,
						"ResourceType": volType,
					},
					VolumeCapability: stdVolCap,
					TargetPath:       targetPath,
				}

				_, err := driver.NodePublishVolume(ctx, req)
				if err == nil {
					t.Fatalf("NodePublishVolume is not failed: %v", err)
				}

				mockCtl.Finish()
			},
		},
		{
			name: "fail: invalid volumePath when mounting a file system",
			testFunc: func(t *testing.T) {
				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()

				mockMetadata := cloudMock.NewMockMetadataService(mockCtl)
				mockMounter := driverMocks.NewMockMounter(mockCtl)

				driver := &nodeService{
					metadata: mockMetadata,
					mounter:  mockMounter,
					inFlight: internal.NewInFlight(),
				}

				ctx := context.Background()
				req := &csi.NodePublishVolumeRequest{
					VolumeId: fsVolumeId,
					VolumeContext: map[string]string{
						"DNSName":      dnsName,
						"VolumePath":   volumePath,
						"ResourceType": fsType,
					},
					VolumeCapability: stdVolCap,
					TargetPath:       targetPath,
				}

				_, err := driver.NodePublishVolume(ctx, req)
				if err == nil {
					t.Fatalf("NodePublishVolume is not failed: %v", err)
				}

				mockCtl.Finish()
			},
		},
		{
			name: "fail: missing target path",
			testFunc: func(t *testing.T) {
				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()

				mockMetadata := cloudMock.NewMockMetadataService(mockCtl)
				mockMounter := driverMocks.NewMockMounter(mockCtl)

				driver := &nodeService{
					metadata: mockMetadata,
					mounter:  mockMounter,
					inFlight: internal.NewInFlight(),
				}

				ctx := context.Background()
				req := &csi.NodePublishVolumeRequest{
					VolumeId: fsVolumeId,
					VolumeContext: map[string]string{
						"DNSName":      dnsName,
						"VolumePath":   rootVolumePath,
						"ResourceType": fsType,
					},
					VolumeCapability: stdVolCap,
				}

				_, err := driver.NodePublishVolume(ctx, req)
				if err == nil {
					t.Fatalf("NodePublishVolume is not failed: %v", err)
				}

				mockCtl.Finish()
			},
		},
		{
			name: "fail: mounter failed to MakeDir",
			testFunc: func(t *testing.T) {
				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()

				mockMetadata := cloudMock.NewMockMetadataService(mockCtl)
				mockMounter := driverMocks.NewMockMounter(mockCtl)

				driver := &nodeService{
					metadata: mockMetadata,
					mounter:  mockMounter,
					inFlight: internal.NewInFlight(),
				}

				ctx := context.Background()
				req := &csi.NodePublishVolumeRequest{
					VolumeId: fsVolumeId,
					VolumeContext: map[string]string{
						"DNSName":      dnsName,
						"VolumePath":   rootVolumePath,
						"ResourceType": fsType,
					},
					VolumeCapability: stdVolCap,
					TargetPath:       targetPath,
				}

				err := fmt.Errorf("failed to MakeDir")
				mockMounter.EXPECT().MakeDir(gomock.Eq(targetPath)).Return(err)

				_, err = driver.NodePublishVolume(ctx, req)
				if err == nil {
					t.Fatalf("NodePublishVolume is not failed: %v", err)
				}

				mockCtl.Finish()
			},
		},
		{
			name: "fail: mounter failed to Mount",
			testFunc: func(t *testing.T) {
				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()

				mockMetadata := cloudMock.NewMockMetadataService(mockCtl)
				mockMounter := driverMocks.NewMockMounter(mockCtl)

				driver := &nodeService{
					metadata: mockMetadata,
					mounter:  mockMounter,
					inFlight: internal.NewInFlight(),
				}

				ctx := context.Background()
				req := &csi.NodePublishVolumeRequest{
					VolumeId: fsVolumeId,
					VolumeContext: map[string]string{
						"DNSName":      dnsName,
						"VolumePath":   rootVolumePath,
						"ResourceType": fsType,
					},
					VolumeCapability: stdVolCap,
					TargetPath:       targetPath,
				}

				source := fmt.Sprintf("%s:%s", dnsName, rootVolumePath)
				err := fmt.Errorf("failed to Mount")
				mockMounter.EXPECT().MakeDir(gomock.Eq(targetPath)).Return(nil)
				mockMounter.EXPECT().IsLikelyNotMountPoint(gomock.Eq(targetPath)).Return(true, nil)
				mockMounter.EXPECT().Mount(gomock.Eq(source), gomock.Eq(targetPath), gomock.Eq("nfs"), gomock.Any()).Return(err)

				_, err = driver.NodePublishVolume(ctx, req)
				if err == nil {
					t.Fatalf("NodePublishVolume is not failed: %v", err)
				}

				mockCtl.Finish()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.testFunc)
	}
}

func TestNodeUnpublishVolume(t *testing.T) {

	var (
		targetPath = "/target/path"
		volumeId   = "fsvol-0efb292807cc770ff"
	)

	testCases := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "success: normal",
			testFunc: func(t *testing.T) {
				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()

				mockMetadata := cloudMock.NewMockMetadataService(mockCtl)
				mockMounter := driverMocks.NewMockMounter(mockCtl)

				driver := &nodeService{
					metadata: mockMetadata,
					mounter:  mockMounter,
					inFlight: internal.NewInFlight(),
				}

				ctx := context.Background()
				req := &csi.NodeUnpublishVolumeRequest{
					VolumeId:   volumeId,
					TargetPath: targetPath,
				}

				mockMounter.EXPECT().IsLikelyNotMountPoint(gomock.Eq(targetPath)).Return(false, nil)
				mockMounter.EXPECT().Unmount(gomock.Eq(targetPath)).Return(nil)

				_, err := driver.NodeUnpublishVolume(ctx, req)
				if err != nil {
					t.Fatalf("NodeUnpublishVolume is failed: %v", err)
				}
			},
		},
		{
			name: "success: target already unmounted",
			testFunc: func(t *testing.T) {
				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()

				mockMetadata := cloudMock.NewMockMetadataService(mockCtl)
				mockMounter := driverMocks.NewMockMounter(mockCtl)

				driver := &nodeService{
					metadata: mockMetadata,
					mounter:  mockMounter,
					inFlight: internal.NewInFlight(),
				}

				ctx := context.Background()
				req := &csi.NodeUnpublishVolumeRequest{
					VolumeId:   volumeId,
					TargetPath: targetPath,
				}

				mockMounter.EXPECT().IsLikelyNotMountPoint(gomock.Eq(targetPath)).Return(true, nil)

				_, err := driver.NodeUnpublishVolume(ctx, req)
				if err != nil {
					t.Fatalf("NodeUnpublishVolume is failed: %v", err)
				}
			},
		},
		{
			name: "fail: targetPath is missing",
			testFunc: func(t *testing.T) {
				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()

				mockMetadata := cloudMock.NewMockMetadataService(mockCtl)
				mockMounter := driverMocks.NewMockMounter(mockCtl)

				driver := &nodeService{
					metadata: mockMetadata,
					mounter:  mockMounter,
					inFlight: internal.NewInFlight(),
				}

				ctx := context.Background()
				req := &csi.NodeUnpublishVolumeRequest{
					VolumeId: "volumeId",
				}

				_, err := driver.NodeUnpublishVolume(ctx, req)
				if err == nil {
					t.Fatalf("NodeUnpublishVolume is not failed: %v", err)
				}
			},
		},
		{
			name: "fail: mounter failed to umount",
			testFunc: func(t *testing.T) {
				mockCtl := gomock.NewController(t)
				defer mockCtl.Finish()

				mockMetadata := cloudMock.NewMockMetadataService(mockCtl)
				mockMounter := driverMocks.NewMockMounter(mockCtl)

				driver := &nodeService{
					metadata: mockMetadata,
					mounter:  mockMounter,
					inFlight: internal.NewInFlight(),
				}

				ctx := context.Background()
				req := &csi.NodeUnpublishVolumeRequest{
					VolumeId:   volumeId,
					TargetPath: targetPath,
				}

				mockMounter.EXPECT().IsLikelyNotMountPoint(gomock.Eq(targetPath)).Return(false, nil)
				mountErr := fmt.Errorf("unmount failed")
				mockMounter.EXPECT().Unmount(gomock.Eq(targetPath)).Return(mountErr)

				_, err := driver.NodeUnpublishVolume(ctx, req)
				if err == nil {
					t.Fatalf("NodeUnpublishVolume is not failed: %v", err)
				}
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, tc.testFunc)
	}
}
