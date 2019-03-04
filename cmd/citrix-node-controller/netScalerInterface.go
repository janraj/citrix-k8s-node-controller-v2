package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"k8s.io/klog"
	"log"
	"net/http"
	"strings"
)

type NitroClient struct {
	url       string
	statsURL  string
	username  string
	password  string
	proxiedNs string
	client    *http.Client
}

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
func NewNitroClient(obj *ControllerInput) *NitroClient {
	c := new(NitroClient)
	c.url = "http://" + strings.Trim(obj.IngressDeviceIP, " /") + "/nitro/v1/config/"
	c.username = obj.IngressDeviceUsername
	c.password = obj.IngressDevicePassword
	c.client = &http.Client{}
	return c
}

type responseHandlerFunc func(resp *http.Response) ([]byte, error)

func createResponseHandler(resp *http.Response) ([]byte, error) {
	switch resp.Status {
	case "201 Created", "200 OK":
		body, _ := ioutil.ReadAll(resp.Body)
		return body, nil
	case "409 Conflict":
		body, _ := ioutil.ReadAll(resp.Body)
		return body, errors.New("failed: " + resp.Status + " (" + string(body) + ")")

	case "207 Multi Status":
		//This happens in case of Bulk operations, which we do not support yet
		body, _ := ioutil.ReadAll(resp.Body)
		return body, nil
	case "400 Bad Request", "401 Unauthorized", "403 Forbidden",
		"404 Not Found", "405 Method Not Allowed", "406 Not Acceptable",
		"503 Service Unavailable", "599 Netscaler specific error":
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("[INFO] go-nitro: error = " + string(body))
		return body, errors.New("failed: " + resp.Status + " (" + string(body) + ")")
	default:
		body, err := ioutil.ReadAll(resp.Body)
		return body, err

	}
}

func (c *NitroClient) createHTTPRequest(method string, url string, buff *bytes.Buffer) (*http.Request, error) {
	req, err := http.NewRequest(method, url, buff)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	if c.proxiedNs == "" {
		req.Header.Set("X-NITRO-USER", c.username)
		req.Header.Set("X-NITRO-PASS", c.password)
	} else {
		req.SetBasicAuth(c.username, c.password)
		req.Header.Set("_MPS_API_PROXY_MANAGED_INSTANCE_IP", c.proxiedNs)
	}
	return req, nil
}
func (c *NitroClient) doHTTPRequest(method string, url string, bytes *bytes.Buffer, respHandler responseHandlerFunc) ([]byte, error) {
	req, err := c.createHTTPRequest(method, url, bytes)
	if err != nil {
		return []byte{}, err
	}
	resp, err := c.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return []byte{}, err
	}
	log.Println("[DEBUG] go-nitro: response Status:", resp.Status)
	return respHandler(resp)
}

func (c *NitroClient) createResource(resourceType string, resourceJSON []byte) ([]byte, error) {
	log.Println("[DEBUG] go-nitro: Creating resource of type ", resourceType)
	url := c.url + resourceType
	log.Println("[TRACE] go-nitro: url is ", url)
	return c.doHTTPRequest("POST", url, bytes.NewBuffer(resourceJSON), createResponseHandler)
}

func (c *NitroClient) AddResource(resourceType string, name string, resourceStruct interface{}) (string, error) {

	nsResource := make(map[string]interface{})
	nsResource[resourceType] = resourceStruct

	resourceJSON, err := JSONMarshal(nsResource)
	if err != nil {
		return "", fmt.Errorf("[ERROR] go-nitro: Failed to create resource of type %s, name=%s, err=%s", resourceType, name, err)
	}

	log.Printf("[TRACE] go-nitro: Resourcejson is " + string(resourceJSON))

	body, err := c.createResource(resourceType, resourceJSON)
	if err != nil {
		return "", fmt.Errorf("[ERROR] go-nitro: Failed to create resource of type %s, name=%s, err=%s", resourceType, name, err)
	}
	_ = body

	return name, nil
}

type Vxlan struct {
	Dynamicrouting     string `json:"dynamicrouting,omitempty"`
	Id                 int    `json:"id,omitempty"`
	Ipv6dynamicrouting string `json:"ipv6dynamicrouting,omitempty"`
	Port               int    `json:"port,omitempty"`
	Td                 int    `json:"td,omitempty"`
	Vlan               int    `json:"vlan,omitempty"`
}
type Vxlan_srcip_binding struct {
	Id        int    `json:"id,omitempty"`
	Ipaddress string `json:"ipaddress,omitempty"`
	Srcip     string `json:"srcip",omitempty"`
	Netmask   string `json:"netmask,omitempty"`
}
type Route struct {
	Network string `json:"network,omitempty"`
	Gateway string `json:"gateway,omitempty"`
	Netmask string `json:"netmask,omitempty"`
}
type Nsip struct {
	Ipaddress string `json:"ipaddress,omitempty"`
	Netmask   string `json:"netmask,omitempty"`
}

