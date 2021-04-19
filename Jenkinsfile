pipeline {
    agent { docker { image 'golang:1.16.3' } }
    environment {
        XDG_CACHE_HOME = '/tmp/.cache'
    }
    stages {
        stage('lint') {
            steps {
                sh './setup-dev-env.sh'
                sh 'bin/golangci-lint run'
            }
        }
        stage('test') {
            steps {
                sh 'go test ./...'
            }
        }
    }
}