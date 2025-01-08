/*
Copyright 2024 NTT Corporation , FUJITSU LIMITED
*/

package splitfilter_test

import (
	"context"
	"fmt"
	"path/filepath"

	v1 "github.com/compsysg/whitebox-k8s-flowctrl/api/v1"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dfStatus "github.com/compsysg/whitebox-k8s-flowctrl/lib/dataflowStatus"
	. "github.com/compsysg/whitebox-k8s-flowctrl/test/lib/testutils"
)

const basicApiVersion = "example.com/v1"

// Normal TopologyInfo not found
func Test1_Nominal(ctx context.Context, k8sClient client.Client) {

	// Get Resource
	commonDir := filepath.Join("..", "..", "resources", "splitfilter", "common")
	paths := []string{
		filepath.Join(commonDir, "functionTargets"),
		filepath.Join(commonDir, "user_requirement.yaml"),
		filepath.Join(commonDir, "strategy.yaml"),
		filepath.Join(commonDir, "dataflow.yaml"),
	}
	_, err := Deploy(ctx, k8sClient, paths...)
	Expect(err).ShouldNot(HaveOccurred())

	// Get the expected value
	testDir := filepath.Join("..", "..", "resources", "splitfilter", "1_nominal", "with_scheduler")
	expYaml := filepath.Join(testDir, "exp.yaml")
	if !Exists(expYaml) {
		CreateExpectYaml(ctx, k8sClient, expYaml, "sample-flow", "default",
			basicApiVersion, dfStatus.WB_CreationInProgress, v1.DataFlow{})
	}
	exp, err := GetResourceFromYaml[*v1.DataFlow](expYaml)
	Expect(err).ShouldNot(HaveOccurred())

	// verification
	var act *v1.DataFlow = &v1.DataFlow{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(
			ctx, client.ObjectKey{Name: exp.Name, Namespace: exp.Namespace}, act)).To(Succeed())
		g.Expect(Check(exp, act, ".status", ".spec")).Should(Succeed())
	}).Should(Succeed())
}

// Normal TopologyInfo available
func Test2_Connection(ctx context.Context, k8sClient client.Client) {

	commonDir := filepath.Join("..", "..", "resources", "splitfilter", "common")
	testDir := filepath.Join("..", "..", "resources", "splitfilter", "2_connection", "with_scheduler")

	// Deploy Resource
	paths := []string{
		filepath.Join(commonDir, "functionTargets"),
		filepath.Join("..", "..", "resources", "splitfilter", "2_connection", "common"),
		filepath.Join(commonDir, "dataflow.yaml"),
	}
	_, err := Deploy(ctx, k8sClient, paths...)
	Expect(err).ShouldNot(HaveOccurred())

	// Get the expected value
	expYaml := filepath.Join(testDir, "exp.yaml")
	if !Exists(expYaml) {
		CreateExpectYaml(ctx, k8sClient, expYaml, "sample-flow", "default",
			basicApiVersion, dfStatus.WB_CreationInProgress, v1.DataFlow{})
	}
	exp, err := GetResourceFromYaml[*v1.DataFlow](expYaml)
	Expect(err).ShouldNot(HaveOccurred())

	// verification
	var act *v1.DataFlow = &v1.DataFlow{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(
			ctx, client.ObjectKey{Name: exp.Name, Namespace: exp.Namespace}, act)).To(Succeed())
		g.Expect(Check(exp, act, ".status", ".spec")).Should(Succeed())
	}).Should(Succeed())

}

