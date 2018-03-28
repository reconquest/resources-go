package resources

type ErrorUndefinedResource struct {
	key string
}

func (err ErrorUndefinedResource) Error() string {
	return "undefined resource: " + err.key
}
