package main
import (
        "fmt"
	"strconv"
	"os"
        "k8s.io/client-go/kubernetes"
        "k8s.io/client-go/tools/clientcmd"
        "k8s.io/klog"
        "path/filepath"
        restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
        "k8s.io/client-go/util/workqueue"
	uruntime "k8s.io/apimachinery/pkg/util/runtime"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
	"encoding/json"
	"strings"
	"encoding/binary"
)

var (
      kubeconfig = filepath.Join(sudo, os.Getenv("HOME"), ".kube", "config",)
      config *restclient.Config = nil
      err error = nil
)
type KubernetesAPIServer struct {
	Suffix string
	Client kubernetes.Interface
}


type Controller struct {
	indexer        cache.Indexer
	queue          workqueue.RateLimitingInterface
	informer       cache.Controller
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
func ConvertPrefixLenToMask (prefixLen string)(string){
	len, _ := strconv.Atoi(prefixLen)
	netmask := (uint32)(^((1 << (32 - (uint32)(len)) - 1)))
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
func CreateK8sApiserverClient()(*KubernetesAPIServer, error){
    api := &KubernetesAPIServer{}
    config, err = clientcmd.BuildConfigFromFlags("", "")
    if err != nil {
        config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
        if err != nil {
            klog.Fatal(err)
        }
    }

    client, err := kubernetes.NewForConfig(config)
    if err != nil {
        klog.Fatal(err)
    }
    klog.Info("Kubernetes Client is created", client)
    api.Client=client
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
		AddFunc: func(obj interface{}){
			CoreHandler(obj, nil,  "ADD", IngressDeviceClient, ControllerInputObj)
		},
		UpdateFunc: func(obj interface{}, newobj interface{}){
			CoreHandler(obj, newobj, "UPDATE", IngressDeviceClient, ControllerInputObj)
		},
		DeleteFunc: func(obj interface{}){
			fmt.Println("Node DELETE", obj)
			CoreHandler(obj, nil, "DELETE", IngressDeviceClient, ControllerInputObj)
		},
	    },
        )
	stop := make(chan struct{})
	defer close(stop)
	go nodecontroller.Run(stop)
	
	select {}
	nodequeue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	fmt.Println("Citrix Node Watcher")
	indexer, informer := cache.NewIndexerInformer(nodeListWatcher, &v1.Node{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				fmt.Println("Add ", obj)
				nodequeue.Add(QueueUpdate{key, false})
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				nodequeue.Add(QueueUpdate{key, false})
			}
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				nodequeue.Add(QueueUpdate{key, false})
			}
		},
	}, cache.Indexers{})

	controller := NewController(nodequeue, indexer, informer)

	fmt.Println("Starting Citrix Node controller")
	stop = make(chan struct{})
	defer close(stop)
	go controller.Run(5, stop)

	//select {}

	/*
        _, nodecontroller := cache.NewInformer(nodeListWatcher, &v1.Node{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: CoreHandler(obj interface{}, k8sclient, ingress_device_client, input_data), 
	    },
        )
	stop := make(chan struct{})
	defer close(stop)
	go nodecontroller.Run(stop)
	
	select {}
	*/
}
func NewController(queue workqueue.RateLimitingInterface, indexer cache.Indexer, informer cache.Controller) *Controller {

	return &Controller{
		informer:       informer,
		indexer:        indexer,
		queue:          queue,
	}
}
func (c *Controller) Run(threadiness int, stopCh chan struct{}) {
        klog.Info("JANRAJ RUN FUNC: Starting Node controller")
	defer uruntime.HandleCrash()

	// Let the workers stop when we are done
	defer c.queue.ShutDown()
	//g

	go c.informer.Run(stopCh)


	// Start a number of worker threads to read from the queue.
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
	klog.Info("Stopping Node controller")
	//glog.Info("Stopping Node controller")
}
func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}
func (c *Controller) processNextItem() bool {
	// Wait until there is a new item in the working queue
	upd, quit := c.queue.Get()
	klog.Info("JANRAJ Data is ", upd)
	if quit {
		return false
	}
	// Tell the queue that we are done with processing this key. This unblocks the key for other workers
	// This allows safe parallel processing because two nodes with the same key are never processed in
	// parallel.
	defer c.queue.Done(upd)

	// Invoke the method containing the business logic
	//err := c.syncToCalico(upd.(QueueUpdate))

	// Handle the error if something went wrong during the execution of the business logic
	//c.handleErr(err, upd)
	return true
}
/*
*************************************************************************************************
*   APIName :  ParseNodeEvents                                                                *
*   Input   :  Takes Node object, IngressDeviceObject and InputData.             		*		
*   Output  :  Every node addition, it creates a Route entry in Ingress Device.	                *
*   Descr   :  This API is for watching the Nodes. API Monitors Kubernetes API server for Nodes *
*            events and store in node cache. Based on the events type, call back functions      *
*	     Will execute and perform the desired tasks.					*
*************************************************************************************************
*/
//TODO: Make it independant of CNI
func ParseNodeEvents(obj interface{}, IngressDeviceClient *NitroClient, ControllerInputObj *ControllerInput) (*Node){
	originalObjJS, err := json.Marshal(obj)
        var originalNode v1.Node
        if err = json.Unmarshal(originalObjJS, &originalNode); err != nil {
                        klog.Errorf("Failed to unmarshal original object: %v", err)
        }
        PodCIDR := originalNode.Spec.PodCIDR
	split_string := strings.Split(PodCIDR, "/")
    	address, masklen := split_string[0], split_string[1]
	backend_data := []byte(obj.(*v1.Node).Annotations["flannel.alpha.coreos.com/backend-data"])
	vtep_mac := make(map[string]string)
	err = json.Unmarshal(backend_data, &vtep_mac)
	if err != nil {
		fmt.Println("Error")
	}
	node := new(Node)
	node.HostName = "janraj"
	node.IPAddr = obj.(*v1.Node).Annotations["flannel.alpha.coreos.com/public-ip"]
	node.PodVTEP = vtep_mac["VtepMAC"]
	node.PodAddress = address
	node.PodNetMask =  ConvertPrefixLenToMask(masklen) 
	node.PodMaskLen = masklen
	node.Type = obj.(*v1.Node).Annotations["flannel.alpha.coreos.com/backend-type"]
	ControllerInputObj.NodesInfo=make(map[string]*Node)
	ControllerInputObj.NodesInfo[node.IPAddr]=node
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
func CoreAddHandler(obj interface{}, IngressDeviceClient *NitroClient, ControllerInputObj *ControllerInput){
	node := ParseNodeEvents(obj, IngressDeviceClient , ControllerInputObj)
	NsInterfaceAddRoute(IngressDeviceClient, ControllerInputObj, node)        
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
func CoreDeleteHandler(obj interface{}, ingressDevice *NitroClient, controllerInput *ControllerInput){
	node := ParseNodeEvents(obj, ingressDevice, controllerInput)
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
func CoreUpdateHandler(obj interface{}, IngressDeviceClient *NitroClient, ControllerInputObj *ControllerInput){
	node := ParseNodeEvents(obj, IngressDeviceClient , ControllerInputObj)
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
func CoreHandler(obj interface{}, newobj interface{}, event string, IngressDeviceClient *NitroClient, ControllerInputObj *ControllerInput){
	if event == "ADD"{
		CoreAddHandler(obj, IngressDeviceClient, ControllerInputObj)
	}
	if event == "DELETE"{
		CoreDeleteHandler(obj, IngressDeviceClient, ControllerInputObj)
	}
	if event == "UPDATE"{
	//	CoreUpdateHandler(obj, IngressDeviceClient, ControllerInputObj)
	}
}
func GetClusterCNI(api *KubernetesAPIServer, controllerInput *ControllerInput) {
	pods, err:=  api.Client.Core().Pods("kube-system").List(metav1.ListOptions{})
	if err != nil {
		klog.Error("Error in Pod Listing", err)
	}
    	for _,pod:= range pods.Items{
		if strings.Contains(pod.Name, "flannel") {
			controllerInput.ClusterCNI = "Flannel"
		}else if strings.Contains(pod.Name, "weave") {
			controllerInput.ClusterCNI = "Weave"
		}else if strings.Contains(pod.Name, "calico") {
			controllerInput.ClusterCNI = "Calico"
		}
    	}
}
func ConfigDecider(api *KubernetesAPIServer, ingressDevice *NitroClient, controllerInput *ControllerInput) {
	GetClusterCNI(api, controllerInput)	
	if (controllerInput.ClusterCNI == "Flannel") {
		InitFlannel(api, ingressDevice, controllerInput)
	} else {
		klog.Info("Network Automation is not supported for other than Flannel")
	}
}

