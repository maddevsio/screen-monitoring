package dashboard

import "github.com/stretchr/testify/mock"

type MigratorMock struct {
	mock.Mock
}

func (o *MigratorMock) Up() ([]error, bool) {
	args := o.Called()
	return args.Get(0).([]error), args.Bool(1)
}

func (o *MigratorMock) Down() ([]error, bool) {
	args := o.Called()
	return args.Get(0).([]error), args.Bool(1)
}
