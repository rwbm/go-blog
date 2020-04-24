package template

import (
	"fmt"
	"go-blog/pkg/util/log"
	"go-blog/pkg/util/model"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// NewProcessor creates a new instance of the template processor
func NewProcessor(database *gorm.DB, logger *log.Log, processedOKLocation string, processedErrorLocation string) *Processor {
	return &Processor{
		database:               database,
		logger:                 logger,
		processedOKLocation:    processedOKLocation,
		processedErrorLocation: processedErrorLocation,
	}
}

// Processor is used to process template files from the templates folder.
// Once processed, the files are moved to OK or Error folders, just for future references.
type Processor struct {
	logger                 *log.Log
	processedOKLocation    string
	processedErrorLocation string
	database               *gorm.DB
}

// ProcessTemplate process a template file, by
func (p *Processor) ProcessTemplate(filePath string) {

	p.logger.Info("processing file "+filePath, nil)

	// read file content
	data, errRead := ioutil.ReadFile(filePath)
	if errRead != nil {
		p.logger.Error("error reading template content", errRead, map[string]interface{}{"file": filePath})
		return
	}

	// parse
	post, errParse := ParseTemplate(string(data))
	if errParse != nil {
		p.logger.Error("error parsing template", errParse, map[string]interface{}{"file": filePath})

		// move file to error folder
		if errMove := p.moveFile(filePath, true); errMove != nil {
			p.logger.Error("error moving template", errMove, map[string]interface{}{"file": filePath})
		}

		return
	}

	// save original file name, for reference
	post.OriginalFileName = path.Base(filePath)

	// set date created and updated if wasn't set in the template
	if post.DateCreated.Year() == 1 {
		post.DateCreated = time.Now()
	}
	if post.DateUpdated.Year() == 1 {
		post.DateUpdated = time.Now()
	}

	// save in the database
	if errSave := p.savePost(&post); errSave != nil {
		p.logger.Error("error saving template to the database", errSave, map[string]interface{}{"file": filePath})

		// move file to error folder
		if errMove := p.moveFile(filePath, true); errMove != nil {
			p.logger.Error("error moving template", errMove, map[string]interface{}{"file": filePath})
		}

		return
	}

	// mark file to processed OK
	if errMove := p.moveFile(filePath, false); errMove != nil {
		p.logger.Error("error moving template", errMove, map[string]interface{}{"file": filePath})
		return
	}

	p.logger.Info("file "+filePath+" processed OK", nil)
}

func (p *Processor) savePost(post *model.Post) (err error) {

	trx := p.database.Begin()

	// save post
	if err = trx.Create(post).Error; err != nil {
		return
	}

	// save categories
	if len(post.Categories) > 0 {
		categories := strings.Split(post.Categories, ",")
		for i := range categories {
			if err = trx.Create(&model.PostCategory{IDPost: post.ID, Name: strings.Trim(categories[i], " ")}).Error; err != nil {
				trx.Rollback()
				return
			}
		}
	}

	// save tags
	if len(post.Tags) > 0 {
		tags := strings.Split(post.Tags, ",")
		for i := range tags {
			if err = trx.Create(&model.PostTag{IDPost: post.ID, Name: strings.Trim(tags[i], " ")}).Error; err != nil {
				trx.Rollback()
				return
			}
		}
	}

	trx.Commit()
	return
}

func (p *Processor) moveFile(srcFile string, failed bool) (err error) {

	// move to OK or error?
	destPath := p.processedOKLocation
	if failed {
		destPath = p.processedErrorLocation
	}

	// prefix new file name with timestamp
	srcFileName := filepath.Base(srcFile)
	newFileName := fmt.Sprintf("%v_%s", time.Now().Unix(), srcFileName)
	destFile := path.Join(destPath, newFileName)

	// copy file to destination
	if err = copyFile(srcFile, destFile); err != nil {
		return
	}

	// remove original file
	err = os.Remove(srcFile)

	return
}

// copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func copyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}

	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		return fmt.Errorf("non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}

	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}

	if err = os.Link(src, dst); err == nil {
		return
	}

	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}

	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
