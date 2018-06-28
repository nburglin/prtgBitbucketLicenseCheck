## Description

This is a simple Go program designed to leverage PRTG's monitoring system and Atlassian Bitbucket's API.
It returns the count of how many licenses are currently being consumed in your bitbucket server.

My company already uses PRTG to monitor systems and processes and we have the SSH Script Advanced Sensor available to us, so I decided to
leverage this system to alert our team if we are getting close to the license max

I am also new to Go, so this is my first Go program written. There probably are better or more efficient ways to do this,
but this does what I needed it to do so I figured I would share it out there to possibly help others with a similar desire

## Access Needed

 * Access to create a new sensor in PRTG for your Bitbucket sensor (or anywhere that has access to your bitbucket url really)
 * Access to copy the executable and settings.json to your PRTG scriptsxml directory

## Configuration

The program expects a 'settings.json' file to be located in the same directory as the executable to provide the following values

 * Username - Username of user with access to licensing info on Bitbucket
 * Password - Password to use the API
 * BaseUrl - The Base API url of your Bitbucket server.

It uses Basic Authentication to authenticate with your bitbucket server url.

Example file contents:

```json
{
    "username": "nburglin",
    "password": "somepassword",
    "baseurl": "https://bitbucket.mycompany.com:7893/rest/api/1.0/"
}
```

There is a samle settings.json in this repo that you can overwrite to use. Special characters like '!' work fine with basic auth

##Create Executable

It's expected that you have GO installed locally. You can execute the command below to retrieve the source

```bash
go get github.com/nburglin/prtgBitbucketLicenseCheck
```

This will put the executable in your $GOPATH/bin, however you must first make sure you build an executable that will be able to run
on your remote server. For instance, if you are going to be running on a linux host with AMD, your build command will look like this:

```bash
env GOOS=linux GOARCH=amd64 go install github.com/nburglin/prtgBitbucketLicenseCheck
```

In the example above, the executable you need will be found in $GOPATH/bin/linux_amd64/prtgBitbucketLIcenseCheck

Now you simply copy this executable to your remote server to PRTG's scriptsxml directory, along with 
a settings.json mentioned above in the Configuration section. Once this is done, you can create the new SSH Script
Advanced Sensor in PRTG and point it to this executable.


## Flow

This serves a pretty specific purpose, so the flow is fairly simple.

1. Read in config file
2. Set up a Go http client
3. Perform a GET against BASEURL with "admin/license" at the end
4. Print out the current license count inside of some HTML tags that PRTG expects. The channel name is just hardcoded to "Bitbucket License Count"

## References

More info on PRTG's SSH Script Advanced Sensor found here:
 - [PRTG SSH Script Advanced Sensory](https://blog.paessler.com/prtg-ssh-script-advanced-sensor)

Bitbucket license API (tested as of 5.9.1)
 - [Bitbucket API Info](https://docs.atlassian.com/bitbucket-server/rest/5.9.1/bitbucket-rest.html#idm93660011056)
