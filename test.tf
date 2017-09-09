provider "teamcity" {
    url = "http://localhost:8112"
    user = "admin"
    password = "admin"
}

resource "teamcity_build_template" "terraform-provider-teamcity" {
  project = "Single"
  name = "terraform-provider-teamcity"
  parameter {
      name = "env.MUH"
      type = "password"
    }
  parameter {
      name = "env.TEST"
      type = "text"
      validation_mode = "not_empty"
      label = "Test framwork"
      description = "Name of the test framework to use"
    }
  parameter_values = {
    "env.TEST" = "Hello"
  }
  step {
    type = "DockerBuild"
    name = ""
    properties = {
      "docker.build.path" = "."
      "env.DOCKER_TAG" = "%env.DOCKER_TAG%"
      "teamcity.step.mode" = "default"
    }
  }
  step {
    type = "DockerTest"
    name = ""
    properties = {
      "deploy.environment.whitelist" = "^BUILD_\n^ELASTIC_BEANSTALK_\n^AWS_\n^AZURE_WA_\n^TRAVIS_"
      "env.DOCKER_TAG" = "%env.DOCKER_TAG%"
      "teamcity.step.mode" = "default"
    }
  }
  step {
      type = "simpleRunner"
      name = "Docker Push"
      properties = {
        "script.content" = <<EOF
set -o errexit
[ -n \"$DEBUG\" ] && set -o xtrace

function retry {
  local wait=$1
  local retry=$2
  shift 2
  local times=1
  while ! "$@" ; do
    if [[ $times -ge $retry ]]; then
      echo "task failed $retry times quitting"
      return 1
    fi
    echo "task failed sleeping $wait seconds before retrying"
    sleep $wait
    ((times=$times+1))
  done
}

export HOME="%system.teamcity.build.tempDir%"
echo "##teamcity[blockOpened name='Save docker login']"
retry "%retry.wait%" "%retry.count%" docker login --email=%docker.login.email% --username=%docker.login.username% --password="%docker.login.password%"
echo "##teamcity[blockClosed name='Save docker login']"

echo "##teamcity[blockOpened name='Push %env.DOCKER_TAG%']"
retry "%retry.wait%" "%retry.count%" docker push %env.DOCKER_TAG%
echo "##teamcity[blockClosed name='Push %env.DOCKER_TAG%']"
EOF
      "use.custom.script" = "true"
      "teamcity.step.mode" = "default"
    }
  }
  step {
    type = "simpleRunner"
    name = "Replace DOCKER_TAG in Dockerrun.aws.json"
    properties = {
      "script.content" = <<EOF
set -o errexit
[ -n "$DEBUG" ] && set -o xtrace

if [[ -f "Dockerrun.aws.json" ]]; then
  REPLACE="$(echo ${DOCKER_TAG} | sed s/\\\\//\\\\\\\\\\\\//g)"
  sed s/\\${DOCKER_TAG}/${REPLACE}/ Dockerrun.aws.json
  cp Dockerrun.aws.json Dockerrun.aws.json.original
  sed s/\\${DOCKER_TAG}/${REPLACE}/ Dockerrun.aws.json.original > Dockerrun.aws.json
else
  echo "WARNING: Skipping because no Dockerrun.aws.json was found"
fi
EOF
      "use.custom.script" = "true"
      "teamcity.step.mode" = "default"
    }
  }
  step {
    type = "simpleRunner"
    name = "Package for Elastic Beanstalk"
    properties = {
      "script.content" = <<EOF
set -o errexit
if [ -d ".ebextensions" -a -f "Dockerrun.aws.json" ]; then
  echo "Creating eb-deploy.zip"
  rm -f eb-deploy.zip
  zip -r eb-deploy.zip .ebextensions Dockerrun.aws.json
elif [ -f "Dockerrun.aws.json" ]; then
  echo "Creating eb-deploy.zip"
  rm -f eb-deploy.zip
  zip -r eb-deploy.zip Dockerrun.aws.json
else
  echo "Nothing to do - skipping step"
fi
EOF
      "use.custom.script" = "true"
      "teamcity.step.mode" = "default"
    }
  }
  step {
    type = "DeployWithUmploy"
    name = ""
    properties = {
      "deploy.environment.whitelist" = "^BUILD_\n^AWS_\n^ELASTIC_BEANSTALK_\n^AZURE_WA_\n^TRAVIS_"
      "deploy.retry.count" = "%retry.count%"
      "deploy.retry.wait" = "%retry.wait%"
      "deploy.source.path" = "%teamcity.build.checkoutDir%"
      "env.BUILD_BRANCH" = "%env.BUILD_BRANCH%"
      "env.BUILD_REF" = "%env.BUILD_REF%"
      "teamcity.step.mode" = "default"
    }
  }
  attached_vcs_root {
    vcs_root = "Single_HttpsGithubComUmweltdkDockerNodeGit"
    checkout_rules = "+:refs/heads/master\n+:refs/heads/develop"
  }
}

resource "teamcity_build_configuration" "terraform-provider-teamcit" {
  project = "Single"
  name = "terraform-provider-teamcit"
  template = "${teamcity_build_template.terraform-provider-teamcity.id}"
  parameter {
    name = "env.DER"
    type = "text"
  }

  step {
    type = "DeployWithUmploy"
    name = "Muhd"
    properties = {
      "deploy.environment.whitelist" = "^BUILD_\n^AWS_\n^ELASTIC_BEANSTALK_\n^AZURE_WA_\n^TRAVIS_"
      "deploy.retry.count" = "%retry.count%"
      "deploy.retry.wait" = "%retry.wait%"
      "deploy.source.path" = "%teamcity.build.checkoutDir%"
      "env.BUILD_BRANCH" = "%env.BUILD_BRANCH%"
      "env.BUILD_REF" = "%env.BUILD_REF%"
      "teamcity.step.mode" = "default"
    }
  }

}