# tesla
Tesla Take Home Assignment

URL: https://bruce.teitoku.net/

Hosted on a server
Written in GoLang
No front-end API just a server where CURL Commands or Postman can be run on

NOTE ON CHECK FOR TEMPERATURE:
Not sure what is happening for the string check on 'Temperature', but I'm guessing for some reason the terminal on my windows machine I've installed doesn't handle quotes properly. I've made a check which handles both cases, but on certain terminals it might let in more Temperature strings than it should, if that's the case I could tighten the check for standard terminal, but will leave as is unless it is a huge problem.
