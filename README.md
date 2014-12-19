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

This will first prompt for your Crucible Base URL. This can be found by typing the address of the Crucible
server in a browser, and watching what url it redirects to. (i.e. http://my.crucible.com might resolve to http://my.crucible.com/viewer)
All trailing slashes should be omitted.

Here is what your tongs.cfg file should look like. You can add as many config sections
as you need, as long as the default and settings are included. Not all fields are required for all sections
but you can currently use project-key, duration, reviewers and title as you see fit. There is no restriction on spacing
around the commas and equal signs. Also titles should be written without quotes of any kind.

Note: you must have at least a default project-key to create reviews. If this project key is invalid for the Crucible instance you are using, the review will not be created. Also, for the baseurl, add any necessary URL paths to the base URL such as /viewer if needed, just omit the trailing slash.

    [default]                       
    project-key=PROJECT-KEY                     <--- Make sure to set the correct key
    duration=3
    reviewers=userid1,userid9,userid12

    [my-team]                                   <--- Change this to anything you want...
    project-key=OTHER-KEY           
    reviewers=userid6,userid1,userid2,userid3,userid4       
    title=My Team Code Review Template

    [java]                                      <--- Create as many sections as you need.
    reviewers=userid6

    [settings]                                  
    crucible-baseurl=http://crucible.company.com/basepath
    crucible-token=                        
    

This will ask you for your username and password and if it is able to successfully
able to connect it will save the username and token to the config file.

Notes about the Crucible Token
-------------------------------
No password is being stored in this application, however currently there is no way to mask
the password entry in this application when running the --get-token command. As such if you
would rather get your token by some other means, here are links to the Crucible API's and
REST services. Note that tokens don't expire so you should not have to enter your password very often.

Example Token:
    
    xw027966:41556:8ff7ddfgdfgl0ec44234fg9e2cb6cb

https://developer.atlassian.com/display/FECRUDEV/Authenticating+REST+Requests

https://docs.atlassian.com/fisheye-crucible/latest/wadl/crucible.html#rest-service:auth-v1:login





Future Development Ideas
--------------------------------------------

 - [ ] Github hosted config files for teams
 - [ ] Github and Subversion integration
 - [x] Ability to add reviewers to existing reviews
 - [ ] Ability to see list of current reviews
 - [x] Ability to see list of templates
 - [ ] Password Masking
 - [ ] Refactor / Code Cleanup (ongoing)

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

