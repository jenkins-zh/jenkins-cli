pipeline {
    agent {
        label 'golang'
    }

    stages {
        stage('Build') {
            container('golang') {
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