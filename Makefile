NAME=$(shell basename `pwd` | cut -d '-' -f 1)
VERSION=0.1

N=23
#$(shell head -n1 /etc/issue | cut -d ' ' -f 3)

EXE=$(shell basename `pwd`)
PKGNAME=$(EXE)
RPMSHORT=$(PKGNAME)-$(VERSION)-1.x86_64.rpm
RPM=package/RPMS/x86_64/$(RPMSHORT)

all: $(EXE)

rpm: $(RPM)

$(EXE):
	go test
	go build
	mkdir -p static
	mkdir -p content
	mkdir -p site
	cp -a $(EXE) site/
	cp -a config site/
	cp -a add-user revoke-auth err static content site/

put: $(RPM)
	rsync -az $(RPM) bpowers.net:.
	ssh bpowers.net -t sudo rpm --force -fvi ./$(RPMSHORT)

$(RPM): $(EXE)
	cp -a site $(PKGNAME)-$(VERSION)
	mkdir -p package/{RPMS,BUILD,SOURCES,BUILDROOT}
	tar -czf package/SOURCES/$(PKGNAME)-$(VERSION).tar.gz $(PKGNAME)-$(VERSION)
	rm -rf $(PKGNAME)-$(VERSION)
	cat server.service.in | sed "s/%NAME%/$(NAME)/g" >package/SOURCES/server.service
	cat server.spec.in | sed "s/%NAME%/$(NAME)/g" | sed "s/%VERSION%/$(VERSION)/g" >server.spec
	rpmbuild --define "_topdir $(PWD)/package" -ba server.spec
	rm -rf package/{BUILD,BUILDROOT}

serve: $(EXE)
	./$(EXE) -dev

clean:
	rm -f server.spec $(EXE)

.PHONY: $(EXE) rpm clean serve put
