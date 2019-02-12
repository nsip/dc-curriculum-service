
cd tools; go build release.go; cd ..
./tools/release dc-curriculum-service > version/version.go
sh build.sh
./tools/release dc-curriculum-service dc-curriculum-service-Mac.zip build/dc-curriculum-service-Mac.zip
./tools/release dc-curriculum-service dc-curriculum-service-Win64.zip build/dc-curriculum-service-Win64.zip
./tools/release dc-curriculum-service dc-curriculum-service-Win32.zip build/dc-curriculum-service-Win32.zip
./tools/release dc-curriculum-service dc-curriculum-service-Linux64.zip build/dc-curriculum-service-Linux64.zip
./tools/release dc-curriculum-service dc-curriculum-service-Linux32.zip build/dc-curriculum-service-Linux32.zip
