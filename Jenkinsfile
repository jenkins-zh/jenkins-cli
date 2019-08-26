library identifier: 'jenkins_zh@', retriever: modernSCM([$class: 'GitSCMSource', credentialsId: '', remote: 'https://github.com/LinuxSuRen/shared-library', traits: [[$class: 'jenkins.plugins.git.traits.BranchDiscoveryTrait']]])

pipeline {
    agent {
        label 'golang-1.12'
    }

    stages {
        stage('Build') {
            parallel {
                stage('MacOS') {
                    steps {
                        script {
                            entry.container_x('golang', 'go version'){
                                sh label: 'make darwin', script: 'make darwin'
                            }
                        }
                    }
                }
                stage('Linux') {
                    steps {
                        script {
                            entry.container_x('golang', 'go version'){
                                sh label: 'make linux', script: 'make linux'
                            }
                        }
                    }
                }
                stage('Windows') {
                    steps {
                        script {
                            entry.container_x('golang', 'go version'){
                                sh label: 'make win', script: 'make win'
                            }
                        }
                    }
                }
            }
        }

        stage('Test') {
            steps {
                script {
                    entry.container_x('golang', 'go version'){
                        sh label: 'go test', script: 'make test'
                    }
                }
            }
        }
    }

    post {
        always {
            junit allowEmptyResults: true, testResults: "*/**/*.xml"
        }
    }
}