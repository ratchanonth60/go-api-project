package service

import "errors"

func wrapError(baseErr, err error) error {
	if err == nil {
		return baseErr // If original error is nil, return base error
	}
	return errors.Join(baseErr, err) // Use errors.Join for Go 1.20+
}
