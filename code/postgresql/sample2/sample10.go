//go:build run
// +build run

package main

import (
	"context"
	"fmt"
	"os"
	"sample/files"
	"sample/gorm/model"
	"sample/orm"

	"github.com/spiegel-im-spiegel/errs"
	"github.com/spiegel-im-spiegel/gocli/exitcode"
	"gorm.io/gorm"
)

func Run() exitcode.ExitCode {
	// create gorm.DB instance for PostgreSQL service
	gormCtx, err := orm.NewGORM()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitcode.Abnormal
	}
	defer gormCtx.Close()

	file1 := "files/file1.txt"
	bin1, err := files.GetBinary(file1)
	if err != nil {
		gormCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
		return exitcode.Abnormal
	}
	file2 := "files/file2.txt"
	bin2, err := files.GetBinary(file2)
	if err != nil {
		gormCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
		return exitcode.Abnormal
	}
	file3 := "files/file3.txt"
	bin3, err := files.GetBinary(file3)
	if err != nil {
		gormCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
		return exitcode.Abnormal
	}
	data1 := &model.User{
		Username: "Alice",
		BinaryFiles: []model.BinaryFile{
			{Filename: file1, Body: bin1},
			{Filename: file2, Body: bin2},
		},
	}
	data2 := &model.User{
		Username: "Bob",
		BinaryFiles: []model.BinaryFile{
			{Filename: file3, Body: bin3},
		},
	}

	if err := gormCtx.GetDb().WithContext(context.TODO()).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(data1).Error; err != nil {
			return errs.Wrap(err) // return any error will rollback
		}
		if err := tx.Create(data2).Error; err != nil {
			return errs.Wrap(err) // return any error will rollback
		}
		return nil // return nil will commit the whole transaction
	}); err != nil {
		gormCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
		return exitcode.Abnormal
	}

	return exitcode.Normal
}

func main() {
	Run().Exit()
}
