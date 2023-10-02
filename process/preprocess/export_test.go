package preprocess

// NewBaseEventsPreProcessor -
func NewBaseEventsPreProcessor(args ArgsEventsPreProcessor) (*baseEventsPreProcessor, error) {
	return newBaseEventsPreProcessor(args)
}
