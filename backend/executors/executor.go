package executors

import "github.com/pocketbase/pocketbase"

type MessageExecutorParameters struct {
	Message string
	Script  string
	PB      *pocketbase.PocketBase
	OwnerID string
}

type MessageExecutor interface {
	Execute(mep *MessageExecutorParameters) string
}

var executors = make(map[string]MessageExecutor)

func RegisterExecutor(name string, executor MessageExecutor) {
	executors[name] = executor
}

func GetExecutor(name string) MessageExecutor {
	return executors[name]
}
