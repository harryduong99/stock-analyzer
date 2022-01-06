package sourcefactory

import (
	"github.com/duongnam99/stock-analyzer/config"
)

type ISourceHandlerFactory interface {
	GetData([]string, int, string)
}

func GetSourceHandlerFactory(sourceType string) ISourceHandlerFactory {

	switch sourceType {
	case config.CAFEF:
		return CafefSourceHandler{}
	case config.VNDIRECT:

	case config.VIETSTOCK:

	}

	return nil
}
