package tasks

type TaskRunner interface {
    GetType() string
    Fetch(params TaskParams) ([]ABook, error)
}
