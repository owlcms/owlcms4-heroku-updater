# Cloud Updater for OWLCMS4 Apps
Updater for cloud installations of owlcms4 and publicresults

#### Rationale

OWLCMS applications such as [OWLCMS4](https://github.com/owlcms/owlcms4-heroku) and its [Public Results Relay](https://github.com/owlcms/publicresults-heroku) can be deployed to Heroku using a simple button.  But since hundreds of people could conceivably deploy that way, Heroku does not provide an automatic mechanism to subscribe to the source and update automatically.  And unfortunately, no "please update me" option either.  This program fills that gap.

#### Updating Existing OWLCMS Cloud Installations

> This program will update installations done using the  `Deploy to Heroku` button.

This program is downloaded to a user's workstation from the [Releases](https://github.com/jflamy/owlcms4-heroku-updater/releases/latest) page.  On Windows, you can simply double-click on the .exe.  On other platforms, the program is run from the command-line.

1. If the user has not used the `heroku` program or this updater on the machine before, a prompt for the Heroku username and password is given. The API token is stored locally so that subsequent updates do not require the password.
2. The program fetches the list of the user's Heroku applications and detects the ones that are for owlcms (currently, `owlcms4` and `publicresults`).
3. Each application is then updated, if needed, to the latest version available (prerelease applications are updated to the latest prerelease, stable applications to the latest stable.)  

![image](https://user-images.githubusercontent.com/678663/74204710-348c2480-4c6c-11ea-82d7-4908fabb296c.png)

#### Command-line options

By default, on Windows, the program opens a new command-line Window.  If you want to use these you will need to use `-createshell false`

| Option&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
| <nobr>`-createshell false`</nobr>                            | true or false. If true, open a new terminal window.  If false, run in the current command-line interface without opening a new window (useful for scripts).  Only works on Windows  (ignored on other platforms) |
| -apikey *keyvalue*                                           | Ignore the Heroku access token currently stored in the home directory `.netrc` (`_netrc` on Windows).  Use instead the token provided. The token can be obtained for a given user from the [User Account](https://dashboard.heroku.com/account) page and using the `Reveal`button at the right of the `API Key` section. |
| -app *appName*                                               | Used together with `-archive`, the name of a single application to be updated.  This is used to revert to a prior version in the event of a glitch. |
| `-archive` *tarball*                                         | The explicit URL of a .tar.gz file to be used to rebuild the application.  Such files are located in the Releases section of owlcms-heroku and publicresults-heroku. |
| `-force false`                                               | If true, ignore the version numbers and update to the latest available version. Used to work around a bug in the semantic versioning library that thinks that rc9 is bigger than rc10. |

