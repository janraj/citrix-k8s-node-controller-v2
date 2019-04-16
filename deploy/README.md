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
       <summary>NS_SNIP</summary>

         NS_SNIP used for configuring as VTEP IP when cluster CNI as flannel and it is also used for creating dummy node in kubernetes cluster  

       </details>
       <details>
       <summary>NS_USER and NS_PASSWORD</summary>

         This is for authenticating with NetScaler if it has non default username and password. We can directly pass username/password or use Kubernetes secrets.

       </details>
       <details>
       <summary>NS_VTEP_MAC</summary>
         NetScaler MAC of VTEP IP, which is required for Flannel CNI.
       </details>
       <details>
       <summary>NS_POD_CIDR</summary>
         Reserve a Pod subnet for Netscaler by providing podcidr.
       </details>
       <details>
       <summary>NODE_CNI_CIDR</summary>
         CIDR of kubernetes cluster nodes from where each node  gets its own pod CIDR.
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
