package clients

type DataNodeClient interface {
	IsHealthy() bool
}
