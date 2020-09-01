package Database

import (
	"TianTing/Logger"
	"TianTing/Settings"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"io"
	"strings"
	"time"
)

type CacheDB struct {
	IDatabase
	clientPool *redis.Pool
}

// 這樣設計的好處就是只要抽換這裡以及 Config 就可以換資料庫了

func ConnectWithCacheDB(config *Settings.CacheDbConf) (*CacheDB, error) {

	Logger.SysLog.Info("[CacheDatabase] Connecting to Cache Service")

	client := &redis.Pool{
		MaxIdle:     config.MaxIdle,
		MaxActive:   config.MaxActive,
		IdleTimeout: time.Duration(config.IdleTimeout) * time.Second,
		Wait:        config.Wait,
		Dial: func() (redis.Conn, error) {
			Logger.SysLog.Debug("[CacheDatabase] Dial Connects To The Cache Server")
			var dialOpt []redis.DialOption
			if config.Password != "" {
				dialOpt = append(dialOpt, redis.DialPassword(config.Password))
			}
			dialOpt = append(dialOpt, redis.DialDatabase(config.Database))
			connectString := fmt.Sprintf("%s:%d", config.Host, config.Port)
			c, err := redis.Dial("tcp", connectString, dialOpt...)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	CacheDbClient := &CacheDB{
		clientPool: client,
	}

	// Checking cache can be connected
	err := CacheDbClient.pingHealth()
	if err != nil {
		return nil, err
	}

	return CacheDbClient, nil
}

func (r *CacheDB) pingHealth() error {
	_, err := r.redo("PING")
	if err != nil {
		return err
	}
	Logger.SysLog.Info("[CacheDatabase] Connected to the Redis Database")
	return nil
}

func (r *CacheDB) GetClient() redis.Conn {
	return r.clientPool.Get()
}

func (r *CacheDB) SetString(key string, data string, time int) error {
	_, err := r.redo("SET", key, data)
	if err != nil {
		return err
	}
	if time > 0 {
		_, err = r.redo("EXPIRE", key, time)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *CacheDB) Exists(key string) bool {
	exists, err := redis.Bool(r.redo("EXISTS", key))
	if err != nil {
		return false
	}
	return exists
}

func (r *CacheDB) Get(key string) ([]byte, error) {
	reply, err := redis.Bytes(r.redo("GET", key))
	if err != nil {
		return nil, err
	}
	return reply, nil
}

func (r *CacheDB) Delete(key string) (bool, error) {
	return redis.Bool(r.redo("DEL", key))
}

func (r *CacheDB) LikeDeletes(key string) error {
	keys, err := redis.Strings(r.redo("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}
	for _, key := range keys {
		_, err = r.Delete(key)
		if err != nil {
			return err
		}
	}
	return nil
}

func isConnError(err error) bool {
	var needNewConn bool

	if err == nil {
		return false
	}

	if err == io.EOF {
		needNewConn = true
	}
	if strings.Contains(err.Error(), "use of closed network connection") {
		needNewConn = true
	}
	if strings.Contains(err.Error(), "connect: connection refused") {
		needNewConn = true
	}
	return needNewConn
}

func (r *CacheDB) redo(command string, opt ...interface{}) (interface{}, error) {
	client := r.clientPool.Get()
	defer client.Close()

	var conn redis.Conn
	var err error
	var maxRetry = 3
	var needNewConn bool

	resp, err := client.Do(command, opt...)
	needNewConn = isConnError(err)
	if needNewConn == false {
		return resp, err
	} else {
		conn, err = r.clientPool.Dial()
	}

	for index := 0; index < maxRetry; index++ {
		if conn == nil && index+1 > maxRetry {
			return resp, err
		}
		if conn == nil {
			conn, err = r.clientPool.Dial()
		}
		if err != nil {
			continue
		}

		resp, err := conn.Do(command, opt...)
		needNewConn = isConnError(err)
		if needNewConn == false {
			return resp, err
		} else {
			conn, err = r.clientPool.Dial()
		}
	}
	conn.Close()
	return "", errors.New("redis error")
}
