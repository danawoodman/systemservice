package systemservice

import "fmt"

/*
ServiceDoesNotExistError is an error return if a given service does
not exist on the system. This is usually returned if the user
attempts to manage a service not yet configured on the system.
*/
type ServiceDoesNotExistError struct {
	serviceName string
}

/*
Error implements the errors.Error interface
*/
func (e *ServiceDoesNotExistError) Error() string {
	return fmt.Sprintf("the service \"%s\" does not exist", e.serviceName)
}
