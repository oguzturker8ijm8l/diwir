package directory

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/containers/image/types"
	"github.com/docker/docker/reference"
)

type dirImageSource struct {
	ref dirReference
}

// newImageSource returns an ImageSource reading from an existing directory.
func newImageSource(ref dirReference) types.ImageSource {
	return &dirImageSource{ref}
}

// IntendedDockerReference returns the Docker reference for this image, _as specified by the user_
// (not as the image itself, or its underlying storage, claims).  Should be fully expanded, i.e. !reference.IsNameOnly.
// This can be used e.g. to determine which public keys are trusted for this image.
// May be nil if unknown.
func (s *dirImageSource) IntendedDockerReference() reference.Named {
	return nil
}

// it's up to the caller to determine the MIME type of the returned manifest's bytes
func (s *dirImageSource) GetManifest(_ []string) ([]byte, string, error) {
	m, err := ioutil.ReadFile(manifestPath(s.ref.path))
	if err != nil {
		return nil, "", err
	}
	return m, "", err
}

func (s *dirImageSource) GetBlob(digest string) (io.ReadCloser, int64, error) {
	r, err := os.Open(layerPath(s.ref.path, digest))
	if err != nil {
		return nil, 0, nil
	}
	fi, err := os.Stat(layerPath(s.ref.path, digest))
	if err != nil {
		return nil, 0, nil
	}
	return r, fi.Size(), nil
}

func (s *dirImageSource) GetSignatures() ([][]byte, error) {
	signatures := [][]byte{}
	for i := 0; ; i++ {
		signature, err := ioutil.ReadFile(signaturePath(s.ref.path, i))
		if err != nil {
			if os.IsNotExist(err) {
				break
			}
			return nil, err
		}
		signatures = append(signatures, signature)
	}
	return signatures, nil
}

func (s *dirImageSource) Delete() error {
	return fmt.Errorf("directory#dirImageSource.Delete() not implmented")
}