func Test3_OnePass(ctx context.Context, k8sClient client.Client) {

	testDir := filepath.Join("..", "..", "resources", "splitfilter", "3_one_pass")

	// Deploy the prerequisite resource
	// Deploy ComputeResource
	cr, err := GetResourceFromYaml[*v1.ComputeResource](filepath.Join(testDir, "computeresource_env3_sm7.yaml"))
	Expect(err).ShouldNot(HaveOccurred())
	st := cr.Status
	Expect(k8sClient.Create(ctx, cr)).Should(Succeed())
	cr.Status = st
	Eventually(func(g Gomega) {
		k8sClient.Status().Update(ctx, cr)
	}).Should(Succeed())

	// Deploy FunctionInfo
	// Deploy FunctionType
	// Deploy FunctionChain
	// Deploy TopologyInfoConfigMap
	paths := []string{
		filepath.Join(testDir, "strategy.yaml"),
		filepath.Join(testDir, "user_requirement.yaml"),
		filepath.Join(testDir, "functioninfos"),
		filepath.Join(testDir, "functionTypes"),
		filepath.Join(testDir, "functionchain.yaml"),
		filepath.Join(testDir, "topology_configmap.yaml"),
	}
	objs, err := Deploy(ctx, k8sClient, paths...)
	Expect(err).ShouldNot(HaveOccurred())

	// Check Deploy the prerequisite resource
	Eventually(func(g Gomega) {
		// Check FunctionChain
		fc := objs[filepath.Join(testDir, "functionchain.yaml")].(*v1.FunctionChain)
		Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: fc.Namespace, Name: fc.Name}, fc)).Should(Succeed())
		// Check TopologyInfo
		tp := &v1.TopologyInfo{}
		Expect(k8sClient.Get(ctx, types.NamespacedName{Namespace: "topologyinfo", Name: "topology"}, tp)).Should(Succeed())
		// Check FunctionTarget
		ftl := &v1.FunctionTargetList{}
		Expect(k8sClient.List(ctx, ftl)).Should(Succeed())
		Expect(ftl.Items).Should(HaveLen(11))
	}).Should(Succeed())

	// Deploy Resource
	_, _ = Deploy(ctx, k8sClient, filepath.Join(testDir, "df-day01.yaml")) // FIXME: need to handle 2nd return value (error)

	// Read expected value
	expYaml := filepath.Join(testDir, "exp.yaml")
	if !Exists(expYaml) {
		deployDataFlow, err := GetResourceFromYaml[*v1.DataFlow](filepath.Join(testDir, "df-day01.yaml"))
		Expect(err).ShouldNot(HaveOccurred())
		CreateExpectYaml(ctx, k8sClient, expYaml, deployDataFlow.Name, deployDataFlow.Namespace,
			basicApiVersion, dfStatus.WB_CreationInProgress, v1.DataFlow{})
	}
	exp, err := GetResourceFromYaml[*v1.DataFlow](expYaml)
	Expect(err).ShouldNot(HaveOccurred())

	// Verify expected value
	var act *v1.DataFlow = &v1.DataFlow{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(
			ctx, client.ObjectKey{Name: exp.Name, Namespace: exp.Namespace}, act)).To(Succeed())
		g.Expect(act.Status.Status).Should(Or(Equal(dfStatus.WB_Created), Equal(dfStatus.WB_CreationInProgress)))
		g.Expect(Check(exp, act, ".spec")).Should(Succeed())
	}).Should(Succeed())

}

// Normal TopologyInfo not found, various DataFlows
func Test4_Nominal_DataFlows(ctx context.Context, k8sClient client.Client, dfType string) {

	// Get Resource
	commonDir := filepath.Join("..", "..", "resources", "splitfilter", "common")
	paths := []string{
		filepath.Join(commonDir, "functionTargets"),
		filepath.Join(commonDir, "user_requirement.yaml"),
		filepath.Join(commonDir, "strategy.yaml"),
		filepath.Join(commonDir, fmt.Sprintf("dataflow_%s.yaml", dfType)),
	}
	_, err := Deploy(ctx, k8sClient, paths...)
	Expect(err).ShouldNot(HaveOccurred())

	// Get the expected value
	testDir := filepath.Join("..", "..", "resources", "splitfilter", "4_nominal_dataflows", "with_scheduler")
	expYaml := filepath.Join(testDir, fmt.Sprintf("exp_%s.yaml", dfType))
	if !Exists(expYaml) {
		CreateExpectYaml(ctx, k8sClient, expYaml, "sample-flow", "default",
			basicApiVersion, dfStatus.WB_CreationInProgress, v1.DataFlow{})
	}
	exp, err := GetResourceFromYaml[*v1.DataFlow](expYaml)
	Expect(err).ShouldNot(HaveOccurred())

	// verification
	var act *v1.DataFlow = &v1.DataFlow{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(
			ctx, client.ObjectKey{Name: exp.Name, Namespace: exp.Namespace}, act)).To(Succeed())
		g.Expect(Check(exp, act, ".status", ".spec")).Should(Succeed())
	}).Should(Succeed())
}

