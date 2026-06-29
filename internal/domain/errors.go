package domain

import "errors"

var ErrNotFound = errors.New("ShortUrl not found")
var ErrInvalidUrl = errors.New("ShortUrl is invalid")
var ErrInvalidJson = errors.New("JSON is invalid")
