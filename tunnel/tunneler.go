package tunnel

//Tunneler interface defines implementaion for creating a connection tunnel
type Tunneler interface {
	Connect() error
	Close() error
}
