package runtime

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"math/rand"
	rs "runtime"
	"time"
)

var (
	name     = "kubor"
	version  = "development"
	revision = ""
	built    = ""

	Runtime Holder
)

type Holder struct {
	Name      string
	Version   string
	Revision  string
	Built     time.Time
	GoVersion string
	GOOS      string
	GOARCH    string
}

func (instance Holder) String() string {
	return fmt.Sprintf(`%s
 Version:      %s
 Git revision: %s
 Built:        %v
 Go version:   %s
 OS/Arch:      %s/%s`,
		instance.Name, instance.Version, instance.Revision, instance.Built, instance.GoVersion, instance.GOOS, instance.GOARCH)
}

func (instance Holder) ShortString() string {
	return fmt.Sprintf(`%s (version: %s, revision: %s)`,
		instance.Name, instance.Version, instance.Revision)
}

func (instance Holder) LongVersion() string {
	return fmt.Sprintf(`%s (revision: %s)`,
		instance.Version, instance.Revision)
}

func init() {
	if built == "" {
		built = time.Now().Format(time.RFC3339)
	}
	var err error
	if Runtime.Built, err = time.Parse(time.RFC3339, built); err != nil {
		panic(fmt.Sprintf("illegal built value '%s': %v", built, err))
	}
	if revision == "" {
		revision = RandomRevision(Runtime.Built)
	}
	Runtime.Name = name
	Runtime.Version = version
	Runtime.Revision = revision
	Runtime.GoVersion = rs.Version()
	Runtime.GOOS = rs.GOOS
	Runtime.GOARCH = rs.GOARCH
}

func RandomRevision(baseOn time.Time) string {
	b := make([]byte, sha1.Size)
	rng := rand.New(rand.NewSource(baseOn.UnixNano()))
	if n, err := rng.Read(b); err != nil {
		panic(err)
	} else if n < len(b) {
		panic(io.ErrShortBuffer)
	}
	return hex.EncodeToString(b)
}
