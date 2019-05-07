# Captive Portal

Simple application that allows to set up a wireless access point with a captive portal to log the credentials. It also provides a way to clone an existing page & use it as your own portal page.

## Dependencies
* go
* hostapd
* dnsmasq
* iptables
* macchanger
* wget _(Only required for page cloning)_
* uchardet _(Only required for page cloning)_

## Usage

To clone a web page:
```sh
./clone.sh <url> <directory>
```
Target page will be downloaded to specified directory within the ``site`` folder.
The script also injects special javascript into all html files which removes hidden form fields and transforms absolute form action URLs into relative URLs. 

To set up access point:
```sh
./portal.sh <evil_iface> <internet_iface> <ssid> <portal_directory> [index_file]
```
Where:
* evil_iface - wireless interface to be used by access point, e.g. ``wlan0``
* internet_iface - network interface used to provide access to the web after authentication, e.g. ``eth0``
* ssid - SSID of the network
* portal_directory - path to directory within the ``site`` folder where static files of portal web page are located
* index_file - file within ``portal_directory`` that will be used as index file, defaults to ``index.html``

Note that all ``POST`` requests to the portal regardless their URL will be treated as authentication attempts. Thus, any ``POST`` request without a form or with a form with empty fields is considered as failed attempt and is redirected back to index page. All the other requests are considered successful, and all form fields are logged into ``creds.txt`` file.
