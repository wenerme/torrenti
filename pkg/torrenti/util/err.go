package util

func FirstErrorFunc(f ...func() error) error {
	for _, v := range f {
		if err := v(); err != nil {
			return err
		}
	}
	return nil
}
