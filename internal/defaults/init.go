package defaults

import (
	_cleanenv "github.com/ilyakaznacheev/cleanenv"
	_zap "go.uber.org/zap"
)

var Conf conf
var Logger *_zap.Logger

func init() {
	_cleanenv.ReadEnv(&Conf)
	Logger, _ = _zap.NewDevelopment()
}
