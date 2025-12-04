package preprocess

// NewBaseEventsPreProcessor -
func NewBaseEventsPreProcessor(args ArgsEventsPreProcessor) (*baseEventsPreProcessor, error) {
	return newBaseEventsPreProcessor(args)
}

// CreateEmptyBlockCreatorContainer -
func CreateEmptyBlockCreatorContainer() (EmptyBlockCreatorContainer, error) {
	return createEmptyBlockCreatorContainer()
}
