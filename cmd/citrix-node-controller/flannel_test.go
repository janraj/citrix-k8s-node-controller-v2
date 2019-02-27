package main
import (
	"testing"	
	"runtime"
	"fmt"
	"k8s.io/klog"
       )

func getClientAndDeviceInfo()(*ControllerInput, *KubernetesAPIServer){
        nsObj := FetchCitrixNodeControllerInput()
        api, err := CreateK8sApiserverClient()
        if err != nil {
                klog.Fatal("K8s Client Error", err)
	}
        return nsObj, api	
}

func TestGetDummyNode(t *testing.T) {
	 _, filename, _, _ := runtime.Caller(0)
    		fmt.Println("Current test filename: " + filename)
        nsObj, api := getClientAndDeviceInfo()
	node := api.GetDummyNode(nsObj)
	if (node == nil){
		 t.Error("Expected Node but its NULL ")		
	}
	nsObj.DummyNodeLabel = "DUMMY"
	node = api.GetDummyNode(nsObj)
	if (node != nil){
		 t.Error("Expected Null Node but it does not ", node)		
	}
}

func TestInitializeNode(t *testing.T){
        nsObj, _ := getClientAndDeviceInfo()
	node := InitializeNode(nsObj)
	if (node == nil){
		 t.Error("Expected Node but its NULL ")		
	}
}
func TestCreateDummyNode(t *testing.T){
        nsObj, api := getClientAndDeviceInfo()
	api.CreateDummyNode(nsObj)
}
