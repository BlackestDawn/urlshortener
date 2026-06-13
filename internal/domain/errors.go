package domain

import "errors"

var ErrNotFound = errors.New("ShortUrl not found")
var ErrInvalidUrl = errors.New("ShortUrl is invalid")
