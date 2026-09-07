package common

// ConvertTimeStampMsToSec will convert unix timestamp from milliseconds to seconds
func ConvertTimeStampMsToSec(timeStamp uint64) uint64 {
	return timeStamp / 1000
}
