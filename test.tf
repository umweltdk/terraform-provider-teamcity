provider "teamcity" {
    url = "http://localhost:8111"
    user = "admin"
    password = "admin"
}

resource "teamcity_build_configuration" "terraform-provider-teamcity" {
  project = "Single"
  name = "terraform-provider-teamcity"
  parameters {
      name = "env.MUH"
      type = "password"
    }
  parameters {
      name = "env.TEST"
      type = "text"
      validation_mode = "not_empty"
      label = "Test framwork"
      description = "Name of the test framework to use"
    }
  parameter_values = {
    "env.TEST" = "Hello"
  }
  steps {
      type = "DockerTest"
      name = ""
      properties = {
        "deploy.environment.whitelist" = "^BUILD_\n^ELASTIC_BEANSTALK_\n^AWS_\n^AZURE_WA_\n^TRAVIS_"
        "env.DOCKER_TAG" = "%env.DOCKER_TAG%"
        "teamcity.step.mode" = "default"
      }
    }
  steps {
      type = "DockerBuild"
      name = ""
      properties = {
        "docker.build.path" = "."
        "env.DOCKER_TAG" = "%env.DOCKER_TAG%"
        "teamcity.step.mode" = "default"
      }
    }
  steps = [{
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
    },
    {
      type = "simpleRunner"
      name = "Replace DOCKER_TAG in Dockerrun.aws.json"
      properties = {
        "script.content" = "set -o errexit\n[ -n \"$DEBUG\" ] && set -o xtrace\n\nif [[ -f \"Dockerrun.aws.json\" ]]; then\n  REPLACE=\"$(echo ${DOCKER_TAG} | sed s/\\\\//\\\\\\\\\\\\//g)\"\n  sed s/\\${DOCKER_TAG}/${REPLACE}/ Dockerrun.aws.json\n  cp Dockerrun.aws.json Dockerrun.aws.json.original\n  sed s/\\${DOCKER_TAG}/${REPLACE}/ Dockerrun.aws.json.original > Dockerrun.aws.json\nelse\n  echo \"WARNING: Skipping because no Dockerrun.aws.json was found\"\nfi"
        "use.custom.script" = "true"
        "teamcity.step.mode" = "default"
      }
    },
    {
      type = "simpleRunner"
      name = "Package for Elastic Beanstalk"
      properties = {
        "script.content" = "set -o errexit\nif [ -d \".ebextensions\" -a -f \"Dockerrun.aws.json\" ]; then\n  echo \"Creating eb-deploy.zip\"\n  rm -f eb-deploy.zip\n  zip -r eb-deploy.zip .ebextensions Dockerrun.aws.json\nelif [ -f \"Dockerrun.aws.json\" ]; then\n  echo \"Creating eb-deploy.zip\"\n  rm -f eb-deploy.zip\n  zip -r eb-deploy.zip Dockerrun.aws.json\nelse\n  echo \"Nothing to do - skipping step\"\nfi"
        "use.custom.script" = "true"
        "teamcity.step.mode" = "default"
      }
    },
    {
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
  ]
  attached_vcs_roots = [
    {
      vcs_root = "Single_HttpsGithubComUmweltdkDockerNodeGit"
      checkout_rules = "+:refs/heads/master\n+:refs/heads/develop"
    }
  ]
}