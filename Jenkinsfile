pipeline {
    agent { docker { image 'golang:1.16.3' } }
    stages {
        stage('build') {
            steps {
                sh 'go version'
            }
        }
    }
}