// Normal TopologyInfo available, various DataFlows
func Test5_Connection_DataFlows(ctx context.Context, k8sClient client.Client, dfType string) {

	commonDir := filepath.Join("..", "..", "resources", "splitfilter", "common")
	testDir := filepath.Join("..", "..", "resources", "splitfilter", "5_connection_dataflows", "with_scheduler")

	// Deploy Resource
	paths := []string{
		filepath.Join(commonDir, "functionTargets"),
		filepath.Join("..", "..", "resources", "splitfilter", "5_connection_dataflows", "common"),
		filepath.Join(commonDir, fmt.Sprintf("dataflow_%s.yaml", dfType)),
	}
	_, err := Deploy(ctx, k8sClient, paths...)
	Expect(err).ShouldNot(HaveOccurred())

	// Get the expected value
	expYaml := filepath.Join(testDir, fmt.Sprintf("exp_%s.yaml", dfType))
	if !Exists(expYaml) {
		CreateExpectYaml(ctx, k8sClient, expYaml, "sample-flow", "default",
			basicApiVersion, dfStatus.WB_CreationInProgress, v1.DataFlow{})
	}
	exp, err := GetResourceFromYaml[*v1.DataFlow](expYaml)
	Expect(err).ShouldNot(HaveOccurred())

	// verification
	var act *v1.DataFlow = &v1.DataFlow{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(
			ctx, client.ObjectKey{Name: exp.Name, Namespace: exp.Namespace}, act)).To(Succeed())
		g.Expect(Check(exp, act, ".status", ".spec")).Should(Succeed())
	}).Should(Succeed())

}

// Dynamic Reconfiguration
func Test6_Dynamic_Reconfiguration(ctx context.Context, k8sClient client.Client, testType string) {

	commonDir := filepath.Join("..", "..", "resources", "splitfilter", "6_dynamic_reconfiguration", testType, "common")
	testDir := filepath.Join("..", "..", "resources", "splitfilter", "6_dynamic_reconfiguration", testType, "with_scheduler")
	paths := []string{
		filepath.Join(commonDir, "functioninfos"),
		filepath.Join(commonDir, "functionTargets"),
		filepath.Join(commonDir, "dataflow.yaml"),
		filepath.Join(commonDir, "strategy.yaml"),
		filepath.Join(commonDir, "user_requirement.yaml"),
		filepath.Join(commonDir, "topology_customresource.yaml"),
	}

	_, err := Deploy(ctx, k8sClient, paths...)
	Expect(err).ShouldNot(HaveOccurred())

	// Get the expected value
	expYaml := filepath.Join(testDir, "exp.yaml")
	if !Exists(expYaml) {
		CreateExpectYaml(ctx, k8sClient, expYaml, "sample-flow", "default",
			basicApiVersion, dfStatus.WB_CreationInProgress, v1.DataFlow{})
	}
	exp, err := GetResourceFromYaml[*v1.DataFlow](expYaml)
	Expect(err).ShouldNot(HaveOccurred())

	// verification
	var act *v1.DataFlow = &v1.DataFlow{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(
			ctx, client.ObjectKey{Name: exp.Name, Namespace: exp.Namespace}, act)).To(Succeed())
		g.Expect(Check(exp, act, ".status", ".spec")).Should(Succeed())
	}).Should(Succeed())
}

// Dynamic Reconfiguration
func Test7_Scheduling_Retry(ctx context.Context, k8sClient client.Client, testType string) {

	commonDir := filepath.Join("..", "..", "resources", "splitfilter", "7_scheduling_retry", testType, "common")
	testDir := filepath.Join("..", "..", "resources", "splitfilter", "7_scheduling_retry", testType, "with_scheduler")
	paths := []string{
		filepath.Join(commonDir, "functionTargets"),
		filepath.Join("..", "..", "resources", "splitfilter", "7_scheduling_retry", testType, "common"),
	}
	_, err := Deploy(ctx, k8sClient, paths...)
	Expect(err).ShouldNot(HaveOccurred())

	// Get the expected value
	expYaml1 := filepath.Join(testDir, "exp1.yaml")
	if !Exists(expYaml1) {
		CreateExpectYaml(ctx, k8sClient, expYaml1, "sample-flow", "default",
			basicApiVersion, dfStatus.SchedulingInProgress, v1.DataFlow{})
	}
	exp1, err := GetResourceFromYaml[*v1.DataFlow](expYaml1)
	Expect(err).ShouldNot(HaveOccurred())

	// verification
	var act1 *v1.DataFlow = &v1.DataFlow{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(
			ctx, client.ObjectKey{Name: exp1.Name, Namespace: exp1.Namespace}, act1)).To(Succeed())
		g.Expect(Check(exp1, act1, ".status", ".spec")).Should(Succeed())
	}).Should(Succeed())

	// Update FunctionTarget
	var ft *v1.FunctionTarget = &v1.FunctionTarget{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(
			ctx, client.ObjectKey{Name: "node1.a100-0.gpu", Namespace: exp1.Namespace}, ft)).To(Succeed())
	}).Should(Succeed())

	ft.Status.CurrentCapacity = func() *int32 { i := int32(0); return &i }()
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Status().Update(ctx, ft)).To(Succeed())
	}).Should(Succeed())
	// time.Sleep(5 * time.Second)

	// Get the expected value
	expYaml2 := filepath.Join(testDir, "exp2.yaml")
	if !Exists(expYaml2) {
		CreateExpectYaml(ctx, k8sClient, expYaml2, "sample-flow", "default",
			basicApiVersion, dfStatus.WB_CreationInProgress, v1.DataFlow{})
	}
	exp, err := GetResourceFromYaml[*v1.DataFlow](expYaml2)
	Expect(err).ShouldNot(HaveOccurred())

	// verification
	var act2 *v1.DataFlow = &v1.DataFlow{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(
			ctx, client.ObjectKey{Name: exp.Name, Namespace: exp.Namespace}, act2)).To(Succeed())
		g.Expect(Check(exp, act2, ".status", ".spec")).Should(Succeed())
	}).Should(Succeed())
}

