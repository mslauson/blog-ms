package iamConst

import "fmt"

var (
	NoCustomersFound = fmt.Errorf("no customers exist")
	NoCustomerFound  = fmt.Errorf("no customer exists with the given information")
)
