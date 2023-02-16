package checker

import "errors"

var (
	ZeroCookiesError = errors.New("cannot extract cookies from txt")
	SaveError        = errors.New("error while saving valid")

	BaseErrpr = errors.New("something went wrong")
)
