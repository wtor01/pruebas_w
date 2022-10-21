  ### Crear hpa en base a pubsub

[Doc gcp](https://cloud.google.com/kubernetes-engine/docs/tutorials/autoscaling-metrics?hl=es-419#step1)

1. Otórgale al usuario la capacidad de crear las funciones de autorización requeridas:

    
    ``kubectl create clusterrolebinding cluster-admin-binding \
        --clusterrole cluster-admin --user "$(gcloud config get-value account)"``


2. Implementa el adaptador del modelo de recursos nuevo en el clúster:

    
    ``kubectl apply -f https://raw.githubusercontent.com/GoogleCloudPlatform/k8s-stackdriver/master/custom-metrics-stackdriver-adapter/deploy/production/adapter_new_resource_model.yaml``


3. Crea un objeto HorizontalPodAutoscaler.

    Hay ejemplos en overlays/dev y overlays/staging