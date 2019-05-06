# **Install Citrix Node Controller**

 1. **Download the "citrix-k8s-node-controller.yaml" from the deployment Directory.**
    ```
      wget  https://raw.githubusercontent.com/janraj/citrix-k8s-node-controller/master/deploy/citrix-k8s-node-controller.yaml?token=AMvewY7ooAOE6KZsmhr07BswqSTAj3Ilks5ceA_rwA%3D%3D
    ```
                        
    This yaml creates following on a new namespace called **citrix**

    * Cluster roles
    * Cluster role bindings
    * Service account
    * Citrix Node Controller service
   
    First three are required for citrix Node controller to monitor k8s events. No changes required.
    Next section defines environment variables required for Citrix Node Controller to configure the Citrix ADC.

 2. **Update the following env variables, for Citrix Node Controller bringup.**

    1. "Mandatory" Arguments:
       <details>
       <summary>NS_IP</summary>

         This is must for Citrix Node Controller to configure the NetScaler appliance. Citrix Node Controller uses NS_IP for configuration needs. NS_IP can be of,
         ```
            NSIP for standalone NetScaler  
            SNIP for HA (Management access has to be enabled) 
            CLIP for Cluster
         
         ```
       </details>
       <details>
       <summary>NS_USER and NS_PASSWORD</summary>

         This is for authenticating with NetScaler if it has non default username and password. We can directly pass username/password or use Kubernetes secrets.
         Please refer our [guide](https://github.com/citrix/citrix-k8s-ingress-controller/blob/master/docs/command-policy.md) for configuring a non default NetScaler username and password.
         
         Given Yaml uses k8s secrets. Following steps helps to create secrets to be used in yaml.

         Create secrets on Kubernetes for NS_USER and NS_PASSWORD
         Kubernetes secrets can be created by using 'kubectl create secret'.  

                 kubectl create secret  generic nslogin --from-literal=username='nsroot' --from-literal=password='nsroot'

         >**Note:** If you are using different secret name rather than nslogin, you have to update the "name" field in the yaml. 

       </details>
       <details>
       <summary>NS_POD_CIDR</summary>
         Reserve a Pod subnet for Netscaler by providing podcidr.
       </details>
       <details>
       <summary>NODE_CNI_CIDR</summary>
         CIDR of kubernetes cluster nodes from where each node  gets its own pod CIDR.
       </details>
    
    2. "Optional" Arguments:

       <details>
       <summary>NS_VTEP_MAC</summary>
         Optional Field. CNC automatically detect from NetScaler. If its failed, CNC has to restart with explict  NS_VTEP_MAC.
       </details>
       <details>
       <summary>NS_VTEP_IP</summary>
         Optional Field. If require different IP as VTEP than NS_IP, then please provide.
       </details>
       <details>
       <summary>NS_VXLAN_ID</summary>
         Optional Field. Default Value is 1. This argument is used in case of flannel. If flannel uses different VXLAN_ID, we have to provide this information.
       </details>
       <details>
       <summary>K8S_VXLAN_PORT</summary>
         Optional Field. If VXLAN port is other than 8472, you have to use this field to pass to Citrix Node Controller.
       </details>
3. **Deploy Citrix Node Controller.**

   Deploy Citrix Node Controller  on kubernetes by using 'kubectl create' command
        
           kubectl create -f citrix-k8s-node-controller.yaml

   This pulls the latest stable  image and brings up the Citrix Node Controller.
                
   Official Citrix Node Controller docker images is <span style="color:red"> `quay.io/citrix/citrix-k8s-node-controller:latest` </span>

4.  **Apply config map input.**
    
    Citrix Node controller is ready for operation and must be waiting for this input. Config map input is manadatory for citrix node controller to work. Config Map is used for controlling the citrix node operation via **operation** data field.

    ```
	kubectl apply -f https://raw.githubusercontent.com/janraj/citrix-k8s-node-controller/master/deploy/config_map.yaml
    ```
    if **operation** data field in config map is **ADD**, Citrix node controller creates routing configuration on netscaler. If its **DELETE**, citrix node controller remove the routing configuration added earlier. 
