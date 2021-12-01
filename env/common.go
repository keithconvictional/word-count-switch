package env

import "os"

func DoExtract() bool {
	return os.Getenv("EXTRACT") == "true"
}

func ExtractMethod() string {
	return os.Getenv("EXTRACT_METHOD")
}

func DoTransform() bool {
	return os.Getenv("TRANSFORM") == "true"
}

func DoLoad() bool {
	return os.Getenv("LOAD") == "true"
}

func LoadMethod() string {
	return os.Getenv("LOAD_METHOD")
}

func Debug() bool {
	return os.Getenv("DEBUG") == "true"
}

