import {execSync} from 'child_process';
import fs from 'fs';
import replace from 'replace-in-file';

const protoFiles = fs.readdirSync('../../proto')
    .map(f => 'proto/' + f);
const protocPluginOpts = '--js_out=import_style=commonjs:. --grpc-web_out=import_style=commonjs+dts,mode=grpcwebtext:.';

const out = execSync(`protomod protoc -- -I../.. ${protocPluginOpts} ${protoFiles.join(' ')}`);
console.log(out.toString());

// update the generated files to replace
// `require('../../../traits/*_pb.js');`
// with `require('@smart-core-os/sc-api-grpc-web/traits/*_pb.js')`

// replace .js imports
replace.sync({
  files: ['proto/**/*_pb.js'],
  from: /require\('(?:\.\.\/)+((?:traits|types|info)\/.+_pb.js)'\)/g,
  to: `require('@smart-core-os/sc-api-grpc-web/$1')`
});
// replace .d.ts imports
replace.sync({
  files: ['proto/**/*_pb.d.ts'],
  from: /from '(?:\.\.\/)+((?:traits|types|info)\/.+_pb)'/g,
  to: `from '@smart-core-os/sc-api-grpc-web/$1'`
});
