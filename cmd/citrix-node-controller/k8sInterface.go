package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
	"os"
	"time"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config     *restclient.Config
	err        error
)

//This is interface for Kubernetes API Server
type KubernetesAPIServer struct {
	Suffix string
	Client kubernetes.Interface
}

type Controller struct {
	indexer  cache.Indexer
	queue    workqueue.RateLimitingInterface
	informer cache.Controller
}

type QueueUpdate struct {
	Key   string
	Force bool
}

/*
*************************************************************************************************
*   APIName :  ConvertPrefixLenToMask                                                           *
*   Input   :  Prefix Length. 								        *
*   Output  :  Return Net Mask in dotted decimal.	                                        *
*   Descr   :  This API takes Prefix length and generate coresponding dotted Decimal            *
*	       notation of net mask						  		*
*************************************************************************************************
 */
func ConvertPrefixLenToMask(prefixLen string) string {
	len, _ := strconv.Atoi(prefixLen)
	netmask := (uint32)(^(1<<(32-(uint32)(len)) - 1))
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, netmask)
	fmt.Println("NETMASK", bytes)
	netmaskdot := fmt.Sprintf("%d.%d.%d.%d", bytes[0], bytes[1], bytes[2], bytes[3])
	return netmaskdot
}

/*
*************************************************************************************************
*   APIName :  CreateK8sApiserverClient                                                         *
*   Input   :  Nil. 								              	*
*   Output  :  Return Kubernetes APIserver session.	                                        *
*   Descr   :  This API creates a session with kube api server which can be used for   		*
*              wathing  different events. Does not take any input as APi Func parameter.	*
*	       This API automatically get API server informations if the binary running  	*
*	       inside the cluster. If Binary is running outside cluster, cluster kube config    *
*              file must have to be in local nodes $HOME/.kube/config  location                 *
*************************************************************************************************
 */
func CreateK8sApiserverClient() (*KubernetesAPIServer, error) {
	klog.Info("[INFO] Creating API Client")
	api := &KubernetesAPIServer{}
	config, err = clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
	 	klog.Error("[WARNING] Citrix Node Controller Runs outside cluster")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
	 	        klog.Error("[ERROR] Did not find valid kube config info")
			klog.Fatal(err)
		}
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
	 	klog.Error("[ERROR] Failed to establish connection")
		klog.Fatal(err)
	}
	klog.Info("[INFO] Kubernetes Client is created")
	api.Client = client
	return api, nil
}

/*
*************************************************************************************************
*   APIName :  NodeWatcher                                                                      *
*   Input   :  Takes API server session called client.             			        *
*   Output  :  Invokes call back functions.	                                                *
*   Descr   :  This API is for watching the Nodes. API Monitors Kubernetes API server for Nodes *
*            events and store in node cache. Based on the events type, call back functions      *
*	     Will execute and perform the desired tasks.					*
*************************************************************************************************
 */
func CitrixNodeWatcher(api *KubernetesAPIServer, IngressDeviceClient *NitroClient, ControllerInputObj *ControllerInput) {

	nodeListWatcher := cache.NewListWatchFromClient(api.Client.Core().RESTClient(), "nodes", v1.NamespaceAll, fields.Everything())
	_, nodecontroller := cache.NewInformer(nodeListWatcher, &v1.Node{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			CoreHandler(api, obj, nil, "ADD", IngressDeviceClient, ControllerInputObj)
		},
		UpdateFunc: func(obj interface{}, newobj interface{}) {
			CoreHandler(api, obj, newobj, "UPDATE", IngressDeviceClient, ControllerInputObj)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("Node DELETE", obj)
			CoreHandler(api, obj, nil, "DELETE", IngressDeviceClient, ControllerInputObj)
		},
	},
	)
	stop := make(chan struct{})
	defer close(stop)
	go nodecontroller.Run(stop)

	select {}

}
/*
*************************************************************************************************
*   APIName :  Generate Next PodCIRIP                                                           *
*   Input   :  Podaddr in dotted decimal notation. 						*
*   Output  :  Return Net Mask in dotted decimal.	                                        *
*   Descr   :  This API takes Prefix length and generate coresponding dotted Decimal            *
*	       notation of net mask						  		*
*************************************************************************************************
 */
func GenerateNextPodAddr(PodAddr string) string{
	oct := strings.Split(PodAddr, ".")
	oct3, _ := strconv.Atoi(oct[3])
	if (oct3 >= 254) {
		klog.Errorf("[ERROR] Cannot increment the last octect of the IP as it is 254")
                return "Error"
        }
	oct3 = oct3 + 1
	nextaddr := fmt.Sprintf("%s.%s.%s.%d", oct[0], oct[1], oct[2], oct3)
	return nextaddr
}
/*
*************************************************************************************************
*   APIName :  GetNodeAddress                                           	                *
*   Input   :  Takes Node object.					             		*
*   Output  :  Return Internal IP, External IP and Hostname.					*
*   Descr   :  This API Gets the Address info of the Node if present 				*
*************************************************************************************************
 */
