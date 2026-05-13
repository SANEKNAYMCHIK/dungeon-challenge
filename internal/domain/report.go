package domain

type ReportHeader string

const (
	ReportHeaderSuccess ReportHeader = "SUCCESS"
	ReportHeaderFail    ReportHeader = "FAIL"
	ReportHeaderDisqual ReportHeader = "DISQUAL"
)
