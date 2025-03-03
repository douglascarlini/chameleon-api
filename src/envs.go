package main

import "os"

var MDB_HOST = os.Getenv("MDB_HOST")
var MDB_PORT = os.Getenv("MDB_PORT")
var MDB_NAME = os.Getenv("MDB_NAME")
var MDB_USER = os.Getenv("MDB_USER")
var MDB_PASS = os.Getenv("MDB_PASS")

var RMQ_USER = os.Getenv("RMQ_USER")
var RMQ_PASS = os.Getenv("RMQ_PASS")
