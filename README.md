tongs
================================
Command line tool for interacting with [Atlassian Crucible](https://www.atlassian.com/software/crucible/overview). Solves the problem of 
remembering who you need to put on that code review and saves you time 
by adding everyone automatically.
Usage
--------------------------------
    tongs create my-team --title "The new awesome server code"

Creates a new code review draft for the users in the ruby reviewers section,
with the title 'The new awesome server code'. The end user then just has to go to
Crucible, add the desired revisions to the review and start it. 

    tongs create

Creates a new code review draft for the users in the default reviewers section,
with the default title. The end user then just has to go to
Crucible, add the desired revisions to the review and start it. 

Other Usage Examples:

    tongs --create-config
    tongs create
    tongs create <template-name>
    tongs create <template-name> --title "Special Title"
    tongs create --title "Default Review, With a not so default title"
    tongs help

Setup
--------------------------------
Tongs is written in [Golang](http://golang.org/). Golang code can be compiled on 
Windows, Mac, Linux etc.. just as long as Go is installed.
([Download Go here](https://code.google.com/p/go/wiki/Downloads) or 'brew install go' on a mac) 
With go installed build tongs by calling:
    
    go build tongs.go
    
This will create a new executable in the same directory. You can then move this executable (file called tongs WITHOUT the .go extention) to a folder location that is in your PATH so that it can be run from anywhere.

If you just want to run the code without building it, you can call:

    go run tongs.go

Configuration
---------------------

In order to use tongs you will need to create a .tongs_config file that 
exists in your home directory on your system. This has currently been verified
on Windows and OSX to write to the correct directory for the respective OS.
To generate a template file in your home directory, run:

    tongs --create-config

Here is what your .tongs_config file should look like. You can add as many config sections 
as you need, as long as the default and settings. Not all fields are required for all sections
but you can currently use project-key, duration, reviewers and title as you see fit. There is no restriction on spacing 
around the commas and equal signs. Also titles should be written without quotes of any kind. 

Note: you must have at least a default project-key to create reviews. If this project key is invalid for the Crucible instance you are using, the review will not be created. Also, for the baseurl, add any nessesary url paths to the base url such as /viewer if needed, just omit the trailing slash.

    [default]
    project-key=PROJECT-KEY
    duration=3
    reviewers=userid1,userid9,userid12

    [my-team]
    project-key=OTHER-KEY
    reviewers=userid6,userid1,userid2,userid3,userid4
    title=My Team Code Review Template

    [java]
    reviewers=userid6

    [settings]
    crucible-baseurl=http://crucible.mycompany.com
    crucible-username=
    crucible-token=
    
Notes about the Crucible Token
-------------------------------

The Crucible token and username will be auto populated on first connect, when you provide your username 
and password. No password is being stored in this application, however currently there is no way to mask 
the password entry in this application. As such if you would rather get your token by some other means, 
here are links to the Crucible API's and REST services. Note that tokens dont expire so you should not 
have to enter your password very often.

Example Token:
    
    xw027966:41556:8ff7ddfgdfgl0ec44234fg9e2cb6cb

https://developer.atlassian.com/display/FECRUDEV/Authenticating+REST+Requests

https://docs.atlassian.com/fisheye-crucible/latest/wadl/crucible.html#rest-service:auth-v1:login





Future Development Ideas
--------------------------------------------

 - [ ] Central, Github hosted, config files for teams. Teams could create one central list to manage who is on various code reviews for different languages.
 - [ ] Github and SVN intergration
 - [ ] Ability to Add reviewers to existing reviews
 - [ ] Ability to see list of current reviews
 - [ ] Refactor code
 - [ ] MASK PASSWORDS

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

