package option

type Receiver interface {
	Receive(...Fn) error
}

// Fn is a function for setting option into package.
type Fn func(Receiver) error

// Receive receives options for a package from athoter package.
func Receive(r Receiver, options ...Fn) error {
	var err error
	for i := range options {
		err = options[i](r)
		if err != nil {
			return err
		}
	}
	return nil
}
