package service

import (
	"context"
	"github.com/KirillMironov/ci/internal/domain"
	"github.com/KirillMironov/ci/pkg/logger"
	"io"
	"os"
	"time"
)

// Poller is a service that can poll a VCS and execute a pipeline.
type Poller struct {
	ciFilename    string
	cloner        cloner
	archiver      archiver
	parser        parser
	executor      executor
	vcsRepository vcsRepository
	logger        logger.Logger
}

type (
	// cloner is a service that can clone a repository.
	cloner interface {
		CloneRepository(url string) (sourceCodePath string, removeSourceCode func() error, err error)
	}
	// archiver is a service that works with archives.
	archiver interface {
		Compress(dir string) (archivePath string, removeArchive func() error, err error)
		FindFile(filename, archivePath string) ([]byte, error)
	}
	// parser is a service that can parse a pipeline.
	parser interface {
		ParsePipeline(b []byte) (domain.Pipeline, error)
	}
	// executor is a service that can execute pipeline steps.
	executor interface {
		ExecuteStep(ctx context.Context, step domain.Step, sourceCodeArchive io.Reader) (logs io.ReadCloser, err error)
	}
	// vcsRepository stores information about a VCS repository.
	vcsRepository interface {
		Put(vcs domain.VCS) error
		GetAll() (arr []domain.VCS, err error)
	}
)

// NewPoller creates a new Poller.
func NewPoller(ciFilename string, cloner cloner, archiver archiver, parser parser, executor executor,
	vcsRepository vcsRepository, logger logger.Logger) *Poller {
	return &Poller{
		ciFilename:    ciFilename,
		cloner:        cloner,
		archiver:      archiver,
		parser:        parser,
		executor:      executor,
		vcsRepository: vcsRepository,
		logger:        logger,
	}
}

// Recover starts polling saved repositories.
func (p Poller) Recover() error {
	arr, err := p.vcsRepository.GetAll()
	if err != nil {
		return err
	}

	for _, vcs := range arr {
		go p.Start(vcs)
	}

	return nil
}

// Start starts VCS polling with a given interval.
func (p Poller) Start(vcs domain.VCS) {
	err := p.vcsRepository.Put(vcs)
	if err != nil {
		p.logger.Errorf("failed to put VCS %q: %v", vcs.URL, err)
		return
	}

	timer := time.NewTimer(vcs.PollingInterval)

	for range timer.C {
		err := p.poll(vcs)
		if err != nil {
			p.logger.Error(err)
		}
		timer.Reset(vcs.PollingInterval)
	}
}

// poll clones a repository and executes a pipeline.
func (p Poller) poll(vcs domain.VCS) error {
	sourceCodePath, removeSourceCode, err := p.cloner.CloneRepository(vcs.URL)
	if err != nil {
		return err
	}
	defer func() {
		if err = removeSourceCode(); err != nil {
			p.logger.Error(err)
		}
	}()

	archivePath, removeArchive, err := p.archiver.Compress(sourceCodePath)
	if err != nil {
		return err
	}
	defer func() {
		if err = removeArchive(); err != nil {
			p.logger.Error(err)
		}
	}()

	yaml, err := p.archiver.FindFile(p.ciFilename, archivePath)
	if err != nil {
		return err
	}

	pipeline, err := p.parser.ParsePipeline(yaml)
	if err != nil {
		return err
	}

	for _, step := range pipeline.Steps {
		err = p.executeStep(step, archivePath)
		if err != nil {
			return err
		}
	}

	return nil
}

// executeStep executes a pipeline step.
func (p Poller) executeStep(step domain.Step, sourceCodeArchivePath string) error {
	archive, err := os.Open(sourceCodeArchivePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	logs, err := p.executor.ExecuteStep(context.Background(), step, archive)
	if err != nil {
		return err
	}
	defer logs.Close()

	_, err = io.Copy(os.Stdout, logs)
	return err
}
