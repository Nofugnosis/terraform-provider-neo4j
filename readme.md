# Terraform Neo4j provider

## Build
Run the following command to build the provider
`
go build -o terraform-provider-neo4j
`
## Run Neo4j



## Testing

 - run a Neo4j server for testing:

```
docker run -d --name neo4j -p 7474:7474 -p 7687:7687 --env=NEO4J_ACCEPT_LICENSE_AGREEMENT=yes --env NEO4J_AUTH=neo4j/password neo4j:enterprise
```

> :warning: **Neo4j Enterprise is a licensed product. Please read the [official license documentation](https://neo4j.com/licensing)**

 - Export the required variables:
```
export TF_LOG=DEBUG
export NEO4J_HOST=neo4j://localhost:7687
export NEO4J_USERNAME=neo4j
export NEO4J_PASSWORD=password
```

 - Run the tests
```
make test
```