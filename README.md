tongs
================================
Utility for creating code reviews quickly
in [Crucible](https://www.atlassian.com/software/crucible/overview) based on predefined templates.

Usage: 
    
    tongs [OPTION] [TEMPLATE] [REVIEW-ID]

Options:
    
    setup
    token
    templates
    templates [TEMPLATE]
    create [TEMPLATE]
    update [TEMPLATE] [REVIEW-ID]

Examples
--------------------------------

    tongs create [TEMPLATE]

    Example: tongs create javateam

Creates a new code review draft for the users in the given config section.
The user then just has to go to Crucible, add the desired revisions
to the review, and start it.

    update [TEMPLATE] [REVIEW-ID]
    
    Example: tongs update javateam CODEREVIEW-1001

Updates an existing code review with the users from the config template
provided. You must be the author of the code review in order to run this
command. Users that were already on the review will be ignored, allowing
the creation of overlapping groups as needed.

Installing
--------------------------------
To install Tongs, simply download the latest Tongs binary file for your 
system from the [Releases](https://github.com/southpawlife/tongs/releases) page and place it 
somewhere on your PATH. To update tongs, simply replace the existing binary file on your system.

Building From Source
--------------------------------
Tongs is written in [Go](http://golang.org/). Go can be compiled on any archecture as long as Go is installed.
([Download Go Here](https://code.google.com/p/go/wiki/Downloads) or 'brew install go' on Mac)

Build Tongs by running:
    
    go build tongs.go
    
This will create an executable in the source directory. Move this tongs executable to a location that is in your PATH so that it can be run from anywhere.

Configuration
---------------------

In order to use tongs you will need to create a tongs.cfg file that
exists in your home directory on your system. This has been verified
on Windows and Mac to write to the correct directory.
To generate a template file in your home directory, run setup.

    tongs setup

This will first prompt for your Crucible base url, username, and password.
The Crucible base url can be found by typing the address of the Crucible
server in a browser, and watching what url it redirects to. 
(i.e. http://my.crucible.com might resolve to http://my.crucible.com/viewer)
All trailing slashes should be omitted.

(See 'Crucible Token' section below for more information.)

Below is what your tongs.cfg file should look like. You can add as many config sections
as you need, as long as the default and settings are included. Not all fields are required 
for all sections but you can currently use project-key, duration, reviewers and title as you 
see fit. There is no restriction on spacing around the commas and equal signs. Also titles 
should be written without quotes of any kind.

Note that you must have at least a default project-key to create reviews. If this project 
key is invalid for the Crucible instance you are using, the review will not be created. 
    
Example tongs.cfg:

    [default]                       
    project-key=PROJECT-KEY                     <--- Make sure to set the correct key
    duration=3
    reviewers=userid1,userid9,userid12

    [my-team]                                   <--- Change this to anything you want
    project-key=OTHER-KEY           
    reviewers=userid6,userid1,userid2,userid3,userid4       
    title=My Team Code Review Template

    [java]                                      <--- Create as many sections as you need
    reviewers=userid6

    [settings]                                  
    crucible-baseurl=http://crucible.company.com/basepath
    crucible-token=                        
    

Remote Configuration
---------------------

In many cases it makes sense to have a Tongs configuration defined once for group of users. 
Tongs supports this by using the optional 'url' option in the local configuration.

Example tongs.cfg:

    [default]                       
    project-key=PROJECT-KEY                     
    duration=3

    [my-team]                               <--- Template name must match remote template name
    url = https://mydomain.com/myconfig/    <--- URL pointing to a remote configuration file

    [settings]                                  
    crucible-baseurl = http://crucible.company.com/basepath
    crucible-token = xw027966:41556:8ff7ddfgdfgl0ec44234fg9e2cb6cb                        

Example remote configuration at https://mydomain.com/myconfig/:
    
    [my-team]                               <--- Remote template name 
    project-key = OTHER-KEY           
    reviewers = userid6,userid1,userid2,userid3,userid4       
    title = My Team Code Review Template

When setting up a remote template, the local template name must match the remote template name 
you wish to consume. 

The order in which the templates are honored is as follows:

1. Load from the remote configuration section with the template name given
2. Load from the remote configuration default section
3. Load from the local configuration section with the template name given
4. Load from the local configuration default section

Crucible Token
-------------------------------
Use the 'token' command to reset your token without re-entering your crucible url.

    tongs token

This will ask you for your username and password and if it is able to successfully
able to connect it will save the token to the config file.

No password is being stored in this application, however currently there is no way to mask
the password entry when running the 'token' and 'setup' commands. As such if you would rather 
get your token by some other means, here are links to the Crucible API's and REST services. 
Note that tokens don't expire so you should not have to enter your password very often.

Example Token:
    
    xw027966:41556:8ff7ddfgdfgl0ec44234fg9e2cb6cb

https://developer.atlassian.com/display/FECRUDEV/Authenticating+REST+Requests

https://docs.atlassian.com/fisheye-crucible/latest/wadl/crucible.html#rest-service:auth-v1:login


Future Development Ideas
--------------------------------------------

 - [x] Github hosted config files for teams
 - [ ] Github and Subversion integration
 - [x] Ability to add reviewers to existing reviews
 - [ ] Ability to see list of current reviews
 - [x] Ability to see list of templates
 - [ ] Password Masking

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

