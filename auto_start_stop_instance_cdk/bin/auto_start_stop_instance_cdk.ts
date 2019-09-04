#!/usr/bin/env node
import 'source-map-support/register';
import cdk = require('@aws-cdk/core');
import { AutoStartStopInstanceCdkStack } from '../lib/auto_start_stop_instance_cdk-stack';



const util = require('util');
const exec = util.promisify(require('child_process').exec);

async function deploy(){
    await exec('go get -v -t -d ./lambdaSources/function/... && GOOS=linux GOARCH=amd64 go build -o ./lambdaSources/source/main ./lambdaSources/function/**.go')
    const app = new cdk.App();
    new AutoStartStopInstanceCdkStack(app, 'AutoStartStopInstanceCdkStack');
    app.synth()
    await  exec('rm ./lambdaSources/source/main')
}

deploy()