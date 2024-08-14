import {execSync} from 'child_process';
import fs from 'fs';
import {replaceInFileSync} from 'replace-in-file';

const protoFiles = fs.readdirSync('../../proto')
    .filter(f => f.endsWith('.proto'))
    .map(f => f);
const protocPluginOpts = '--js_out=import_style=commonjs:proto --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:proto';

const out = execSync(`protomod protoc -- -I../../proto ${protocPluginOpts} ${protoFiles.join(' ')}`);
console.log(out.toString());

// update the generated files to replace
// `require('../../../traits/*_pb.js');`
// with `require('@smart-core-os/sc-api-grpc-web/traits/*_pb.js')`

// replace .js imports
replaceInFileSync({
  files: ['proto/**/*_pb.js'],
  from: /require\('(?:\.\/)+((?:traits|types|info)\/.+_pb.js)'\)/g,
  to: `require('@smart-core-os/sc-api-grpc-web/$1')`
});
// replace .d.ts imports
replaceInFileSync({
  files: ['proto/**/*_pb.d.ts'],
  from: /from '(?:\.\/)+((?:traits|types|info)\/.+_pb)'/g,
  to: `from '@smart-core-os/sc-api-grpc-web/$1'`
});