func GetNodeAddress(node v1.Node) (string, string, string){
        var InternalIP, ExternalIP, HostName string
        for _, addr := range node.Status.Addresses {
		if (addr.Type == "InternalIP"){
			InternalIP = addr.Address
        		klog.Info("[INFO] Internal IP of Node", InternalIP)
		}else if (addr.Type == "Hostname"){
			HostName = addr.Address
        		klog.Info("[INFO] Host Name of Node", HostName)
		}else if (addr.Type == "ExternalIP"){
			ExternalIP = addr.Address
        		klog.Info("[INFO] External IP  of Node", ExternalIP)
		}
	}
	return InternalIP, ExternalIP, HostName
}
/*
*************************************************************************************************
*   APIName :  ParseNodeEvents                                                                  *
*   Input   :  Takes Node object, IngressDeviceObject and InputData.             		*
*   Output  :  Return Node Object.						                *
*   Descr   :  This API  Parses the object and prepare node object. 				*
*************************************************************************************************
 */
func ParseNodeEvents(api *KubernetesAPIServer, obj interface{}, IngressDeviceClient *NitroClient, ControllerInputObj *ControllerInput) *Node {
	node := new(Node)
	originalObjJS, err := json.Marshal(obj)
	if err != nil {
		klog.Errorf("[ERROR] Failed to Marshal original object: %v", err)
	}
	var originalNode v1.Node
	if err = json.Unmarshal(originalObjJS, &originalNode); err != nil {
		klog.Errorf("[ERROR] Failed to unmarshal original object: %v", err)
	}
	if (originalNode.Labels["com.citrix.nodetype"] == "citrixadc"){ 
		node.Label = "citrixadc"
		klog.Info("[INFO] Processing Citrix Dummy Node")
	}
	PodCIDR := originalNode.Spec.PodCIDR
        InternalIP, ExternalIP, HostName := GetNodeAddress(originalNode)
	node.IPAddr = InternalIP
        node.HostName = HostName
        node.ExternalIPAddr = ExternalIP
        if (PodCIDR != ""){
		klog.Infof("[INFO] PodCIDR Information is Present: PodiCIDR=%v", PodCIDR)
		splitString := strings.Split(PodCIDR, "/")
		address, masklen := splitString[0], splitString[1]
		backendData := []byte(obj.(*v1.Node).Annotations["flannel.alpha.coreos.com/backend-data"])
		vtepMac := make(map[string]string)
		err = json.Unmarshal(backendData, &vtepMac)
		if err != nil {
			klog.Error("[ERROR] Issue with Json unmarshel", err)
		}
		if (node.HostName != ""){
			node.HostName = "Citrix"
		}
		if (node.IPAddr != ""){
			node.IPAddr = obj.(*v1.Node).Annotations["flannel.alpha.coreos.com/public-ip"]
		}
		node.PodVTEP = vtepMac["VtepMAC"]
		node.PodAddress = address
		NextPodAddress := GenerateNextPodAddr(address)
		if (NextPodAddress != "Error"){
			node.NextPodAddress = NextPodAddress
		}else{
			node.NextPodAddress = address
		}
		node.PodNetMask = ConvertPrefixLenToMask(masklen)
		node.PodMaskLen = masklen
		node.Type = obj.(*v1.Node).Annotations["flannel.alpha.coreos.com/backend-type"]
		ControllerInputObj.NodesInfo[node.IPAddr] = node
	}else{
		klog.Errorf("[WARNING] Does not have PodCIDR Information")
		klog.Info("[INFO] Generating PODCIDR and Node Information")
		if originalNode.Labels["NodeIP"] == "" {
    			originalNode.Labels = make(map[string]string)
		}
        	originalNode.Labels["NodeIP"] = node.IPAddr
        	if _, err = api.Client.CoreV1().Nodes().Update(&originalNode); err != nil {  
            		klog.Error("Failed to update label " + err.Error())
        	}
		pod := &v1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name: "citrixdummypod",
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					{
						Name:  "citrixdummypod",
						Image: "fakeimage",
					},
				},
			},
		}
		nodeSelector :=  make(map[string]string)
		nodeSelector["NodeIP"] = node.IPAddr
		pod.Spec.NodeSelector = nodeSelector
        	if _, err = api.Client.CoreV1().Pods("default").Create(pod); err != nil {  
            		klog.Error("Failed to Create a Pod " + err.Error())
        	}
                time.Sleep(10 * time.Second) //TODO, We have to wait till Node is available.
		pod, err = api.Client.CoreV1().Pods("default").Get(pod.Name, metav1.GetOptions{})
		//if err != nil {
		//	return pod, fmt.Errorf("pod Get API error: %v", err)
		//}
		klog.Info("PODS INFO", pod.Status.PodIP)
		node.PodVTEP = "00:11:11:00:01:10"
		node.PodAddress = pod.Status.PodIP
		node.PodNetMask = ConvertPrefixLenToMask("24")
	} 
	return node
}

