package model

// having some basic functions for result and results

// Returns true if the results are true
func (r *Result) IsValid() bool {
	return r != nil && len(*r) > 0
}

// Retuens true if the results not empty
func (r *Results) IsEmpty() bool {
	return r != nil && len(*r) > 0
}
