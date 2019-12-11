.SILENT: ;               # no need for @

#SWAGGERYAMLURL = $(shell bash -c 'read -p "Enter swagger ui url (Path where this utility will host the swagger yaml, example: http://myhost.mycompany.com:8080/swagger.yaml): " surl; echo $$surl')
GATEWAYURL = $(shell bash -c 'read -p "Enter open_faas gateway url. (without the trailing slash, example: http://mygateway.mycompany.com:8080): " gurl; echo $$gurl')
REPONAME = $(shell bash -c 'read -p "Enter repository name. (example: myrepo.mycompany.com/acct): " rname; echo $$rname')

.prompt-yesno:
	@exec 9<&0 0</dev/tty; \
	echo "$(message) [Y]:"; \
	read -rs -n 1 yn; \
	exec 0<&9 9<&-; \
	[[ -z $$yn ]] || [[ $$yn == [yY] ]] && echo Y >&2 || (echo N >&2 && exit 1)

.backup:
	mkdir -p .deploy
	cp -R swaggerui .deploy/
	cp k8.yaml .deploy/

all: .backup
	  echo; \
	  echo "Please review what you have entered."; \
	  if ! make .prompt-yesno message="Do you want to continue?" 2> /dev/null; then \
	    echo "Aborted on request."; \
	    exit; \
	  fi; \
	echo "Let's go"
	#sed "s,http://localhost:8080/swagger.yaml,$(SWAGGERYAMLURL)," swaggerui/index.html > .deploy/swaggerui/index.html
	#echo "Check changed lines"
	#-diff .deploy/swaggerui/index.html swaggerui/index.html
	-g=$(GATEWAYURL); r=$(REPONAME); \
	sed -e "s,http://gateway:8080,$$g," -e "s,http://localhost:8080,$$g," -e "s,<changereponame>,$$r," k8.yaml > .deploy/k8.yaml; \
	diff .deploy/k8.yaml k8.yaml; \
	docker build --rm -t $$r/faas-swagger:latest .; \
