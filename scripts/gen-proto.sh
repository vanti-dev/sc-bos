#!/usr/bin/env bash

pushd ../pkg/gen
go generate
if [ -d ../trait/gen ]; then
  # old version fixes to files
  # Remove imports of the gen package: `gen "github.com/smart-core-os/sc-bos/pkg/gen"`.
  sed -i.bak '/gen "github.com\/smart-core-os\/sc-bos\/pkg\/gen"/d' ../trait/gen/*.pb.go
  # Remove any `gen.` qualifiers on statements: `gen.Foo` => `Foo`.
  sed -i.bak 's/\bgen\.//g' ../trait/gen/*.pb.go
  sed -i.bak 's/[[:<:]]gen\.//g' ../trait/gen/*.pb.go
  mv -f ../trait/gen/*.pb.go .
fi
if [ -d ../trait/genpb ]; then
  # new version fixes to files
  # Remove imports of the gen package: `gen "github.com/smart-core-os/sc-bos/pkg/gen"`.
  sed -i.bak '/gen "github.com\/smart-core-os\/sc-bos\/pkg\/gen"/d' ../trait/genpb/*.pb.go
  # Remove any `gen.` qualifiers on statements: `gen.Foo` => `Foo`.
  sed -i.bak 's/\bgen\.//g' ../trait/genpb/*.pb.go
  sed -i.bak 's/[[:<:]]gen\.//g' ../trait/genpb/*.pb.go
  # Rename the package from genpb to gen.
  sed -i.bak 's/package genpb/package gen/g' ../trait/genpb/*.pb.go
  mv -f ../trait/genpb/*.pb.go .
fi
rm -rf ../trait
popd

pushd ../ui/ui-gen
yarn run gen
popd
