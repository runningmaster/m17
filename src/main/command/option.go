package command

type logger interface {
	Printf(string, ...interface{})
}

type Option struct {
	log logger
}

var defaultOption = &Option{}

func (o *Option) setLogger(l logger) error {
	o.log = l
	return nil
}

func (o *Option) override(options ...func(*Option) error) error {
	var err error
	for i := range options {
		err = options[i](o)
		if err != nil {
			return err
		}
	}
	return nil
}

func Logger(l logger) func(*Option) error {
	return func(o *Option) error {
		return o.setLogger(l)
	}
}
