# Cloud Updater for OWLCMS4 Apps
Updater for cloud installations of owlcms4 and publicresults

#### Rationale

OWLCMS applications such as [OWLCMS4](https://github.com/owlcms/owlcms4-heroku) and its [Public Results Relay](https://github.com/owlcms/publicresults-heroku) can be deployed to Heroku using a simple button.  But hundreds of people could conceivably deploy that way, there is no automation that would cause an update to the master to update the deployments.  Heroku does not provide an automatic "redeploy" button either.

#### Updating Existing OWLCMS Cloud Installations

> This program will update installations done using the  `Deploy to Heroku` button *since release 4.5.* For older installations, proceed as before (`heroku` command-line or uninstall/reinstall).

This program is downloaded to a user's workstation from the [Releases](https://github.com/jflamy/owlcms4-heroku-updater/releases/latest) page.  On Windows, you can double-click on the .exe.  On other platforms, the program is run from the command-line.

1. If the user has not used the `heroku` program or this updater on the machine before, a prompt for the Heroku username and password is given.
2. The program queries Heroku for the user's applications and detects the ones that are owlcms modules (currently, `owlcms4` and `publicresults`).  This works because the OWLCMS deployment button, starting with version 4.5, adds a configuration variable that defines where the program came from)
3. Each application is then updated to the latest version available from its source (prerelease applications are updated to the latest prerelease, stable applications to the latest stable.)  
4. The API token is stored locally so that subsequent updates do not require the password.

![image](https://user-images.githubusercontent.com/678663/74204710-348c2480-4c6c-11ea-82d7-4908fabb296c.png)

#### Command-line options

By default, on Windows, the program opens a new command-line Window.  Using these options requires that you also give `-createshell false`

| Option&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; | Description                                                  |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
| <nobr>`-createshell false`</nobr>                            | true or false. If true, open a new terminal window (default value on Windows).  If false, run in the current command-line interface without opening a new window (useful for scripts).  Only works on Windows  (ignored on other platforms) |
| -apikey *keyvalue*                                           | Ignore the Heroku access token currently stored in the home directory `.netrc` (`_netrc` on Windows).  Use instead the token provided. The token can be obtained for a given user from the [User Account](https://dashboard.heroku.com/account) page and using the `Reveal`button at the right of the `API Key` section. |
| -app *appName*                                               | Used together with `-archive`, the name of a single application to be updated.  This is used to revert to a prior version in the event of a glitch. |
| `-archive` *tarball*                                         | The explicit URL of a .tar.gz file to be used to rebuild the application.  These `.tar.gz` files are found in the Releases section of the [OWLCMS4](https://github.com/owlcms/owlcms4-heroku) and [Public Results Relay](https://github.com/owlcms/publicresults-heroku) repositories.  Use the files with `-heroku.tar.gz` in the name. |
|                                                              |                                                              |

