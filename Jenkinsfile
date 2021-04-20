def linuxBuildTargets = [
    "386",
    "amd64",
    "arm",
    "arm64"
]

def windowsBuildTargets = [
    "386",
    "amd64"
]

def builtImage

pipeline {
    agent none
    environment {
        XDG_CACHE_HOME = '/tmp/.cache'
    }
    stages {
        stage('Validate code') {
            parallel {
                stage('Lint') {
                    agent { docker { image 'golang:1.16.3' } }
                    steps {
                        sh './setup-dev-env.sh'
                        sh 'bin/golangci-lint run'
                    }
                }
                stage('Test') {
                    agent { docker { image 'golang:1.16.3' } }
                    steps {
                        sh 'go test ./...'
                    }
                }
            }
        }
        stage('Build code') {
            agent { docker { image 'golang:1.16.3' } }
            steps {
                script {
                    def builds = [:]
                    linuxBuildTargets.each { target ->
                        builds["linux-"+target] = {
                            stage("Build Linux - ${target}") {
                                sh "GOOS=linux GOARCH=${target} go build -o eve-industry-linux-${target} ./..."
                            }
                        }
                    }
                    windowsBuildTargets.each { target ->
                        builds["windows-"+target] = {
                            stage("Build Windows - ${target}") {
                                sh "GOOS=windows GOARCH=${target} go build -o eve-industry-windows-${target} ./..."
                            }
                        }
                    }
                    parallel builds
                }
            }
        }
        stage('Build docker image') {
            agent {
                label 'docker-build'
            }
            steps {
                script {
                    builtImage = docker.build("eve-industry:${env.BUILD_ID}")
                }
            }
        }
        stage('Acceptance test') {
            agent {
                label 'docker-build'
            }
            tools {
                go '1.16.3'
            }
            steps {
                script {
                    builtImage.withRun('-p 18080:18080') {
                        sh 'go test --tags=acceptance ./...'
                    }
                }
            }
        }
    }
    post {
        always {
            node('docker-build') {
                sh "docker rmi ${builtImage.id}"
            }
        }
    }
}