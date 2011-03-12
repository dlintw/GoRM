include $(GOROOT)/src/Make.inc

TARG=sdegutis/gorm
GOFILES=\
	gorm.go \
	util.go

include $(GOROOT)/src/Make.pkg
