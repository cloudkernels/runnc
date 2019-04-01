package profile

import (
	"fmt"
	"os"
	"time"
)

type Profile struct {
	// Name of the Profile
	name string

	// Currently measuring
	on bool

	// The last timestamp we took
	last time.Time

	// Duration of this profile so far
	elapsed time.Duration
}

var profiles []*Profile
var logFile *os.File

func NewProfile(name string) *Profile {
	ret := &Profile{name: name}
	profiles = append(profiles, ret)

	str := fmt.Sprintf("Creating new profile: '%s'\n", name)
	logFile.WriteString(str)

	return ret
}

func (p *Profile) Start() {
	if p.on {
		str := fmt.Sprintf("Profile %s already started\n", p.name)
		logFile.WriteString(str)
		return
	}

	p.last = time.Now()
	p.on = true
}

func (p *Profile) Stop() {
	if !p.on {
		str := fmt.Sprintf("Profile %s not started\n", p.name)
		logFile.WriteString(str)
		return
	}

	p.elapsed += time.Now().Sub(p.last)
	p.on = false
}

func (p *Profile) String() string {
	return fmt.Sprintf("[%s]: %v", p.name, p.elapsed)
}

func Init(profFileName string) error {
	f, err := os.OpenFile(profFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	logFile = f

	return nil
}

func WriteProfiles() {
	for _, p := range profiles {
		if p.on {
			p.Stop()
		}
		logFile.WriteString(p.String() + "\n")
	}
	logFile.Close()
}
