def linuxBuildTargets = [
    "386",
    "amd64",
    "arm",
    "arm64"
]

pipeline {
    agent { docker { image 'golang:1.16.3' } }
    environment {
        XDG_CACHE_HOME = '/tmp/.cache'
    }
    stages {
        stage('Validate code') {
            parallel {
                stage('Lint') {
                    steps {
                        sh './setup-dev-env.sh'
                        sh 'bin/golangci-lint run'
                    }
                }
                stage('Test') {
                    steps {
                        sh 'go test ./...'
                    }
                }
            }
        }
        stage('Build code') {
            parallel {
                stage('Linux') {
                    steps {
                        parallel script {
                            linuxBuildTargets.each { target ->
                                sh "GOOS=linux GOARCH=${target} go build -o eve-industry-linux-${target} ./..."
                            }
                        }
                    }
                }
            }
        }
    }
}