package domain

import "net/url"

func ValidateURL(input string) (bool, error) {
	_, err := url.ParseRequestURI(input)
	return err == nil, err
}
