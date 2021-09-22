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
	"gorm.io/gorm/clause"
)

func Run() exitcode.ExitCode {
	// create gorm.DB instance for PostgreSQL service
	gormCtx, err := orm.NewGORM()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitcode.Abnormal
	}
	defer gormCtx.Close()

	file4 := "files/file4.txt"
	bin4, err := files.GetBinary(file4)
	if err != nil {
		gormCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
		return exitcode.Abnormal
	}

	// select data for 'Bob 2nd' (with preload)
	data := []model.User{}
	tx := gormCtx.GetDb().WithContext(context.TODO()).Preload(clause.Associations).Where(&model.User{Username: "Bob 2nd"}).Find(&data)
	if tx.Error != nil {
		gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
		return exitcode.Abnormal
	}

	// update data in binary_files table
	tx = gormCtx.GetDb().WithContext(context.TODO()).Model(&data[0].BinaryFiles[0]).Updates(model.BinaryFile{Filename: file4, Body: bin4})
	if tx.Error != nil {
		gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
		return exitcode.Abnormal
	}

	return exitcode.Normal
}

func main() {
	Run().Exit()
}
