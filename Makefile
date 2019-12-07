.SILENT: ;               # no need for @

SWAGGERYAMLURL = $(shell bash -c 'read -p "Enter swagger ui url (example: http://myhost.mycompany.com:8080/swagger.yaml): " surl; echo $$surl')
GATEWAYURL = $(shell bash -c 'read -p "Enter open_faas gateway url. (example: http://mygateway.mycompany.com:8080): " gurl; echo $$gurl')
REPONAME = $(shell bash -c 'read -p "Enter repository name. (example: myrepo.mycompany.com/acct): " rname; echo $$rname')

.prompt-yesno:
	@exec 9<&0 0</dev/tty; \
	echo "$(message) [Y]:"; \
	read -rs -n 1 yn; \
	exec 0<&9 9<&-; \
	[[ -z $$yn ]] || [[ $$yn == [yY] ]] && echo Y >&2 || (echo N >&2 && exit 1)

.backup:
	-[[ ! -f "swaggerui/index.html.bak" ]] && cp swaggerui/index.html swaggerui/index.html.bak
	-[[ ! -f "k8.yaml.bak" ]] && cp k8.yaml k8.yaml.bak

.statik:
	-[[ ! -f "$${HOME}/go/bin/statik" ]] && echo no statik && go get github.com/rakyll/statik

all: .backup .statik
	  echo; \
	  echo "Please review what you have entered. And do docker login first before continue."; \
	  if ! make .prompt-yesno message="Do you want to continue?" 2> /dev/null; then \
	    echo "Aborted on request."; \
	    exit; \
	  fi; \
	echo "Let's go"
	sed "s,http://localhost:8080/swagger.yaml,$(SWAGGERYAMLURL)," swaggerui/index.html.bak > swaggerui/index.html
	echo "Check changed lines"
	-diff swaggerui/index.html.bak swaggerui/index.html
	-g=$(GATEWAYURL); r=$(REPONAME); \
	sed -e "s,http://gateway:8080,$$g," -e "s,url: http://localhost:8080,url: $$g," -e "s,<changereponame>,$$r," k8.yaml.bak > k8.yaml; \
	diff k8.yaml.bak k8.yaml; \
	$${HOME}/go/bin/statik -f -src=swaggerui/; \
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o release/faas-swagger .; \
	docker build --rm -t $$r/faas-swagger:latest .; \
	docker push $$r/faas-swagger:latest
