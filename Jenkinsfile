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
        IPV4_REGEX = '^([0-9]{1,3}\\.){3}[0-9]{1,3}$'

        XDG_CACHE_HOME = '/tmp/.cache'
        VM_IMAGE_NAME = 'evevulcan'
        SERVER_NAME = 'evevulcan'
        // Openstack
        OS_AUTH_URL = 'https://auth.cloud.ovh.net/v3/'
        OS_IDENTITY_API_VERSION=3
        OS_USER_DOMAIN_NAME="Default"
        OS_PROJECT_DOMAIN_NAME="Default"
        OS_REGION_NAME="GRA5"
        // Test server parameters
        VM_TEST_SERVER_NAME = 'evevulcan-ci-test'
        VM_TEST_SERVER_FLAVOR = 's1-2'
    }
    stages {
        stage('Validate code') {
            parallel {
                stage('Lint Go code') {
                    agent any
                    tools {
                        go '1.16.3'
                    }
                    steps {
                        sh './setup-dev-env.sh'
                        sh 'bin/golangci-lint run'
                    }
                }
                stage('Lint Vue code') {
                    agent any
                    tools {
                        node '14.16.1'
                    }
                    steps {
                        dir("ui/web") {
                            sh 'npm install'
                            sh 'npm run lint'
                        }
                    }
                }
                stage('Lint Packer configuration') {
                    agent any
                    steps {
                        dir(".deployment") {
                            withCredentials([usernamePassword(credentialsId: 'OpenstackOVH', usernameVariable: 'OS_USERNAME', passwordVariable: 'OS_PASSWORD')]) {
                                sh 'packer validate openstack.pkr.hcl'
                            }
                        }
                    }
                }
                stage('Lint Ansible configurations') {
                    agent any
                    steps {
                        sh 'ansible-lint */***'
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
        stage('Generate code') {
            agent any
            tools {
                go '1.16.3'
                node '14.16.1'
            }
            steps {
                sh 'go generate ./...'
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
                                sh "GOOS=linux GOARCH=${target} go build -o evevulcan-linux-${target} ./..."
                                stash name: "evevulcan-linux-${target}", allowEmpty: false, includes: "evevulcan-linux-${target}"
                            }
                        }
                    }
                    windowsBuildTargets.each { target ->
                        builds["windows-"+target] = {
                            stage("Build Windows - ${target}") {
                                sh "GOOS=windows GOARCH=${target} go build -o evevulcan-windows-${target}.exe ./..."
                                stash name: "evevulcan-windows-${target}", allowEmpty: false, includes: "evevulcan-windows-${target}.exe"
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
                    builtDockerImage = docker.build("normegil/evevulcan:${env.BUILD_ID}")
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
                                sh 'mkdir -p bin/'
                                sh 'wget https://github.com/github-release/github-release/releases/download/v0.10.0/linux-amd64-github-release.bz2 -O bin/github-release.bz2'
                                sh 'bzip2 -d bin/github-release.bz2 || true'
                                sh 'chmod +x bin/github-release'

                                sh 'ls'

                                sh 'bin/github-release delete --user ${GITHUB_USER} --repo evevulcan --tag latest || true'
                                sh 'bin/github-release release --user ${GITHUB_USER} --repo evevulcan --tag latest --name latest'
                                linuxBuildTargets.each { target ->
                                    unstash name: "evevulcan-linux-${target}"
                                    sh "bin/github-release upload --user ${GITHUB_USER} --repo evevulcan --tag latest --name evevulcan-linux-${target} --file evevulcan-linux-${target}"
                                }
                                windowsBuildTargets.each { target ->
                                unstash name: "evevulcan-windows-${target}"
                                    sh "bin/github-release upload --user ${GITHUB_USER} --repo evevulcan --tag latest --name evevulcan-windows-${target}.exe --file evevulcan-windows-${target}.exe"
                                }
                            }
                        }
                    }
                }
            }
        }
        stage('Launch heavy tests ?') {
            agent none
            steps {
                input message: "Launch integration tests & performance tests ?"
                milestone 1
            }
        }
        stage('Create VM Image') {
            agent any
            steps {
                sh 'ansible-galaxy collection install community.general'
                sh 'ansible-galaxy collection install community.docker'
                dir(".deployment") {
                    withCredentials([usernamePassword(credentialsId: 'OpenstackOVH', usernameVariable: 'OS_USERNAME', passwordVariable: 'OS_PASSWORD'), file(credentialsId: 'AnsibleVaultPasswordFile', variable: 'AUSIBLE_VAULT_PASSWORD_PATH')]) {
                        sh "packer build -var=\"image_name=${env.VM_IMAGE_NAME}-${env.BUILD_NUMBER}\" -var=\"vault_password_file=${env.AUSIBLE_VAULT_PASSWORD_PATH}\" openstack.pkr.hcl"
                    }
                }
            }
        }
        stage('Test VM: Launch') {
            agent any
            steps {
                withCredentials([usernamePassword(credentialsId: 'OpenstackOVH', usernameVariable: 'OS_USERNAME', passwordVariable: 'OS_PASSWORD')]) {
                    sh "openstack server create --flavor ${env.VM_TEST_SERVER_FLAVOR} --image ${env.VM_IMAGE_NAME}-${env.BUILD_NUMBER} --wait ${env.VM_TEST_SERVER_NAME}-${env.BUILD_NUMBER}"
                }
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
        stage('Test VM: Delete') {
            agent any
            steps {
                script {
                    withCredentials([usernamePassword(credentialsId: 'OpenstackOVH', usernameVariable: 'OS_USERNAME', passwordVariable: 'OS_PASSWORD')]) {
                        sh "openstack server delete ${env.VM_TEST_SERVER_NAME}-${env.BUILD_NUMBER}"
                    }
                }
            }
        }
        stage('Release to production ?') {
            agent none
            steps {
               input message: "Release new code to production ?"
               milestone 2
            }
        }
        stage('Release') {
            agent any
            steps {
                script {
                    withCredentials([usernamePassword(credentialsId: 'OpenstackOVH', usernameVariable: 'OS_USERNAME', passwordVariable: 'OS_PASSWORD'), sshUserPrivateKey(credentialsId: 'JenkinsSSHKey', keyFileVariable: 'JENKINS_PRIVATE_KEY')]) {
                        sh "openstack server create --key-name JENKINS_KEY --flavor s1-2 --image ${env.VM_IMAGE_NAME}-${env.BUILD_NUMBER} --wait ${env.SERVER_NAME}-${env.BUILD_NUMBER}"

                        PRODUCTION_EXIST_STR = sh (
                            script: "openstack server list | grep ${env.SERVER_NAME}-production | wc -l",
                            returnStdout: true
                        ).trim()
                        echo PRODUCTION_EXIST_STR
                        PRODUCTION_EXIST = PRODUCTION_EXIST_STR as Integer

                        STAGING_IP = sh (
                            script: ".deployment/openstack-server-private-ipv4.sh ${env.SERVER_NAME}-${env.BUILD_NUMBER}",
                            returnStdout: true
                        ).trim()
                        assert STAGING_IP =~ IPV4_REGEX

                        PRODUCTION_IP = ""
                        if (PRODUCTION_EXIST > 0) {
                            PRODUCTION_IP = sh (
                                script: ".deployment/openstack-server-private-ipv4.sh ${env.SERVER_NAME}-production",
                                returnStdout: true
                            ).trim()
                            assert PRODUCTION_IP =~ IPV4_REGEX
                        }

                        dir(".deployment/ansible/") {
                            sh "ansible-playbook release.yml -i inventory.yml --extra-vars \"release_env_ip=${STAGING_IP}\""
                        }

                        if (PRODUCTION_EXIST > 0) {
                            // Wait for no connections to current production machine
                            NUMBER_OF_CONNECTIONS = sh (
                                script: "ssh -i ${env.JENKINS_PRIVATE_KEY} ubuntu@${PRODUCTION_IP} netstat -an | grep -E \":443|:80\" | grep -v \":8080\" | grep -E \"ESTABLISHED|CLOSING\" | wc -l",
                                returnStdout: true
                            ).trim()
                            while(NUMBER_OF_CONNECTIONS > 0) {
                                NUMBER_OF_CONNECTIONS = sh (
                                    script: "ssh -i ${env.JENKINS_PRIVATE_KEY} ubuntu@${PRODUCTION_IP} netstat -an | grep -E \":443|:80\" | grep -v \":8080\" | grep -E \"ESTABLISHED|CLOSING\" | wc -l",
                                    returnStdout: true
                                ).trim()
                                sleep (time:1)
                            }

                            // Delete previous production
                            sh "openstack server delete ${env.SERVER_NAME}-production"
                            sh "openstack image set --property name=${env.VM_IMAGE_NAME}-previous ${env.VM_IMAGE_NAME}-production"
                        }

                        // Promote staging to production
                        sh "openstack server set --name ${env.SERVER_NAME}-production ${env.SERVER_NAME}-${env.BUILD_NUMBER}"
                        sh "openstack image set --property name=${env.VM_IMAGE_NAME}-production ${env.VM_IMAGE_NAME}-${env.BUILD_NUMBER}"
                    }
                }
            }
        }
    }
    post {
        always {
            node('docker-build') {
                sh "docker rmi ${builtDockerImage.id}"
                script {
                    withCredentials([usernamePassword(credentialsId: 'OpenstackOVH', usernameVariable: 'OS_USERNAME', passwordVariable: 'OS_PASSWORD')]) {
                        try {
                            // Delete server if release step failed
                            sh "openstack server delete ${env.SERVER_NAME}-${env.BUILD_NUMBER}"
                        } catch (Exception e) {
                            echo e.getMessage()
                        }
                        try {
                            // Delete test server
                            sh "openstack server delete ${env.VM_TEST_SERVER_NAME}-${env.BUILD_NUMBER}"
                        } catch (Exception e) {
                            echo e.getMessage()
                        }
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
        unstable {
            mail from: 'jenkins@ci.normegil.be', to: "mail@normegil.be", subject: "Unstable: Job ${JOB_NAME} - ${env.BUILD_NUMBER}", body: "The job ${JOB_NAME} (${env.BUILD_NUMBER}) is unstable. Please check ${env.JENKINS_URL}/blue/organizations/jenkins/${JOB_NAME}/detail/${env.BRANCH_NAME}/${env.BUILD_NUMBER}/pipeline."
        }
        failure {
            mail from: 'jenkins@ci.normegil.be', to: "mail@normegil.be", subject: "Error: Job ${JOB_NAME} - ${env.BUILD_NUMBER}", body: "The job ${JOB_NAME} (${env.BUILD_NUMBER}) is in error. Please check ${env.JENKINS_URL}/blue/organizations/jenkins/${JOB_NAME}/detail/${env.BRANCH_NAME}/${env.BUILD_NUMBER}/pipeline."
        }
    }
}