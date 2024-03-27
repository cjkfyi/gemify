#!/bin/bash

echo -e "\n🚀  Launching 'scripts/scrub.sh'\n" 

PB_GO_FILES="backend/api/gen/*.go"
rm -rf $PB_GO_FILES 
echo -e "💦  Cleared up Go Protobuf files!\n"

PB_JS_FILES="extension/src/gen/*.js"
rm -rf $PB_JS_FILES 
echo -e "💦  Cleared up JS Protobuf files!\n"

GEN_JS_FILE="extension/build/*.js"
rm -rf $GEN_JS_FILE 
echo -e "💦  Cleared generated JS file!\n"

BUILT_EXT="extension/dist/*.vsix"
rm -rf $BUILT_EXT
echo -e "💦  Cleared up built ext pkgs!\n" 

EXT_MODULES="extension/node_modules"
rm -rf $EXT_MODULES 
echo -e "💦  Cleared up ext 'node_modules'!\n"

ROOT_MODULES="node_modules"
rm -rf $ROOT_MODULES 
echo -e "💦  Cleared up root 'node_modules'!\n"

echo -e "👻  Ending 'scripts/scrub.sh'\n" 
