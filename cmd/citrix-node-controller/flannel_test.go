package main
import (
	"testing"	
	"runtime"
	"fmt"
        "k8s.io/client-go/kubernetes/fake"
       )

func getClientAndDeviceInfo()(*ControllerInput, *KubernetesAPIServer){
        nsObj := FetchCitrixNodeControllerInput()
        fake := fake.NewSimpleClientset()
        api := &KubernetesAPIServer{
           Suffix: "Test", 
           Client: fake,
        } 
        return nsObj, api	
}

func TestGetDummyNode(t *testing.T) {
	 _, filename, _, _ := runtime.Caller(0)
    		fmt.Println("Current test filename: " + filename)
        nsObj, api := getClientAndDeviceInfo()
	nsObj.DummyNodeLabel = "DUMMY"
	node := api.GetDummyNode(nsObj)
	if (node == nil){
	         node1 := api.CreateDummyNode(nsObj)
                 if (node1 == nil){
		     t.Error("Created Node and Get nodes are different ")		
                 }
	}else {
	         node=api.CreateDummyNode(nsObj)
                 if (node == nil){
		     t.Error("Expected Nil since there is already a nOde with same Label ")		
		 }
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
