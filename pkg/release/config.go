package release

import (
	"errors"
	"github.com/helmwave/helmwave/pkg/pubsub"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/storage/driver"
	"time"
)

type Config struct {
	Chart           chart
	Name            string
	Namespace       string
	CreateNamespace bool

	Devel                    bool
	dryRun                   bool
	SkipCRDs                 bool
	Timeout                  time.Duration
	Wait                     bool
	WaitForJobs              bool
	DisableHooks             bool
	Force                    bool
	ResetValues              bool
	ReuseValues              bool
	Recreate                 bool
	MaxHistory               int
	Atomic                   bool
	CleanupOnFail            bool
	SubNotes                 bool
	Description              string
	DisableOpenAPIValidation bool

	// Helmwave
	Tags         []string
	Values       []string
	Store        map[string]interface{}
	DependsOn    []string `yaml:"depends_on"`
	dependencies map[string]<-chan pubsub.ReleaseStatus

	cfg *action.Configuration
}

func (rel *Config) DryRun(b bool) *Config {
	rel.dryRun = b
	return rel
}

type chart struct {
	Name                    string
	action.ChartPathOptions `yaml:",inline"`
}

func (rel *Config) newInstall() *action.Install {
	client := action.NewInstall(rel.cfg)

	// Only Install
	client.CreateNamespace = rel.CreateNamespace
	client.ReleaseName = rel.Name

	// Common Part
	client.DryRun = rel.dryRun
	client.Devel = rel.Devel
	client.Namespace = rel.Namespace
	client.ChartPathOptions = rel.Chart.ChartPathOptions
	client.DisableHooks = rel.DisableHooks
	client.SkipCRDs = rel.SkipCRDs
	client.Timeout = rel.Timeout
	client.Wait = rel.Wait
	client.Atomic = rel.Atomic
	client.DisableOpenAPIValidation = rel.DisableOpenAPIValidation
	client.SubNotes = rel.SubNotes
	client.Description = rel.Description

	if client.DryRun {
		client.Replace = true
	}

	return client
}

func (rel *Config) newUpgrade() *action.Upgrade {
	client := action.NewUpgrade(rel.cfg)
	// Only Upgrade
	client.CleanupOnFail = rel.CleanupOnFail
	client.MaxHistory = rel.MaxHistory
	client.Recreate = rel.Recreate
	client.ReuseValues = rel.ReuseValues
	client.ResetValues = rel.ReuseValues

	// Common Part
	client.DryRun = rel.dryRun
	client.Devel = rel.Devel
	client.Namespace = rel.Namespace
	client.ChartPathOptions = rel.Chart.ChartPathOptions
	client.DisableHooks = rel.DisableHooks
	client.SkipCRDs = rel.SkipCRDs
	client.Timeout = rel.Timeout
	client.Wait = rel.Wait
	client.Atomic = rel.Atomic
	client.DisableOpenAPIValidation = rel.DisableOpenAPIValidation
	client.SubNotes = rel.SubNotes
	client.Description = rel.Description

	return client
}

var (
	ErrNotFound      = driver.ErrReleaseNotFound
	ErrFoundMultiple = errors.New("found multiple releases o_0")
	ErrEmpty         = errors.New("releases are empty")
	ErrDepFailed     = errors.New("dependency failed")
)

// UniqName redis@my-namespace
func (rel *Config) UniqName() string {
	return rel.Name + "@" + rel.Namespace
}

// In check that 'x' found in 'array'
func (rel *Config) In(a []*Config) bool {
	for _, r := range a {
		if rel == r {
			return true
		}
	}
	return false
}
