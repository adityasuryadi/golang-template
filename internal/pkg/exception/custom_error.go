package exception

const (
	ERRBADREQUEST Status = 1
	ERRSERVER     Status = 2
	ERRUNKNOWN    Status = 3
	ERRAUTHORIZED Status = 4
	ERRFORBIDDEN  Status = 5
)

type Status int

type CustomError struct {
	Status Status
	Error  error
}
