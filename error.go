package main

import "github.com/Sirupsen/logrus"

var log = logrus.New()

func check(err error) {
	if err != nil {
		log.Error(err)
	}
}
