
### Para nueva instalacion en el cluster


``kubectl apply -f ./k8s/mongo/storage.yaml``

cambiar el values.yaml del entorno y aplicar el comando de helm

``helm install mongo-poc -f ./k8s/mongo/values.yaml bitnami/mongodb-sharded --namespace measures``

### Para modificar config en el cluster


cambiar el values.yaml del entorno y aplicar el comando de helm

`` 
helm upgrade mongo-poc -f ./k8s/mongo/values.yaml bitnami/mongodb-sharded --namespace measures``