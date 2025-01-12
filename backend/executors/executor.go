package executors

type MessageExecutor interface {
	Execute(message string, script string) string
}

var executors = make(map[string]MessageExecutor)

func RegisterExecutor(name string, executor MessageExecutor) {
	executors[name] = executor
}

func GetExecutor(name string) MessageExecutor {
	return executors[name]
}
