pipeline {
  agent {
    dockerfile {
      filename 'Dockerfile'
    }

  }
  stages {
    stage('install') {
      steps {
        sh '''make clean
make build
make test'''
      }
    }

  }
}