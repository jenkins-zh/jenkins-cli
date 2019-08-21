pipeline {
    agent {
        label 'golang-1.12'
    }

    stages {
        stage('Build') {
            parallel {
                stage('MacOS') {
                    steps {
                        container('golang') {
                            sh label: 'make darwin', script: 'make darwin'
                        }
                    }
                }
                stage('Linux') {
                    steps {
                        container('golang') {
                            sh label: 'make linux', script: 'make linux'
                        }
                    }
                }
                stage('Windows') {
                    steps {
                        container('golang') {
                            sh label: 'make win', script: 'make win'
                        }
                    }
                }
            }
        }
    }
}