# Deploy the Citrix k8s node controller
  This creates Citrix Node Controller on Kubernetes and establish the route between Citrix ADC and kubernetes Nodes.

Perform the following:

1.  Download the `citrix-k8s-node-controller.yaml` deployment file using the following command:

        wget  https://raw.githubusercontent.com/janraj/citrix-k8s-node-controller-v2/master/deploy/citrix-k8s-node-controller.yaml

    The deployment file contains definitions for the following:

    -  Cluster Role (`ClusterRole`)

    -  Cluster Role Bindings (`ClusterRoleBinding`)

    -  Service Account (`ServiceAccount`)

    -  Citrix Node Controller service (`citrix-node-controller`)

    You don't have to modify the definitions for `ClusterRole`, `ClusterRoleBinding`, and `ServiceAccount` definitions. The definitions are used by Citrix node controller to monitor Kubernetes events. But, in the `citrix-node-controller` definition you have to provide the values for the environment variables that is required for Citrix k8s node controller to configure the Citric ADC.

    You must provide values for the following environment variables in the Citrix k8s node controller service definition:

    | Environment Variable | Mandatory or Optional | Description |
    | -------------------- | --------------------- | ----------- |
    | NS_IP | Mandatory | Citrix k8s node controller uses this IP address to configure the Citrix ADC. The NS_IP can be anyone of the following: </br></br> - **NSIP** for standalone Citrix ADC </br>- **SNIP** for high availability deployments (Ensure that management access is enabled) </br> - **CLIP** for Cluster deployments |
    | NS_USER and NS_PASSWORD | Mandatory | The user name and password of Citrix ADC. Citrix k8s node controller uses these credentials to authenticate with Citrix ADC. You can either provide the user name and password or Kubernetes secrets. If you want to use a non-default Citrix ADC user name and password, you can [create a system user account in Citrix ADC](https://developer-docs.citrix.com/projects/citrix-k8s-ingress-controller/en/latest/deploy/deploy-cic-yaml/#create-system-user-account-for-citrix-ingress-controller-in-citrix-adc). </br></br> The deployment file uses Kubernetes secrets, create a secret for the user name and password using the following command: </br></br> `kubectl create secret  generic nslogin --from-literal=username='nsroot' --from-literal=password='nsroot'` </br></br> **Note**: If you want to use a different secret name other than `nslogin`, ensure that you update the `name` field in the `citrix-node-controller` definition. |
    | NETWORK | Mandatory | The IP address range (for example, `192.128.1.0/24`) that Citrix node controller uses to configure the VTEP overlay end points on the Kubernetes nodes.|
    | VNID | Mandatory | A unique VXLAN VNID to create a VXLAN overlays between kubernetes cluster and the ingress devices. </br></br>**Note:** Ensure that the VXLAN VNID that you use does not conflict with the Kubernetes cluster or Citrix ADC VXLAN VNID.|
    | VXLAN_PORT | Mandatory | The VXLAN port that you want to use for the overlay.|
    | REMOTE_VTEPIP | Mandatory | The Ingress Citrix ADC SNIP.|

1.  After you have updated the Citrix k8s node controller deployment YAML file, deploy it using the following command:

        kubectl create -f citrix-k8s-node-controller.yaml

1.  Create the configmap using the following command:

        kubectl apply -f https://raw.githubusercontent.com/janraj/citrix-k8s-node-controller-v2/master/deploy/config_map.yaml

# Delete the Citrix K8s node conroller 


1.  Delete the [config map](config_map.yaml) using the following command:
	
	When we delete the configmap, citrix node controller clean up teh configuration created on Citrix ADC.

        kubectl delete -f https://raw.githubusercontent.com/janraj/citrix-k8s-node-controller-v2/master/deploy/config_map.yaml


1.  Delete the citrix node controller  using the following command:

        kubectl delete -f citrix-k8s-node-controller.yaml
