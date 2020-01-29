#!/bin/sh
# -*- coding: utf-8 -*-

go mod download
go build -o /usr/bin/app .

exec app ${@}
