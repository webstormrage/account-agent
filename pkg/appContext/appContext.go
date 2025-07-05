package appContext

import (
	"github.com/joho/godotenv"
	"os"
)

type Context struct {
	DataSourceName string
}

var context *Context

func Init()error{
    err := godotenv.Load()
	if err != nil {
		return err
	}
	context = &Context{
		DataSourceName: os.Getenv("DATA_SOURCE_NAME"),
	}
	return nil
}

func Get()Context{
	return *context
}