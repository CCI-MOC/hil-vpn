
# Variables for various system directories; see:
#
# https://www.gnu.org/prep/standards/standards.html#Directory-Variables
#
# ...which describes a convention that is typical on unix-like systems.
# Note that the bit about preferring info files to man pages does not
# really reflect the reality of most systems (even GNU/Linux systems),
# but the rest of it is good.
#
# These can be overridden by doing e.g:
#
#    make PREFIX=/usr SYSCONFDIR=/etc ...
#
# Note that:
#
# 1. We use uppercase variable names, as these are a bit more idiomatic in
#    makefiles (as opposed to autotools stuff).
# 2. hil-vpn doesn't actually use most of these, and probably never will,
#    but we specify the whole set for the sake of making this orthogonal
#    to the rest of the code.
PREFIX         ?= /usr/local
EXEC_PREFIX    ?= $(PREFIX)
BINDIR         ?= $(EXEC_PREFIX)/bin
SBINDIR        ?= $(EXEC_PREFIX)/sbin
LIBEXECDIR     ?= $(EXEC_PREFIX)/libexec
DATAROOTDIR    ?= $(PREFIX)/share
DATADIR        ?= $(DATAROOTDIR)
SYSCONFDIR     ?= $(PREFIX)/etc
SHAREDSTATEDIR ?= $(PREFIX)/com
LOCALSTATEDIR  ?= $(PREFIX)/var
RUNSTATEDIR    ?= $(LOCALSTATEDIR)/run
INCLUDEDIR     ?= $(PREFIX)/include
OLDINCLUDEDIR  ?= /usr/include
DOCDIR         ?= $(DATAROOTDIR)/doc/hil-vpn
INFODIR        ?= $(DATAROOTDIR)/info
HTMLDIR        ?= $(DOCDIR)
DVIDIR         ?= $(DOCDIR)
PDFDIR         ?= $(DOCDIR)
PSDIR          ?= $(DOCDIR)
LIBDIR         ?= $(EXEC_PREFIX)/lib
LISDIR         ?= $(DATAROOTDIR)/emacs/site-lisp
LOCALEDIR      ?= $(DATAROOTDIR)/locale
MANDIR         ?= $(DATAROOTDIR)/man
# We don't bother with all of the man page directories for now.


# Shorthand for the static config package, in which we override several
# variables at link-time
CONFIGPKG := github.com/CCI-MOC/hil-vpn/internal/staticconfig

GO_LDFLAGS := \
	-X $(CONFIGPKG).Prefix=$(PREFIX) \
	-X $(CONFIGPKG).Execprefix=$(EXEC_PREFIX) \
	-X $(CONFIGPKG).Bindir=$(BINDIR) \
	-X $(CONFIGPKG).Sbindir=$(SBINDIR) \
	-X $(CONFIGPKG).Libexecdir=$(LIBEXECDIR) \
	-X $(CONFIGPKG).Datarootdir=$(DATAROOTDIR) \
	-X $(CONFIGPKG).Datadir=$(DATADIR) \
	-X $(CONFIGPKG).Sysconfdir=$(SYSCONFDIR) \
	-X $(CONFIGPKG).Sharedstatedir=$(SHAREDSTATEDIR) \
	-X $(CONFIGPKG).Localstatedir=$(LOCALSTATEDIR) \
	-X $(CONFIGPKG).Runstatedir=$(RUNSTATEDIR) \
	-X $(CONFIGPKG).Includedir=$(INCLUDEDIR) \
	-X $(CONFIGPKG).Oldincludedir=$(OLDINCLUDEDIR) \
	-X $(CONFIGPKG).Docdir=$(DOCDIR) \
	-X $(CONFIGPKG).Infodir=$(INFODIR) \
	-X $(CONFIGPKG).Htmldir=$(HTMLDIR) \
	-X $(CONFIGPKG).Dvidir=$(DVIDIR) \
	-X $(CONFIGPKG).Pdfdir=$(PDFDIR) \
	-X $(CONFIGPKG).Psdir=$(PSDIR) \
	-X $(CONFIGPKG).Libdir=$(LIBDIR) \
	-X $(CONFIGPKG).Lisdir=$(LISDIR) \
	-X $(CONFIGPKG).Localedir=$(LOCALEDIR) \
	-X $(CONFIGPKG).Mandir=$(MANDIR)

all:
	@echo BUILD hil-vpnd
	@cd cmd/hil-vpnd       ; go build -ldflags "$(GO_LDFLAGS)"
	@echo BUILD hil-vpn-privop
	@cd cmd/hil-vpn-privop ; go build -ldflags "$(GO_LDFLAGS)"
install:
	install -Dm755 ./cmd/hil-vpnd/hil-vpnd             $(SBINDIR)/
	install -Dm755 ./cmd/hil-vpn-privop/hil-vpn-privop $(LIBEXECDIR)/
	install -Dm755 ./openvpn-hooks/hil-vpn-hook-up     $(LIBEXECDIR)/

.PHONY: all install
