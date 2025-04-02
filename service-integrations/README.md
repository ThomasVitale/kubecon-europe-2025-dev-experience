# Running

```
helm install postgresql bitnami/postgresql --version 12.5.7 --set "image.debug=true" --set "primary.initdb.user=postgres" --set "primary.initdb.password=postgres" --set "primary.initdb.scriptsConfigMap=appointments-init-sql" --set "global.postgresql.auth.postgresPassword=postgres"
```

```
helm upgrade --install dapr dapr/dapr \                                    
--version=1.15 \      
--namespace dapr-system \
--create-namespace \
--wait
```

```
helm install kafka bitnami/kafka --version 22.1.5 --set "provisioning.topics[0].name=events-topic" --set "provisioning.topics[0].partitions=1"
```