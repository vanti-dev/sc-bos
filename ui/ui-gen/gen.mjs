import {execSync} from 'child_process';
import fs from 'fs';

const path = fs.readdirSync('.');
console.log('paths', path);
const protoFiles = fs.readdirSync('../../proto');

execSync('protoc -I ../../proto --js_out=import_style=commonjs:src --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:src ' + protoFiles.join(' '));
