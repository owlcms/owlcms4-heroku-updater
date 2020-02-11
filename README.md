# Cloud Deployment Updater for OWLCMS4 Apps
Updater for cloud installations of owlcms4 and publicresults

#### Rationale

OWLCMS applications such as [OWLCMS4](https://github.com/owlcms/owlcms4-heroku) and its [Public Results Relay](https://github.com/owlcms/publicresults-heroku) can be deployed to Heroku using a simple button.  But hundreds of people could conceivably deploy that way, there is no automation that would cause an update to the master to update the deployments.  Heroku does not provide an automatic "redeploy" button either.

#### On-demand Updating

This program is downloaded to a user's workstation. 

1. If the user has not used the `heroku` program or this updater on the machine before, a prompt for the Heroku username and password is given.
2. The program queries Heroku for the user's applications and detects the ones that are owlcms modules (currently, owlcms4 and publicresults).  This works because the deployment button adds a configuration variable that defines where the program came from)
3. Each application is then updated to the latest version available from its source (prerelease applications are updated to the latest prerelease, stable applications to the latest stable.)  
4. The API token is stored locally so that subsequent updates do not require the password.

![image](https://user-images.githubusercontent.com/678663/74204710-348c2480-4c6c-11ea-82d7-4908fabb296c.png)

#### Command-line options

By default, on Windows, the program opens a new command-line Window.  Using these options requires that you also give `-createshell false`

| Option&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;                               | Description                                                  |
| ---------------------------------------- | ------------------------------------------------------------ |
| <nobr>`-createshell false`</nobr> | true or false.  If false, open a new terminal window.  Only works on Windows  (ignored on other platforms) |
| -apikey *keyvalue*                       | Ignore the Heroku access token currently stored in the home directory `.netrc` (`_netrc` on Windows).  Use instead the provided token. The token can be obtained for a given user from the [User Account](https://dashboard.heroku.com/account) page. |
| -app *appName*                           | Used together with `-archive`, the name of a single application to be updated.  This is used to revert to a prior version in the event of a glitch. |
| `-archive` *tarball*                     | The explicit URL of a .tar.gz file to be used to rebuild the application.  These `.tar.gz` files are found in the Releases section of the [OWLCMS4](https://github.com/owlcms/owlcms4-heroku) and [Public Results Relay](https://github.com/owlcms/publicresults-heroku) repositories.  Use the files with `-heroku.tar.gz` in the name. |
|                                          |                                                              |

