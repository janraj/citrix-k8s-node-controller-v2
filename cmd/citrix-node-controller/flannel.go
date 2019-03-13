package main

import (
	"fmt"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
        "time"
)

/*
*************************************************************************************************
*   APIName :  InitializeNode                                                                   *
*   Input   :  Nil.					             			        *
*   Output  :  Nil.				                                                *
*   Descr   :  This API initialize a node and return it.					*
*************************************************************************************************
 */
func InitializeNode(obj *ControllerInput) *v1.Node {
	klog.Info("[INFO] Initializing a Dummy Node")
	backend_data := fmt.Sprintf("{VtepMAC:%s}", obj.IngressDeviceVtepMAC)
	NewNode := &v1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: "citrixadc",
		},
		//Spec: v1.NodeSpec{
		//        PodCIDR: obj.IngressDevicePodCIDR,
		//},
	}
	NewNode.Labels = make(map[string]string)
	NewNode.Labels["com.citrix.nodetype"] = obj.DummyNodeLabel
	NewNode.Annotations = make(map[string]string)
	NewNode.Annotations["flannel.alpha.coreos.com/kube-subnet-manager"] = "true"
	NewNode.Annotations["flannel.alpha.coreos.com/backend-type"] = "vxlan"
	NewNode.Annotations["flannel.alpha.coreos.com/public-ip"] = obj.IngressDeviceVtepIP
	NewNode.Annotations["flannel.alpha.coreos.com/backend-data"] = backend_data
	return NewNode
}

/*
*************************************************************************************************
*   APIName :  CreateDummyNode                                                                  *
*   Input   :  Takes API server session called client.             			        *
*   Output  :  Nil.				                                                *
*   Descr   :  This API  Creates a Dummy Node on K8s CLuster.					*
*************************************************************************************************
 */
func (api KubernetesAPIServer) CreateDummyNode(obj *ControllerInput) *v1.Node {
	klog.Info("[INFO] Creating Citrix ADC Node")
	NsAsDummyNode := InitializeNode(obj)
	node, err := api.Client.CoreV1().Nodes().Create(NsAsDummyNode)
	if err != nil {
		klog.Error("[ERROR] Node Creation has failed", err)
		return node
	}
	klog.Info("[INFO] Created Citrix ADC Node of name=", node.GetObjectMeta().GetName())
	return node
}

/*
*************************************************************************************************
*   APIName :  GetDummyNode	                                                                *
*   Input   :  Takes API server session called client.             			        *
*   Output  :  Node Object if it present else retun Nil.				        *
*   Descr   :  This API  Get the Citrix Adc node if its present in the Cluster.			*
*************************************************************************************************
 */
func (api KubernetesAPIServer) GetDummyNode(obj *ControllerInput) *v1.Node {
	opts := metav1.GetOptions{}
	node, err := api.Client.CoreV1().Nodes().Get(obj.DummyNodeLabel, opts)
	if err != nil {
		return nil
	}
	klog.Info("[INFO] Get Node: Node Name:", node.GetObjectMeta().GetName())
	return node
}

/*
*************************************************************************************************
*   APIName :  vxlanConfig	                                                                *
*   Input   :  Takes ingress Device session.		             			        *
*   Output  :  Node Object if it present else retun Nil.				        *
*   Descr   :  This API  Get the Citrix Adc node if its present in the Cluster.			*
*************************************************************************************************
 */
func CreateVxlanConfig(ingressDevice *NitroClient, controllerInput *ControllerInput, node *Node) {

	configPack := ConfigPack{}
	vxlan := Vxlan{
		Id:   controllerInput.IngressDeviceVxlanID,
		Port: controllerInput.ClusterCNIPort,
	}
	configPack.Set("vxlan", &vxlan)
	vxlanbind := Vxlan_srcip_binding{
		Id:    controllerInput.IngressDeviceVxlanID,
		Srcip: controllerInput.IngressDeviceVtepIP,
	}
	configPack.Set("vxlan_srcip_binding", &vxlanbind)
   
	nsip := Nsip{
		Ipaddress: node.NextPodAddress,
		Netmask:   node.PodNetMask,
	}
	configPack.Set("nsip", &nsip)
	AddIngressDeviceConfig(&configPack, ingressDevice)
}

/*
*************************************************************************************************
*   APIName :  InitFlannel	                                                                *
*   Input   :  Takes Api, ingress Device session and controller input.		           	*
*   Output  :  Retun Nil.								        *
*   Descr   :  This API  Initialize flannel Config by creating Dummy Node Vxlan Config.		*
*************************************************************************************************
 */
func InitFlannel(api *KubernetesAPIServer, ingressDevice *NitroClient, controllerInput *ControllerInput) {
	klog.Info("[INFO] Initializing Flannel Config")
	dummyNode := api.GetDummyNode(controllerInput)
	ingressDevice.GetVxlanConfig(controllerInput)
	if dummyNode == nil {
		klog.Info("[INFO] Creating Citrix ADC node \n")
		api.CreateDummyNode(controllerInput)
                time.Sleep(60 * time.Second) //TODO, We have to wait till Node is available.
		dummyNode = api.GetDummyNode(controllerInput)
	}
	node := ParseNodeEvents(dummyNode, ingressDevice, controllerInput)
	node.PodNetMask = "255.255.0.0" //Automate to find next highest number
	CreateVxlanConfig(ingressDevice, controllerInput, node)
}
