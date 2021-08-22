FROM openjdk:8-jdk
LABEL Version="1.0-SNAPSHOT"Description="Upgraded jenkins core and plugins in a YAML specification"Vendor="jenkins-cli"
ADD target/jcli-yaml-demo-1.0-SNAPSHOT.war /usr/share/jenkins/jenkins.war
ENTRYPOINT ["tini", "--", "/usr/local/bin/jenkins.sh"]
