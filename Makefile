BINARY_NAME=MeDo.app
APP_NAME=MeDo
VERSION=0.0.1
BUILD_NO=1

## Run with local test DB
run:
	env DB_PATH="./sql.db" go run -v .

runp:
	go run -v .

fynepackage:
	@rm -rf ${BINARY_NAME}
	@#fyne package -appVersion ${VERSION} -appBuild ${BUILD_NO} -name ${APP_NAME} -release
	fyne package -name ${APP_NAME} -release

bundle:
	@echo "Generating bundled.go from resource"
	go generate






