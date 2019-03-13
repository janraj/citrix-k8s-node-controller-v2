package main

import (
	"k8s.io/klog"
	"os"
	"strconv"
)

type Node struct {
	HostName   string `json:"hostname,omitempty"`
	IPAddr     string `json:"address,omitempty"`
	ExternalIPAddr     string `json:"address,omitempty"`
	PodCIDR    string `json:"podcidr,omitempty"`
	PodVTEP    string `json:"podvtep,omitempty"`
	PodNetMask string `json:"podcidr,omitempty"`
	PodAddress string `json:"podaddress,omitempty"`
	NextPodAddress string `json:"nextpodaddress,omitempty"`
	PodMaskLen string `json:"prefix, omitempty"`
	Type       string `json:"podvtep,omitempty"`
	VxlanPort  string `json:"prefix, omitempty"`
	Count      int    `json:"count,omitempty"`
	Label     string `json:"label,omitempty"`
}

type ControllerInput struct {
	IngressDeviceIP       string
	IngressDeviceVtepMAC  string
	IngressDeviceUsername string
	IngressDevicePassword string
	IngressDeviceVtepIP  string
	IngressDevicePodCIDR  string
	IngressDeviceVxlanID  int
	IngressDeviceVxlanIDs  string
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
		klog.Error("[ERROR] Ingress Device IP (NS_IP) is required")
		configError = 1
	}
	InputDataBuff.IngressDeviceVtepMAC = os.Getenv("NS_VTEP_MAC")
	if len(InputDataBuff.IngressDeviceVtepMAC) == 0 {
		klog.Error("[ERROR] Ingress Device VtepMAC (NS_VTEP_MAC) is  required")
		configError = 1
	}
	InputDataBuff.IngressDeviceUsername = os.Getenv("NS_LOGIN")
	if len(InputDataBuff.IngressDeviceUsername) == 0 {
		klog.Error("[ERROR] Ingress Device USER_NAME (NS_LOGIN) is  required")
		configError = 1
	}
	InputDataBuff.IngressDevicePassword = os.Getenv("NS_PASSWORD")
	if len(InputDataBuff.IngressDevicePassword) == 0 {
		klog.Error("[ERROR] Ingress Device PASSWORD (NS_PASSWORD) is  required")
		configError = 1
	}
	InputDataBuff.IngressDeviceVtepIP = os.Getenv("NS_SNIP")
	if len(InputDataBuff.IngressDeviceVtepIP) == 0 {
		klog.Info("[ERROR] Ingress Device VTEP IP (NS_SNIP)  is empty")
		configError = 1
	}
	if configError == 1 {
		klog.Error("Unable to get the above mentioned input from YAML")
		panic("[ERROR] Killing Container.........Please restart Citrix Node Controller with Valid Inputs")
	}
	InputDataBuff.ClusterCNI = os.Getenv("K8S_CNI")
	if len(InputDataBuff.ClusterCNI) == 0 {
		klog.Infof("[INFO] Cluster CNI information is Empty")
	}
	InputDataBuff.IngressDevicePodCIDR = os.Getenv("NS_POD_CIDR")
	if len(InputDataBuff.IngressDevicePodCIDR) == 0 {
		klog.Infof("[INFO] IngressDevicePodCIDR is Empty")
	}
	InputDataBuff.DummyNodeLabel = "citrixadc"
        InputDataBuff.IngressDeviceVxlanIDs = os.Getenv("NS_VXLAN_ID")
	InputDataBuff.IngressDeviceVxlanID, _ = strconv.Atoi(InputDataBuff.IngressDeviceVxlanIDs)
	if InputDataBuff.IngressDeviceVxlanID == 0 {
		klog.Info("[INFO] VXLAN ID has Not Given, taking 5000 as default VXLAN_ID")
		InputDataBuff.IngressDeviceVxlanID = 5000
		InputDataBuff.IngressDeviceVxlanIDs = "5000"
	}
	InputDataBuff.ClusterCNIPort, _ = strconv.Atoi(os.Getenv("K8S_VXLAN_PORT"))
	if InputDataBuff.ClusterCNIPort == 0 {
		klog.Info("[INFO] K8S_VXLAN_PORT has Not Given, taking default 8472 as Vxlan Port")
		InputDataBuff.ClusterCNIPort = 8472
	}
	InputDataBuff.NodesInfo = make(map[string]*Node)
	return &InputDataBuff
}
