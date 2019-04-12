#/bin/bash

go build -o build/bin/application
cd build
zip -r ../app.zip *
