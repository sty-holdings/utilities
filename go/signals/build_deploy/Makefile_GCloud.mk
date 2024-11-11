# Makefile_GCloud.mk
#
# This will push the build of the SavUp HTTP Server to the GCloud compute instance.
#
# Check that the correct CONFIG_NAME is uncommented before running all.
#
# Requirements:
# - The GC_SERVER must already be provisioned and running
# - Your .ssh directory must contain the google_compute_engine and google_compute_engine.pub file so you can ssh to the instance.
#    - Contact the System Administrator to be added to Google IAM, if you do not have the files.
#
# Relative Path:
# - Change the ROOT_DIRECTORY export to a dot - DO NOT PUSH THIS CHANGE TO GIT REMOTE REPOSITORY
# - Change the INSTALL_ROOT_DIRECTORY export to a dot - DO NOT PUSH THIS CHANGE TO GIT REMOTE REPOSITORY
#
# NOTE: If you have an improvement or a correct, please make the change and check it into the repository.
#
###########################
# User controlled parameter - This is not a CLI argument because is doesn't change frequently
###########################
export GC_INSTANCE_NAME=savup-local-0030
export GC_PROJECT_ID=savup-development
export GC_REGION=us-central1-c
export GC_REMOTE_LOGIN=scott_yacko_sty_holdings_com@${GC_INSTANCE_NAME}

####################################
# DO NOT CHANGE BELOW THIS LINE
####################################
# FileNames
export SERVER_NAME=signals

# SOURCE DIRECTORIES
export ROOT_DIRECTORY=/Users/syacko
export SOURCE_DIRECTORY=${ROOT_DIRECTORY}/workspace/styh-dev/src
export BINARY_NAME=${SOURCE_DIRECTORY}/albert/Utilities/${SERVER_NAME}/bin/${SERVER_NAME}
export TEMPLATE_FILE_DIRECTORY=${SOURCE_DIRECTORY}/albert/Utilities/${SERVER_NAME}/build_deploy/templates

# TARGET DIRECTORIES
export INSTALL_ROOT_DIRECTORY=/home/scott_yacko_sty_holdings_com

all: settings setGoogleProject preInstall buildLinux installingFiles installDaemon clean displayRunNotes

settings:
  	$(info GC INSTANCE NAME:          $(GC_INSTANCE_NAME))
  	$(info GC PROJECT ID:             $(GC_PROJECT_ID))
  	$(info GC REMOTE LOGIN:           ${GC_REMOTE_LOGIN})
  	$(info SERVER NAME:               $(SERVER_NAME))
  	$(info -------------------------)
  	$(info Here are the pre-set or defined variables:)
  	$(info BINARY NAME:             ${SOURCE_DIRECTORY}/albert/${SERVER_NAME}/bin/${SERVER_NAME})
  	$(info GC REGION:               $(GC_REGION))
  	$(info INSTALL ROOT DIRECTORY:  $(INSTALL_ROOT_DIRECTORY))
  	$(info ROOT DIRECTORY:          $(ROOT_DIRECTORY))
  	$(info SOURCE DIRECTORY:        $(SOURCE_DIRECTORY))
  	$(info )

preInstall:
	$(info )
	$(info preInstall)
	$(info -------------------------)
	-gcloud compute ssh --zone ${GC_REGION} ${GC_REMOTE_LOGIN} --command "rm  ${INSTALL_ROOT_DIRECTORY}/bin/${SERVER_NAME};"

buildLinux:
	$(info )
	$(info buildLinux)
	$(info -------------------------)
	env GOOS=linux GOARCH=amd64 go build -o ${BINARY_NAME} ${SOURCE_DIRECTORY}/albert/Utilities/${SERVER_NAME}/main.go

setGoogleProject:
	$(info )
	$(info setGoogleProject)
	$(info -------------------------)
	gcloud config set project "${GC_PROJECT_ID}";


installingFiles:
	$(info )
	$(info installingFiles)
	$(info -------------------------)
	gcloud compute scp --recurse --zone ${GC_REGION} ${BINARY_NAME} ${GC_REMOTE_LOGIN}:${INSTALL_ROOT_DIRECTORY}/bin/${SERVER_NAME}

installDaemon:
	$(info )
	$(info installDaemon)
	$(info -------------------------)
	envsubst '$${INSTALL_ROOT_DIRECTORY},$${SERVER_NAME}' < ${TEMPLATE_FILE_DIRECTORY}/${SERVER_NAME}.servicefile.template > /tmp/${SERVER_NAME}.servicefile.tmp
	gcloud compute scp --recurse --zone ${GC_REGION} /tmp/${SERVER_NAME}.servicefile.tmp ${GC_REMOTE_LOGIN}:${INSTALL_ROOT_DIRECTORY}/.config/${SERVER_NAME}.servicefile

	envsubst < ${TEMPLATE_FILE_DIRECTORY}/${SERVER_NAME}-install-daemon.sh.template > /tmp/${SERVER_NAME}-install-daemon.sh.tmp
	gcloud compute scp --recurse --zone ${GC_REGION} /tmp/${SERVER_NAME}-install-daemon.sh.tmp ${GC_REMOTE_LOGIN}:${INSTALL_ROOT_DIRECTORY}/scripts/${SERVER_NAME}-install-daemon.sh

	gcloud compute ssh --zone ${GC_REGION} ${GC_REMOTE_LOGIN} --command "sudo sh ${INSTALL_ROOT_DIRECTORY}/scripts/${SERVER_NAME}-install-daemon.sh"

displayRunNotes:
	$(info )
	$(info displayRunNotes)
	$(info -------------------------)
	$(info Daemon Commands:)
	$(info -  sudo systemctl status ${SERVER_NAME}.service)
	$(info -  sudo systemctl start ${SERVER_NAME}.service)
	$(info -  sudo systemctl stop ${SERVER_NAME}.service)
	$(info -  sudo systemctl restart ${SERVER_NAME}.service)
	$(info -  sudo journalctl -u ${SERVER_NAME}.service -n 50)

clean:
	$(info )
	$(info clean)
	$(info -------------------------)
	go clean
