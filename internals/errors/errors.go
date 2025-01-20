package dmerrors

import "errors"

type Error struct {
	libErr error
	appErr error
}

func (e Error) LibError() error {
	return e.libErr
}

func (e Error) AppError() error {
	return e.appErr
}

func DMError(apperror, liberr error) error {
	return Error{
		libErr: liberr,
		appErr: apperror,
	}
}

func (e Error) Error() string {
	return errors.Join(e.libErr, e.appErr).Error()
}

// func DMErrorChain(err)