// Normal TopologyInfo not found
func UnitTest1_Nominal(ctx context.Context, k8sClient client.Client) {

	// Deploy Resource
	commonDir := filepath.Join("..", "..", "resources", "splitfilter", "common")
	paths := []string{
		filepath.Join(commonDir, "functionTargets"),
		filepath.Join(commonDir, "user_requirement.yaml"),
		filepath.Join(commonDir, "strategy.yaml"),
		filepath.Join(commonDir, "dataflow.yaml"),
		filepath.Join(commonDir, "scheduling_data.yaml"),
	}
	_, err := Deploy(ctx, k8sClient, paths...)
	Expect(err).ShouldNot(HaveOccurred())

	// Get the expected value
	testDir := filepath.Join("..", "..", "resources", "splitfilter", "1_nominal", "unit")
	expYaml := filepath.Join(testDir, "exp.yaml")
	if !Exists(expYaml) {
		CreateExpectYaml(ctx, k8sClient, expYaml, "sample-flow", "default", basicApiVersion, "Finish", v1.SchedulingData{})
	}
	exp, err := GetResourceFromYaml[*v1.SchedulingData](expYaml)
	Expect(err).ShouldNot(HaveOccurred())

	// verification
	var act *v1.SchedulingData = &v1.SchedulingData{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(
			ctx, client.ObjectKey{Name: exp.Name, Namespace: exp.Namespace}, act)).To(Succeed())
		g.Expect(Check(exp, act, ".status", ".spec")).Should(Succeed())
	}).Should(Succeed())

}

// Normal TopologyInfo available
func UnitTest2_Connection(ctx context.Context, k8sClient client.Client) {

	commonDir := filepath.Join("..", "..", "resources", "splitfilter", "common")
	testDir := filepath.Join("..", "..", "resources", "splitfilter", "2_connection", "unit")
	paths := []string{
		filepath.Join(commonDir, "functionTargets"),
		filepath.Join("..", "..", "resources", "splitfilter", "2_connection", "common"),
		filepath.Join(commonDir, "dataflow.yaml"),
		filepath.Join(testDir, "scheduling_data.yaml"),
	}
	_, err := Deploy(ctx, k8sClient, paths...)
	Expect(err).ShouldNot(HaveOccurred())

	// Get the expected value
	expYaml := filepath.Join(testDir, "exp.yaml")
	if !Exists(expYaml) {
		CreateExpectYaml(ctx, k8sClient, expYaml, "sample-flow",
			"default", basicApiVersion, "Finish", v1.SchedulingData{})
	}
	exp, err := GetResourceFromYaml[*v1.SchedulingData](expYaml)
	Expect(err).ShouldNot(HaveOccurred())

	// verification
	var act *v1.SchedulingData = &v1.SchedulingData{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(
			ctx, client.ObjectKey{Name: exp.Name, Namespace: exp.Namespace}, act)).To(Succeed())
		g.Expect(Check(exp, act, ".status", ".spec")).Should(Succeed())
	}).Should(Succeed())

}

