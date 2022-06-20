package common

import "github.com/ilyakaznacheev/cleanenv"

var Conf conf

func init() {
	cleanenv.ReadEnv(&Conf)
}
