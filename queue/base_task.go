package queue

type BaseTask interface {
	Execute() error
}
