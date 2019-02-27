package main
import (
	"testing"	
	"runtime"
	"fmt"
	"k8s.io/klog"
       )

func TestConvertPrefixLenToMask(t *testing.T) {
	 _, filename, _, _ := runtime.Caller(0)
    		fmt.Println("Current test filename: " + filename)
	t.Log("Prefix Length 24")
	mask := ConvertPrefixLenToMask("24")
	if mask != "255.255.255.0" {
    		t.Error("Expected 255.255.255.0, got ", mask)
  	}
	mask = ConvertPrefixLenToMask("8")
	if mask != "255.0.0.0" {
    		t.Error("Expected 255.0.0.0, got ", mask)
  	}
	mask = ConvertPrefixLenToMask("16")
	if mask != "255.255.0.0" {
    		t.Error("Expected 255.255.0.0, got ", mask)
  	}
	mask = ConvertPrefixLenToMask("30")
	if mask != "255.255.255.252" {
    		t.Error("Expected 255.255.255.252, got ", mask)
  	}
	mask = ConvertPrefixLenToMask("29")
	if mask != "255.255.255.248" {
    		t.Error("Expected 255.255.255.248, got ", mask)
  	}
	mask = ConvertPrefixLenToMask("25")
	if mask != "255.255.255.128" {
    		t.Error("Expected 255.255.255.128, got ", mask)
  	}
	mask = ConvertPrefixLenToMask("19")
	if mask != "255.255.224.0" {
    		t.Error("Expected 255.255.224.0, got ", mask)
  	}
	mask = ConvertPrefixLenToMask("17")
	if mask != "255.255.128.0" {
    		t.Error("Expected 255.255.128.0, got ", mask)
  	}
	mask = ConvertPrefixLenToMask("30")
	if mask != "255.255.255.252" {
    		t.Error("Expected 255.255.255.252, got ", mask)
  	}
}

func TestConfigDecider(t *testing.T){
	 _, filename, _, _ := runtime.Caller(0)
    		fmt.Println("Current test filename: " + filename)
        ControllerInputObj := FetchCitrixNodeControllerInput()
        k8sclient, err := CreateK8sApiserverClient()
        if err != nil {
                klog.Fatal("K8s Client Error", err)
	}
        IngressDeviceClient := createIngressDeviceClient(ControllerInputObj)
	
	ConfigDecider(k8sclient, IngressDeviceClient, ControllerInputObj) 
}
func TestCreateK8sApiserverClient(t *testing.T){
	client, err := CreateK8sApiserverClient()
	if client == nil{
		t.Error("Expected valid K8s Client")
	}
	if err != nil{
		t.Error("Expected valid K8s Client without error\n client=", client, "error=", err)
	}
		
}
