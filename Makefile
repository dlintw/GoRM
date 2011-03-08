include $(GOROOT)/src/Make.inc

TARG=gorm
GOFILES=\
	gorm.go \
	util.go

include $(GOROOT)/src/Make.pkg
