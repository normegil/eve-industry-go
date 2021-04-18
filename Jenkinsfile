pipeline {
    agent { docker { image 'golang:1.16.3' } }
    stages {
        stage('lint') {
            steps {
                sh './setup-dev-env.sh'
                sh 'bin/golangci-lint --version'
            }
        }
    }
}