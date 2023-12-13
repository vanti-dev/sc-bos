#!/usr/bin/env bash

pushd ../pkg/gen
go generate
# Remove imports of the gen package: `gen "github.com/vanti-dev/sc-bos/pkg/gen"`.
sed -i.bak '/gen "github.com\/vanti-dev\/sc-bos\/pkg\/gen"/d' ../trait/gen/*.pb.go
# Remove any `gen.` qualifiers on statements: `gen.Foo` => `Foo`.
sed -i.bak 's/\bgen\.//g' ../trait/gen/*.pb.go
sed -i.bak 's/[[:<:]]gen\.//g' ../trait/gen/*.pb.go
mv -f ../trait/gen/*.pb.go .
rm -rf ../trait
popd

pushd ../ui/ui-gen
yarn run gen
popd
