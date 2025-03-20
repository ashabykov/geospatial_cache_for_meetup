## Prepare

1. Install Docker and start Colima:
    ```sh
    make colima-start
    ```

2. Bring up the Docker containers:
    ```sh
    make up
    ```

3. Install dependencies:
    ```sh
    make deps
    ```

4. Run tests to ensure everything is set up correctly:
    ```sh
    make test
    ```

## Prepare Kafka UI
1. Pull Kafka UI image:

    ```sh
    docker pull provectuslabs/kafka-ui:latest
    ```

2. Run Kafka UI container:
    ```sh
    docker run -d --rm -p 8080:8080 \
        -e KAFKA_CLUSTERS_0_NAME=local \
        -e KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS=kafka:9092 \
        -e KAFKA_CLUSTERS_0_ZOOKEEPER=zookeeper:2181 \
        provectuslabs/kafka-ui:latest
    ```

3. Access Kafka UI at [http://localhost:8082](http://localhost:8082)
4. Set up Kafka UI:
    - Cluster Name: `local`
    - Bootstrap Servers: `kafka:29092`
    - Click `Validate`
    - If everything is correct, click `Submit`

## Run main application
1. Start the main application:
    ```sh
    make run
    ```