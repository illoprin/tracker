package storage

// getContentTypeByExtension returns MIME type by file extension
func GetAudioContentTypeByExtension(ext string) string {
	switch ext {
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".flac":
		return "audio/flac"
	case ".ogg":
		return "audio/ogg"
	case ".aac":
		return "audio/aac"
	case ".m4a":
		return "audio/mp4"
	default:
		return "audio/mpeg"
	}
}
