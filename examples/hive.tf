resource "aws_elastic_beanstalk_application" "hive" {
  name = "umwelt-hive"
  description = "Providing the backbone of Umwelts internal services"
}

resource "aws_elastic_beanstalk_environment" "hive-dev" {
  name = "umw-hive-dev"
  application = "${aws_elastic_beanstalk_application.hive.name}"
  description = "Development environment"
  solution_stack_name = "64bit Amazon Linux 2016.03 v2.1.6 running Docker 1.11.2"
  wait_for_ready_timeout = "10m"
}