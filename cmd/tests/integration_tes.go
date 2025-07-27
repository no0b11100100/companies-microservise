package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var kafkaBootstrap string

func TestMain(m *testing.M) {
	ctx := context.Background()

	// --- Setup MySQL container ---
	mysqlContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mysql:8",
			ExposedPorts: []string{"3306/tcp"},
			Env: map[string]string{
				"MYSQL_ROOT_PASSWORD": "test",
				"MYSQL_DATABASE":      "testdb",
				"MYSQL_USER":          "test",
				"MYSQL_PASSWORD":      "test",
			},
			WaitingFor: wait.ForListeningPort("3306/tcp").WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		panic(err)
	}
	defer mysqlContainer.Terminate(ctx)

	mysqlHost, _ := mysqlContainer.Host(ctx)
	mysqlPort, _ := mysqlContainer.MappedPort(ctx, "3306")

	// --- Setup Zookeeper container ---
	zookeeperContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "confluentinc/cp-zookeeper:7.3.0",
			ExposedPorts: []string{"2181/tcp"},
			Env: map[string]string{
				"ZOOKEEPER_CLIENT_PORT": "2181",
				"ZOOKEEPER_TICK_TIME":   "2000",
			},
			WaitingFor: wait.ForListeningPort("2181/tcp").WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		panic(err)
	}
	defer zookeeperContainer.Terminate(ctx)

	zooHost, _ := zookeeperContainer.Host(ctx)
	zooPort, _ := zookeeperContainer.MappedPort(ctx, "2181/tcp")
	zookeeperAddress := fmt.Sprintf("%s:%s", zooHost, zooPort.Port())

	// --- Setup Kafka container ---
	kafkaContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "confluentinc/cp-kafka:7.3.0",
			ExposedPorts: []string{"9092/tcp"},
			Env: map[string]string{
				"KAFKA_BROKER_ID":                        "1",
				"KAFKA_ZOOKEEPER_CONNECT":                zookeeperAddress,
				"KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR": "1",
				"KAFKA_ADVERTISED_LISTENERS":             "PLAINTEXT://localhost:9092",
				"KAFKA_LISTENERS":                        "PLAINTEXT://0.0.0.0:9092",
			},
			WaitingFor: wait.ForListeningPort("9092/tcp").WithStartupTimeout(60 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		panic(err)
	}
	defer kafkaContainer.Terminate(ctx)

	kafkaHost, _ := kafkaContainer.Host(ctx)
	kafkaPort, _ := kafkaContainer.MappedPort(ctx, "9092")
	kafkaBootstrap = fmt.Sprintf("%s:%s", kafkaHost, kafkaPort.Port())

	// --- Set ENV variables for your app ---
	os.Setenv("DB_USER", "test")
	os.Setenv("DB_PASSWORD", "test")
	os.Setenv("DB_HOST", mysqlHost)
	os.Setenv("DB_PORT", mysqlPort.Port())
	os.Setenv("DB_NAME", "testdb")

	os.Setenv("KAFKA_BROKER", kafkaBootstrap)

	os.Setenv("HTTP_HOST", "localhost")
	os.Setenv("HTTP_PORT", "8080")

	time.Sleep(10 * time.Second) // wait for app startup

	os.Exit(m.Run())
}
