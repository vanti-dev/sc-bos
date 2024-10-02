@echo off
REM Build the Docker image for Windows
docker build -f Dockerfile-windows -t go-windows-builder .

REM Run the Docker container and copy the executable to the current directory
docker run --rm -v "%cd%:/output" go-windows-builder cmd /C "copy /app/windows-seeder.exe /output/windows-seeder.exe"

REM Execute the Windows binary with the provided configuration files
output\windows-seeder.exe --appconf C:\path\to\git\sc-bos\example\config\vanti-ugs\app.conf.json --sysconf C:\path\to\git\sc-bos\example\config\vanti-ugs\system.conf.json
