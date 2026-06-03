/*
All Rights Reversed (ɔ)
*/

package engine

type ExecuteSummary struct {
	Actions          []ActionEvent
	OperationSummary *Summary
	Label            string
}

type Summary struct {
	Stowed   int
	Unstowed int
	Replaced int
	Backed   int
	Adopted  int
	Skipped  int
	UpToDate int
}
