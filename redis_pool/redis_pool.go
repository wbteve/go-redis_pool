//
// Redis Connection Pool written in GO
//

package redis_pool

import "fmt"
import "errors"
import "github.com/fzzy/radix/redis"
import "github.com/alecthomas/log4go"
import cp "github.com/gnagel/go-connection_pool/connection_pool"

type ConnectionMode int

const (
	// Populate the connection pool, but don't actually connect to redis
	LAZY = itoa
	// Populate the connection pool, and ping each one to verify it is alive
	AGRESSIVE
)

//
// Connection Pool wrapper
//
type RedisConnectionPool struct {
	Mode   ConnectionMode            "How should we prepare the connection pool?"
	Size   int                       "(Max) Pool size"
	Urls   []string                  "Redis URLs to connect to"
	Logger log4go.Logger             "Logger we are using in the connection pool"
	myPool *cp.ConnectionPoolWrapper "Connection Pool wrapper"
}

//
// Open the connection pool
//
func (p *RedisConnectionPool) Open() error {
	Close()

	// Lambda to iterate the urls
	nextUrl := findNextUrl(p.Urls)

	// Lambda for creating the factories
	var initfn cp.InitFunction
	switch p.Mode {
	// Create the factory
	// DON'T Connect to Redis
	// DON'T Test the connection
	case LAZY:
		initfn = func() (interface{}, error) {
			return newLazyFactory(nextUrl(), p.Logger)
		}
		// Create the factory
		// AND Connect to Redis
		// AND Test the connection
	case AGRESSIVE:
		initfn = func() (interface{}, error) {
			return newAgressiveFactory(nextUrl(), p.Logger)
		}
		// No mode specified!
	default:
		return errors.New(fmt.Sprintf("Invalid connection mode: %v", p.Mode))
	}

	// Create the new pool
	pool, err := cp.NewPool(p.Size, initfn)

	// Error creating the pool?
	if nil != err {
		return nil, err
	}

	// Save the pointer to the pool
	p.myPool = pool

	// Return nil
	return nil
}

//
// Helper to iterate urls
//
func findNextUrl(urls []string) func() string {
	i := 0
	return func() string {
		url := p.Urls[i%len(p.Urls)]
		i++
		return url
	}
}

//
// Close the connection pool
//
func (p *RedisConnectionPool) Close() error {
	// If the pool is not nil,
	// Then close all the connections and release the pointer
	if nil != p.myPool {
		for i := 0; i < p.Size; i++ {
			// Pop a connection from the pool
			c := Pop()

			// Skip nils
			if nil == c {
				continue
			}

			// Close the connection
			f := c.(*RedisConnectionFactory)
			f.Close()
		}
	}

	// Release the connection pool
	p.myPool = nil
}

//
// Get a RedisConnectionFactory from the pool
//
func (p *RedisConnectionPool) Pop() (*RedisConnectionFactory, error) {
	// Pop a connection from the pool
	c := p.myPool.GetConnection()

	// Return the connection
	if c != nil {
		return c.(*RedisConnectionFactory)
	}

	// Return an error when all connections are exhausted
	return errors.New("No RedisConnectionFactory availble")
}

//
// Return a RedisConnectionFactory
//
func (p *RedisConnectionPool) Push(connection *RedisConnectionFactory) {
	c.myPool.ReleaseConnection(c)
}
