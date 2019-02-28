# Citrix-k8s-node-controller
Citrix Node Controller (CNC) is a micro service which creates the network between cluster and ingress device.  Citrix node controller can run inside the cluster as a pod or outside the cluster. In case of outside the cluster, it requires a proper Kube config file to run successfully.

## Description
When services on Kubernetes expose to external access via the Ingress device, there has to be proper networking between the Kubernetes nodes and ingress device to route the traffic into the cluster.   This is important because the pods will be having private IP’s based on the CNI framework.  These Private IP’s will not be able to directly access from ingress device without proper network configurations. Manual configuration to ensure such reachability is troublesome in Kubernetes world.

# **Questions**
For questions and support the following channels are available:
* [Citrix Discussion Forum](https://discussions.citrix.com/forum/1657-netscaler-cpx/). 
* [NetScaler Slack Channel](https://citrixadccloudnative.slack.com/)

# **Issues**
Describe the Issue in Details, Collects the logs and  Use the forum mentioned below
```
   https://discussions.citrix.com/forum/1657-netscaler-cpx/
```

# **Code of Conduct**
This project adheres to the [Kubernetes Community Code of Conduct](https://github.com/kubernetes/community/blob/master/code-of-conduct.md). By participating in this project you agree to abide by its terms.

# **License**
[Apache License 2.0](./license/LICENSE)
