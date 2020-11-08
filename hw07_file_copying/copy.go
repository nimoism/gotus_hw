package main

import (
	"errors"
	"io"
	"os"

	"github.com/nimoism/gotus_hw/hw07_file_copying/progress"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrSameFiles             = errors.New("from and to files are same")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	var from, to *os.File
	var fromStat, toStat os.FileInfo
	var err error
	if from, err = os.Open(fromPath); err != nil {
		return err
	}
	defer from.Close()
	if to, err = os.Create(toPath); err != nil {
		return err
	}
	defer to.Close()
	if fromStat, err = from.Stat(); err != nil {
		return err
	}
	if toStat, err = to.Stat(); err != nil {
		return err
	}
	if os.SameFile(fromStat, toStat) {
		return ErrSameFiles
	}
	if limit, err = calcLimit(fromStat, offset, limit); err != nil {
		return err
	}
	if offset > 0 {
		if _, err = from.Seek(offset, io.SeekStart); err != nil {
			return err
		}
	}

	var w io.Writer = to

	width, _, err := progress.Term{}.Size()
	// Continue without progress bar on getting terminal width error
	if err == nil {
		w = io.MultiWriter(w, progress.NewBar(limit, os.Stdout, width))
	}

	if _, err = io.CopyN(w, from, limit); err != nil {
		return err
	}

	return nil
}

func calcLimit(stat os.FileInfo, offset, limit int64) (int64, error) {
	if !stat.Mode().IsRegular() {
		if limit < 1 {
			return 0, ErrUnsupportedFile
		}
		return limit, nil
	}
	size := stat.Size()
	if size < offset {
		return 0, ErrOffsetExceedsFileSize
	}
	if limit < 1 || size < offset+limit {
		return size - offset, nil
	}
	return limit, nil
}
