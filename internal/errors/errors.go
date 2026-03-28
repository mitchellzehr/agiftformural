package errors

import stderrors "errors"

// ErrNotFound is returned by repositories when a requested row does not exist.
var ErrNotFound = stderrors.New("not found")
