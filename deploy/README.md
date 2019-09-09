# Deploy the Citrix k8s node controller

Citrix k8s node controller is controlled using a [config map](https://github.com/janraj/citrix-k8s-node-controller/blob/master/deploy/config_map.yaml). The [config map](https://github.com/janraj/citrix-k8s-node-controller/blob/master/deploy/config_map.yaml) file contains a `data.operation:` field that you can use to define Citrix k8s node controller to automatically create, apply, and delete routing configuration on Citrix ADC. You can use the following values for the `data.operation:` field:

| **Value** | **Description** |
| ----- | ----------- |
| ADD | Citrix k8s node controller creates a routing configuration on the Citrix ADC instance. |
| DELETE | Citrix k8s node controller deletes the routing configuration on the Citrix ADC instance. |

[config_map.yaml](https://github.com/janraj/citrix-k8s-node-controller/blob/master/deploy/config_map.yaml):

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: citrix
  labels:
    name: citrix
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: citrix-node-controller
  namespace: citrix
data:
  operation: "ADD"
```

## Deploy the Citrix k8s node controller

Perform the following:

1.  Download the `citrix-k8s-node-controller.yaml` deployment file using the following command:

        wget  https://raw.githubusercontent.com/janraj/citrix-k8s-node-controller/master/deploy/citrix-k8s-node-controller.yaml

    The deployment file contains definitions for the following:

    -  Cluster Role (`ClusterRole`)

    -  Cluster Role Bindings (`ClusterRoleBinding`)

    -  Service Account (`ServiceAccount`)

    -  Citrix Node Controller service (`citrix-node-controller`)

    You don't have to modify the definitions for `ClusterRole`, `ClusterRoleBinding`, and `ServiceAccount` definitions. The definitions are used by Citrix node controller to monitor Kubernetes events. But, in the`citrix-node-controller` definition you have to provide the values for the environment variables that is required for Citrix k8s node controller to configure the Citric ADC.

    You must provide values for the following environment variables in the Citrix k8s node controller service definition:

    | Environment Variable | Mandatory or Optional | Description |
    | -------------------- | --------------------- | ----------- |
    | NS_IP | Mandatory | Citrix k8s node controller uses this IP address to configure the Citrix ADC. The NS_IP can be anyone of the following: </br> - NSIP for standalone Citrix ADC </br>- SNIP for high availability deployments (Ensure that management access is enabled) </br> - CLIP for Cluster deployments |
    | NS_USER and NS_PASSWORD | Mandatory | The user name and password of Citrix ADC. Citrix k8s node controller uses these credentials to authenticate with Citrix ADC. You can either provide the user name and password or Kubernetes secrets. If you want to use a non-default Citrix ADC user name and password, you can [create a system user account in Citrix ADC](https://developer-docs.citrix.com/projects/citrix-k8s-ingress-controller/en/latest/deploy/deploy-cic-yaml/#create-system-user-account-for-citrix-ingress-controller-in-citrix-adc). </br> The deployment file uses Kubernetes secrets, create a secret for the user name and password using the following command: </br> `kubectl create secret  generic nslogin --from-literal=username='nsroot' --from-literal=password='nsroot'` </br> **Note**: If you want to use a different secret name other than `nslogin`, ensure that you update the `name` field in the `citrix-node-controller` definition. |
    | ADDRESS | Mandatory | kube-router uses this address to configure the VTEP overlay end points on nodes.| 
    | VNID | Mandatory | A unique VNID tp create a VXLAn overlays between kubernetes nodes and ingress devices.|
    | CNI_NAME | Mandatory | Provide the CNI name used in the cluster[flannel, calico, openshift-azure, etc].|
    | K8S_VXLAN_PORT | Mandatory | VXLAN PORT for overlays.|
    | REMOTE_VTEPIP | Mandatory | Ingress device VTEP IP|

1.  After you have updated the Citrix k8s node controller deployment YAML file, deploy it using the following command:

        kubectl create -f citrix-k8s-node-controller.yaml

1.  Apply the [config map](https://github.com/janraj/citrix-k8s-node-controller/blob/master/deploy/config_map.yaml) using the following command:

        kubectl apply -f https://raw.githubusercontent.com/janraj/citrix-k8s-node-controller-v2/master/deploy/config_map.yaml
