/*
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/

package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/syndtr/goleveldb/leveldb/iterator"

	"github.com/syndtr/goleveldb/leveldb/util"

	"github.com/syndtr/goleveldb/leveldb/opt"

	"chainmaker.org/chainmaker/tcip-bcos/v2/module/conf"
	"chainmaker.org/chainmaker/tcip-bcos/v2/module/logger"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/zap"
)

// DbHandle 数据库操作对象
type DbHandle struct {
	db  *leveldb.DB
	log *zap.SugaredLogger
}

// Db 数据库全局对象
var Db *DbHandle

// NewDbHandle 初始化数据库
func NewDbHandle() {
	err := createDirIfNotExist(conf.Config.DbPath)
	if err != nil {
		panic(fmt.Sprintf("Error create dir %s by DbHandler: %s", conf.Config.DbPath, err))
	}
	db, err := leveldb.OpenFile(conf.Config.DbPath, &opt.Options{})
	if err != nil {
		panic(err.Error())
	}
	Db = &DbHandle{
		log: logger.GetLogger(logger.ModuleDb),
		db:  db,
	}
}

// Get returns the value for the given key, or returns nil if none exists
//  @receiver d
//  @param key
//  @return []byte
//  @return error
func (d *DbHandle) Get(key []byte) ([]byte, error) {
	value, err := d.db.Get(key, nil)
	if err == leveldb.ErrNotFound {
		value = nil
		err = nil
	}
	if err != nil {
		msg := fmt.Sprintf("[Get] getting DbHandle key [%#v], err:%s", key, err.Error())
		d.log.Errorf(msg)
		return nil, errors.New(msg)
	}
	d.log.Debug("[Get] key: %s, value: %s", string(key), string(value))
	return value, nil
}

// Put saves the key-values
//  @receiver d
//  @param key
//  @param value
//  @return error
func (d *DbHandle) Put(key []byte, value []byte) error {
	msg := fmt.Sprintf("[Put] writing leveldbprovider key [%#v] with nil value", key)
	if value == nil {
		d.log.Warn(msg)
		return errors.New(msg)
	}

	err := d.db.Put(key, value, &opt.WriteOptions{Sync: true})
	if err != nil {
		d.log.Errorf(msg)
		return errors.New(msg)
	}
	d.log.Debugf("[Put] key: %s, value: %s", string(key), string(value))
	return err
}

// Has return true if the given key exist, or return false if none exists
//  @receiver d
//  @param key
//  @return bool
//  @return error
func (d *DbHandle) Has(key []byte) (bool, error) {
	exist, err := d.db.Has(key, nil)
	if err != nil {
		d.log.Errorf("getting leveldbprovider key [%#v], err:%s", key, err.Error())
		return false, errors.New("error getting leveldbprovider key [%#v]")
	}
	return exist, nil
}

// Delete Delete deletes the given key
//  @receiver d
//  @param key
//  @return error
func (d *DbHandle) Delete(key []byte) error {
	wo := &opt.WriteOptions{Sync: true}
	err := d.db.Delete(key, wo)
	if err != nil {
		d.log.Errorf("deleting leveldbprovider key [%#v]", key)
		return errors.New("error deleting leveldbprovider key [%#v]")
	}
	return err
}

// NewIteratorWithRange returns an iterator that contains all the key-values between given key
//	 					ranges start is included in the results and limit is excluded.
//  @receiver d
//  @param startKey
//  @param limitKey
//  @return iterator.Iterator
//  @return error
func (d *DbHandle) NewIteratorWithRange(startKey []byte, limitKey []byte) (iterator.Iterator, error) {
	if len(startKey) == 0 || len(limitKey) == 0 {
		return nil, fmt.Errorf("iterator range should not start(%s) or limit(%s) with empty key",
			string(startKey), string(limitKey))
	}
	keyRange := &util.Range{Start: startKey, Limit: limitKey}
	iter := d.db.NewIterator(keyRange, nil)
	return iter, nil
}

// Close closes the leveldb
//  @receiver d
//  @return error
func (d *DbHandle) Close() error {
	return d.db.Close()
}

// createDirIfNotExist 创建文件夹
//  @param path
//  @return error
func createDirIfNotExist(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		// 创建文件夹
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
