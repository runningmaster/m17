package option

// Fn is a function on the options for a package.
type Fn func(interface{}) error

// Receive receives options for a package from athoter package.
func Receive(receiver interface{}, options ...Fn) error {
	var err error
	for i := range options {
		err = options[i](receiver)
		if err != nil {
			return err
		}
	}
	return nil
}