// Normal TopologyInfo not found, various DataFlows
func UnitTest3_Nominal_DataFlows(ctx context.Context, k8sClient client.Client, dfType string) {

	// Deploy Resource
	commonDir := filepath.Join("..", "..", "resources", "splitfilter", "common")
	paths := []string{
		filepath.Join(commonDir, "functionTargets"),
		filepath.Join(commonDir, "user_requirement.yaml"),
		filepath.Join(commonDir, "strategy.yaml"),
		filepath.Join(commonDir, fmt.Sprintf("dataflow_%s.yaml", dfType)),
		filepath.Join(commonDir, "scheduling_data.yaml"),
	}
	_, err := Deploy(ctx, k8sClient, paths...)
	Expect(err).ShouldNot(HaveOccurred())

	// Get the expected value
	testDir := filepath.Join("..", "..", "resources", "splitfilter", "4_nominal_dataflows", "unit")
	expYaml := filepath.Join(testDir, fmt.Sprintf("exp_%s.yaml", dfType))
	if !Exists(expYaml) {
		CreateExpectYaml(ctx, k8sClient, expYaml, "sample-flow", "default", basicApiVersion, "Finish", v1.SchedulingData{})
	}
	exp, err := GetResourceFromYaml[*v1.SchedulingData](expYaml)
	Expect(err).ShouldNot(HaveOccurred())

	// verification
	var act *v1.SchedulingData = &v1.SchedulingData{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(
			ctx, client.ObjectKey{Name: exp.Name, Namespace: exp.Namespace}, act)).To(Succeed())
		g.Expect(Check(exp, act, ".status", ".spec")).Should(Succeed())
	}).Should(Succeed())

}

// Normal TopologyInfo available, various DataFlows
func UnitTest4_Connection_DataFlows(ctx context.Context, k8sClient client.Client, dfType string) {

	commonDir := filepath.Join("..", "..", "resources", "splitfilter", "common")
	testDir := filepath.Join("..", "..", "resources", "splitfilter", "5_connection_dataflows", "unit")
	paths := []string{
		filepath.Join(commonDir, "functionTargets"),
		filepath.Join("..", "..", "resources", "splitfilter", "5_connection_dataflows", "common"),
		filepath.Join(commonDir, fmt.Sprintf("dataflow_%s.yaml", dfType)),
		filepath.Join(testDir, "scheduling_data.yaml"),
	}
	_, err := Deploy(ctx, k8sClient, paths...)
	Expect(err).ShouldNot(HaveOccurred())

	// Get the expected value
	expYaml := filepath.Join(testDir, fmt.Sprintf("exp_%s.yaml", dfType))
	if !Exists(expYaml) {
		CreateExpectYaml(ctx, k8sClient, expYaml, "sample-flow",
			"default", basicApiVersion, "Finish", v1.SchedulingData{})
	}
	exp, err := GetResourceFromYaml[*v1.SchedulingData](expYaml)
	Expect(err).ShouldNot(HaveOccurred())

	// verification
	var act *v1.SchedulingData = &v1.SchedulingData{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(
			ctx, client.ObjectKey{Name: exp.Name, Namespace: exp.Namespace}, act)).To(Succeed())
		g.Expect(Check(exp, act, ".status", ".spec")).Should(Succeed())
	}).Should(Succeed())

}

// Dynamic_Reconfiguration
func UnitTest5_Dynamic_Reconfiguration(ctx context.Context, k8sClient client.Client, testType string) {

	commonDir := filepath.Join("..", "..", "resources", "splitfilter", "6_dynamic_reconfiguration", testType, "common")
	testDir := filepath.Join("..", "..", "resources", "splitfilter", "6_dynamic_reconfiguration", testType, "unit")
	paths := []string{
		filepath.Join(commonDir, "functioninfos"),
		filepath.Join(commonDir, "functionTargets"),
		filepath.Join(commonDir, "dataflow.yaml"),
		filepath.Join(commonDir, "strategy.yaml"),
		filepath.Join(commonDir, "user_requirement.yaml"),
		filepath.Join(commonDir, "topology_customresource.yaml"),
		filepath.Join(testDir, "scheduling_data.yaml"),
	}
	_, err := Deploy(ctx, k8sClient, paths...)
	Expect(err).ShouldNot(HaveOccurred())

	// Get the expected value
	expYaml := filepath.Join(testDir, "exp.yaml")
	if !Exists(expYaml) {
		CreateExpectYaml(ctx, k8sClient, expYaml, "sample-flow",
			"default", basicApiVersion, "Finish", v1.SchedulingData{})
	}
	exp, err := GetResourceFromYaml[*v1.SchedulingData](expYaml)
	Expect(err).ShouldNot(HaveOccurred())

	// verification
	var act *v1.SchedulingData = &v1.SchedulingData{}
	Eventually(func(g Gomega) {
		g.Expect(k8sClient.Get(
			ctx, client.ObjectKey{Name: exp.Name, Namespace: exp.Namespace}, act)).To(Succeed())
		g.Expect(Check(exp, act, ".status", ".spec")).Should(Succeed())
	}).Should(Succeed())

}
