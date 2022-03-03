package main

import (
	"context"
	"errors"
	"fmt"
	"log"
)

type StorageService interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Put(ctx context.Context, key string, val interface{}) error
	Delete(ctx context.Context, key string) error
}

type StorageChannel struct {
	commands chan Command
}

type CommandType int

const (
	Get = iota
	Put
	Delete
)

type Result struct {
	val interface{}
	err error
}

type Command struct {
	ctx   context.Context
	ty    CommandType
	key   string
	val   interface{}
	reply chan Result
}

func (channel *StorageChannel) Get(ctx context.Context, key string) (interface{}, error) {
	replyChan := make(chan Result)
	channel.commands <- Command{ctx: ctx, ty: Get, key: key, reply: replyChan}

	result := <-replyChan
	return result.val, result.err
}

func (channel *StorageChannel) Put(ctx context.Context, key string, val interface{}) error {
	replyChan := make(chan Result)

	channel.commands <- Command{ctx: ctx, ty: Put, key: key, val: val, reply: replyChan}

	result := <-replyChan
	return result.err
}

func (channel *StorageChannel) Delete(ctx context.Context, key string) error {
	replyChan := make(chan Result)
	channel.commands <- Command{ctx: ctx, ty: Delete, key: key, reply: replyChan}

	result := <-replyChan
	return result.err
}

type KVStorage interface {
	Get(ctx context.Context, key string) (interface{}, error)
	Put(ctx context.Context, key string, val interface{}) error
	Delete(ctx context.Context, key string) error
}

type Storage struct {
	Repository map[string]interface{}
}

func startStorageManager() StorageService {
	storage := &Storage{make(map[string]interface{})}

	storageChannel := &StorageChannel{make(chan Command)}
	commands := storageChannel.commands
	go func() {
		for command := range commands {
			switch command.ty {
			case Get:
				if val, err := storage.Get(command.ctx, command.key); err == nil {
					command.reply <- Result{val, err}
				} else {
					fmt.Println(err)
					command.reply <- Result{val, err}
				}
			case Put:
				if err := storage.Put(command.ctx, command.key, command.val); err == nil {
					command.reply <- Result{nil, err}
				} else {
					fmt.Println(err)
					command.reply <- Result{nil, err}
				}
			case Delete:
				if err := storage.Delete(command.ctx, command.key); err == nil {
					command.reply <- Result{nil, err}
				} else {
					fmt.Println(err)
					command.reply <- Result{nil, err}
				}
			default:
				log.Fatal("command not supported")
			}
		}
	}()
	var service StorageService = storageChannel
	return service
}

func (storage *Storage) Get(ctx context.Context, key string) (interface{}, error) {
	if val, ok := storage.Repository[key]; ok {
		return val, nil
	} else {
		return nil, errors.New("value not found")
	}

}

func (storage *Storage) Put(ctx context.Context, key string, val interface{}) error {
	storage.Repository[key] = val
	return nil
}

func (Storage *Storage) Delete(ctx context.Context, key string) error {
	delete(Storage.Repository, key)
	return nil
}