/*
*************************************************************************************************
*   APIName :  core_add_handler                                                                 *
*   Input   :  Takes Node object, IngressDeviceObject and InputData.             		*
*   Output  :  Every node addition, it creates a Route entry in Ingress Device.	                *
*   Descr   :  This API being Invoked when an Add node event comes.				*
*	       It parses the Node event object and calls route addition for the new Node.	*
*************************************************************************************************
 */
func CoreAddHandler(api *KubernetesAPIServer, obj interface{}, IngressDeviceClient *NitroClient, ControllerInputObj *ControllerInput) {
	node := ParseNodeEvents(api, obj, IngressDeviceClient, ControllerInputObj)
	if (node.Label != "citrixadc"){
		NsInterfaceAddRoute(IngressDeviceClient, ControllerInputObj, node)
	}else {
		klog.Info("[INFO] Skipping Route addition for Dummy Node")
	}
}

/*
*************************************************************************************************
*   APIName :  CoreDeleteHandler                                                                 *
*   Input   :  Takes Node object, IngressDeviceObject and InputData.             		*
*   Output  :  Every node addition, it creates a Route entry in Ingress Device.	                *
*   Descr   :  This API is for watching the Nodes. API Monitors Kubernetes API server for Nodes *
*            events and store in node cache. Based on the events type, call back functions      *
*	     Will execute and perform the desired tasks.					*
*************************************************************************************************
 */
func CoreDeleteHandler(api *KubernetesAPIServer, obj interface{}, ingressDevice *NitroClient, controllerInput *ControllerInput) {
	node := ParseNodeEvents(api, obj, ingressDevice, controllerInput)
	NsInterfaceDeleteRoute(ingressDevice, controllerInput, node)
}

/*
*************************************************************************************************
*   APIName :  CoreUpdateHandler                                                              *
*   Input   :  Takes Node object, IngressDeviceObject and InputData.             		*
*   Output  :  Every node addition, it creates a Route entry in Ingress Device.	                *
*   Descr   :  This API being Invoked when an Add node event comes.				*
*	       It parses the Node event object and calls route addition for the new Node.	*
*************************************************************************************************
 */
func CoreUpdateHandler(api *KubernetesAPIServer, obj interface{}, IngressDeviceClient *NitroClient, ControllerInputObj *ControllerInput) {
	node := ParseNodeEvents(api, obj, IngressDeviceClient, ControllerInputObj)
	fmt.Println("UPDATE HANDLER", node)
}

/*
*************************************************************************************************
*   APIName :  CoreHandler                                                                     *
*   Input   :  Takes API server session called client.             			        *
*   Output  :  Invokes call back functions.	                                                *
*   Descr   :  This API is for watching the Nodes. API Monitors Kubernetes API server for Nodes *
*            events and store in node cache. Based on the events type, call back functions      *
*	     Will execute and perform the desired tasks.					*
*************************************************************************************************
 */
func CoreHandler(api *KubernetesAPIServer, obj interface{}, newobj interface{}, event string, IngressDeviceClient *NitroClient, ControllerInputObj *ControllerInput) {
	//create a slice of ops

	if event == "ADD" {
		CoreAddHandler(api, obj, IngressDeviceClient, ControllerInputObj)
	}
	if event == "DELETE" {
		CoreDeleteHandler(api, obj, IngressDeviceClient, ControllerInputObj)
	}
	if event == "UPDATE" {
		//	CoreUpdateHandler(obj, IngressDeviceClient, ControllerInputObj)
	}
}
func GetClusterCNI(api *KubernetesAPIServer, controllerInput *ControllerInput) {
	pods, err := api.Client.Core().Pods("kube-system").List(metav1.ListOptions{})
	if err != nil {
		klog.Error("[ERROR] Error in Pod Listing", err)
	}
	for _, pod := range pods.Items {
		if strings.Contains(pod.Name, "flannel") {
			controllerInput.ClusterCNI = "Flannel"
		} else if strings.Contains(pod.Name, "weave") {
			controllerInput.ClusterCNI = "Weave"
		} else if strings.Contains(pod.Name, "calico") {
			controllerInput.ClusterCNI = "Calico"
		} else {
			controllerInput.ClusterCNI = "Flannel"
                }
	}
}
func ConfigDecider(api *KubernetesAPIServer, ingressDevice *NitroClient, controllerInput *ControllerInput) {
	GetClusterCNI(api, controllerInput)
	if controllerInput.ClusterCNI == "Flannel" {
		InitFlannel(api, ingressDevice, controllerInput)
	} else {
		klog.Info("[INFO] Network Automation is not supported for other than Flannel")
	}
}
