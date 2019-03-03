package main

import (
	"k8s.io/klog"
)
func  InitCitrixNodeController ()(error){
	klog.InitFlags(nil)
	klog.Info("Initializing CNC")
        return nil
}
func StartCitrixNodeController(){
        controllerInput := FetchCitrixNodeControllerInput()
        api, err := CreateK8sApiserverClient()
        if err != nil {
                klog.Fatal("K8s Client Error", err)
	}
        ingressDevice := createIngressDeviceClient(controllerInput)
	ConfigDecider(api, ingressDevice, controllerInput)
        CitrixNodeWatcher(api, ingressDevice, controllerInput) 
}
func main() {
	InitCitrixNodeController()	
	StartCitrixNodeController()
}