type Arp struct {
	Ipaddress string `json:"ipaddress,omitempty"`
	Mac       string `json:"mac,omitempty"`
	Vxlan     string `json:"vxlan,omitempty"`
	Vtep      string `json:"vtep,omitempty"`
}

// Key the key of the dictionary
type Key interface{}

// Value the content of the dictionary
type Value interface{}

// ValueDictionary the set of Items
type ConfigPack struct {
	items map[Key]Value
	keys  []string
}

type SameSubnet struct {
	items map[Key]Value
}

// Set adds a new item to the dictionary
func (d *ConfigPack) Set(k Key, v Value) {
	if d.items == nil {
		d.items = make(map[Key]Value)
	}
	d.items[k] = v
	d.keys = append(d.keys, k.(string))
}

//func main() {
func NsInterface(client *NitroClient, obj *ControllerInput) {
	flannel := ConfigPack{}
	vxlan := Vxlan{
		Id:   5556,
		Port: 3336,
	}
	flannel.Set("vxlan", &vxlan)
	vxlanbind := Vxlan_srcip_binding{
		Id:    5556,
		Srcip: "10.102.169.244",
	}
	flannel.Set("vxlan_srcip_binding", &vxlanbind)

	nsip := Nsip{
		Ipaddress: "10.244.250.250",
		Netmask:   "255.255.0.0",
	}
	flannel.Set("nsip", &nsip)

	route := Route{
		Network: "10.244.1.0",
		Netmask: "255.255.255.0",
		Gateway: "10.244.1.0",
	}
	flannel.Set("route", &route)

	arp := Arp{
		Ipaddress: "10.244.1.0",
		Mac:       "2e:3c:5a:0c:64:68",
		Vxlan:     "5556",
		Vtep:      "10.102.33.198",
	}
	flannel.Set("arp", &arp)

	for ind, _ := range flannel.keys {
		result, err := client.AddResource(flannel.keys[ind], "flannel", flannel.items[flannel.keys[ind]])
		if err != nil {
			fmt.Println("Result err ", result, err)
		}
	}
}
func createIngressDeviceClient(input *ControllerInput) *NitroClient {
	client := NewNitroClient(input)
	return client
}
func AddIngressDeviceConfig(config *ConfigPack, client *NitroClient) {
	for ind, _ := range config.keys {
		result, err := client.AddResource(config.keys[ind], "ADD", config.items[config.keys[ind]])
		if err != nil {
			fmt.Println("Result  err ", result, err)
		}
	}
}

/*
*************************************************************************************************
*   APIName :  InitializeNode                                                                   *
*   Input   :  Nil.					             			        *
*   Output  :  Nil.				                                                *
*   Descr   :  This API initialize a node and return it.					*
*************************************************************************************************
 */
func NsInterfaceAddRoute(client *NitroClient, input *ControllerInput, node *Node) {
	configPack := ConfigPack{}
	route := Route{
		Network: node.PodAddress,
		Netmask: node.PodNetMask,
		Gateway: node.PodAddress,
	}
	configPack.Set("route", &route)

	arp := Arp{
		Ipaddress: node.PodAddress,
		Mac:       node.PodVTEP,
		Vxlan:     input.IngressDeviceVxlanIDs,
		Vtep:      node.IPAddr,
	}
	configPack.Set("arp", &arp)
	AddIngressDeviceConfig(&configPack, client)
}

func NsInterfaceDeleteRoute(client *NitroClient, obj *ControllerInput, nodeinfo *Node) {
	var argsBundle = map[string]string{"network": nodeinfo.PodAddress, "netmask": nodeinfo.PodNetMask, "gateway": nodeinfo.PodAddress}
	err2 := client.DeleteResourceWithArgsMap("route", "", argsBundle)
	if err2 != nil {
		fmt.Println(err2)
	}
	argsBundle = map[string]string{"Ipaddress": nodeinfo.PodAddress}
	err2 = client.DeleteResourceWithArgsMap("arp", "", argsBundle)
	if err2 != nil {
		fmt.Println(err2)
	}

}
func (ingressDevice *NitroClient) GetVxlanConfig(controllerInput *ControllerInput) {
	klog.Info("GetVxlanConfig")
}

