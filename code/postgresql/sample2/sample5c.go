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
)

func Run() exitcode.ExitCode {
	// create gorm.DB instance for PostgreSQL service
	gormCtx, err := orm.NewGORM()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return exitcode.Abnormal
	}
	defer gormCtx.Close()

	file3 := "files/file3.txt"
	bin3, err := files.GetBinary(file3)
	if err != nil {
		gormCtx.GetLogger().Error().Interface("error", errs.Wrap(err)).Send()
		return exitcode.Abnormal
	}

	data := &model.User{
		Username: "Bob",
		BinaryFiles: []model.BinaryFile{
			{Filename: file3, Body: bin3},
		},
	}

	// insert data
	tx := gormCtx.GetDb().WithContext(context.TODO()).Create(data)
	if tx.Error != nil {
		gormCtx.GetLogger().Error().Interface("error", errs.Wrap(tx.Error)).Send()
		return exitcode.Abnormal
	}

	return exitcode.Normal
}

func main() {
	Run().Exit()
}
