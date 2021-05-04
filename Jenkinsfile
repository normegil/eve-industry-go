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

def builtDockerImage

pipeline {
    agent none
    environment {
        XDG_CACHE_HOME = '/tmp/.cache'
        VM_IMAGE_NAME = 'eve-industry'
        // Openstack
        OS_AUTH_URL = 'https://auth.cloud.ovh.net/v3/'
        OS_IDENTITY_API_VERSION=3
        OS_USER_DOMAIN_NAME="Default"
    }
    stages {
        stage('Validate code') {
            parallel {
                stage('Lint') {
                    agent any
                    tools {
                        go '1.16.3'
                    }
                    steps {
                        sh './setup-dev-env.sh'
                        sh 'bin/golangci-lint run'
                    }
                }
                stage('Test') {
                    agent any
                    tools {
                        go '1.16.3'
                    }
                    steps {
                        sh 'go test ./...'
                    }
                }
            }
        }
        stage('Build code') {
            agent any
            tools {
                go '1.16.3'
            }
            steps {
                script {
                    def builds = [:]
                    linuxBuildTargets.each { target ->
                        builds["linux-"+target] = {
                            stage("Build Linux - ${target}") {
                                sh "GOOS=linux GOARCH=${target} go build -o eve-industry-linux-${target} ./..."
                                stash name: "eve-industry-linux-${target}", allowEmpty: false, includes: "eve-industry-linux-${target}"
                            }
                        }
                    }
                    windowsBuildTargets.each { target ->
                        builds["windows-"+target] = {
                            stage("Build Windows - ${target}") {
                                sh "GOOS=windows GOARCH=${target} go build -o eve-industry-windows-${target}.exe ./..."
                                stash name: "eve-industry-windows-${target}", allowEmpty: false, includes: "eve-industry-windows-${target}.exe"
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
                    builtDockerImage = docker.build("normegil/eve-industry:${env.BUILD_ID}")
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
                    builtDockerImage.withRun('-p 18080:18080') {
                        sh 'go test --tags=acceptance ./...'
                    }
                }
            }
        }
        stage('Publish artefacts') {
            parallel {
                stage('Publish docker') {
                    agent {
                        label 'docker-build'
                    }
                    steps {
                        script {
                            docker.withRegistry('https://index.docker.io/v1/', 'DockerHub') {
                                builtDockerImage.push('latest')
                            }
                        }
                    }
                }
                stage('Publish binaries') {
                    agent any
                    tools {
                        go '1.16.3'
                    }
                    steps {
                        withCredentials([usernamePassword(credentialsId: 'GithubToken', usernameVariable: 'GITHUB_USER', passwordVariable: 'GITHUB_TOKEN')]) {
                            script {
                                sh 'wget https://github.com/github-release/github-release/releases/download/v0.10.0/linux-amd64-github-release.bz2 -O bin/github-release.bz2'
                                sh 'bzip2 -d bin/github-release.bz2 || true'
                                sh 'chmod +x bin/github-release'

                                sh 'ls'

                                sh 'bin/github-release delete --user ${GITHUB_USER} --repo eve-industry-go --tag latest || true'
                                sh 'bin/github-release release --user ${GITHUB_USER} --repo eve-industry-go --tag latest --name latest'
                                linuxBuildTargets.each { target ->
                                    unstash name: "eve-industry-linux-${target}"
                                    sh "bin/github-release upload --user ${GITHUB_USER} --repo eve-industry-go --tag latest --name eve-industry-linux-${target} --file eve-industry-linux-${target}"
                                }
                                windowsBuildTargets.each { target ->
                                unstash name: "eve-industry-windows-${target}"
                                    sh "bin/github-release upload --user ${GITHUB_USER} --repo eve-industry-go --tag latest --name eve-industry-windows-${target}.exe --file eve-industry-windows-${target}.exe"
                                }
                            }
                        }
                    }
                }
            }
        }
        stage('Create VM Image') {
            agent any
            steps {
                withCredentials([usernamePassword(credentialsId: 'OpenstackOVH', usernameVariable: 'OS_USERNAME', passwordVariable: 'OS_PASSWORD')]) {
                    sh 'packer validate .deployment/openstack.pkr.hcl'
                    sh "packer build -var=\"image_name=${env.VM_IMAGE_NAME}-${env.BUILD_NUMBER}\" .deployment/openstack.pkr.hcl"
                }
            }
        }
        stage('Test VM Image') {
            agent none
            stages {
                stage('launch server with image') {
                    agent any
                    steps {
                        sh 'echo "launch server"'
                    }
                }
                stage('Integration tests') {
                    agent any
                    tools {
                        go '1.16.3'
                    }
                    steps {
                        sh 'go test --tags=integration ./...'
                    }
                }
                stage('Performance tests') {
                    agent any
                    steps {
                        sh 'echo "Performance tests"'
                    }
                }
            }
        }
        stage('Publish VM Image') {
            agent any
            steps {
                sh 'echo "Remove old staging openstack image"'
                sh 'echo "Rename openstack image to staging"'
            }
        }
        stage('Release') {
            agent any
            steps {
                sh 'echo "Create server from staging image"'
                sh 'echo "Switch Load balancer"'
                sh 'echo "Switch image names"'
                sh 'echo "Remove old server"'
            }
        }
    }
    post {
        always {
            node('docker-build') {
                sh "docker rmi ${builtDockerImage.id}"
                sh "echo 'Remove openstack running server'"
                script {
                    withCredentials([usernamePassword(credentialsId: 'OpenstackOVH', usernameVariable: 'OS_USERNAME', passwordVariable: 'OS_PASSWORD')]) {
                        try {
                            sh "openstack image delete ${env.VM_IMAGE_NAME}-${env.BUILD_NUMBER}"
                        } catch (Exception e) {
                            echo e.getMessage()
                        }
                    }
                }
                cleanWs()
            }
        }
    }
}