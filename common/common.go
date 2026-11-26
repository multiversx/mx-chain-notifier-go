package common

// ConvertTimeStampMsToSec will convert unix timestamp from milliseconds to seconds
func ConvertTimeStampMsToSec(timeStamp uint64) uint64 {
	return timeStamp / 1000
}

// DefaultMaxPayloadSizeInBytes is the max payload size used when none is configured. It matches
// the default rabbitMQ max_message_size (16 MiB), so that the broker does not reject messages.
const DefaultMaxPayloadSizeInBytes uint64 = 16 * 1024 * 1024

// GetMaxPayloadSize returns the configured max payload size, or the default one if not configured
func GetMaxPayloadSize(configuredMaxPayloadSize uint64) uint64 {
	if configuredMaxPayloadSize == 0 {
		return DefaultMaxPayloadSizeInBytes
	}

	return configuredMaxPayloadSize
}
