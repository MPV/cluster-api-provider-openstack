/*
Copyright 2018 The Kubernetes Authors.

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

package v1alpha1

import (
	"fmt"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/json"
	"sigs.k8s.io/yaml"

	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/runtime/scheme"
)

const GroupName = "openstackproviderconfig"

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: "openstackproviderconfig.k8s.io", Version: "v1alpha1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}
)

func ClusterSpecAndStatusFromProviderSpec(cluster *clusterv1.Cluster) (*OpenstackClusterProviderSpec, *OpenstackClusterProviderStatus, error) {
	clusterProviderSpec, err := ClusterSpecFromProviderSpec(cluster.Spec.ProviderSpec)
	if err != nil {
		return nil, nil, errors.Errorf("failed to load cluster provider spec: %v", err)
	}
	clusterProviderStatus, err := ClusterStatusFromProviderStatus(cluster.Status.ProviderStatus)
	if err != nil {
		return nil, nil, errors.Errorf("failed to load cluster provider status: %v", err)
	}
	return clusterProviderSpec, clusterProviderStatus, nil
}

// ClusterConfigFromProviderSpec unmarshals a provider config into an OpenStack Cluster type
func ClusterSpecFromProviderSpec(providerSpec clusterv1.ProviderSpec) (*OpenstackClusterProviderSpec, error) {
	if providerSpec.Value == nil {
		return nil, errors.New("no such providerSpec found in manifest")
	}

	var config OpenstackClusterProviderSpec
	if err := yaml.Unmarshal(providerSpec.Value.Raw, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// ClusterStatusFromProviderStatus unmarshals a provider status into an OpenStack Cluster Status type
func ClusterStatusFromProviderStatus(extension *runtime.RawExtension) (*OpenstackClusterProviderStatus, error) {
	if extension == nil {
		return &OpenstackClusterProviderStatus{}, nil
	}

	status := new(OpenstackClusterProviderStatus)
	if err := yaml.Unmarshal(extension.Raw, status); err != nil {
		return nil, err
	}

	return status, nil
}

// This is the same as ClusterSpecFromProviderSpec but we
// expect there to be a specific Spec type for Machines soon
func MachineSpecFromProviderSpec(providerSpec clusterv1.ProviderSpec) (*OpenstackProviderSpec, error) {
	if providerSpec.Value == nil {
		return nil, errors.New("no such providerSpec found in manifest")
	}

	var config OpenstackProviderSpec
	if err := yaml.Unmarshal(providerSpec.Value.Raw, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func EncodeClusterSpecAndStatus(cluster *clusterv1.Cluster, clusterProviderSpec *OpenstackClusterProviderSpec, clusterProviderStatus *OpenstackClusterProviderStatus) (*runtime.RawExtension, *runtime.RawExtension, error) {
	rawSpec, err := EncodeClusterSpec(clusterProviderSpec)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to encod cluster spec for cluster %q in namespace %q: %v", cluster.Name, cluster.Namespace, err)
	}
	rawStatus, err := EncodeClusterStatus(clusterProviderStatus)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to update cluster status for cluster %q in namespace %q: %v", cluster.Name, cluster.Namespace, err)
	}
	return rawSpec, rawStatus, nil
}

func EncodeClusterSpec(spec *OpenstackClusterProviderSpec) (*runtime.RawExtension, error) {
	if spec == nil {
		return &runtime.RawExtension{}, nil
	}

	var rawBytes []byte
	var err error

	//  TODO: use apimachinery conversion https://godoc.org/k8s.io/apimachinery/pkg/runtime#Convert_runtime_Object_To_runtime_RawExtension
	if rawBytes, err = json.Marshal(spec); err != nil {
		return nil, err
	}

	return &runtime.RawExtension{
		Raw: rawBytes,
	}, nil
}

func EncodeClusterStatus(status *OpenstackClusterProviderStatus) (*runtime.RawExtension, error) {
	if status == nil {
		return &runtime.RawExtension{}, nil
	}

	var rawBytes []byte
	var err error

	//  TODO: use apimachinery conversion https://godoc.org/k8s.io/apimachinery/pkg/runtime#Convert_runtime_Object_To_runtime_RawExtension
	if rawBytes, err = json.Marshal(status); err != nil {
		return nil, err
	}

	return &runtime.RawExtension{
		Raw: rawBytes,
	}, nil
}

func EncodeMachineStatus(status *OpenstackMachineProviderStatus) (*runtime.RawExtension, error) {
	if status == nil {
		return &runtime.RawExtension{}, nil
	}

	var rawBytes []byte
	var err error

	// TODO: use apimachinery conversion https://godoc.org/k8s.io/apimachinery/pkg/runtime#Convert_runtime_Object_To_runtime_RawExtension
	if rawBytes, err = json.Marshal(status); err != nil {
		return nil, err
	}

	return &runtime.RawExtension{
		Raw: rawBytes,
	}, nil
}