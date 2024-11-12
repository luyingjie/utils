package redis

import (
	"errors"
	"reflect"
	"time"

	"github.com/luyingjie/utils/pkg/container/vvar"
	"github.com/luyingjie/utils/pkg/conv"
	"github.com/luyingjie/utils/util/json"

	"github.com/gomodule/redigo/redis"
)

// Do sends a command to the server and returns the received reply.
// It uses json.Marshal for struct/slice/map type values before committing them to redis.
// The timeout overrides the read timeout set when dialing the connection.
func (c *Conn) do(timeout time.Duration, commandName string, args ...interface{}) (reply interface{}, err error) {
	var (
		reflectValue reflect.Value
		reflectKind  reflect.Kind
	)
	for k, v := range args {
		reflectValue = reflect.ValueOf(v)
		reflectKind = reflectValue.Kind()
		if reflectKind == reflect.Ptr {
			reflectValue = reflectValue.Elem()
			reflectKind = reflectValue.Kind()
		}
		switch reflectKind {
		case
			reflect.Struct,
			reflect.Map,
			reflect.Slice,
			reflect.Array:
			// Ignore slice type of: []byte.
			if _, ok := v.([]byte); !ok {
				if args[k], err = json.Marshal(v); err != nil {
					return nil, err
				}
			}
		}
	}
	if timeout > 0 {
		conn, ok := c.Conn.(redis.ConnWithTimeout)
		if !ok {
			return vvar.New(nil), errors.New(`current connection does not support "ConnWithTimeout"`)
		}
		return conn.DoWithTimeout(timeout, commandName, args...)
	}
	return c.Conn.Do(commandName, args...)
}

// Do sends a command to the server and returns the received reply.
// It uses json.Marshal for struct/slice/map type values before committing them to redis.
func (c *Conn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	return c.do(0, commandName, args...)
}

// DoWithTimeout sends a command to the server and returns the received reply.
// The timeout overrides the read timeout set when dialing the connection.
func (c *Conn) DoWithTimeout(timeout time.Duration, commandName string, args ...interface{}) (reply interface{}, err error) {
	return c.do(timeout, commandName, args...)
}

// DoVar retrieves and returns the result from command as vvar.Var.
func (c *Conn) DoVar(commandName string, args ...interface{}) (*vvar.Var, error) {
	return resultToVar(c.Do(commandName, args...))
}

// DoVarWithTimeout retrieves and returns the result from command as vvar.Var.
// The timeout overrides the read timeout set when dialing the connection.
func (c *Conn) DoVarWithTimeout(timeout time.Duration, commandName string, args ...interface{}) (*vvar.Var, error) {
	return resultToVar(c.DoWithTimeout(timeout, commandName, args...))
}

// ReceiveVar receives a single reply as vvar.Var from the Redis server.
func (c *Conn) ReceiveVar() (*vvar.Var, error) {
	return resultToVar(c.Receive())
}

// ReceiveVarWithTimeout receives a single reply as vvar.Var from the Redis server.
// The timeout overrides the read timeout set when dialing the connection.
func (c *Conn) ReceiveVarWithTimeout(timeout time.Duration) (*vvar.Var, error) {
	conn, ok := c.Conn.(redis.ConnWithTimeout)
	if !ok {
		return vvar.New(nil), errors.New(`current connection does not support "ConnWithTimeout"`)
	}
	return resultToVar(conn.ReceiveWithTimeout(timeout))
}

// resultToVar converts redis operation result to vvar.Var.
func resultToVar(result interface{}, err error) (*vvar.Var, error) {
	if err == nil {
		if result, ok := result.([]byte); ok {
			return vvar.New(conv.UnsafeBytesToStr(result)), err
		}
		// It treats all returned slice as string slice.
		if result, ok := result.([]interface{}); ok {
			return vvar.New(conv.Strings(result)), err
		}
	}
	return vvar.New(result), err
}
