/*
Copyright 2014 The Kubernetes Authors.
Copyright 2020 Authors of Arktos - file modified.

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

package defaults

import (
	"k8s.io/klog"

	"k8s.io/apimachinery/pkg/util/sets"
	utilfeature "k8s.io/apiserver/pkg/util/feature"

	"k8s.io/kubernetes/pkg/features"
	"k8s.io/kubernetes/pkg/scheduler/algorithm/predicates"
	"k8s.io/kubernetes/pkg/scheduler/algorithm/priorities"
	"k8s.io/kubernetes/pkg/scheduler/factory"
)

const (
	// ClusterAutoscalerProvider defines the default autoscaler provider
	ClusterAutoscalerProvider = "ClusterAutoscalerProvider"
)

func init() {
	registerAlgorithmProvider(defaultPredicates(), defaultPriorities())
}

func defaultPredicates() sets.String {
	return sets.NewString(
		predicates.NoVolumeZoneConflictPred,
		predicates.MaxEBSVolumeCountPred,
		predicates.MaxGCEPDVolumeCountPred,
		predicates.MaxAzureDiskVolumeCountPred,
		predicates.MaxCSIVolumeCountPred,
		predicates.MatchInterPodAffinityPred,
		predicates.NoDiskConflictPred,
		predicates.GeneralPred,
		predicates.CheckNodeDiskPressurePred,
		predicates.CheckNodePIDPressurePred,
		predicates.CheckNodeConditionPred,
		predicates.CheckVolumeBindingPred,
	)
}

// ApplyFeatureGates applies algorithm by feature gates.
func ApplyFeatureGates() {
	if utilfeature.DefaultFeatureGate.Enabled(features.TaintNodesByCondition) {
		// Remove "CheckNodeCondition", CheckNodePIDPressure"
		// and "CheckNodeDiskPressure" predicates
		factory.RemoveFitPredicate(predicates.CheckNodeConditionPred)
		factory.RemoveFitPredicate(predicates.CheckNodeDiskPressurePred)
		factory.RemoveFitPredicate(predicates.CheckNodePIDPressurePred)
		// Remove key "CheckNodeCondition", CheckNodePIDPressure" and "CheckNodeDiskPressure"
		// from ALL algorithm provider
		// The key will be removed from all providers which in algorithmProviderMap[]
		// if you just want remove specific provider, call func RemovePredicateKeyFromAlgoProvider()
		factory.RemovePredicateKeyFromAlgorithmProviderMap(predicates.CheckNodeConditionPred)
		factory.RemovePredicateKeyFromAlgorithmProviderMap(predicates.CheckNodeDiskPressurePred)
		factory.RemovePredicateKeyFromAlgorithmProviderMap(predicates.CheckNodePIDPressurePred)

		// Fit is determined based on whether a pod can tolerate unschedulable of node
		factory.RegisterMandatoryFitPredicate(predicates.CheckNodeUnschedulablePred, predicates.CheckNodeUnschedulablePredicate)
		// Insert Key "CheckNodeUnschedulable" To All Algorithm Provider
		// The key will insert to all providers which in algorithmProviderMap[]
		// if you just want insert to specific provider, call func InsertPredicateKeyToAlgoProvider()
		factory.InsertPredicateKeyToAlgorithmProviderMap(predicates.CheckNodeUnschedulablePred)

		klog.Infof("TaintNodesByCondition is enabled, PodToleratesNodeTaints predicate is mandatory")
	}
}

func registerAlgorithmProvider(predSet, priSet sets.String) {
	// Registers algorithm providers. By default we use 'DefaultProvider', but user can specify one to be used
	// by specifying flag.
	factory.RegisterAlgorithmProvider(factory.DefaultProvider, predSet, priSet)
	// Cluster autoscaler friendly scheduling algorithm.
	factory.RegisterAlgorithmProvider(ClusterAutoscalerProvider, predSet,
		copyAndReplace(priSet, priorities.LeastRequestedPriority, priorities.LeastRequestedPriority))
}

func defaultPriorities() sets.String {
	return sets.NewString(
		priorities.LeastRequestedPriority,
		priorities.BalancedResourceAllocation,
	)
}

func copyAndReplace(set sets.String, replaceWhat, replaceWith string) sets.String {
	result := sets.NewString(set.List()...)
	if result.Has(replaceWhat) {
		result.Delete(replaceWhat)
		result.Insert(replaceWith)
	}
	return result
}
