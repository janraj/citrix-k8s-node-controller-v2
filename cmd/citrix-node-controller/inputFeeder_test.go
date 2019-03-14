package main

import (
	"testing"
	"os"
)

func TestFetchCitrixNodeControllerInput(t *testing.T) {
	IngressDeviceIP := os.Getenv("NS_IP")
	os.Setenv("NS_IP", "")
	IngressDeviceVtepMAC := os.Getenv("NS_VTEP_MAC")
	os.Setenv("NS_VTEP_MAC", "")
	IngressDeviceUsername := os.Getenv("NS_LOGIN")
	os.Setenv("NS_LOGIN", "")
	IngressDevicePassword := os.Getenv("NS_PASSWORD")
	os.Setenv("NS_PASSWORD", "")
	IngressDeviceVtepIP := os.Getenv("NS_SNIP")
	os.Setenv("NS_SNIP", "")
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("FetchCitrixNodeControllerInput should have panicked!")
			}
		}()
		// This function should cause a panic 
		FetchCitrixNodeControllerInput()
	}()

	os.Setenv("NS_IP", IngressDeviceIP)
	os.Setenv("NS_VTEP_MAC", IngressDeviceVtepMAC)
	os.Setenv("NS_LOGIN", IngressDeviceUsername)
	os.Setenv("NS_PASSWORD", IngressDevicePassword)
	os.Setenv("NS_SNIP", IngressDeviceVtepIP)
	FetchCitrixNodeControllerInput()
}
