package watcher

import (
	"go-blog/pkg/util/log"
	"os"
	"path"
	"path/filepath"
	"time"
)

// NewWatcher creates a new watcher instance
func NewWatcher(path, templateExtension string, checkCycleDuration time.Duration, logger *log.Log, fileHandler func(string)) *Watcher {
	return &Watcher{
		logger:             logger,
		templatesExtension: templateExtension,
		pathToWatch:        path,
		checkCycleDuration: checkCycleDuration,
		fileHandler:        fileHandler,
	}
}

// Watcher is able to watch a folder in order to process when new files are created
type Watcher struct {
	logger             *log.Log
	templatesExtension string
	pathToWatch        string
	quitChannel        chan bool
	checkCycleDuration time.Duration
	fileHandler        func(string)
}

// Start begins with the watching process
func (w *Watcher) Start() {

	w.logger.Info("starting watcher on "+w.pathToWatch, nil)

	go w.watch()

	// keep running until it's stopped
	w.quitChannel = make(chan bool)
	<-w.quitChannel

	w.logger.Info("stopping watcher on "+w.pathToWatch, nil)
}

// Stop ends the watching process
func (w *Watcher) Stop() {
	if w.quitChannel != nil {
		w.quitChannel <- true
	}
}

// get files in pathToLook, filter by the indicated file extension
func (w *Watcher) listExsitingFiles(pathToLook string, extension string) (currentFiles []string, err error) {
	err = filepath.Walk(pathToLook, func(filepath string, info os.FileInfo, err error) error {
		if path.Ext(filepath) == extension {
			currentFiles = append(currentFiles, filepath)
		}
		return nil
	})

	return
}

func (w *Watcher) watch() {

	// process existing files
	for {
		w.logger.Debug("checking for new files", nil)
		existingFiles, err := w.listExsitingFiles(w.pathToWatch, w.templatesExtension)
		if err != nil {
			w.logger.Error("error reading existing files in folder to watch", err, map[string]interface{}{"path": w.pathToWatch})
			return
		}

		processedCount := 0
		if len(existingFiles) > 0 {
			for i := range existingFiles {
				if path.Dir(existingFiles[i]) == w.pathToWatch {
					w.logger.Info("found existing template; sending to be processed", map[string]interface{}{"file": existingFiles[i]})
					go w.fileHandler(existingFiles[i])
					processedCount++
				}
			}
		}

		// if processedCount == 0 {
		// 	w.logger.Debug("no new files found", nil)
		// }

		// wait before checking again
		time.Sleep(w.checkCycleDuration)
	}

}
