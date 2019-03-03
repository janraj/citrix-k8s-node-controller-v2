package main

import (
	"k8s.io/klog"
	"os"
	"strconv"
)

type Node struct {
	HostName   string `json:"hostname,omitempty"`
	IPAddr     string `json:"address,omitempty"`
	PodCIDR    string `json:"podcidr,omitempty"`
	PodVTEP    string `json:"podvtep,omitempty"`
	PodNetMask string `json:"podcidr,omitempty"`
	PodAddress string `json:"address,omitempty"`
	PodMaskLen string `json:"prefix, omitempty"`
	Type       string `json:"podvtep,omitempty"`
	VxlanPort  string `json:"prefix, omitempty"`
	Count      int    `json:"count,omitempty"`
}

type ControllerInput struct {
	IngressDeviceIP       string
	IngressDeviceVtepMAC  string
	IngressDeviceUsername string
	IngressDevicePassword string
	IngressDeviceVtepIP  string
	IngressDeviceVxlanID  int
	ClusterCNI            string
	ClusterCNIPort        int
	DummyNodeLabel        string
	NodesInfo             map[string]*Node
}

func FetchCitrixNodeControllerInput() *ControllerInput {
	InputDataBuff := ControllerInput{}
	InputDataBuff.IngressDeviceIP = os.Getenv("NS_IP")
	configError := 0
	if len(InputDataBuff.IngressDeviceIP) == 0 {
		klog.Error("Ingress Device IP is required")
		configError = 1
	}
	InputDataBuff.IngressDeviceVtepMAC = os.Getenv("NS_VTEP_MAC")
	if len(InputDataBuff.IngressDeviceVtepMAC) == 0 {
		klog.Error("Ingress Device VtepMAC is  required")
		configError = 1
	}
	InputDataBuff.IngressDeviceUsername = os.Getenv("NS_LOGIN")
	if len(InputDataBuff.IngressDeviceUsername) == 0 {
		klog.Error("Ingress Device USER_NAME is  required")
		configError = 1
	}
	InputDataBuff.IngressDevicePassword = os.Getenv("NS_PASSWORD")
	if len(InputDataBuff.IngressDevicePassword) == 0 {
		klog.Error("Ingress Device PASSWORD is  required")
		configError = 1
	}
	InputDataBuff.IngressDeviceVtepIP = os.Getenv("NS_SNIP")
	if len(InputDataBuff.IngressDeviceVtepIP) == 0 {
		klog.Info("Ingress Device VTEP IP is empty")
		configError = 1
	}
	if configError == 1 {
		klog.Error("Unable to get the above mentioned input from YAML")
		panic("Killing Container.........Please restart Citrix Node Controller with Valid Inputs")
	}
	InputDataBuff.ClusterCNI = os.Getenv("K8S_CNI")
	if len(InputDataBuff.ClusterCNI) == 0 {
		klog.Info("Cluster CNI information is Empty")
	}
	InputDataBuff.DummyNodeLabel = "citrixadc"
	InputDataBuff.IngressDeviceVxlanID, _ = strconv.Atoi(os.Getenv("NS_VXLAN_ID"))
	if InputDataBuff.IngressDeviceVxlanID == 0 {
		klog.Info("VXLAN ID has Not Given, taking 5000 as default VXLAN_ID")
		InputDataBuff.IngressDeviceVxlanID = 5000
	}
	InputDataBuff.ClusterCNIPort, _ = strconv.Atoi(os.Getenv("K8S_VXLAN_PORT"))
	if InputDataBuff.ClusterCNIPort == 0 {
		klog.Info("K8S_VXLAN_PORT has Not Given, taking default 8472 as Vxlan Port")
		InputDataBuff.ClusterCNIPort = 8472
	}
	return &InputDataBuff
}
