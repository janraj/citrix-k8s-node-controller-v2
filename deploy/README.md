# **Install Citrix Node Controller**

 1. Download the "citrix-k8s-node-controller.yaml" from the deployment Directory.
    ```
      wget  https://raw.githubusercontent.com/janraj/citrix-k8s-node-controller/master/deploy/citrix-k8s-node-controller.yaml?token=AMvewY7ooAOE6KZsmhr07BswqSTAj3Ilks5ceA_rwA%3D%3D
    ```
                        
    This yaml has four section, in which first three is for cluster role creation and service account creation and the 
    next one is for citrix node controller. 
    * Cluster roles
    * Cluster role bindings
    * Service account
    * Citrix Node Controller
   
    First three are required for citrix Node controller to monitor k8s events. No changes required.
    Next section defines environment variables required for Citrix Node Controller to configure the Citrix ADC.

 2. Update the following env variables, for Citrix Node Controller bringup.

       <details>
       <summary>NS_IP</summary>

         This is must for Citrix Node Controller to configure the NetScaler appliance. Provide,
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
       <summary>NS_VTEP_MAC</summary>

       </details>
    
3. Deploy Citrix Node Controller. 

   Deploy Citrix Node Controller  on kubernetes by using 'kubectl create' command
        
           kubectl create -f citrix-k8s-node-controller.yaml

   This pulls the latest image and brings up the Citrix Node Controller.
                
   Official Citrix Node Controller docker images is <span style="color:red"> `quay.io/citrix/citrix-k8s-node-controller:latest` </span> 
