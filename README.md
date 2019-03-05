# structwalk

`structwalk` is a go package to recursively visit fields in a (nested) struct. You supply a function that is called in every visited field. It returns a bool that indicates if the field should be recursed into and an error that indicates that recursion should be halted. A couple of examples are provided with the tests.