//DeleteResourceWithArgsMap deletes a resource of supplied type and name. Args are supplied as map of key value
func (c *NitroClient) DeleteResourceWithArgsMap(resourceType string, resourceName string, args map[string]string) error {

	_, err := c.listResourceWithArgsMap(resourceType, resourceName, args)
	if err == nil { // resource exists
		log.Printf("[INFO] go-nitro: DeleteResource found resource of type %s: %s", resourceType, resourceName)
		_, err = c.deleteResourceWithArgsMap(resourceType, resourceName, args)
		if err != nil {
			log.Printf("[ERROR] go-nitro: Failed to delete resourceType %s: %s, err=%s", resourceType, resourceName, err)
			return err
		}
	} else {
		log.Printf("[INFO] go-nitro: Resource %s already deleted ", resourceName)
	}
	return nil
}
func (c *NitroClient) deleteResourceWithArgs(resourceType string, resourceName string, args []string) ([]byte, error) {
	log.Println("[DEBUG] go-nitro: Deleting resource of type ", resourceType, "with args ", args)
	var url string
	if resourceName != "" {
		url = c.url + fmt.Sprintf("%s/%s?args=", resourceType, resourceName)
	} else {
		url = c.url + fmt.Sprintf("%s?args=", resourceType)
	}
	url = url + strings.Join(args, ",")
	log.Println("[TRACE] go-nitro: url is ", url)

	return c.doHTTPRequest("DELETE", url, bytes.NewBuffer([]byte{}), deleteResponseHandler)

}

func (c *NitroClient) deleteResourceWithArgsMap(resourceType string, resourceName string, argsMap map[string]string) ([]byte, error) {
	args := make([]string, len(argsMap))
	i := 0
	for key, value := range argsMap {
		args[i] = fmt.Sprintf("%s:%s", key, value)
		i++
	}
	return c.deleteResourceWithArgs(resourceType, resourceName, args)

}
func (c *NitroClient) listResourceWithArgsMap(resourceType string, resourceName string, argsMap map[string]string) ([]byte, error) {
	args := make([]string, len(argsMap))
	i := 0
	for key, value := range argsMap {
		args[i] = fmt.Sprintf("%s:%s", key, value)
		i++
	}
	return c.listResourceWithArgs(resourceType, resourceName, args)

}
func (c *NitroClient) listResourceWithArgs(resourceType string, resourceName string, args []string) ([]byte, error) {
	log.Println("[DEBUG] go-nitro: listing resource of type ", resourceType, ", name: ", resourceName, ", args:", args)
	var url string

	if resourceName != "" {
		url = c.url + fmt.Sprintf("%s/%s", resourceType, resourceName)
	} else {
		url = c.url + fmt.Sprintf("%s", resourceType)
	}
	strArgs := strings.Join(args, ",")
	url2 := url + "?args=" + strArgs
	log.Println("[TRACE] go-nitro: url is ", url)

	data, err := c.doHTTPRequest("GET", url2, bytes.NewBuffer([]byte{}), readResponseHandler)
	if err != nil {
		log.Println("[DEBUG] go-nitro: error listing with args, trying filter")
		url2 = url + "?filter=" + strArgs
		return c.doHTTPRequest("GET", url2, bytes.NewBuffer([]byte{}), readResponseHandler)
	}
	return data, err

}
func readResponseHandler(resp *http.Response) ([]byte, error) {
	switch resp.Status {
	case "200 OK":
		body, _ := ioutil.ReadAll(resp.Body)
		return body, nil
	case "404 Not Found":
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("[DEBUG] go-nitro: read: 404 not found")
		return body, errors.New("go-nitro: read: 404 not found: ")
	case "400 Bad Request", "401 Unauthorized", "403 Forbidden",
		"405 Method Not Allowed", "406 Not Acceptable",
		"409 Conflict", "503 Service Unavailable", "599 Netscaler specific error":
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("[INFO] go-nitro: read: error = " + string(body))
		return body, errors.New("[INFO] go-nitro: failed read: " + resp.Status + " (" + string(body) + ")")
	default:
		body, err := ioutil.ReadAll(resp.Body)
		log.Println("[INFO] go-nitro: read error = " + string(body))
		return body, err

	}
}
func deleteResponseHandler(resp *http.Response) ([]byte, error) {
	switch resp.Status {
	case "200 OK", "404 Not Found":
		body, _ := ioutil.ReadAll(resp.Body)
		return body, nil

	case "400 Bad Request", "401 Unauthorized", "403 Forbidden",
		"405 Method Not Allowed", "406 Not Acceptable",
		"409 Conflict", "503 Service Unavailable", "599 Netscaler specific error":
		body, _ := ioutil.ReadAll(resp.Body)
		log.Println("[INFO] go-nitro: delete: error = " + string(body))
		return body, errors.New("[INFO] delete failed: " + resp.Status + " (" + string(body) + ")")
	default:
		body, err := ioutil.ReadAll(resp.Body)
		return body, err

	}
}
