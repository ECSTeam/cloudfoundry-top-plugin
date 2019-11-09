# top-plugin by Kurt Kellner of CGI
https://www.cgi.com/en/media/brochure/cloud-native-solutions

This is a Cloud Foundry command-line cf interactive plugin for showing live statistics of the targeted Cloud Foundry foundation.
The live statistics include application statistics and route statistics among others.
The primary source of information that the top plugin uses is via monitoring the Cloud Foundry firehose.

[Cloud Foundry Summit 2017 session about TOP](https://www.youtube.com/watch?v=XDY64HKB7CI&t=7m48s)

The plugin will run in one of two modes, privileged or non-privileged depending on your Cloud Foundry user permission.
If you are a foundation operator you will want to use top in privileged mode.  This is done automatically if the
correct permissions are granted to your Cloud Foundry login (or if you are logged in via `admin` account).  See
[Assign Permissions](#assign-permissions-if-privileged-mode-is-needed) for more information on assigning permissions.


[Installation Instructions](#installation) 

![Screenshot](screenshots/screencast2.gif?raw=true)

## Screenshots

More [screenshots here](screenshots/screenshots.md)

# Usage Documentation

After installation be sure to view the [full `top` documentation](docs/doc.md) as
well as the [Frequently Asked Questions (FAQ)](docs/faq.md) page.

# Installation
There are two options for installation; use the plugin repo (recommended) or manual installation.

## Install from plugin repository (recommended)
NOTE: This installation method requires that your client computer has access to the internet.
If internet access is not available from client computer use the manual method.

Verify you have a repo named `CF-Community` registered in your cf client.

```
cf list-plugin-repos
```
If the above command does not show `CF-Community` you can add the repo via:

```
cf add-plugin-repo CF-Community http://plugins.cloudfoundry.org/
```
Now that we have the cloud foundry community repo registered, install `top`:

```
cf install-plugin -r CF-Community "top"
```


## Manual installation method
* **Download the binary file for your target OS from the latest [release](https://github.com/ecsteam/cloudfoundry-top-plugin/releases/latest)**
* If you've already installed the plugin and are updating, you must first run `cf uninstall-plugin top`
* Then install the plugin with `cf install-plugin top-plugin-darwin`  (or `top-plugin-linux` or `top-plugin.exe`)
* If you get a permission error run: `chmod +x top-plugin-darwin` (or `top-plugin-linux`) on the binary
* Verify the plugin installed by looking for it with `cf plugins`

## Upgrade to latest version
To upgrade to the lastest version of top plugin, uninstall plugin and install again.
```
cf uninstall-plugin top
cf install-plugin -r CF-Community "top"     (or use manual install method described above)
```

## Assign scope if privileged mode is needed

The `top` plugin will run without special scope (permissions) however it determines at runtime
what scopes you have and displays the appropriate functionality based on those
scopes.  If you are a foundation operator you will want the additional functionality
that top provides to privileged users.

If you are logged in with the Cloud Foundry `admin` account, no additional scopes
are needed, the `admin` account has everything it needs to run top with full functionality.

For non-admin accounts, to run top in privileged mode you need to assign two scopes
to an existing Cloud Foundry user (or LDAP group).  To assign needed scopes:

Install the uaac client CLI if you do not already have it:
```
gem install cf-uaac
```

Login and add 2 or 3 scopes as showed below.  Note that the UAA password is NOT the
"Admin Credentials", the password is found in the PAS tile under Credentials tab,
look for password for "Admin Client Credentials".

```
uaac target https://login.system.YOUR.DOMAIN --skip-ssl-validation
uaac token client get admin -s [UAA Admin Client Credentials]  
```

### To assign scopes to a LDAP group (recommended if connected to LDAP/Active Directory).
Read-only admin API is all that is needed for this plugin but both options given below.

Read-only API access:
```
uaac group map --name cloud_controller.admin_read_only [FULL DN to LDAP group]
uaac group map --name scim.read  [FULL DN to LDAP group]
uaac group map --name doppler.firehose [FULL DN to LDAP group]
```
-or-

Full API access:
```
uaac group map --name cloud_controller.admin [FULL DN to LDAP group]
uaac group map --name doppler.firehose [FULL DN to LDAP group]
```

### To assign scopes directly to a user. 
Read-only admin API is all that is needed for this plugin but both options given below.

Read-only API access:
```
uaac member add cloud_controller.admin_read_only [username]
uaac member add scim.read [username]
uaac member add doppler.firehose [username]
```
-or-

Full API access:
```
uaac member add cloud_controller.admin [username]
uaac member add doppler.firehose [username]
```

Note: The change in permissions does not take effect until user username performs
a logout and login.


# Usage

Although `top` does not *require* any special permissions, foundation operators 
will want to run `top` in privileged mode as described in the
[Assign permissions](#Assign-permissions-if-privileged-mode-is-needed)
section above.  The plugin does not require arguments.  Simply run:
```
cf top
```

## Options

List top live statistics for CF.

```
NAME:
   top - Displays top stats - by Kurt Kellner of ECS Team (now part of CGI)

USAGE:
   cf top

OPTIONS:
   -debug              -d, enable debugging
   -no-top-check       -ntc, do not check if there are other instances of top running
   -nozzles            -n, specify the number of nozzle instances (default: 2)
   -cygwin             -c, force run under cygwin (Use this to run: 'cmd /c start cf top -cygwin' )
```
