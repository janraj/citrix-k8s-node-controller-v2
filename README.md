[![Build Status](https://travis-ci.com/janraj/citrix-k8s-node-controller.svg?token=GfEuWKxn7TJJesWboygR&branch=master)](https://travis-ci.com/janraj/citrix-k8s-node-controller)
[![codecov](https://codecov.io/gh/janraj/citrix-k8s-node-controller/branch/master/graph/badge.svg?token=9c5R8ukQGY)](https://codecov.io/gh/janraj/citrix-k8s-node-controller)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](./license/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/janraj/citrix-k8s-node-controller)](https://goreportcard.com/report/github.com/janraj/citrix-k8s-node-controller)
[![Docker Repository on Quay](https://quay.io/repository/citrix/citrix-k8s-node-controller/status "Docker Repository on Quay")](https://quay.io/repository/citrix/citrix-k8s-node-controller)
[![GitHub stars](https://img.shields.io/github/stars/janraj/citrix-k8s-node-controller.svg)](https://github.com/janraj/citrix-k8s-node-controller/stargazers)
[![HitCount](http://hits.dwyl.com/janraj/citrix-k8s-node-controller.svg)](http://hits.dwyl.com/janraj/citrix-k8s-node-controller)

# **Citrix-k8s-node-controller**
Citrix Node Controller  is a micro service which creates the network between cluster and ingress device. <span style="color:red">**Citrix Node Controller currently works only with flannel as CNI. The scope of Citrix node controller can be extended to other CNI.**</span>

## Table of Contents

- [Introduction](#introduction)
- [Architecture](#architecture)
- [How it Works](#how-it-works)
- [Getting Started](#getting-started)
- [Questions](#questions)
- [Issues](#issues)
- [Code of Conduct](#Code-of-conduct)
- [License](#api-reference)


## **Introduction**
When services on Kubernetes expose to external access via the Ingress device, there has to be proper networking between the Kubernetes nodes and ingress device to route the traffic into the cluster.   This is important because the pods will be having private IP’s based on the CNI framework.  These Private IP’s will not be able to directly access from ingress device without proper network configurations. Manual configuration to ensure such reachability is troublesome in Kubernetes world.

## **Architecture**
This is the high-level preview of Citrix node controller architecture. Following are the main components.	


![](./images/CitrixControllerArchitecture.png)
       <details>
       <summary>**Ingress Interface**</summary>
	    Ingress Interface is responsible for interacting with Citrix ADC via nitro rest API. It maintains the nitro session and invokes it when required. 
       </details>
       <details>
       <summary>**K8s Interface**</summary>
	    This module interacts with Kube API server via K8s Go Client. It ensures the availability of client and maintains a healthy client session.
       </details>
       <details>
       <summary>**Node Watcher**</summary>
	    The node watcher unit is used to watch the node events via K8s Interface. It responds to the node events such as node addition, deletion or modification with its call            back functions.
       </details>
       <details>
       <summary>**Input Feeder**</summary>
	    It provides inputs to the config decider. Some of the inputs are auto detect and the rest are taken from the CNC deployment yaml. 
       </details>
       <details>
       <summary>**Config Decider**</summary>
	    This segment takes inputs from both the node watcher and the input feeder and decides the best network automation required between cluster and NetScaler.
       </details>
       <details>
       <summary>**Core**</summary>
	    The core module interacts with node watcher and updates the corresponding config engine.  It is responsible for starting the best config engine for the corresponding             cluster.
       </details>

## **How it Works**

Citrix Node controller monitor the node events and establish a route between the node to Citrix ADC via VXLAN. Citrix Node Controller adds route on Citrix ADC when a new node joins to the cluster. Similarly when node leaves, citrix node controller removes the route from Citrix ADC. Citrix Node Controller uses VXLAN overlay between kubernetes cluster and Citrix ADC for service routing. 



## **Getting Started**

Citrix Node controller can be used in two flavour. 

	1) In cluster CNC Configuration [As a Process].
	2) Out of cluster CNC Configuration [As a Micro Service]

In cluster configuration is recomended for production. Out of cluster configuration can be used for easy development.
  
#### **As Processs***
```

        1) Download/Clone the citrix-k8s-node-controller.

        2) Perform "make run" from build folder.
                This starts the citrix node controller. Go binary has to be installed for running MIC. You have set the few inputs as Enviornment variable.

        3) Deploy Config MAP.
		kubectl apply -f https://raw.githubusercontent.com/janraj/citrix-k8s-node-controller/master/deploy/config_map.yaml
```
#### **As Micro Service***
	Please refer [deployment](deploy/README.md) page for running CNC as a micro service inside the cluster.

## **Questions**
For questions and support the following channels are available:
* [Citrix Discussion Forum](https://discussions.citrix.com/forum/1657-netscaler-cpx/). 
* [NetScaler Slack Channel](https://citrixadccloudnative.slack.com/)

## **Issues**
Describe the Issue in Details, Collects the logs and  Use the forum mentioned below
```
   https://discussions.citrix.com/forum/1657-netscaler-cpx/
```

## **Code of Conduct**
This project adheres to the [Kubernetes Community Code of Conduct](https://github.com/kubernetes/community/blob/master/code-of-conduct.md). By participating in this project you agree to abide by its terms.

## **License**
[Apache License 2.0](./license/LICENSE)
