# owlcms4-heroku-updater
Updater for cloud installations of owlcms4 and publicresults

This program is downloaded to a user's workstation.  It asks Heroku for the list of applications owned by the user, and
detects the ones that are owlcms modules (currently, owlcms4 and publicresults). Each application is then updated to the latest
version available from its source (prerelease applications are updated to the latest prerelease, stable applications to the latest
stable.

