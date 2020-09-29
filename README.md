# Pi Webcam
Web service interface written in go to capture images of a raspberry camera module. Implemented on top off [raspistill](https://github.com/raspberrypi/userland/blob/master/host_applications/linux/apps/raspicam/RaspiStill.c) and [raspiyuv](https://github.com/raspberrypi/userland/blob/master/host_applications/linux/apps/raspicam/RaspiStillYUV.c) commands. Please have a look at [main.go](main.go) for a list of available query parameters.

### Typical usage
```
# Get a jpg image over the still interface
wget --user webcam --password webcam -O test.jpg  http://localhost/still?timeout=1000

# Get yuv image as raw
wget --user webcam --password webcam -O test.bin  http://localhost/yuv?timeout=1000&awb=off
```

### Installing
- Import the repository key  
`wget -O - http://www.bayceer.uni-bayreuth.de/repos/apt/conf/bayceer_repo.gpg.key |apt-key add -`
- Add the BayCEER Debian repository  
echo "deb http://www.bayceer.uni-bayreuth.de/repos/apt/ $(lsb_release -c -s) main" | tee /etc/apt/sources.list.d/bayceer.list
- Update your repository cache  
`apt-get update`
- Install the package  
`apt-get install pi-webcam`

## Authors 
* **Oliver Archner** - *Developer* - [BayCEER, University of Bayreuth](https://www.bayceer.uni-bayreuth.de)

## History

### Version 1.0.0, Sep, 2020
- Initial release with basic capture capabilities

## License
GNU GENERAL PUBLIC LICENSE, Version 3, 29 June 2007