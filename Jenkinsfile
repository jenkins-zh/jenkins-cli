pipeline {
    agent {
        label 'golang'
    }

    stages {
        container('golang') {
            stage('Build') {
                parallel {
                    stage('MacOS') {
                        sh label: 'make darwin', script: 'make darwin'
                    }
                    stage('Linux') {
                        sh label: 'make linux', script: 'make linux'
                    }
                    stage('Windows') {
                        sh label: 'make win', script: 'make win'
                    }
                }
            }
        }
    }
}