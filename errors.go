// Copyright (c) 2012-present The upper.io/db authors. All rights reserved.
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package db

import (
	"errors"
)

// Error messages
var (
	ErrMissingAdapter           = errors.New(`db: missing adapter`)
	ErrAlreadyWithinTransaction = errors.New(`db: already within a transaction`)
	ErrCollectionDoesNotExist   = errors.New(`db: collection does not exist`)
	ErrExpectingNonNilModel     = errors.New(`db: expecting non nil model`)
	ErrExpectingPointerToStruct = errors.New(`db: expecting pointer to struct`)
	ErrGivingUpTryingToConnect  = errors.New(`db: giving up trying to connect: too many clients`)
	ErrInvalidCollection        = errors.New(`db: invalid collection`)
	ErrMissingCollectionName    = errors.New(`db: missing collection name`)
	ErrMissingConditions        = errors.New(`db: missing selector conditions`)
	ErrMissingConnURL           = errors.New(`db: missing DSN`)
	ErrMissingDatabaseName      = errors.New(`db: missing database name`)
	ErrNoMoreRows               = errors.New(`db: no more rows in this result set`)
	ErrNotConnected             = errors.New(`db: not connected to a database`)
	ErrNotImplemented           = errors.New(`db: call not implemented`)
	ErrQueryIsPending           = errors.New(`db: can't execute this instruction while the result set is still open`)
	ErrQueryLimitParam          = errors.New(`db: a query can accept only one limit parameter`)
	ErrQueryOffsetParam         = errors.New(`db: a query can accept only one offset parameter`)
	ErrQuerySortParam           = errors.New(`db: a query can accept only one order-by parameter`)
	ErrSockerOrHost             = errors.New(`db: you may connect either to a UNIX socket or a TCP address, but not both`)
	ErrTooManyClients           = errors.New(`db: can't connect to database server: too many clients`)
	ErrUndefined                = errors.New(`db: value is undefined`)
	ErrUnknownConditionType     = errors.New(`db: arguments of type %T can't be used as constraints`)
	ErrUnsupported              = errors.New(`db: action is not supported by the DBMS`)
	ErrUnsupportedDestination   = errors.New(`db: unsupported destination type`)
	ErrUnsupportedType          = errors.New(`db: type does not support marshaling`)
	ErrUnsupportedValue         = errors.New(`db: value does not support unmarshaling`)
	ErrNilRecord                = errors.New(`db: invalid item (nil)`)
	ErrRecordIDIsZero           = errors.New(`db: item ID is not defined`)
	ErrMissingPrimaryKeys       = errors.New(`db: collection %q has no primary keys`)
	ErrWarnSlowQuery            = errors.New(`db: slow query`)
	ErrTransactionAborted       = errors.New(`db: transaction was aborted`)
	ErrNotWithinTransaction     = errors.New(`db: not within transaction`)
	ErrNotSupportedByAdapter    = errors.New(`db: not supported by adapter`)
)
