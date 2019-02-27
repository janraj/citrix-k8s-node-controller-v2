package main
import (
	"testing"	
	"runtime"
	"k8s.io/klog"
       )

func TestcreateIngressDeviceClient(t *testing.T) {
	 _, filename, _, _ := runtime.Caller(0)
    		klog.Info("Current test filename: " + filename)
	ControllerInput := FetchCitrixNodeControllerInput()
	client := createIngressDeviceClient(ControllerInput)
	if (client == nil){
    		t.Error("Expected a Valid Client ")
		
	}
}

