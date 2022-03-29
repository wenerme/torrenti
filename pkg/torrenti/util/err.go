package util

import "go.uber.org/multierr"

func FirstErrorFunc(f ...func() error) error {
	for _, v := range f {
		if err := v(); err != nil {
			return err
		}
	}
	return nil
}

func CombineErrorFunc(f ...func() error) error {
	var errs []error
	for _, v := range f {
		if err := v(); err != nil {
			errs = append(errs, err)
		}
	}
	return multierr.Combine(errs...)
}
