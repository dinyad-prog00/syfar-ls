package fs

import (
	"errors"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"
)

func uriToPath(uri string) (string, error) {
	s := strings.ReplaceAll(uri, "%5C", "/")
	parsed, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	if parsed.Scheme != "file" {
		return "", errors.New("URI was not a file:// URI")
	}

	if runtime.GOOS == "windows" {
		// In Windows "file:///c:/tmp/foo.md" is parsed to "/c:/tmp/foo.md".
		// Strip the first character to get a valid path.
		if strings.Contains(parsed.Path[1:], ":") {
			// url.Parse() behaves differently with "file:///c:/..." and "file://c:/..."
			return parsed.Path[1:], nil
		} else {
			// if the windows drive is not included in Path it will be in Host
			return parsed.Host + "/" + parsed.Path[1:], nil
		}
	}
	return parsed.Path, nil
}

func canonical(path string) (string, error) {
	path = filepath.Clean(path)

	resolvedPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path, err
	} else {
		path = resolvedPath
	}

	return path, nil
}

func normalizePath(pathOrUri string) (string, error) {
	path, err := uriToPath(pathOrUri)
	if err != nil {
		return "", err
	}
	return canonical(path)
}
