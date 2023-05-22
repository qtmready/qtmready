package shared

func StackWorkflowID(id string) string {
	return Temporal().Queue(CoreQueue).CreateWorkflowID("stack", id)
}
