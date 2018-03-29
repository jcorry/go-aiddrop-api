package controllers

import (
	"database/sql"

	"github.com/coopernurse/gorp"
	"github.com/revel/revel"
)

var (
	// Dbm is the DB object
	Dbm *gorp.DbMap
)

// GorpController is the controller that resource controllers
// extend from
type GorpController struct {
	*revel.Controller
	Txn *gorp.Transaction
}

// Begin starts a transaction
func (c *GorpController) Begin() revel.Result {
	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

// Commit commits a transaction
func (c *GorpController) Commit() revel.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

// Rollback rolls back a transation
func (c *GorpController) Rollback() revel.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}
